package gotyphon

import (
	"net/url"
)

// Creating a new type called `BotData` that is a struct.
type BotData struct {
	// A struct that is used to store the data that is sent to the bot.
	HttpRequest struct {
		Input struct {
			Username  string `json:"username"`
			Password  string `json:"password"`
			Proxy     string `json:"proxy"`
			UseProxy  bool   `json:"useproxy"`
			Ja3       string `json:"ja3"`
			RandomJa3 bool   `json:"randomja3"`
		}
		Response struct {
			Location        url.URL           `json:"location"`
			Status          string            `json:"status"`
			StatusCode      int               `json:"statuscode"`
			Headers         map[string]string `json:"headers"`
			Cookies         map[string]string `json:"cookies"`
			Protocol        string            `json:"protocol"`
			ContentLength   int64             `json:"contentlength"`
			Source          string            `json:"source"`
			ProtoMajor      int               `json:"protomajor"`
			ProtoMinor      int               `json:"protominor"`
			WasUncompressed bool              `json:"isuncompressed"`
		}
		Request struct {
			URL          string            `json:"url"`
			ReadResponse bool              `json:"readresponse"`
			Method       string            `json:"method"`
			Headers      map[string]string `json:"headers"`
			Payload      string            `json:"payload"`
			// Default 2.0
			Protocol         float32  `json:"protocol"`
			Timeout          int      `json:"timeout"`
			MaxRedirects     int      `json:"maxredirects"`
			DisableRedirects bool     `json:"disableredirects"`
			DecompressGzip   bool     `json:"decompressgzip"`
			HeaderOrderKey   []string `json:"headerorderkey"`
			PHeaderOrderKey  []string `json:"pheaderorderkey"`
			HTTP1TRANSPORT   struct {
				DisableKeepAlives      bool  `json:"disablekeepal"`
				DisableCompression     bool  `json:"disablecompression"`
				Timeout                int   `json:"timeout"`
				IdleConnTimeout        int   `json:"idleconntimeout"`
				TLSHandshakeTimeout    int   `json:"tlshandshaketimeout"`
				ExpectContinueTimeout  int   `json:"expectcontinuetimeout"`
				ForceAttemptHTTP2      bool  `json:"forceattempthttp2"`
				MaxConnsPerHost        int   `json:"maxconnsperhost"`
				MaxIdleConns           int   `json:"maxidleconns"`
				MaxIdleConnsPerHost    int   `json:"maxidleconnsperhost"`
				MaxResponseHeaderBytes int64 `json:"maxresponseheaderbytes"`
				ResponseHeaderTimeout  int   `json:"responseheadertimeout"`
			}
			HTTP2TRANSPORT struct {
				PingTimeout                int    `json:"pingtimeout"`
				ReadIdleTimeout            int    `json:"readidletimeout"`
				AllowHTTP                  bool   `json:"allowhttp"`
				DisableCompression         bool   `json:"disablecompression"`
				HeaderTableSize            uint32 `json:"headertablesize"`
				InitialWindowSize          uint32 `json:"initialwindowsize"`
				MaxHeaderListSize          uint32 `json:"maxheaderlistsize"`
				StrictMaxConcurrentStreams bool   `json:"strictmaxconcurrentstreams"`
				MaxFrameSize               uint32 `json:"maxframesize"`
				EnablePush                 uint32 `json:"enablepush"`
				MaxConcurrentStreams       uint32 `json:"maxconcurrentstreams"`
			}
		}
	}
}
