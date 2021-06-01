package discovery

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"
)

var (
	NOT_FOUND = errors.New("Not Found")
)

const (
	RenewInterval               = 30 * time.Second   //client heart beat interval
	CheckEvictInterval          = 60 * time.Second   //evict task interval
	InstanceExpireDuration      = 90 * time.Second   //instance's renewTimestamp after this will be canceled
	InstanceMaxExpireDuration   = 3600 * time.Second //instance's renewTimestamp after this will be canceled
	ResetGuardNeedCountInterval = 15 * time.Minute   //ticker reset guard need count
	SelfProtectThreshold        = 0.85
)

type Registry struct {
	apps map[string]*Application
	gd   *Guard
	lock sync.RWMutex
}

type Application struct {
	appid           string
	instances       map[string]*Instance
	latestTimestamp int64
	lock            sync.RWMutex
}

func (app *Application) AddInstance(in *Instance, latestTimestamp int64) (*Instance, bool) {
	app.lock.Lock()
	defer app.lock.Unlock()
	appIns, ok := app.instances[in.Hostname]
	if ok {
		in.UpTimestamp = appIns.UpTimestamp
		if in.DirtyTimestamp < appIns.DirtyTimestamp {
			log.Println("register exist dirty timestamp")
			in = appIns
		}
	}
	app.instances[in.Hostname] = in
	app.upLatestTimestamp(latestTimestamp)
	returnIns := new(Instance)
	*returnIns = *in
	return returnIns, !ok
}

func (app *Application) upLatestTimestamp(latestTimestamp int64) {
	if latestTimestamp <= app.latestTimestamp {
		latestTimestamp = app.latestTimestamp + 1
	}
	app.latestTimestamp = latestTimestamp
}

/*
Env 服务环境标识，如 online、dev、test
AppId  应用服务的唯一标识

Hostname 服务实例的唯一标识

Addrs 服务实例的地址，可以是 http 或 rpc 地址，多个地址可以维护数组

Version 服务实例版本

Status 服务实例状态，用于控制上下线

xxTimestamp 依次记录服务实例注册时间戳，上线时间戳，最近续约时间戳，脏时间戳（后面解释），最后更新时间戳
*/

type Instance struct {
	Env      string   `json:"env"`
	AppId    string   `json:"appid"`
	Hostname string   `json:"hostname"`
	Addrs    []string `json:"addrs"`
	Version  string   `json:"version"`
	Status   uint32   `json:"status"`

	RegTimestamp    int64 `json:"reg_timestamp"`
	UpTimestamp     int64 `json:"up_timestamp"`
	RenewTimestamp  int64 `json:"renew_timestamp"`
	DirtyTimestamp  int64 `json:"dirty_timestamp"`
	LatestTimestamp int64 `json:"latest_timestamp"`
}

func NewRegistry() *Registry {
	registry := &Registry{
		apps: make(map[string]*Application),
		gd:   new(Guard),
	}
	go registry.evictTask()
	return registry
}

func NewInstance(req *RequestRegister) *Instance {
	now := time.Now().UnixNano()
	instance := &Instance{
		Env:             req.Env,
		AppId:           req.AppId,
		Hostname:        req.Hostname,
		Addrs:           req.Addrs,
		Version:         req.Version,
		Status:          req.Status,
		RegTimestamp:    now,
		UpTimestamp:     now,
		RenewTimestamp:  now,
		DirtyTimestamp:  now,
		LatestTimestamp: now,
	}
	return instance
}

func (r *Registry) Register(instance *Instance, latestTimestamp int64) (*Application, error) {
	key := getKey(instance.AppId, instance.Env)
	r.lock.RLock()
	app, ok := r.apps[key]
	r.lock.RUnlock()
	if !ok {
		app = NewApplication(instance.AppId)
	}
	_, isNew := app.AddInstance(instance, latestTimestamp)
	if isNew {
		r.gd.incrNeed()
	}
	r.lock.Lock()
	r.apps[key] = app
	r.lock.Unlock()
	return app, nil
}
func (r *Registry) Fetch(env, appid string, status uint32, latestTime int64) (*FetchData, error) {
	app, ok := r.getApplication(appid, env)
	if !ok {
		return nil, errors.New("Not found")
	}
	return app.GetInstance(status, latestTime)
}

func (r *Registry) getApplication(appid, env string) (*Application, bool) {
	key := getKey(appid, env)
	r.lock.RLock()
	app, ok := r.apps[key]
	r.lock.RUnlock()
	return app, ok
}

func (r *Registry) Cancel(env, appid string, hostname string, latestTimestamp int64) (*Instance, error) {
	app, ok := r.getApplication(appid, env)
	if !ok {
		return nil, errors.New("Not found")
	}
	instance, ok, insLen := app.Cancel(hostname, latestTimestamp)
	if !ok {
		return nil, errors.New("Not found")
	}
	if insLen == 0 {
		r.lock.Lock()
		delete(r.apps, getKey(appid, env))
		r.lock.Unlock()
	}
	r.gd.decrNeed()
	return instance, nil
}

func (r *Registry) Renew(env, appid, hostname string) (*Instance, error) {
	app, ok := r.getApplication(appid, env)
	if !ok {
		return nil, NOT_FOUND
	}
	in, ok := app.Renew(hostname)
	if !ok {
		return nil, NOT_FOUND
	}
	r.gd.incrCount()
	return in, nil
}

func (r *Registry) evictTask() {
	ticker := time.Tick(CheckEvictInterval)
	resetTikcer := time.Tick(ResetGuardNeedCountInterval)
	for {
		select {
		case <-ticker:
			r.gd.storeLastCount()
			r.evict()
		case <-resetTikcer:
			var count int64
			for _, app := range r.getAllApplications() {
				count += int64(app.GetInstanceLen())
			}
			r.gd.setNeed(count)
		}
	}
}

func (r *Registry) evict() {
	now := time.Now().UnixNano()
	var expireInstances []*Instance
	apps := r.getAllApplications()
	var registryLen int
	protectStatus := r.gd.selfProtectStatus()
	for _, app := range apps {
		registryLen += app.GetInstanceLen()
		allInstances := app.GetAllInstances()
		for _, instance := range allInstances {
			delta := now - instance.RenewTimestamp
			if !protectStatus && delta > int64(InstanceExpireDuration) || delta > int64(InstanceMaxExpireDuration) {
				expireInstances = append(expireInstances, instance)
			}
		}
	}
	evictionLimit := registryLen - int(float64(registryLen)*SelfProtectThreshold)
	expireLen := len(expireInstances)
	if expireLen > evictionLimit {
		expireLen = evictionLimit
	}
	if expireLen == 0 {
		return
	}
	for i := 0; i < expireLen; i++ {
		j := i + rand.Intn(len(expireInstances)-i)
		expireInstances[i], expireInstances[j] = expireInstances[j], expireInstances[i]
		expiredInstance := expireInstances[i]
		r.Cancel(expiredInstance.Env, expiredInstance.AppId, expiredInstance.Hostname, now)
	}
}

func (r *Registry) getAllApplications() []*Application {
	r.lock.Lock()
	defer r.lock.Unlock()
	apps := make([]*Application, 0, len(r.apps))
	for _, app := range r.apps {
		apps = append(apps, app)
	}
	return apps
}

func (app *Application) GetInstance(status uint32, latestTime int64) (*FetchData, error) {
	app.lock.RLock()
	defer app.lock.RUnlock()
	if latestTime >= app.latestTimestamp {
		return nil, errors.New("app not modified")
	}
	fetchData := FetchData{
		Instances:       make([]*Instance, 0),
		LatestTimestamp: app.latestTimestamp,
	}
	var exists bool
	for _, instance := range app.instances {
		if status&instance.Status > 0 {
			exists = true
			newInstance := copyInstance(instance)
			fetchData.Instances = append(fetchData.Instances, newInstance)
		}
	}
	if !exists {
		return nil, errors.New("Not found")
	}
	return &fetchData, nil
}

func (app *Application) Cancel(hostname string, latestTimestamp int64) (*Instance, bool, int) {
	newInstance := new(Instance)
	app.lock.Lock()
	defer app.lock.Unlock()
	appIn, ok := app.instances[hostname]
	if !ok {
		return nil, ok, 0
	}
	delete(app.instances, hostname)
	appIn.LatestTimestamp = latestTimestamp
	app.upLatestTimestamp(latestTimestamp)
	*newInstance = *appIn
	return newInstance, true, len(app.instances)
}

func (app *Application) Renew(hostname string) (*Instance, bool) {
	app.lock.Lock()
	defer app.lock.Unlock()
	appIn, ok := app.instances[hostname]
	if !ok {
		return nil, ok
	}
	appIn.RenewTimestamp = time.Now().UnixNano()
	return copyInstance(appIn), true
}

func (app *Application) GetInstanceLen() int {
	app.lock.Lock()
	instanceLen := len(app.instances)
	app.lock.Unlock()
	return instanceLen
}

func (app *Application) GetAllInstances() []*Instance {
	app.lock.RLock()
	defer app.lock.RUnlock()
	rs := make([]*Instance, 0, len(app.instances))
	for _, instance := range app.instances {
		newInstance := new(Instance)
		*newInstance = *instance
		rs = append(rs, newInstance)
	}
	return rs
}

func copyInstance(src *Instance) *Instance {
	dst := new(Instance)
	*dst = *src
	dst.Addrs = make([]string, len(src.Addrs))
	for i, addr := range src.Addrs {
		dst.Addrs[i] = addr
	}
	return dst
}

func NewApplication(appid string) *Application {
	return &Application{
		appid:     appid,
		instances: make(map[string]*Instance),
	}
}

func getKey(appid, env string) string {
	return fmt.Sprintf("%s-%s", appid, env)
}
