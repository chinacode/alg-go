package util

import (
	"fmt"
	"sync"
	"time"
)

var (
	lock     *sync.Mutex = &sync.Mutex{}
	instance *TimeCut
)

type TimeCut struct {
	timeHistoryMap map[string][]int64
}

func initTag(tag string) {
	if instance == nil {
		lock.Lock()
		instance = &TimeCut{}
		instance.timeHistoryMap = map[string][]int64{}
		defer lock.Unlock()
	}
	if instance.timeHistoryMap[tag] == nil {
		instance.timeHistoryMap[tag] = []int64{}
	}
}

func Start(tag string, subTag string) {
	initTag(tag)
	t := time.Now().UnixNano()
	instance.timeHistoryMap[tag] = append(instance.timeHistoryMap[tag], t)
	fmt.Printf("[%s] - [%s] start %d \n", tag, subTag, t)
}

func Cut(tag string, subTag string) {
	initTag(tag)
	t := time.Now().UnixNano()
	instance.timeHistoryMap[tag] = append(instance.timeHistoryMap[tag], t)
	if len(instance.timeHistoryMap[tag]) < 2 {
		return
	}
	tagHistory := instance.timeHistoryMap[tag]
	useTime := tagHistory[len(tagHistory)-1] - tagHistory[len(tagHistory)-2]
	fmt.Printf("[%s] - [%s] use time %d %d(ms)  \n", tag, subTag, useTime, useTime/1000000)
}
