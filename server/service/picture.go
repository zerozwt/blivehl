package service

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var gPicCacheDir string
var gSaveLock sync.Mutex

const (
	PIC_DIR = "pics"
)

func InitPictureCache(wwwroot string) (string, error) {
	gPicCacheDir = filepath.Join(wwwroot, PIC_DIR)
	return gPicCacheDir, os.MkdirAll(gPicCacheDir, 0644)
}

func SaveFileAndGetLocalUrl(uid int64, picurl string) (string, error) {
	res, err := url.Parse(picurl)
	if err != nil {
		return "", err
	}

	dotIdx := strings.LastIndex(res.Path, ".")
	if dotIdx < 0 {
		return "", fmt.Errorf("picture %s has no ext", picurl)
	}
	ext := res.Path[dotIdx+1:]

	sum := md5.Sum([]byte(picurl))
	hash := hex.EncodeToString(sum[:])

	localName := fmt.Sprintf("%d_%s.%s", uid, hash, ext)
	localUrl := "/" + PIC_DIR + "/" + localName
	localFile := filepath.Join(gPicCacheDir, localName)

	if _, err := os.Stat(localFile); err == nil {
		return localUrl, nil
	}

	rsp, err := http.Get(picurl)
	if err != nil {
		return "", err
	}
	defer rsp.Body.Close()

	picData, err := io.ReadAll(rsp.Body)
	if err != nil {
		return "", err
	}

	gSaveLock.Lock()
	defer gSaveLock.Unlock()

	picFile, err := os.OpenFile(localFile, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return "", err
	}
	defer picFile.Close()

	_, err = picFile.Write(picData)
	if err != nil {
		return "", err
	}

	return localUrl, nil
}
