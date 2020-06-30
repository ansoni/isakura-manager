package isakura
import (
  "time"
  "strconv"
  "encoding/json"
  "fmt"
)


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
  Name string `json:"name"`
  Folder string `json:"folder"`
  Pattern string `json:"pattern"`
  Watch bool `json:"delete_on_download"`
  PrependDate bool `json:"delete_on_download"`
  PrependTime bool `json:"delete_on_download"`
  AppendDate bool `json:"delete_on_download"`
  AppendTime bool `json:"delete_on_download"`
  Modifiers []MonitorModifiers `json:"modifiers"`
}

type MonitorModifiers struct {
  Search string `json:"search"`
  Replace string `json:"replace"`
}

type Isakura struct {
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
  Downloaded bool `json:"downloaded"`
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
  Monitor Monitor
}

type LocalProgram struct {
  ID string `json:"id"`
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
