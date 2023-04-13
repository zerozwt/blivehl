package engine

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/zerozwt/blivehl/server/logger"
)

type Server struct {
	apiPrefix string
	apiMap    map[string]http.HandlerFunc
	apiLock   sync.Mutex

	fileServer spaFileServer
}

var defaultServer *Server = NewServer("/api")

func NewServer(apiPrfix string) *Server {
	return &Server{
		apiPrefix: apiPrfix,
		apiMap:    map[string]http.HandlerFunc{},
	}
}

func (s *Server) RegisterApi(api string, handler HandlerFunc, middlewares ...HandlerFunc) {
	s.apiLock.Lock()
	defer s.apiLock.Unlock()
	s.apiMap[s.apiPrefix+api] = func(w http.ResponseWriter, r *http.Request) {
		ctx := makeContext(w, r, append(middlewares, handler)...)
		ctx.Next()
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(r.URL.Path, s.apiPrefix) {
		s.serveApi(w, r)
		return
	}
	s.fileServer.ServeHTTP(w, r)
}

func (s *Server) Serve(wwwdir string, port int) {
	s.fileServer = spaFileServer(wwwdir)
	svr := http.Server{
		Addr:    fmt.Sprintf("localhost:%d", port),
		Handler: s,
	}
	svr.ListenAndServe()
}

func (s *Server) serveApi(w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	defer func() {
		logger.INFO("process API call to %s cost time: %v", r.URL.Path, time.Since(now))
		if err := recover(); err != nil {
			logger.ERROR("request to %s panic: %v", r.URL.Path, err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprint(err)))
		}
	}()
	handler, ok := s.apiMap[r.URL.Path]
	if !ok {
		logger.ERROR("requst to API %s failed: API not found", r.URL.Path)
		http.NotFound(w, r)
		return
	}
	handler(w, r)
}

func RegisterRawApi(api string, handler HandlerFunc, middlewares ...HandlerFunc) {
	defaultServer.RegisterApi(api, handler, middlewares...)
}

func RegisterApi[InType, OutType any](api string, handler func(*Context, *InType) (*OutType, error), middlewares ...HandlerFunc) {
	defaultServer.RegisterApi(api, MakeAPIHandler(handler), middlewares...)
}

func Serve(wwwdir string, port int) {
	defaultServer.Serve(wwwdir, port)
}
