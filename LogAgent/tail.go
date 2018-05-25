package main

import (
	"fmt"
	"sync"

	"github.com/astaxie/beego/logs"
	"github.com/hpcloud/tail"
)

// TailObj 一个文件跟踪对象
type TailObj struct {
	filename string
	offset   int64
	tailf    *tail.Tail
}

// TailMgr 管理所有的文件跟踪对象
type TailMgr struct {
	tailMap map[string]*TailObj
	lock    sync.RWMutex
}

// NewTailMgr 生成新的tailMgr对象 管理tail对象
func NewTailMgr(filelist []string) *TailMgr {
	tailMgr := &TailMgr{
		tailMap: make(map[string]*TailObj, 10),
	}

	for _, file := range filelist {
		tailMgr.AddTail(file)
	}
	return tailMgr
}

func (tobj *TailObj) traceLog() {
	for msg := range tobj.tailf.Lines {
		if msg.Err != nil {
			logs.Error(msg.Err)
			continue
		}
		kafkaSend.AddMsg(msg.Text)
	}
	wg.Done()
}

// AddTail 生成新的tail对象追加到map中
func (t *TailMgr) AddTail(filename string) (err error) {
	t.lock.Lock()
	defer t.lock.Unlock()

	_, ok := t.tailMap[filename]
	if ok {
		logs.Warn("duplicate file [%s],tail file process already exist.", filename)
		err = fmt.Errorf("duplicate file [%s],tail file process already exist", filename)
		return
	}

	tails, err := tail.TailFile(filename, tail.Config{
		ReOpen:    true,
		Follow:    true,
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2},
		MustExist: false,
		Poll:      true,
	})
	if err != nil {
		logs.Error("tail file [%s] failed", filename)
		return
	}
	tailObj := &TailObj{
		filename: filename,
		offset:   0,
		tailf:    tails,
	}

	t.tailMap[filename] = tailObj
	return
}

// Process 追踪日志对象
func (t *TailMgr) Process() {
	for _, obj := range t.tailMap {
		wg.Add(1)
		go obj.traceLog()
	}
}
