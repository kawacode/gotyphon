package gotyphon

import (
	"context"
	"net"
	"time"

	fhttp "github.com/Danny-Dasilva/fhttp"
	fhttp2 "github.com/Danny-Dasilva/fhttp/http2"
	utls "github.com/Danny-Dasilva/utls"
)

// It creates a client based on the protocol specified in the bot's request
func CreateClient(bot *BotData) (*fhttp.Client, error) {
	var client *fhttp.Client
	if bot.HttpRequest.Request.Protocol == 1.1 {
		var err error
		client, err = CreateHttp1Client(bot)
		if err != nil {
			return nil, err
		}
	} else {
		var err error
		client, err = CreateHttp2Client(bot)
		if err != nil {
			return nil, err
		}
	}
	return client, nil
}

// It creates an HTTP/1.1 client with the ability to use a proxy, disable redirects, and set a timeout, that uses a custom TLS dialer that uses a custom JA3 fingerprint
func CreateHttp1Client(bot *BotData) (*fhttp.Client, error) {
	http1transport := fhttp.Transport{
		DisableCompression:     bot.HttpRequest.Request.HTTP1TRANSPORT.DisableCompression,
		DisableKeepAlives:      bot.HttpRequest.Request.HTTP1TRANSPORT.DisableKeepAlives,
		ExpectContinueTimeout:  time.Duration(time.Duration(bot.HttpRequest.Request.HTTP1TRANSPORT.ExpectContinueTimeout) * time.Second),
		ForceAttemptHTTP2:      bot.HttpRequest.Request.HTTP1TRANSPORT.ForceAttemptHTTP2,
		IdleConnTimeout:        time.Duration(time.Duration(bot.HttpRequest.Request.HTTP1TRANSPORT.IdleConnTimeout) * time.Second),
		MaxConnsPerHost:        bot.HttpRequest.Request.HTTP1TRANSPORT.MaxConnsPerHost,
		MaxIdleConns:           bot.HttpRequest.Request.HTTP1TRANSPORT.MaxIdleConns,
		MaxIdleConnsPerHost:    bot.HttpRequest.Request.HTTP1TRANSPORT.MaxIdleConnsPerHost,
		MaxResponseHeaderBytes: bot.HttpRequest.Request.HTTP1TRANSPORT.MaxResponseHeaderBytes,
		ResponseHeaderTimeout:  time.Duration(time.Duration(bot.HttpRequest.Request.HTTP1TRANSPORT.ResponseHeaderTimeout) * time.Second),
		TLSHandshakeTimeout:    time.Duration(time.Duration(bot.HttpRequest.Request.HTTP1TRANSPORT.TLSHandshakeTimeout) * time.Second),
		DialTLSContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			utls.EnableWeakCiphers()
			var conn net.Conn
			if bot.HttpRequest.Input.UseProxy {
				dialer, err := CreateHTTPProxyDialer(bot.HttpRequest.Input.Proxy)
				if err != nil {
					return nil, err
				}
				con, err := dialer.Dial(network, addr)
				if err != nil {
					return nil, err
				}
				conn = con
			} else {
				var err error
				conn, err = net.Dial(network, addr)
				if err != nil {
					return nil, err
				}
			}
			host, _, err := net.SplitHostPort(addr)
			if err != nil {
				return nil, err
			}
			config := &utls.Config{ServerName: host, InsecureSkipVerify: true}
			var uconn *utls.UConn
			if bot.HttpRequest.Input.Ja3 == "" || bot.HttpRequest.Input.RandomJa3 {
				uconn = utls.UClient(conn, config, utls.HelloCustom)
				bot.HttpRequest.Input.Ja3 = RandomJA3()
				tlsspec, err := ParseJA3(bot.HttpRequest.Input.Ja3, bot.HttpRequest.Request.Protocol)
				if err != nil {
					return nil, err
				}
				if err := uconn.ApplyPreset(tlsspec); err != nil {
					return nil, err
				}
				if err := uconn.Handshake(); err != nil {
					return nil, err
				}
			} else {
				uconn = utls.UClient(conn, config, utls.HelloCustom)
				tlsspec, err := ParseJA3(bot.HttpRequest.Input.Ja3, bot.HttpRequest.Request.Protocol)
				if err != nil {
					return nil, err
				}
				if err := uconn.ApplyPreset(tlsspec); err != nil {
					return nil, err
				}
				if err := uconn.Handshake(); err != nil {
					return nil, err
				}
			}
			return uconn, nil
		},
	}

	client := &fhttp.Client{
		Transport: &http1transport,
		Timeout:   time.Duration(time.Duration(bot.HttpRequest.Request.Timeout) * time.Second),
	}
	if bot.HttpRequest.Request.DisableRedirects {
		client.CheckRedirect = func(req *fhttp.Request, via []*fhttp.Request) error {
			return fhttp.ErrUseLastResponse
		}
	} else {
		client.CheckRedirect = func(req *fhttp.Request, via []*fhttp.Request) error {
			if len(via) >= bot.HttpRequest.Request.MaxRedirects {
				return fhttp.ErrUseLastResponse
			}
			return nil
		}
	}
	return client, nil
}

// It creates an HTTP2 client with the ability to use a proxy, disable redirects, and set a timeout, that uses a custom TLS dialer that uses a custom JA3 fingerprint
func CreateHttp2Client(bot *BotData) (*fhttp.Client, error) {
	var client *fhttp.Client
	http2transport := fhttp2.Transport{
		PingTimeout:        time.Duration(time.Duration(bot.HttpRequest.Request.HTTP2TRANSPORT.PingTimeout) * time.Second),
		AllowHTTP:          bot.HttpRequest.Request.HTTP2TRANSPORT.AllowHTTP,
		ReadIdleTimeout:    time.Duration(time.Duration(bot.HttpRequest.Request.HTTP2TRANSPORT.ReadIdleTimeout) * time.Second),
		DisableCompression: bot.HttpRequest.Request.HTTP2TRANSPORT.DisableCompression,
		DialTLS: func(network, addr string, cfg *utls.Config) (net.Conn, error) {
			utls.EnableWeakCiphers()
			var conn net.Conn
			if bot.HttpRequest.Input.UseProxy {
				dialer, err := CreateHTTPProxyDialer(bot.HttpRequest.Input.Proxy)
				if err != nil {
					return nil, err
				}
				con, err := dialer.Dial(network, addr)
				if err != nil {
					return nil, err
				}
				conn = con
			} else {
				var err error
				conn, err = net.Dial(network, addr)
				if err != nil {
					return nil, err
				}
			}
			host, _, err := net.SplitHostPort(addr)
			if err != nil {
				return nil, err
			}
			config := &utls.Config{ServerName: host, InsecureSkipVerify: true}
			var uconn *utls.UConn
			if bot.HttpRequest.Input.Ja3 == "" || bot.HttpRequest.Input.RandomJa3 {
				uconn = utls.UClient(conn, config, utls.HelloCustom)
				bot.HttpRequest.Input.Ja3 = RandomJA3()
				tlsspec, err := ParseJA3(bot.HttpRequest.Input.Ja3, bot.HttpRequest.Request.Protocol)
				if err != nil {
					return nil, err
				}
				if err := uconn.ApplyPreset(tlsspec); err != nil {
					return nil, err
				}
				if err := uconn.Handshake(); err != nil {
					return nil, err
				}
			} else {
				uconn = utls.UClient(conn, config, utls.HelloCustom)
				tlsspec, err := ParseJA3(bot.HttpRequest.Input.Ja3, bot.HttpRequest.Request.Protocol)
				if err != nil {
					return nil, err
				}
				if err := uconn.ApplyPreset(tlsspec); err != nil {
					return nil, err
				}
				if err := uconn.Handshake(); err != nil {
					return nil, err
				}
			}
			return uconn, nil
		},
	}
	if bot.HttpRequest.Request.HTTP2TRANSPORT.InitialWindowSize != 0 {
		http2transport.InitialWindowSize = bot.HttpRequest.Request.HTTP2TRANSPORT.InitialWindowSize
	} else {
		http2transport.InitialWindowSize = uint32(RandomInt(131072, 1310720))
	}
	if bot.HttpRequest.Request.HTTP2TRANSPORT.HeaderTableSize != 0 {
		http2transport.HeaderTableSize = bot.HttpRequest.Request.HTTP2TRANSPORT.HeaderTableSize
	} else {
		http2transport.HeaderTableSize = uint32(RandomInt(65536, 999999))
	}
	if bot.HttpRequest.Request.HTTP2TRANSPORT.MaxConcurrentStreams != 0 {
		http2transport.Settings = append(http2transport.Settings, fhttp2.Setting{ID: fhttp2.SettingMaxConcurrentStreams, Val: bot.HttpRequest.Request.HTTP2TRANSPORT.MaxConcurrentStreams})
	} else {
		http2transport.Settings = append(http2transport.Settings, fhttp2.Setting{ID: fhttp2.SettingMaxConcurrentStreams, Val: uint32(RandomInt(1000, 100000))})
	}
	if bot.HttpRequest.Request.HTTP2TRANSPORT.EnablePush != 0 {
		http2transport.Settings = append(http2transport.Settings, fhttp2.Setting{ID: fhttp2.SettingEnablePush, Val: bot.HttpRequest.Request.HTTP2TRANSPORT.EnablePush})
	} else {
		http2transport.Settings = append(http2transport.Settings, fhttp2.Setting{ID: fhttp2.SettingEnablePush, Val: 0})
	}
	if bot.HttpRequest.Request.HTTP2TRANSPORT.MaxFrameSize != 0 {
		http2transport.Settings = append(http2transport.Settings, fhttp2.Setting{ID: fhttp2.SettingMaxFrameSize, Val: bot.HttpRequest.Request.HTTP2TRANSPORT.MaxFrameSize})
	} else {
		http2transport.Settings = append(http2transport.Settings, fhttp2.Setting{ID: fhttp2.SettingMaxFrameSize, Val: uint32(RandomInt(16384, 9999999))})
	}
	if bot.HttpRequest.Request.HTTP2TRANSPORT.MaxHeaderListSize != 0 {
		http2transport.Settings = append(http2transport.Settings, fhttp2.Setting{ID: fhttp2.SettingMaxHeaderListSize, Val: bot.HttpRequest.Request.HTTP2TRANSPORT.MaxHeaderListSize})
	} else {
		http2transport.Settings = append(http2transport.Settings, fhttp2.Setting{ID: fhttp2.SettingMaxHeaderListSize, Val: uint32(RandomInt(262144, 9999999))})
	}
	client = &fhttp.Client{
		Transport: &http2transport,
		Timeout:   time.Duration(time.Duration(bot.HttpRequest.Request.Timeout) * time.Second),
	}
	if bot.HttpRequest.Request.DisableRedirects {
		client.CheckRedirect = func(req *fhttp.Request, via []*fhttp.Request) error {
			return fhttp.ErrUseLastResponse
		}
	} else {
		client.CheckRedirect = func(req *fhttp.Request, via []*fhttp.Request) error {
			if len(via) >= bot.HttpRequest.Request.MaxRedirects {
				return fhttp.ErrUseLastResponse
			}
			return nil
		}
	}

	return client, nil
}
