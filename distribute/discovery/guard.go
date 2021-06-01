package discovery

import (
	"sync"
	"sync/atomic"
)

/*
renewCount  记录所有服务续约次数，每执行一次 renew 加 1
lastRenewCount  记录上一次检查周期（默认 60 秒）服务续约统计次数

needRenewCount 记录一个周期总计需要的续约数，按一次续约 30 秒，一周期 60 秒，一个实例就需要 2 次，所以服务注册时 + 2，服务取消时 - 2

threshold  通过 needRenewCount  和阈值比例 （0.85）确定触发自我保护的值
*/

type Guard struct {
	renewCount     int64
	lastRenewCount int64
	needRenewCount int64
	threshold      int64
	lock           sync.RWMutex
}

func (gd *Guard) incrNeed() {
	gd.lock.Lock()
	defer gd.lock.Unlock()
	gd.needRenewCount += int64(CheckEvictInterval / RenewInterval)
	gd.threshold = int64(float64(gd.needRenewCount) * SelfProtectThreshold)
}

func (gd *Guard) decrNeed() {
	gd.lock.Lock()
	defer gd.lock.Unlock()
	gd.needRenewCount -= int64(CheckEvictInterval / RenewInterval)
	gd.threshold = int64(float64(gd.needRenewCount) * SelfProtectThreshold)
}

func (gd *Guard) setNeed(count int64) {
	gd.lock.Lock()
	defer gd.lock.Unlock()
	gd.needRenewCount = count * int64(CheckEvictInterval/RenewInterval)
	gd.threshold = int64(float64(gd.needRenewCount) * SelfProtectThreshold)
}

func (gd *Guard) incrCount() {
	atomic.AddInt64(&gd.renewCount, 1)
}

func (gd *Guard) storeLastCount() {
	atomic.StoreInt64(&gd.lastRenewCount, atomic.SwapInt64(&gd.needRenewCount, 0))
}

func (gd *Guard) selfProtectStatus() bool {
	return atomic.LoadInt64(&gd.lastRenewCount) < atomic.LoadInt64(&gd.threshold)
}
