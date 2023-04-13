package engine

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/zerozwt/blivehl/server/logger"
)

var APIPrefix string = "/api"

var apiMap map[string]http.HandlerFunc = make(map[string]http.HandlerFunc)
var regLock sync.Mutex

type myServer struct {
	fileServer spaFileServer
}

func RegisterRawApi(api string, handler HandlerFunc, middlewares ...HandlerFunc) {
	regLock.Lock()
	defer regLock.Unlock()
	apiMap[APIPrefix+api] = func(w http.ResponseWriter, r *http.Request) {
		ctx := makeContext(w, r, append(middlewares, handler)...)
		ctx.Next()
	}
}

func RegisterApi[InType, OutType any](api string, handler func(*Context, *InType) (*OutType, error), middlewares ...HandlerFunc) {
	regLock.Lock()
	defer regLock.Unlock()
	apiMap[APIPrefix+api] = makeAPIHandler(handler, middlewares...)
}

func Serve(wwwdir string, port int) {
	svr := http.Server{
		Addr:    fmt.Sprintf("localhost:%d", port),
		Handler: &myServer{fileServer: spaFileServer(wwwdir)},
	}
	svr.ListenAndServe()
}

func (s *myServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(r.URL.Path, APIPrefix) {
		s.serveApi(w, r)
		return
	}
	s.fileServer.ServeHTTP(w, r)
}

func (s *myServer) serveApi(w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	defer func() {
		logger.INFO("process API call to %s cost time: %v", r.URL.Path, time.Since(now))
		if err := recover(); err != nil {
			logger.ERROR("request to %s panic: %v", r.URL.Path, err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprint(err)))
		}
	}()
	handler, ok := apiMap[r.URL.Path]
	if !ok {
		logger.ERROR("requst to API %s failed: API not found", r.URL.Path)
		http.NotFound(w, r)
		return
	}
	handler(w, r)
}
