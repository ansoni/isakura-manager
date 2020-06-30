package isakura


import (
 "fmt"
  "net/http"
  "net/url"
  "time"
  "io"
  "io/ioutil"
  "encoding/json"
  "path/filepath"
  "os"
  "sync"
  "regexp"
  "log"
  "strings"
  uuid "github.com/satori/go.uuid"
)

var saveMutex sync.Mutex

func (isakura *Isakura) refresh() error {
  return fmt.Errorf("Unsupported")
}

func (isakura *Isakura) login() error {
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

func (isakura *Isakura) MaintainSession() {
  if isakura.Session.AccessToken == "" {
    isakura.login()
  } else {
    fmt.Printf("Using existing Session\n")
  }

  isakura.Save()
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

func (isakura *Isakura) MaintainChannelGuide(downloads chan Download) {
  for {
    fmt.Printf("Download Channels\n")
    isakura.retrieveChannelGuide()
    fmt.Printf("Downloaded Channels\n")
    isakura.Save()
    fmt.Printf("Scan for new Content\n")
    isakura.scan(downloads)
    fmt.Printf("Scanned for new Content\n")
    time.Sleep(time.Second * time.Duration(60 * 60))
  }
}

func (isakura *Isakura) scan(downloads chan Download) {
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
        if ! program.Downloaded {
          if r.MatchString(program.Title) {
            fmt.Printf("%v matches\n", program.Title)
            downloads <- Download{Program: program, Channel: channel, Monitor: monitor}
          }
        }
      }
    }
  }
}

func (isakura *Isakura) retrieveChannelGuide() {
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
  //fmt.Printf(string(body))
  if err != nil {
    panic(err)
  }
  err = json.Unmarshal(body, &isakura.Channels)
  if err != nil {
    panic(err)
  }

  // invert channel listing (oldie is goldie!)
  channels := isakura.Channels.Channels
  for i := len(channels)/2-1; i >= 0; i-- {
    opp := len(channels)-1-i
    channels[i], channels[opp] = channels[opp], channels[i]
  }
}

// Stream a preview
func (isakura *Isakura) Preview(download Download, r *http.Request, w io.Writer) error {
  fmt.Printf("Download: %v from %v\n", download.Program.Title, download.Channel.Name)
  if download.Program.Path == "" {
    return fmt.Errorf("Empty Playpath, can't download\n")
  }
  pq, err := isakura.query(download)
  if err != nil {
    return fmt.Errorf("Error Querying Package: %v, hopefully we get it on the next scan\n", err)
  }

  v := url.Values{ "__download": {"1"}, "sc_tk": {pq.ScTk} }
  downloadUrl := fmt.Sprintf("%v/video.fpvsegments?%v", pq.Substreams[0].HttpUrl, v.Encode())
  fmt.Printf("Download Url: %v\n", downloadUrl)
  req, _ := http.NewRequest("GET", downloadUrl, nil) 
  req.Header.Add("Content-Type", "video/x-flv")
  if r.Header.Get("Range") != "" {
    log.Printf("Got Range Header")
    req.Header.Set("Range", r.Header.Get("Range"))
  }
  //req.Header.Add("Range", "bytes=0-5000000")
  //req.Header.Add("Range", "bytes=5000000-10000000")
  //req.Header.Add("Range", "bytes=10000000-15000000")
  var client http.Client
  resp, err := client.Do(req)
  if  err != nil {
    return fmt.Errorf("Error - %v\n", err)
  }
  defer resp.Body.Close()
  io.Copy(w, resp.Body)

  return nil
}

func (isakura *Isakura) FileName(download Download) string {
  program := download.Program
  title := program.Title
  monitor := download.Monitor
  log.Printf("Download - %+v", program)
  log.Printf("Starting Title - %v", title)
  time := time.Unix(download.Program.Time, 0)
  
  finalName := ""
  if monitor.PrependDate {
    finalName+=time.Format("MM-DD-YYYY") + " "
  }
  if monitor.PrependTime {
    finalName+=time.Format("HH:mm") + " "
  }
  for _, modifier := range monitor.Modifiers {
    log.Printf("Modify %v\n", title)
    title = strings.Replace(title, modifier.Search, modifier.Replace, -1)
    log.Printf("Search: %v, Replace: %v, New Title: %v\n", modifier.Search, modifier.Replace, title)
  } 
  finalName+=title
  if monitor.AppendDate {
    finalName+=" " + time.Format("MM-DD-YYYY")
  }
  if monitor.AppendTime {
    finalName+=" " + time.Format("HH:mm") 
  }
  log.Printf("Ending Title - %v", title)
  return title
}

func (isakura *Isakura) Delete(contentId string) error {
  for i, download := range isakura.Downloads {
     if download.ID == contentId {
       err := os.Remove(download.LocalPath)
       if err != nil { 
         return err
       }
       isakura.Downloads = append(isakura.Downloads[:i], isakura.Downloads[i+1:]...)
       isakura.Save()
       return nil
     } 
  }
  return fmt.Errorf("Not Found")
}

func (isakura *Isakura) Download(downloads chan Download) {
  for {
    download := <- downloads
    filename := isakura.FileName(download)
    folder := download.Monitor.Folder
    folderPath := filepath.Join(isakura.Home, "videos", folder)
    err := os.MkdirAll(folderPath, 0777)
    if (err != nil) {
      log.Printf("%v", err)
      continue
    }

    savePath := filepath.Join(isakura.Home, "videos", folder, fmt.Sprintf("%v._flv", filename) )
    finalPath := filepath.Join(isakura.Home, "videos", folder, fmt.Sprintf("%v.flv", filename) )

    fmt.Printf("Download: %v from %v\n", download.Program.Title, download.Channel.Name)
    if download.Program.Path == "" {
       fmt.Printf("Empty Playpath, can't download\n")
       continue
    }

    // Do not to forward if File exists
    if _, err = os.Stat(finalPath); ! os.IsNotExist(err) {
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
    newUuid, _ := uuid.NewV4()
    isakura.Downloads=append(isakura.Downloads, LocalProgram{ ID: newUuid.String(), Program: download.Program, LocalPath: finalPath, Monitor: download.Monitor, DownloadDate: time.Now()  })
    download.Program.Downloaded=true
    isakura.Save()
  }
}

func (isakura *Isakura) query(download Download) (*ProgramQuery, error) {
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

func (isakura *Isakura) Load() {
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

func (isakura *Isakura) Save() {
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
