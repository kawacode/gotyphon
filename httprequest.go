package gotyphon

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/url"

	fhttp "github.com/Danny-Dasilva/fhttp"
)

// It makes an HTTP request and returns the response
// srcbot is the bot where datas of response or request are stored/saved
// dstbot it is the bot from the Request is being used by HttpRequest
func HttpRequest(srcbot *BotData) (fhttp.Response, error) {
	var URL, err = url.Parse(srcbot.HttpRequest.Request.URL)
	if err != nil {
		return fhttp.Response{}, err
	}
	if URL.String() == "" {
		return fhttp.Response{}, errors.New("please provide a URL parameter at srcbot.HttpRequest.Request.URL")
	}
	if srcbot.HttpRequest.Request.Method == "" {
		return fhttp.Response{}, errors.New("please provide a method parameter at srcbot.HttpRequest.Request.Method")
	}
	req, err := fhttp.NewRequest(srcbot.HttpRequest.Request.Method, URL.String(), bytes.NewBuffer([]byte(srcbot.HttpRequest.Request.Payload)))
	if err != nil {
		return fhttp.Response{}, err
	}
	req.Header = MapStringToMapStringSlice(srcbot.HttpRequest.Request.Headers, srcbot)
	if err != nil {
		return fhttp.Response{}, err
	}
	client, err := CreateClient(srcbot)
	if err != nil {
		return fhttp.Response{}, err
	}
	res, err := client.Do(req)
	if err != nil {
		return fhttp.Response{}, err
	}
	location, err := res.Location()
	if err == nil {
		srcbot.HttpRequest.Response.Location = *location
	}
	cookies := make(map[string]string)
	for _, cookie := range res.Cookies() {
		cookies[cookie.Name] = cookie.Value
	}
	srcbot.HttpRequest.Response.Cookies = cookies
	srcbot.HttpRequest.Response.Status = res.Status
	srcbot.HttpRequest.Response.StatusCode = res.StatusCode
	srcbot.HttpRequest.Response.Headers = MapStringSliceToMapString(res.Header)
	srcbot.HttpRequest.Response.Protocol = res.Proto
	srcbot.HttpRequest.Response.ContentLength = res.ContentLength
	if srcbot.HttpRequest.Request.ReadResponse {
		resp, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return fhttp.Response{}, err
		}
		if srcbot.HttpRequest.Request.DecompressGzip {
			source, err := DecompressGzip(string(resp))
			if err != nil {
				srcbot.HttpRequest.Response.Source = string(resp)
			}
			srcbot.HttpRequest.Response.Source = source
		} else {
			srcbot.HttpRequest.Response.Source = string(resp)
		}
	}
	srcbot.HttpRequest.Response.ProtoMajor = res.ProtoMajor
	srcbot.HttpRequest.Response.ProtoMinor = res.ProtoMinor
	srcbot.HttpRequest.Response.WasUncompressed = res.Uncompressed
	defer res.Body.Close()
	client.CloseIdleConnections()
	return *res, nil
}
