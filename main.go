package main

//go:generate swagger generate server -A IsakuraManager -f ./swagger.yaml --exclude-main

import (
  "fmt"
  "flag"
  "net/http"
  "net/url"
  "time"
  "io"
  "io/ioutil"
  "encoding/json"
  _ "reflect"
  "strconv"
  "path/filepath"
  "os"
  "sync"
  "regexp"
  "log"

  "github.com/ansoni/isakura-manager/restapi"
  "github.com/ansoni/isakura-manager/restapi/operations"
  loads "github.com/go-openapi/loads"
  flags "github.com/jessevdk/go-flags"
)

var saveMutex sync.Mutex

type ProductConfig struct {
  VmsHost string `json:"vms_host"`
  VmsRtmfpHost string `json:"vms_rtmfp_host"`
  VmsUid string `json:"vms_uid"`
  VmsLiveCid string `json:"vms_live_cid"`
  VmsVodHost string `json:"vms_vod_host"`
  VmsVodRtmfp string `json:"vms_vod_rtmfp"`
  VmsMarqueeText string `json:"vms_marquee_text"`
  Web string `json:"web"`
  VmsReferer string `json:"vms_referer"`
  Single string `json:"single"`
}

type ProductConfigWrapper struct {
  ProductConfig
}

func (pc *ProductConfig) UnmarshalJSON(data []byte) error {
  type PCAlias ProductConfig
  var pcAlias PCAlias
  jsonData, err := strconv.Unquote(string(data)); 
  if err != nil {
    jsonData = string(data) // try to load the data without unquote!
  } 
  if err := json.Unmarshal([]byte(jsonData), &pcAlias); err != nil {
    panic(err)
  }
  *pc = ProductConfig(pcAlias)
  return nil
}

type ISakuraSession struct {
  AccessToken string `json:"access_token"`
  RefreshToken string `json:"refresh_token"`
  Expired bool `json:"expired"`
  Disabled bool `json:"disabled"`
  Confirmed bool`json:"confirmed"`
  Cid string `json:"cid"`
  Type string `json:"type"`
  Trial int `json:"trial"`
  CreateTime int64 `json:"create_time"`
  ExpireTime int64 `json:"expire_time"`
  ProductConfig *ProductConfig `json:"product_config"`
  ServerTime int64 `json:"server_time"`
  Code string `json:"code"`
}

type ISakuraSaveData struct {
  Session ISakuraSession
  Channels Channels
  Monitors []Monitor
  Downloads []LocalProgram
}

type Monitor struct {
  Pattern string `json:"pattern"`
}

type isakura struct {
  Home string
  AuthHostUrl string
  ApiHostUrl string
  UserAgent string
  Username string
  Password string
  ISakuraSaveData
};

type Channels struct {
  Code string `json:"code"`
  PageCount int `json:"page_count"`
  Channels []Channel `json:"result"`
}

type ChannelHeader struct {
  Name string `json:"name"`
  Description string `json:"description"`
  Playpath string `json:"playpath"`
}

type ChannelPrograms struct {
  Programs []Program `json:"record_epg"`
}

/**
 * Used for de-serializing embedded JSON string
 */
type ChannelProgramsPrivate struct {
  Programs string `json:"record_epg"`
}

type Channel struct {
  ChannelHeader
  ChannelPrograms
}

type Program struct {
  Time int64 `json:"time"`
  Title string `json:"title"`
  Path string `json:"path"`
};

type ProgramQuery struct {
  Name string `json:"name"`
  Description string `json:"description"`
  ScTk string `json:"sc_tk"`
  Substreams []ProgramSubstream `json:"substreams"`
}

type ProgramSubstream struct {
  Path string `json:"path"`
  Width int `json:"width"`
  Height int `json:"height"`
  EncryptCode string `json:"encrypt_code"`
  HttpUrl string `json:"http_url"`
}

type Download struct {
  Channel Channel
  Program Program
}

type LocalProgram struct {
  Program
  LocalPath string `json:"local_path"`
  Monitor Monitor `json:"monitor"`
  DownloadDate time.Time `json:"download_date"`
}

func (c *Channel) UnmarshalJSON(data []byte) error {
  var cHeader ChannelHeader
  var channel Channel
  *c = channel
  if err := json.Unmarshal([]byte(data), &cHeader); err != nil {
    panic(err)
  }
  c.ChannelHeader = cHeader

  var cProgramsPrivate ChannelProgramsPrivate
  fmt.Printf("Unmarshall %v channel bytes \n", len(data))
  if err := json.Unmarshal([]byte(data), &cProgramsPrivate); err != nil {
    // not a string... see if we can parse it normally
    //fmt.Printf("%v\n", err)
    var cPrograms ChannelPrograms
    err = json.Unmarshal([]byte(data), &cPrograms)
    if err != nil { // yikes!
      panic(err)
    }
    c.ChannelPrograms = cPrograms
  
  } else {
    //fmt.Printf("Data Read: %v\n", cProgramsPrivate.Programs)
    //jsonData, err := strconv.Unquote(cProgramsPrivate.Programs); 
    //if err != nil {
    //  fmt.Printf("No Programs read for %v, Error: %v\n", c.Name, err)
    //} else {
      var programs []Program
      err = json.Unmarshal([]byte(cProgramsPrivate.Programs), &programs)
      if err != nil {
        panic(err)
      }
      //fmt.Printf("%v\n", programs)
       
      c.Programs = programs
    //}
  }
  return nil
}

func serve(downloads chan Download) {
  

}

func (isakura *isakura) refresh() error {
  return fmt.Errorf("Unsupported")
}

func (isakura *isakura) login() error {
  auth_host_url := fmt.Sprintf("%s/logon.sjs", isakura.AuthHostUrl)
  resp, err := http.PostForm(auth_host_url, url.Values{"device_id": {"3079898AA6E3EA6D0999F4DE748B5595"}, "cid": {isakura.Username}, "password": {isakura.Password}, "redirect_url": {"http://webtv.jptvpro.net/play.html"}})
  if err != nil {
    panic(err)
  }
  body, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    panic(err)
  }
  err = json.Unmarshal(body, &isakura.Session) 
  if err != nil {
    panic(err)
  }
  //fmt.Printf("It Worked! isakura initialized!\n %+v\n", isakura)
  return nil
}

func (isakura *isakura) maintainSession() {
  if isakura.Session.AccessToken == "" {
    isakura.login()  
  } else {
    fmt.Printf("Using existing Session\n")
  }

  isakura.save()
  refreshInXSeconds := isakura.Session.ExpireTime-isakura.Session.ServerTime-300
  fmt.Printf("Forced Token Refresh in %v seconds\n", refreshInXSeconds)
  backoff := 1
  for {
    refreshInXSeconds := isakura.Session.ExpireTime-isakura.Session.ServerTime-300
    if refreshInXSeconds < 0 || isakura.Session.Expired {
      err := isakura.refresh() // try to fresh
      if err != nil {
        err = isakura.login() // can't refresh... just re-login
        if err != nil {
          if backoff < 60 {
            backoff += 1 // if we can't login.. we should backoff
          }
        } else {
          backoff = 1
        }
      }
    }
    time.Sleep(time.Second * time.Duration(30 * backoff)) 
  }
}

func (isakura *isakura) maintainChannelGuide(downloads chan Download) {
  for {
    fmt.Printf("Download Channels\n")
    isakura.retrieveChannelGuide()
    fmt.Printf("Downloaded Channels\n")
    isakura.save()
    fmt.Printf("Scan for new Content\n")
    isakura.scan(downloads)
    fmt.Printf("Scanned for new Content\n")
    time.Sleep(time.Second * time.Duration(60 * 60))
  }
}

func (isakura *isakura) scan(downloads chan Download) {
  for _, monitor := range isakura.Monitors {
    fmt.Printf("Scan Guide for %v\n", monitor.Pattern)
    r, err := regexp.Compile(monitor.Pattern)
    if (err != nil) {
      fmt.Printf("Regex Error:-( %v\n", err)
      continue 
    }
    // TODO: FIX Channels.Channels
    for _, channel := range isakura.Channels.Channels {
      fmt.Printf("Scan Channel %v\n", channel.Name)
      for _, program := range channel.Programs {
        if r.MatchString(program.Title) {
          downloads <- Download{Program: program, Channel: channel}
          fmt.Printf("%v matches\n", program.Title)
        }
      } 
    }  
  }
}

func (isakura *isakura) retrieveChannelGuide() {
  for isakura.Session.AccessToken == "" { // we are not ready!
    time.Sleep(time.Second)
  }  
  v := url.Values{ "action": {"listLives"}, "cid": {"2E2FAA0BF6E84FE0C34955CA0DFB6AAD"}, "details": {"1"}, "page_size": {"200"}, "sort": {"created_time desc"}, "type": {"video"}, "uid": {"C2D9261F3D5753E74E97EB28FE2D8B26"}, "referer": {"http://isakura.tv"}}
  auth_host_url := fmt.Sprintf("%s/api?%v", isakura.Session.ProductConfig.VmsHost, v.Encode())
  resp, err := http.Get(auth_host_url)
  if err != nil {
    panic(err)
  }
  body, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    panic(err)
  }
  err = json.Unmarshal(body, &isakura.Channels) 
  if err != nil {
    panic(err)
  }
}

func (isakura *isakura) download(downloads chan Download) {
  for {
    download := <- downloads
    savePath := filepath.Join(isakura.Home, "videos", fmt.Sprintf("%v._flv", download.Program.Title) )
    finalPath := filepath.Join(isakura.Home, "videos", fmt.Sprintf("%v.flv", download.Program.Title) )

    fmt.Printf("Download: %v from %v\n", download.Program.Title, download.Channel.Name)
    if download.Program.Path == "" {
       fmt.Printf("Empty Playpath, can't download\n")
       continue
    }

    // Do not to forward if File exists
    if _, err := os.Stat(finalPath); ! os.IsNotExist(err) {
       fmt.Printf("Not Downloading since we already have it -> %v\n", finalPath)
       continue
    }

    pq, err := isakura.query(download)
    if err != nil {
      fmt.Printf("Error Querying Package: %v, hopefully we get it on the next scan\n", err)
      continue
    }

    v := url.Values{ "__download": {"1"}, "sc_tk": {pq.ScTk} }
    downloadUrl := fmt.Sprintf("%v/video.fpvsegments?%v", pq.Substreams[0].HttpUrl, v.Encode())
    fmt.Printf("Downlaod Url: %v\n", downloadUrl)
    resp, err := http.Get(downloadUrl)
    if  err != nil {
      fmt.Printf("Error - %v\n", err)
    }
    defer resp.Body.Close()
    
    out, err := os.Create(savePath)
    if err != nil {
      panic(err)
    }
    defer out.Close()
    io.Copy(out, resp.Body)
    os.Rename(savePath, finalPath)
  }
}

func (isakura *isakura) query(download Download) (*ProgramQuery, error) {
  v := url.Values{ "type": {"vod"}, "access_token": {isakura.Session.AccessToken}, "pageUrl": {"http://webtv.jptvpro.net/isakura/main.html"}, "refresh_token": {isakura.Session.RefreshToken}, "expires_in": {fmt.Sprintf("%v",isakura.Session.ExpireTime)}}
  queryPath := fmt.Sprintf("%s%s.json?%v", isakura.Session.ProductConfig.VmsHost, download.Program.Path, v.Encode())
  fmt.Printf("Download Package Info using: \n%v\n", queryPath)
  resp, err := http.Get(queryPath)
  if err != nil {
    return nil, err
  }
  body, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    panic(err)
  }
  fmt.Printf("ProgramQuery: %v\n", string(body))
  var programQuery ProgramQuery
  err = json.Unmarshal(body, &programQuery) 
  if err != nil {
    panic(err)
  }
  fmt.Printf("ProgramQuery: %+v\n", programQuery)
  return &programQuery, nil

}

func (isakura *isakura) serve(downloads chan Download) {
  for {
    time.Sleep(time.Second * 30)
  }

	swaggerSpec, err := loads.Embedded(restapi.SwaggerJSON, restapi.FlatSwaggerJSON)
	if err != nil {
		log.Fatalln(err)
	}

	api := operations.NewIsakuraManagerAPI(swaggerSpec)
	server := restapi.NewServer(api)
	defer server.Shutdown()

	parser := flags.NewParser(server, flags.Default)
	parser.ShortDescription = "isakura-manager"
	parser.LongDescription = "Isakura-manager"

	server.ConfigureFlags()
	for _, optsGroup := range api.CommandLineOptionsGroups {
		_, err := parser.AddGroup(optsGroup.ShortDescription, optsGroup.LongDescription, optsGroup.Options)
		if err != nil {
			log.Fatalln(err)
		}
	}

	if _, err := parser.Parse(); err != nil {
		code := 1
		if fe, ok := err.(*flags.Error); ok {
			if fe.Type == flags.ErrHelp {
				code = 0
			}
		}
		os.Exit(code)
	}

	server.ConfigureAPI()

	if err := server.Serve(); err != nil {
		log.Fatalln(err)
	}
}

func (isakura *isakura) load() {
  if _, err := os.Stat(isakura.Home); os.IsNotExist(err) {
    fmt.Printf("No existing Configuration found in %v\n", isakura.Home)
    return
  }
  savePath := filepath.Join(isakura.Home, "config")
  data, err := ioutil.ReadFile(savePath)
  if err != nil {
    fmt.Printf("Unable to read Configuration %v, Error: %v\n", savePath, err)
  }
  err = json.Unmarshal(data, &isakura.ISakuraSaveData)
  if err != nil {
    fmt.Printf("Unable to read Configuration %v, Error: %v\n", savePath, err)
  }
}

func (isakura *isakura) save() {
  saveMutex.Lock()
  defer saveMutex.Unlock()
  videosPath := filepath.Join(isakura.Home, "videos")
  if _, err := os.Stat(isakura.Home); os.IsNotExist(err) {
    err = os.Mkdir(isakura.Home, 0700)
    if err != nil {
      panic(err)
    }
    err = os.Mkdir(videosPath, 0700)
    if err != nil {
      panic(err)
    }
  }
  savePath := filepath.Join(isakura.Home, "config.tmp")
  finalPath := filepath.Join(isakura.Home, "config")
  fmt.Printf("Saving Configuration to %v\n", savePath)
  b, err := json.MarshalIndent(isakura.ISakuraSaveData, "", "\t")
  if err != nil {
    fmt.Println("error saving:", err)
  }
  err = ioutil.WriteFile(savePath, b, 0600)
  if err != nil {
    panic(fmt.Sprintf("Could not save to %v, Error: %v!\n", savePath, err))
  }
  os.Rename(savePath, finalPath)
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
  isakura_client := isakura{Home: *home, AuthHostUrl: *authHostUrl, ApiHostUrl: *apiHostUrl, UserAgent: *userAgent, Username: *username, Password: *password}

  if len(isakura_client.Monitors) == 0 {
    isakura_client.Monitors = append(isakura_client.Monitors, Monitor{ Pattern: "ピタゴラスイッチ"})
  }

  isakura_client.load()

  //fmt.Printf("ISakura Manager Config\n#############\n%+v\n#############\n", isakura_client)

  downloads := make(chan Download)

  // Authenticate and maintain a session
  go isakura_client.maintainSession()

  // give 1 second to startup
  time.Sleep(time.Second) 

  // Maintain a list of our channels. Refresh hourly
  go isakura_client.maintainChannelGuide(downloads)

  // Monitor our download channel and retrieve
  go isakura_client.download(downloads)

  // Serve up our webpage
  isakura_client.serve(downloads)
  
}
