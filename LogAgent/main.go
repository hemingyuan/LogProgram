package main

import (
	"fmt"
	"sync"

	"github.com/astaxie/beego/logs"
)

var (
	appConf   = &AppConfig{}     // 存放全局配置
	kafkaSend = NewKafkaSender() // kafkaMgr 管理类
	wg        sync.WaitGroup
)

func main() {
	// 初始化config模块
	err := initConfig("./conf/config.ini")
	if err != nil {
		panic(fmt.Sprintf("init config failed, err: %v", err))
	}
	fmt.Println("init config module success")
	fmt.Printf("%#v\n", appConf)

	// 初始化log模块
	err = initLog(appConf.logPath, appConf.logLevel)
	if err != nil {
		panic(fmt.Sprintf("init log failed, err: %v", err))
	}
	logs.Info("init log module success")

	// 初始化kafka客户端
	err = initKafka(appConf.kafkaAddr)
	if err != nil {
		panic(fmt.Sprintf("init kafka client failed, err: %v", err))
	}
	logs.Info("init kafka client success")

	// 初始化tail
	tailMgr := NewTailMgr(appConf.collectLogs)
	tailMgr.Process()

	wg.Wait()
}
