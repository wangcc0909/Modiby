package discovery

type FetchData struct {
	Instances       []*Instance `json:"instances"`
	LatestTimestamp int64       `json:"latest_timestamp"`
}
