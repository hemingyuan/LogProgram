package config

import (
	"bufio"
	"io"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/astaxie/beego/logs"
)

// Conf Storage config
type Conf struct {
	filename       string
	lastModifyTime int64
	notifyList     []Notifyer
	confMap        map[string]string
	lock           sync.RWMutex
}

// NewConf init a conf object
func NewConf(filename string) (conf *Conf, err error) {
	conf = &Conf{
		filename: filename,
		confMap:  make(map[string]string, 256),
	}

	// 首次解析加载配置
	m, err := conf.firstParse()
	if err != nil {
		return
	}

	conf.lock.Lock()
	conf.confMap = m
	conf.lock.Unlock()

	// 定时监测配置文件是否变更
	go conf.reload()
	return
}

func (c *Conf) firstParse() (m map[string]string, err error) {
	f, err := os.Open(c.filename)
	if err != nil {
		logs.Error("load config file failed", err)
		return
	}
	defer f.Close()

	finfo, err := f.Stat()
	if err == nil {
		c.lastModifyTime = finfo.ModTime().Unix()
	}

	m, err = c.parse(f)
	return
}

// 解析配置文件并保存到map中
func (c *Conf) parse(f *os.File) (m map[string]string, err error) {
	m = make(map[string]string, 128)
	reader := bufio.NewReader(f)
	var lineNum int

	for {
		line, _, err := reader.ReadLine()

		lineNum++

		if err == io.EOF {
			break
		}
		if err != nil {
			logs.Warn("parse config line:[%d] failed, err [%v]", lineNum, err)
			continue
		}

		lineStr := strings.TrimSpace(string(line))
		if len(lineStr) == 0 || lineStr[0] == ';' || lineStr[0] == '#' {
			continue
		}

		fields := strings.Split(lineStr, "=")
		if len(fields) != 2 {
			logs.Warn("parse config line:[%d] failed.", lineNum)
			continue
		}

		key, value := strings.TrimSpace(fields[0]), strings.TrimSpace(fields[1])
		if len(key) == 0 || len(value) == 0 {
			logs.Warn("parse config line:[%d] failed.", lineNum)
			continue
		}

		m[key] = value
	}
	return
}

func (c *Conf) reload() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		func() {
			f, err := os.Open(c.filename)
			if err != nil {
				logs.Error("load config file:[%v] failed. err:%v", c.filename, err)
				return
			}
			defer f.Close()

			fInfo, err := f.Stat()
			if err != nil {
				logs.Error("get file [%s] stat failed. err:%v", c.filename, err)
				return
			}

			modifyTime := fInfo.ModTime().Unix()

			if c.lastModifyTime < modifyTime {
				m, err := c.parse(f)
				if err != nil { // 解析配置文件失败 不做任何操作
					return
				}

				c.lock.Lock()
				c.confMap = m
				c.lock.Unlock()

				c.lastModifyTime = modifyTime

				for _, v := range c.notifyList {
					v.Callback(c)
				}
			}
		}()
	}
}

// AddObserver 添加到配置更新通知列表
func (c *Conf) AddObserver(n Notifyer) {
	c.notifyList = append(c.notifyList, n)
}
