# gotyphon
Gotyphonlib is a library which allows you to request using any proxy type with http1 and http2, it also allows you to customize ja3 string and comes with a lot more feature. This project is created as alternative to Cycletls

# 🚀 Features

## Block
- Map_string_slice_int_tools
- OStools
- Randomtools (Ja3 ranoomizer, header randomizer, Order randomizer etc)
- Sufflleslice
- Httprequest
- Delay function
- Parse between two strings

## HTTPRequest Funtcion
- [High-performance](#-performance) Built-in goroutine pool used for handling asynchronous requests
- Custom header ordering via [fhttp](https://github.com/useflyent/fhttp)
- Proxy support (SOCKS4/SOCKS4a, HTTP/s, SOCKS5)
- Ja3 Customization
- HTTP/2 and HTTP1 Support
## How to install it?
```
go get github.com/kawacode/gotyphon
```
## How to use it?
### HTTP request using gotyphon without proxy or Ja3 customization
```go
package main

import (
	"github.com/kawacode/gotyphon"
)

func main() {
	//Ja3 is by default random
	var bot gotyphon.BotData                                    // Creating struct
	bot.HttpRequest.Request.URL = "https://tls.peet.ws/api/all" // URL
	bot.HttpRequest.Request.Method = "GET"                      // All method that are supported by fhttp is here supported too
	bot.HttpRequest.Request.Protocol = 2.0                      // 2.0 and 1.1 supported
	bot.HttpRequest.Request.ReadResponse = true                 // readresponse is by default false
	gotyphon.HttpRequest(&bot)                                  // Request
	println(bot.HttpRequest.Response.Source)                    // Print response source code of the request
}

```
### HTTP request using gotyphon with custom and ordered headers
```go
package main

import (
	"github.com/kawacode/gotyphon"
)

func main() {
	var bot gotyphon.BotData                                    // Creating a variable called `bot` and setting it to the type `gotyphon.BotData`.
	bot.HttpRequest.Request.URL = "https://tls.peet.ws/api/all" // Setting the URL of the request.
	bot.HttpRequest.Request.Method = "GET"                      // Used to set the method of the request.
	bot.HttpRequest.Request.Protocol = 2.0                      // Used to set the protocol version of the request.
	bot.HttpRequest.Request.ReadResponse = true                 // Used to read the response from the server.
	bot.HttpRequest.Request.Headers = map[string]string{        // Used to add headers to the request.
		"Content-Type": "application/json",
		"TestToken":    "123492190391239102301293",
		"user-agent":   "Go-http-client/2.0",
	}
	bot.HttpRequest.Request.HeaderOrderKey = []string{ // Used to order the headers in the request.
		"content-type",
		"user-agent",
		"testtoken",
	}
	bot.HttpRequest.Request.PHeaderOrderKey = []string{ // Used to order the Pseudo headers in the request.
		":authority",
		":path",
		":scheme",
		":method",
	}
	gotyphon.HttpRequest(&bot)               // A function that is used to send a request to the server.
	println(bot.HttpRequest.Response.Source) // Printing the response from the server.
}

```
### HTTP Request using gotyphon with proxy
```go
package main

import (
	"github.com/kawacode/gotyphon"
)

func main() {
	var bot gotyphon.BotData                                            // Creating a variable called `bot` and setting it to the type `gotyphon.BotData`.
	bot.HttpRequest.Request.URL = "https://tls.peet.ws/api/all"         // Setting the URL of the request.
	bot.HttpRequest.Request.Method = "GET"                              // Used to set the method of the request.
	bot.HttpRequest.Request.Protocol = 2.0                              // Used to set the protocol version of the request.
	bot.HttpRequest.Request.ReadResponse = true                         // Used to read the response from the server.
	bot.HttpRequest.Input.Proxy = "type://user:password@127.0.0.1:1080" // Setting the proxy of the request. Supported types: "http", "https", "socks4", "socks4a", "socks5"
	bot.HttpRequest.Input.UseProxy = true                               // Used to tell the program to use the proxy.
	gotyphon.HttpRequest(&bot)                                          // A function that is used to send a request to the server.
	println(bot.HttpRequest.Response.Source)                            // Printing the response from the server.
}

```
### HTTP Request using gotyphon with custom Ja3
```go
package main

import (
	"github.com/kawacode/gotyphon"
)

func main() {
	var bot gotyphon.BotData                                    // Creating a variable called `bot` and setting it to the type `gotyphon.BotData`.
	bot.HttpRequest.Request.URL = "https://tls.peet.ws/api/all" // Setting the URL of the request.
	bot.HttpRequest.Request.Method = "GET"                      // Used to set the method of the request.
	bot.HttpRequest.Request.Protocol = 2.0                      // Used to set the protocol version of the request.
	bot.HttpRequest.Request.ReadResponse = true                 // Used to read the response from the server.
	bot.HttpRequest.Input.Ja3 = "<Your_Ja3>"                    // Setting the JA3 fingerprint of the request.
	gotyphon.HttpRequest(&bot)                                  // A function that is used to send a request to the server.
	println(bot.HttpRequest.Response.Source)                    // Printing the response from the server.
}

```
### HTTP Request using gotyphon with MaxRedirects or disable redirects
```go
package main

import (
	"github.com/kawacode/gotyphon"
)

func main() {
	var bot gotyphon.BotData                                    // Creating a variable called `bot` and setting it to the type `gotyphon.BotData`.
	bot.HttpRequest.Request.URL = "https://tls.peet.ws/api/all" // Setting the URL of the request.
	bot.HttpRequest.Request.Method = "GET"                      // Used to set the method of the request.
	bot.HttpRequest.Request.Protocol = 2.0                      // Used to set the protocol version of the request.
	bot.HttpRequest.Request.ReadResponse = true                 // Used to read the response from the server.
	bot.HttpRequest.Request.DisableRedirects = true             // Used to disable redirects from server.
	bot.HttpRequest.Request.MaxRedirects = 10                   // Used to set the max number of redirects allowed from server.
	gotyphon.HttpRequest(&bot)                                  // A function that is used to send a request to the server.
	println(bot.HttpRequest.Response.Source)                    // Printing the response from the server.
}

```
### HTTP Request using gotyphon with Payload
```go
package main

import (
	"github.com/kawacode/gotyphon"
)

func main() {
	var bot gotyphon.BotData                                    // Creating a variable called `bot` and setting it to the type `gotyphon.BotData`.
	bot.HttpRequest.Request.URL = "https://tls.peet.ws/api/all" // Setting the URL of the request.
	bot.HttpRequest.Request.Method = "POST"                     // Used to set the method of the request.
	bot.HttpRequest.Request.Protocol = 2.0                      // Used to set the protocol version of the request.
	bot.HttpRequest.Request.ReadResponse = true                 // Used to read the response from the server.
	bot.HttpRequest.Request.Payload = "user=joe&pass=ok"        // Setting the payload of the request.
	gotyphon.HttpRequest(&bot)                                  // A function that is used to send a request to the server.
	println(bot.HttpRequest.Response.Source)                    // Printing the response from the server.
}
```
## LICENSE
### GPL3 LICENSE SYNOPSIS

**_TL;DR_*** Here's what the GPL3 license entails:

```markdown
1. Anyone can copy, modify and distribute this software.
2. You have to include the license and copyright notice with each and every distribution.
3. You can use this software privately.
4. You can use this software for commercial purposes.
5. Source code MUST be made available when the software is distributed.
6. Any modifications of this code base MUST be distributed with the same license, GPLv3.
7. This software is provided without warranty.
8. The software author or license can not be held liable for any damages inflicted by the software.
```

More information on about the [LICENSE can be found here](http://choosealicense.com/licenses/gpl-3.0/)
