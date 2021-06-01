package breaker

import (
	"errors"
	"fmt"
	"runtime"
	"sync"
	"testing"
	"time"
)

func initBreaker() *ServiceBreaker {
	tripOp := TripStrategyOption{
		Strategy: FailRateTrip,
		FailRate: 0.6,
		MinCall:  3,
	}
	option := Option{
		Name:            "breaker1",
		WindowInterval:  5 * time.Second,
		HalfMaxCalls:    3,
		SleepTimeout:    6 * time.Second,
		TripStrategy:    tripOp,
		StateChangeHook: stateChangeHook,
	}
	breaker, _ := NewServiceBreaker(option)
	return breaker
}

func stateChangeHook(name string, fromState State, toState State) {
	fmt.Printf("熔断器%v 触发状态变更：%v --> %v\n", name, fromState, toState)
}

func TestServiceBreaker(t *testing.T) {
	breaker := initBreaker()
	for i := 0; i < 30; i++ {
		breaker.Call(func() (interface{}, error) {
			if i <= 2 || i >= 8 {
				fmt.Println("请求执行成功")
				return nil, nil
			} else {
				fmt.Println("请求执行出错")
				return nil, errors.New("error")
			}
		})
		time.Sleep(1 * time.Second)
	}
}

func TestServiceBreakerInParallel(t *testing.T) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	breaker := initBreaker()
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(group *sync.WaitGroup) {
			defer group.Done()
			for i := 0; i < 30; i++ {
				breaker.Call(func() (interface{}, error) {
					if i <= 2 || i >= 8 {
						fmt.Println("请求执行成功")
						return nil, nil
					} else {
						fmt.Println("请求执行出错")
						return nil, errors.New("error")
					}
				})
				time.Sleep(1 * time.Second)
			}
		}(&wg)
	}
	wg.Wait()
}
