package server

import (
	"airpush/auction"
	"fmt"
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
	"net/url"
	"os"
	"time"
)

// settings model
type ServerSettings struct {
	ServerName string
	ServerAddr string
	ReadBufferSize int
	WriteBufferSize int
	ReadTimeout time.Duration
	WriteTimeout time.Duration
	Concurrency int
	DisableKeepalive bool
}

// settings setter
type ServerSetOption func(*Server)

// set server name
// for debug app in prod for indicate physical node
func SetServerName(name string) ServerSetOption {
	return func(s *Server) {
		s.settings.ServerName = name
	}
}

// set server name
// for debug app in prod for indicate physical node
func SetServerAddr(addr string) ServerSetOption {
	return func(s *Server) {
		s.settings.ServerAddr = addr
	}
}

// set max per request buf size
func SetReadBufferSize(size int) ServerSetOption {
	return func(s *Server) {
		s.settings.ReadBufferSize = size
	}
}

// set max per response buf size
func SetWriteBufferSize(size int) ServerSetOption {
	return func(s *Server) {
		s.settings.WriteBufferSize = size
	}
}

// set time how long server wait read request
func SetReadTimeout(ms int) ServerSetOption {
	return func(s *Server) {
		s.settings.ReadTimeout = time.Duration(ms) * time.Millisecond
	}
}

// set time how long server wait write response
func SetWriteTimeout(ms int) ServerSetOption {
	return func(s *Server) {
		s.settings.WriteTimeout = time.Duration(ms) * time.Millisecond
	}
}

// set keepalive logic, by default server maintains client connection
func SetDisableKeepalive(val bool) ServerSetOption {
	return func(s *Server) {
		s.settings.DisableKeepalive = val
	}
}

// set max concurrency requests
func SetConcurrency(c int) ServerSetOption {
	return func(s *Server) {
		s.settings.Concurrency = c
	}
}

// set custom logger implement fasthttp logger interface
func SetLogger(logger fasthttp.Logger) ServerSetOption {
	return func(s *Server) {
		s.logger = logger
	}
}

// default server settings
func defaultSettings() ServerSettings {

	s := ServerSettings{
		ServerName: "simple rtb",
		ServerAddr: ":8080",

		ReadBufferSize: 4096,
		WriteBufferSize: 4096,

		ReadTimeout: time.Duration(100) * time.Millisecond,
		WriteTimeout: time.Duration(100) * time.Millisecond,

		DisableKeepalive: true,
		Concurrency: fasthttp.DefaultConcurrency,
	}

	s.ServerName, _ = os.Hostname()

	return s
}

// cors middleware
func defaultMiddleWare(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return fasthttp.RequestHandler(func(ctx *fasthttp.RequestCtx) {

		origin := string(ctx.Request.Header.Peek("Origin"))
		if origin == "" {
			ref := string(ctx.Referer())
			if ref != "" {
				u, err := url.Parse(ref)
				if err == nil {
					origin = fmt.Sprintf("%s://%s", u.Scheme, u.Host)
				}
			}
		}

		if origin != "" {
			ctx.Response.Header.Set("Access-Control-Allow-Credentials", "true")
			ctx.Response.Header.Set("Access-Control-Allow-Origin", origin)
			ctx.Response.Header.Set("Vary", "Origin")
			ctx.SetContentType("application/json; charset=utf-8")
		} else {
			ctx.Response.Header.Set("Access-Control-Allow-Credentials", "true")
			ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")
			ctx.Response.Header.Set("Vary", "Origin")
			ctx.SetContentType("application/json; charset=utf-8")
		}

		// no content for cr request
		if ctx.IsOptions() {
			ctx.SetStatusCode(fasthttp.StatusNoContent)
			return
		}

		next(ctx)
	})
}

// Server root struct
type Server struct {
	settings ServerSettings
	server *fasthttp.Server
	auction *auction.Auction
	logger fasthttp.Logger
}

// init server
func New(opts ...ServerSetOption) (proto *Server, err error) {

	proto = &Server{
		settings: defaultSettings(), // set default settings
	}

	// set custom server params
	for _, opt := range opts {
		opt(proto)
	}

	routing := router.New()

	// monitoring app route
	routing.GET("/ping", proto.PingRoute)

	// monitoring app
	routing.GET("/", proto.BidRoute)

	// определяем сервер
	proto.server = &fasthttp.Server{
		ReadTimeout:      proto.settings.ReadTimeout,
		WriteTimeout:     proto.settings.WriteTimeout,

		ReadBufferSize:   proto.settings.ReadBufferSize,
		WriteBufferSize:   proto.settings.WriteBufferSize,

		Concurrency:      proto.settings.Concurrency,
		DisableKeepalive: proto.settings.DisableKeepalive,

		Name:         proto.settings.ServerName,
		Handler:      defaultMiddleWare(routing.Handler),
		Logger:       proto.logger,
	}

	return
}

// services route for monitoring app state
func (s *Server) PingRoute(ctx *fasthttp.RequestCtx) {
	ctx.SetStatusCode(fasthttp.StatusNoContent)
}

// services route for monitoring app state
func (s *Server) BidRoute(ctx *fasthttp.RequestCtx) {
	ctx.SetStatusCode(fasthttp.StatusNoContent)

	//run auction
	//bid, err := s.auction.Do()

	//return bet bid
}

// loop server
func (s *Server) Start() (err error) {
	s.logger.Printf("listen server on: %s\n", s.settings.ServerAddr)
	return s.server.ListenAndServe(s.settings.ServerAddr)
}

// stop server
func (s *Server) Close() (err error) {

	s.logger.Printf("stop server")
	err = s.server.Shutdown()
	if err != nil {
		s.logger.Printf("with err: %s", err)
	}

	return
}
