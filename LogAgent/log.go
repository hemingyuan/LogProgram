package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/astaxie/beego/logs"
)

func getLogLevel(level string) (res int) {
	level = strings.ToLower(level)
	switch level {
	case "debug", "trace":
		res = logs.LevelDebug
	case "info":
		res = logs.LevelInfo
	case "notice":
		res = logs.LevelNotice
	case "warn", "warning":
		res = logs.LevelWarn
	case "error":
		res = logs.LevelError
	case "critical":
		res = logs.LevelCritical
	case "alert":
		res = logs.LevelAlert
	case "emergency":
		res = logs.LevelEmergency
	default:
		fmt.Println("loglevel configuration not found, use default log level")
		res = logs.LevelInfo
	}
	return
}

func initLog(logPath, logLevel string) (err error) {
	logConfig := make(map[string]interface{})
	logConfig["filename"] = logPath
	logConfig["level"] = getLogLevel(logLevel)

	jsonData, err := json.Marshal(logConfig)
	if err != nil {
		fmt.Println("json marshal log config failed")
		return
	}
	logConf := string(jsonData)

	err = logs.SetLogger(logs.AdapterFile, logConf)
	logs.SetLogger(logs.AdapterConsole, logConf)
	return
}
