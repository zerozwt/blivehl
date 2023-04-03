package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/zerozwt/blivehl/server/db"
	"github.com/zerozwt/blivehl/server/engine"
	"github.com/zerozwt/blivehl/server/handler"
	"github.com/zerozwt/blivehl/server/logger"
	"github.com/zerozwt/blivehl/server/service"
)

var gPort int
var gWebDir string
var gLogLevel string
var gDbFile string

func initLog() {
	switch strings.ToLower(gLogLevel) {
	case "debug":
		logger.SetLogLevel(logger.LOG_LEVEL_DEBUG)
	case "info":
		logger.SetLogLevel(logger.LOG_LEVEL_INFO)
	case "warn":
		logger.SetLogLevel(logger.LOG_LEVEL_WARN)
	case "error":
		logger.SetLogLevel(logger.LOG_LEVEL_ERROR)
	default:
		logger.SetLogLevel(logger.LOG_LEVEL_DEBUG)
	}
}

func main() {
	flag.IntVar(&gPort, "port", 4080, "web UI port")
	flag.StringVar(&gWebDir, "wwwroot", "-", "www file root dir")
	flag.StringVar(&gLogLevel, "log_level", "info", "log level (debug/info/warn/error)")
	flag.StringVar(&gDbFile, "db", "-", "database file")
	flag.Parse()
	initLog()

	if gPort <= 0 || gPort > 65535 {
		logger.ERROR("Invalid port: %d", gPort)
		return
	}

	if gWebDir == "-" {
		dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			logger.ERROR("get current exe dir failed: %v", err)
			return
		}
		gWebDir = filepath.Join(dir, "dist")
	}

	logger.INFO("initiating web server, port is %d, wwwroot is %s ...", gPort, gWebDir)

	info, err := os.Stat(gWebDir)
	if err != nil {
		logger.ERROR("www dir stat error: %v", err)
		return
	}
	if !info.IsDir() {
		logger.ERROR("www dir %s not a directory", gWebDir)
		return
	}

	if gDbFile == "-" {
		dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
		gDbFile = filepath.Join(dir, "blivehl.db")
	}

	logger.INFO("initiating database, db file: %s ...", gDbFile)
	if err := db.InitDB(gDbFile); err != nil {
		logger.ERROR("init database file %s failed: %v", gDbFile, err)
		return
	}

	logger.INFO("initiating picture cache ...")
	cacheDir, err := service.InitPictureCache(gWebDir)
	if err != nil {
		logger.ERROR("init picture cache error: %v", err)
		return
	}
	logger.INFO("cached pictures will be stored at %s ...", cacheDir)

	handler.InitHandlers()

	fmt.Printf("Please open your web browser and visit http://localhost:%d\n", gPort)
	engine.Serve(gWebDir, gPort)
}
