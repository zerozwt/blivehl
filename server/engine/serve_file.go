package engine

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type spaFileServer string

func (dir spaFileServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if strings.Contains(r.URL.Path, "..") || strings.Contains(r.URL.Path, "/.") {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	if r.URL.Path != "/" && len(r.URL.Path) > 0 {
		target := filepath.Join(string(dir), r.URL.Path[1:])
		info, err := os.Stat(target)
		if err != nil || info.IsDir() {
			r.URL.Path = "/"
		}
	}
	http.FileServer(http.Dir(dir)).ServeHTTP(w, r)
}
