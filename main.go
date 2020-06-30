package main

//go:generate swagger generate server -A IsakuraManager -f ./swagger.yaml --exclude-main
//go:generate bash build_ui.sh
//go:generate statik -f -src=./ui/dist

import (
  "flag"
  "time"
  "path/filepath"
  "os"
  "log"
  "strings"
  "net/http"
  "io/ioutil"
  "io"
  "strconv"
  "regexp"

  "github.com/go-openapi/runtime"
  "github.com/ansoni/isakura-manager/isakura"
  "github.com/ansoni/isakura-manager/restapi"
  "github.com/ansoni/isakura-manager/models"
  "github.com/ansoni/isakura-manager/restapi/operations"
  "github.com/go-openapi/strfmt" 
  _ "github.com/ansoni/isakura-manager/statik"
  fs "github.com/rakyll/statik/fs"

  loads "github.com/go-openapi/loads"
  middleware "github.com/go-openapi/runtime/middleware"
)

type flushWriter struct {
  f http.Flusher
  w io.Writer
}

func (fw *flushWriter) Write(p []byte) (n int, err error) {
  n, err = fw.w.Write(p)
  if fw.f != nil {
    fw.f.Flush()
  }
  return
}


func main() {
  userHomeDir, _ := os.UserHomeDir()
  home := flag.String("home", filepath.Join(userHomeDir, ".isakura"), "directory where we save our data") 
  authHostUrl := flag.String("auth_host_url", "http://webtv.jptvpro.net:9001", "url for authentication") 
  apiHostUrl := flag.String("api_host_url", "http://livepro.fjp2017.com:9083", "url for host") 
  userAgent := flag.String("user_agent", "Mozilla/5.0 (iPhone; U; CPU iPhone OS 4_3_3 like Mac OS X; en-us) AppleWebKit/533.17.9 (KHTML, like Gecko) Version/5.0.2 Mobile/8J2 Safari/6533.18.5", "user-agent for requests") 
  username := flag.String("username", "", "username for logging in")
  password := flag.String("password", "", "password for logging in")
  flag.Parse()

  // Initialize isakura
  isakura_client := isakura.Isakura{Home: *home, AuthHostUrl: *authHostUrl, ApiHostUrl: *apiHostUrl, UserAgent: *userAgent, Username: *username, Password: *password}

  isakura_client.Load()

  //fmt.Printf("ISakura Manager Config\n#############\n%+v\n#############\n", isakura_client)

  downloads := make(chan isakura.Download)

  // Authenticate and maintain a session
  go isakura_client.MaintainSession()

  // give 1 second to startup
  time.Sleep(time.Second) 

  // Maintain a list of our channels. Refresh hourly
  go isakura_client.MaintainChannelGuide(downloads)

  // Monitor our download channel and retrieve
  go isakura_client.Download(downloads)

  // Serve up our webpage
  swaggerSpec, err := loads.Embedded(restapi.SwaggerJSON, restapi.FlatSwaggerJSON)
  if err != nil {
    log.Fatalln(err)
  }

  api := operations.NewIsakuraManagerAPI(swaggerSpec)

  api.GetRootRedirectHandler = operations.GetRootRedirectHandlerFunc(func(params operations.GetRootRedirectParams) middleware.Responder {
      return middleware.ResponderFunc(
        func(rw http.ResponseWriter, pr runtime.Producer) {
          rw.Header().Add("Location", "/ui/index.html")
          rw.WriteHeader(302) 
        })
    })

  statikFS, err := fs.New()
  if err != nil {
    log.Fatal(err)
  }
  api.GetUIContentHandler = operations.GetUIContentHandlerFunc(func(params operations.GetUIContentParams) middleware.Responder {
    resource := params.Resource
    log.Printf("Read %v\n", resource)
    return middleware.ResponderFunc(
        func(rw http.ResponseWriter, pr runtime.Producer) {
          resourceFile, err := statikFS.Open("/" + resource)
          if err != nil {
            rw.WriteHeader(404)
            return
          }

          if strings.HasSuffix(resource, "js") {
            rw.Header().Add("Content-Type", "application/javascript")
          } else if strings.HasSuffix(resource, "css") {
            rw.Header().Add("Content-Type", "text/css")
          } else {
            rw.Header().Add("Content-Type", "text/html")
          }
          buf, err := ioutil.ReadAll(resourceFile) 
          if err != nil {
            rw.WriteHeader(404)
            return
          }
          rw.Write([]byte(buf)) 
      })
  })

  api.GetChannelsHandler = operations.GetChannelsHandlerFunc(func(params operations.GetChannelsParams) middleware.Responder {
    channels := make([]*models.Channel, 0)
    for _, channel := range isakura_client.Channels.Channels {
      channels = append(channels, &models.Channel{ChannelName: channel.Name, BroadcastType: "BS"})
    }

    return operations.NewGetChannelsOK().WithPayload(channels)
  })

  api.GetChannelsGuideHandler = operations.GetChannelsGuideHandlerFunc(func(params operations.GetChannelsGuideParams) middleware.Responder {
    channels := make([]*models.ChannelGuide, 0)
    search := *params.Search
    var regex_err error
    var regex *regexp.Regexp
    if params.Search != nil {
      regex, regex_err = regexp.Compile(search)
      if regex_err != nil { 
        log.Printf("%v did not regex compile %v", search, regex_err)
      } else {
        log.Printf("%v did regex compile %+v", search, regex)
      }
    }
    for _, channel := range isakura_client.Channels.Channels {
      guide := make([]*models.Guide, 0)
      for _, program := range channel.Programs {
        if params.Search == nil || search == "" || strings.Contains(program.Title, search) || (regex_err == nil && regex.MatchString(program.Title)) {
          airdate := time.Unix(program.Time, 0)
          guide = append(guide, &models.Guide{ID: program.Time, Name: program.Title, Airdate: strfmt.DateTime(airdate)})
        }
      }

      if len(guide) > 0 {
        channels = append(channels, &models.ChannelGuide{ChannelName: channel.Name, Guide: guide})
      }
    }
    return operations.NewGetChannelsGuideOK().WithPayload(channels)
  })

  api.GetContentHandler = operations.GetContentHandlerFunc(func(params operations.GetContentParams) middleware.Responder {
    contents := make([]*models.Content, 0)
    for _, download := range isakura_client.Downloads {
      contents = append(contents, &models.Content{Name: download.ID, LocalPath: download.LocalPath, DownloadDate: strfmt.DateTime(download.DownloadDate)})
    } 
    return operations.NewGetContentOK().WithPayload(contents)
  })

  api.GetContentPreviewHandler = operations.GetContentPreviewHandlerFunc(func(params operations.GetContentPreviewParams) middleware.Responder {
    requestedChannel := params.Channel
    requestedContent,_ := strconv.ParseInt(params.Content, 10, 64)
    log.Printf("Reqeust: %+v", params.HTTPRequest)
    for _, channel := range isakura_client.Channels.Channels {
      if requestedChannel == channel.Name { 
        for _, program := range channel.Programs {
          if requestedContent == program.Time {
            return middleware.ResponderFunc(func(rw http.ResponseWriter, p runtime.Producer) {
              f, _ := rw.(http.Flusher)
              rw.WriteHeader(200)
              _ = isakura_client.Preview(isakura.Download{Channel: channel, Program: program}, params.HTTPRequest, &flushWriter{f: f, w: rw})
            })
          }
        }
      }
    }
    return operations.NewGetContentPreviewNotFound()
  })

  api.GetContentFoldersHandler = operations.GetContentFoldersHandlerFunc(func(params operations.GetContentFoldersParams) middleware.Responder {
    folders := []string{}
    for _, monitor := range isakura_client.Monitors {
      monitorFolder := monitor.Folder
      for _, folder := range folders {
        if folder == monitorFolder {
          continue 
        }
      }
      folders = append(folders, monitorFolder)
    }
    return operations.NewGetContentFoldersOK().WithPayload(folders)
  })

  api.DeleteContentHandler = operations.DeleteContentHandlerFunc(func(params operations.DeleteContentParams) middleware.Responder {
    contentName := params.ContentName
    err := isakura_client.Delete(contentName)
    if err != nil {
      log.Printf("%v", err)
      return operations.NewDeleteScheduleNotFound()
    }
    return operations.NewDeleteScheduleOK()
  })

  api.DeleteScheduleHandler = operations.DeleteScheduleHandlerFunc(func(params operations.DeleteScheduleParams) middleware.Responder{
    scheduleName := params.ScheduleName
    for i, monitor := range isakura_client.Monitors {
      if monitor.Name == scheduleName {
        isakura_client.Monitors = append(isakura_client.Monitors[:i], isakura_client.Monitors[i+1:]...)
        isakura_client.Save()
        return operations.NewDeleteScheduleOK()
      }
    }
    return operations.NewDeleteScheduleOK()
  })

  api.GetSchedulesHandler = operations.GetSchedulesHandlerFunc(func(params operations.GetSchedulesParams) middleware.Responder {
    schedules := make([]*models.Schedule, 0)
    for _, monitor := range isakura_client.Monitors {
      schedule := models.Schedule{Name: monitor.Name, PrependDate: monitor.PrependDate, PrependTime: monitor.PrependTime, AppendDate: monitor.AppendDate, AppendTime: monitor.AppendTime, Watch: monitor.Watch, Folder: monitor.Folder, Filter: monitor.Pattern}
      for _, search := range monitor.Modifiers {
        schedule.Searches = append(schedule.Searches, &models.ScheduleSearchesItems0{Search: search.Search, Replace: search.Replace})
      }
      schedules = append(schedules, &schedule)
    }
    return operations.NewGetSchedulesOK().WithPayload(schedules)
  })
  api.CreateScheduleHandler = operations.CreateScheduleHandlerFunc(func(params operations.CreateScheduleParams) middleware.Responder {
    log.Printf("%+v", params.Body)
    schedule := params.Body
    monitor := isakura.Monitor{ Name: schedule.Name, Pattern: schedule.Filter, PrependDate: schedule.PrependDate, PrependTime: schedule.PrependTime, AppendDate: schedule.AppendDate, AppendTime: schedule.AppendTime, Folder: schedule.Folder, Watch: schedule.Watch}
    for _, search := range schedule.Searches {
      monitor.Modifiers = append(monitor.Modifiers, isakura.MonitorModifiers{ Search: search.Search, Replace: search.Replace})
    }
    isakura_client.Monitors = append(isakura_client.Monitors, monitor)
    isakura_client.Save()
    return operations.NewCreateScheduleOK()
  })

  api.GetChannelGuideHandler = operations.GetChannelGuideHandlerFunc(func(params operations.GetChannelGuideParams) middleware.Responder {
    requestedChannel := params.Channel
    search := params.Search
    var regex_err error
    var regex *regexp.Regexp
    if search !=nil {
      regex, regex_err = regexp.Compile(*search)
      if regex_err != nil { 
        log.Printf("%v did not regex compile %v", search, regex_err)
      }
    }
    guide := make([]*models.Guide, 0)
    for _, channel := range isakura_client.Channels.Channels {
      if requestedChannel == channel.Name {
        for _, program := range channel.Programs {
          if search == nil || *search == "" || strings.Contains(program.Title, *search) || (regex_err == nil && regex.Match([]byte(program.Title))) {
            airdate := time.Unix(program.Time, 0)
            guide = append(guide, &models.Guide{ID: program.Time, Name: program.Title, Airdate: strfmt.DateTime(airdate)})
          }
        }
        channel := &models.ChannelGuide{ChannelName: channel.Name, Guide: guide}
        return operations.NewGetChannelGuideOK().WithPayload(channel)
      }
    }
    return operations.NewGetChannelGuideNotFound()
  })

  server := restapi.NewServer(api)
  defer server.Shutdown()

  server.ConfigureAPI()
  server.Host = "0.0.0.0"
  server.EnabledListeners = []string{"http",}
  server.Port = 8080

  if err := server.Serve(); err != nil {
    log.Fatalln(err)
  }
}
