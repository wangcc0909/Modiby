package timewheel

import (
	"container/list"
	"errors"
	"log"
	"time"
)

//时间轮算法
type TimeWheel struct {
	ticker       *time.Ticker
	interval     time.Duration
	buckets      []*list.List
	bucketSize   int
	currentPos   int
	callbackFunc func(interface{})
	stopChannel  chan bool
}

func New(interval time.Duration, bucketSize int, callbackFunc func(interface{})) (*TimeWheel, error) {
	if interval <= 0 || bucketSize <= 0 || callbackFunc == nil {
		return nil, errors.New("create timewheel instance fail")
	}
	tw := &TimeWheel{
		interval:     interval,
		buckets:      make([]*list.List, bucketSize),
		bucketSize:   bucketSize,
		currentPos:   0,
		callbackFunc: callbackFunc,
		stopChannel:  make(chan bool),
	}
	for i := 0; i < bucketSize; i++ {
		tw.buckets[i] = list.New()
	}
	return tw, nil
}

type Task struct {
	Id     interface{}
	Data   interface{}
	Delay  time.Duration
	Circle int // task position in timewheel
}

//add task
func (tw *TimeWheel) AddTask(task *Task) {
	delaySeconds := int(task.Delay.Seconds())
	intervalSeconds := int(tw.interval.Seconds())
	circle := int(delaySeconds / intervalSeconds / tw.bucketSize)
	pos := int(tw.currentPos+delaySeconds/intervalSeconds) % tw.bucketSize
	task.Circle = circle
	tw.buckets[pos].PushBack(task)
}

func (tw *TimeWheel) Start() {
	tw.ticker = time.NewTicker(tw.interval)
	go func() {
		for {
			select {
			case <-tw.ticker.C:
				log.Println("1 tick")
				tw.tickHandler()
			case <-tw.stopChannel:
				tw.ticker.Stop()
				return
			}
		}
	}()
}

func (tw *TimeWheel) tickHandler() {
	bucket := tw.buckets[tw.currentPos]
	for e := bucket.Front(); e != nil; {
		task := e.Value.(*Task)
		if task.Circle > 0 {
			task.Circle--
			e = e.Next()
			continue
		}
		go tw.callbackFunc(task.Data)
		next := e.Next()
		bucket.Remove(e)
		e = next
	}
	if tw.currentPos == tw.bucketSize-1 {
		log.Println("new circle")
		tw.currentPos = 0
	} else {
		tw.currentPos++
	}
}
