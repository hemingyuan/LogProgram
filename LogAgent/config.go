package main

import (
	"LogProgram/config"
	"strings"
)

// AppConfig 用于存放读取的配置
type AppConfig struct {
	logPath     string
	logLevel    string
	kafkaAddr   string
	kafkaThread int
	collectLogs []string
}

func initConfig(confFile string) (err error) {
	conf, err := config.NewConf(confFile)
	if err != nil {
		return
	}

	logPath := conf.GetDefaultString("logPath", "./logs/server.log")

	logLevel := conf.GetDefaultString("logLevel", "info")

	kafkaAddr, err := conf.GetString("kafkaAddr")
	if err != nil {
		return
	}

	kafkaThread := conf.GetDefaultInt("kafkaThread", 8)

	collectLogs, err := conf.GetString("collectLogs")
	if err != nil {
		return
	}

	logFiles := strings.Split(collectLogs, ",")
	for _, v := range logFiles {
		v = strings.TrimSpace(v)
		if len(v) == 0 {
			continue
		}
		appConf.collectLogs = append(appConf.collectLogs, v)
	}

	appConf.logPath = logPath
	appConf.logLevel = logLevel
	appConf.kafkaAddr = kafkaAddr
	appConf.kafkaThread = kafkaThread
	return
}
