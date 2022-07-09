package gotyphon

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"errors"
	"net"
	"net/http"
	"net/url"
	"time"

	"golang.org/x/net/proxy"
)

// `HttpProxy` is a struct with fields `Host`, `isAuth`, `Username`, `Password`, `Scheme`, and
// `Forward`.
// @property {string} Host - The hostname of the proxy server.
// @property {bool} isAuth - Whether or not the proxy requires authentication
// @property {string} Username - The username to use for authentication
// @property {string} Password - The password for the proxy server.
// @property {string} Scheme - The scheme of the proxy. This is either http or https.
// @property Forward - This is the proxy that will be used to forward the request to the destination.
type HttpProxy struct {
	Host     string
	isAuth   bool
	Username string
	Password string
	Scheme   string
	Forward  proxy.Dialer
}

// function that is used to dial the proxy server.
func (HTTPPROXY *HttpProxy) Dial(network, addr string) (net.Conn, error) {
	c, err := HTTPPROXY.Forward.Dial("tcp", HTTPPROXY.Host)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("CONNECT", HTTPPROXY.Scheme+"://"+addr, nil)
	if err != nil {
		defer c.Close()
		return nil, err
	}
	req.Close = false
	if HTTPPROXY.isAuth {
		req.Header.Add("Proxy-Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(HTTPPROXY.Username+":"+HTTPPROXY.Password)))
		req.Header.Add("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(HTTPPROXY.Username+":"+HTTPPROXY.Password)))
	}
	err = req.Write(c)
	if err != nil {
		defer c.Close()
		return nil, err
	}
	resp, err := http.ReadResponse(bufio.NewReader(c), nil)
	if err != nil {
		defer c.Close()
		return nil, err
	}
	defer resp.Body.Close()
	return c, nil
}

// It creates a new HTTPProxy struct, sets the Host, Forward, Scheme, Username, Password, and isAuth
// fields, and returns the struct
func NewHTTPProxyDialer(uri *url.URL, forward proxy.Dialer) (proxy.Dialer, error) {
	HTTPPROXY := new(HttpProxy)
	HTTPPROXY.Host = uri.Host
	HTTPPROXY.Forward = forward
	HTTPPROXY.Scheme = uri.Scheme
	if uri.User != nil {
		HTTPPROXY.isAuth = true
		HTTPPROXY.Username = uri.User.Username()
		HTTPPROXY.Password, _ = uri.User.Password()
	}
	return HTTPPROXY, nil
}

//SOCKS4 or SOCKS4A Handler
type Socks4_a_Proxy struct {
	Host     string
	Scheme   string
	isAuth   bool
	Username string
	Password string
}

// A function that is used to send the request to the proxy server and receive the response.
func (socks4_a *Socks4_a_Proxy) sendReceive(conn net.Conn, req []byte) (resp []byte, err error) {
	_, err = conn.Write(req)
	if err != nil {
		return
	}
	resp, err = socks4_a.readAll(conn)
	return
}

// Reading the response from the proxy server.
func (socks4_a *Socks4_a_Proxy) readAll(conn net.Conn) (resp []byte, err error) {
	resp = make([]byte, 1024)
	n, err := conn.Read(resp)
	resp = resp[:n]
	return
}

// A function that is used to dial the proxy server.
func (socks4_a *Socks4_a_Proxy) Dial(network, addr string) (net.Conn, error) {
	proxy := socks4_a.Host

	// dial TCP
	conn, err := net.DialTimeout("tcp", proxy, time.Duration(30*time.Second))
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			conn.Close()
		}
	}()

	// connection request
	host, port, err := SplitHostPort(addr)
	if err != nil {
		return nil, err
	}
	ip := net.IPv4(0, 0, 0, 1).To4()
	if socks4_a.Scheme == "socks4" {
		ip, err = LookupIPv4(host)
		if err != nil {
			return nil, err
		}
	}
	req := []byte{
		4,                          // version number
		1,                          // command CONNECT
		byte(port >> 8),            // higher byte of destination port
		byte(port),                 // lower byte of destination port (big endian)
		ip[0], ip[1], ip[2], ip[3], // special invalid IP address to indicate the host name is provided
		0, // user id is empty, anonymous proxy only
	}
	if socks4_a.Scheme == "socks4a" {
		req = append(req, []byte(host+"\x00")...)
	}

	resp, err := socks4_a.sendReceive(conn, req)
	if err != nil {
		return nil, err
	} else if len(resp) != 8 {
		return nil, errors.New("server does not respond properly")
	}
	switch resp[1] {
	case 90:
		// request granted
	case 91:
		return nil, errors.New("socks connection request rejected or failed")
	case 92:
		return nil, errors.New("socks connection request rejected because SOCKS server cannot connect to identd on the client")
	case 93:
		return nil, errors.New("socks connection request rejected because the client program and identd report different user-ids")
	default:
		return nil, errors.New("socks connection request failed, unknown error")
	}
	// clear the deadline before returning
	if err := conn.SetDeadline(time.Time{}); err != nil {
		return nil, err
	}
	return conn, nil
}

// It creates a new SOCKS4_A_ProxyDialer object and returns it.
func NewSOCKS4_A_ProxyDialer(uri *url.URL, forward proxy.Dialer) (proxy.Dialer, error) {
	SOCKSPROXY := new(Socks4_a_Proxy)
	SOCKSPROXY.Host = uri.Host
	SOCKSPROXY.Scheme = uri.Scheme
	if uri.User != nil {
		SOCKSPROXY.isAuth = true
		SOCKSPROXY.Username = uri.User.Username()
		SOCKSPROXY.Password, _ = uri.User.Password()
	}
	return SOCKSPROXY, nil
}

// `SocksProxy` is a struct with fields `Host`, `Scheme`, `isAuth`, `Username`, and `Password`.
// @property {string} Host - The hostname or IP address of the SOCKS proxy server.
// @property {string} Scheme - The scheme of the proxy. This can be either http or socks5.
// @property {bool} isAuth - Whether or not the proxy requires authentication.
// @property {string} Username - The username to use for authentication.
// @property {string} Password - The password for the proxy server.
type SocksProxy struct {
	Host     string
	Scheme   string
	isAuth   bool
	Username string
	Password string
}

// It's a bytes.Buffer that has a method called Build.
// @property  - `method` - The HTTP method to use for the request.
type requestBuilder struct {
	bytes.Buffer
}

// It's a bytes.Buffer that has a method called Build.
func (req *requestBuilder) add(data ...byte) {
	_, _ = req.Write(data)
}

// It's a function that is used to send the request to the proxy server and receive the response.
func (socks5 *SocksProxy) sendReceive(conn net.Conn, req []byte) (resp []byte, err error) {
	_, err = conn.Write(req)
	if err != nil {
		return
	}
	resp, err = socks5.readAll(conn)
	return
}

// It's reading the response from the proxy server.
func (socks5 *SocksProxy) readAll(conn net.Conn) (resp []byte, err error) {
	resp = make([]byte, 1024)
	n, err := conn.Read(resp)
	resp = resp[:n]
	return
}

// It's a function that is used to dial the proxy server.
func (socks5 *SocksProxy) Dial(network, addr string) (net.Conn, error) {
	proxy := socks5.Host

	conn, err := net.DialTimeout("tcp", proxy, time.Duration(30*time.Second))
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			conn.Close()
		}
	}()

	var req requestBuilder

	version := byte(5)
	method := byte(0)
	if socks5.isAuth {
		method = 2
	}
	req.add(
		version,
		1,
		method,
	)

	resp, err := socks5.sendReceive(conn, req.Bytes())
	if err != nil {
		return nil, err
	} else if len(resp) != 2 {
		return nil, errors.New("server does not respond properly")
	} else if resp[0] != 5 {
		return nil, errors.New("server does not support Socks 5")
	} else if resp[1] != method {
		return nil, errors.New("socks method negotiation failed")
	}
	if socks5.isAuth {
		version := byte(1)
		req.Reset()
		req.add(
			version,
			byte(len(socks5.Username)),
		)
		req.add([]byte(socks5.Username)...)
		req.add(byte(len(socks5.Password)))
		req.add([]byte(socks5.Password)...)
		resp, err := socks5.sendReceive(conn, req.Bytes())
		if err != nil {
			return nil, err
		} else if len(resp) != 2 {
			return nil, errors.New("server does not respond properly")
		} else if resp[0] != version {
			return nil, errors.New("server does not support user/password version 1")
		} else if resp[1] != 0 { // not success
			return nil, errors.New("user/password login failed")
		}
	}

	// detail request
	host, port, err := SplitHostPort(addr)
	if err != nil {
		return nil, err
	}
	req.Reset()
	req.add(
		5,
		1,
		0,
		3,
		byte(len(host)),
	)
	req.add([]byte(host)...)
	req.add(
		byte(port>>8),
		byte(port),
	)
	resp, err = socks5.sendReceive(conn, req.Bytes())
	if err != nil {
		return nil, err
	} else if len(resp) != 10 {
		return nil, errors.New("server does not respond properly")
	} else if resp[1] != 0 {
		return nil, errors.New("can't complete SOCKS5 connection")
	}

	return conn, nil
}

// It creates a new SOCKS5 proxy dialer.
func NewSOCKS5ProxyDialer(uri *url.URL, forward proxy.Dialer) (proxy.Dialer, error) {
	SOCKSPROXY := new(SocksProxy)
	SOCKSPROXY.Host = uri.Host
	SOCKSPROXY.Scheme = uri.Scheme
	if uri.User != nil {
		SOCKSPROXY.isAuth = true
		SOCKSPROXY.Username = uri.User.Username()
		SOCKSPROXY.Password, _ = uri.User.Password()
	}
	return SOCKSPROXY, nil
}

// It registers the HTTPProxyDialer as the dialer for all proxy types
func RegisterProxyDialers() {
	proxy.RegisterDialerType("https", NewHTTPProxyDialer)
	proxy.RegisterDialerType("http", NewHTTPProxyDialer)
	proxy.RegisterDialerType("socks4", NewHTTPProxyDialer)
	proxy.RegisterDialerType("socks4a", NewHTTPProxyDialer)
	proxy.RegisterDialerType("socks5", NewHTTPProxyDialer)
}

// It creates a proxy dialer and returns the Dialer
func CreateHTTPProxyDialer(proxystring string) (proxy.Dialer, error) {
	if proxystring != "" {
		proxyurl, err := url.Parse(proxystring)
		if err != nil {
			return nil, err
		}
		dialer, err := proxy.FromURL(proxyurl, proxy.Direct)
		if err != nil {
			return nil, err
		}
		return dialer, nil
	} else {
		return nil, errors.New("please make sure you added proxy")
	}
}
