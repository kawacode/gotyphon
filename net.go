package gotyphon

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"net"
	"strconv"
)

// It takes a string, converts it to a byte array, creates a new gzip reader, reads the gzip reader,
// and returns the result as a string
func DecompressGzip(Gzip string) (string, error) {
	res, err := gzip.NewReader(bytes.NewReader([]byte(Gzip)))
	if err != nil {
		return Gzip, nil
	}
	defer res.Close()
	read, err := ioutil.ReadAll(res)
	if err != nil {
		return "", err
	}
	return string(read), nil
}

// It takes a hostname, looks up its IP addresses, and returns the first IPv4 address it finds
func LookupIPv4(host string) (net.IP, error) {
	ips, err := net.LookupIP(host)
	if err != nil {
		return nil, err
	}
	for _, ip := range ips {
		ipv4 := ip.To4()
		if ipv4 == nil {
			continue
		}
		return ipv4, nil
	}
	return nil, fmt.Errorf("no IPv4 address found for host: %s", host)
}

// It splits a string of the form "host:port" into its two components, and returns an error if the port
// is not a valid number
func SplitHostPort(addr string) (host string, port uint16, err error) {
	host, portStr, err := net.SplitHostPort(addr)
	if err != nil {
		return "", 0, err
	}
	portInt, err := strconv.ParseUint(portStr, 10, 16)
	if err != nil {
		return "", 0, err
	}
	port = uint16(portInt)
	return
}
