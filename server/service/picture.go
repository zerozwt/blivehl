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
	"time"

	"github.com/zerozwt/blivehl/server/utils"
)

var gPicCacheDir string

const (
	PIC_DIR = "pics"
)

func InitPictureCache(wwwroot string) (string, error) {
	gPicCacheDir = filepath.Join(wwwroot, PIC_DIR)
	return gPicCacheDir, os.MkdirAll(gPicCacheDir, 0755)
}

type pictureRequestUnit struct {
	uid int64
	url string
}

var picFetcher *utils.Fetcher[pictureRequestUnit, string] = utils.NewFetcher(time.Microsecond*500, convertPictureUrl)

func convertPictureUrl(req pictureRequestUnit) (*string, error) {
	res, err := url.Parse(req.url)
	if err != nil {
		return nil, err
	}

	dotIdx := strings.LastIndex(res.Path, ".")
	if dotIdx < 0 {
		return nil, fmt.Errorf("picture %s has no ext", req.url)
	}
	ext := res.Path[dotIdx+1:]

	sum := md5.Sum([]byte(req.url))
	hash := hex.EncodeToString(sum[:])

	sUID := fmt.Sprint(req.uid)
	localName := fmt.Sprintf("%s.%s", hash, ext)
	localUrl := "/" + PIC_DIR + "/" + sUID + "/" + localName
	localFile := filepath.Join(gPicCacheDir, sUID, localName)

	if err := os.MkdirAll(filepath.Join(gPicCacheDir, sUID), 0755); err != nil {
		return nil, err
	}

	if _, err := os.Stat(localFile); err == nil {
		return &localUrl, nil
	}

	rsp, err := http.Get(req.url)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()

	picData, err := io.ReadAll(rsp.Body)
	if err != nil {
		return nil, err
	}

	picFile, err := os.OpenFile(localFile, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}
	defer picFile.Close()

	_, err = picFile.Write(picData)
	if err != nil {
		return nil, err
	}

	return &localUrl, nil
}

func SaveFileAndGetLocalUrl(uid int64, picurl string) (string, error) {
	localUrl, err := picFetcher.Get(pictureRequestUnit{
		uid: uid,
		url: picurl,
	})
	if err != nil {
		return "", err
	}
	return *localUrl, nil
}
