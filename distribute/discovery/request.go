package discovery

type RequestRegister struct {
	Env             string   `form:"env"`
	AppId           string   `form:"appid"`
	Hostname        string   `form:"hostname"`
	Addrs           []string `form:"addrs[]"`
	Status          uint32   `form:"status"`
	Version         string   `form:"version"`
	LatestTimestamp int64    `form:"latest_timestamp"`
	DirtyTimestamp  int64    `form:"dirty_timestamp"`
	Replication     bool     `form:"replication"`
}
