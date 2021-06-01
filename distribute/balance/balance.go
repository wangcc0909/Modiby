package balance

type LoadBalance interface {
	Add(...string) error
	Get() (string, error)
}

const (
	Random = iota
	RoundRobin
	Weight
)

func LoadBalanceFactory(lbType int) LoadBalance {
	switch lbType {
	case Random:
		return new(RandomBalance)
	case RoundRobin:
		return new(RoundRobinBalance)
	case Weight:
		return new(WeightRoundRobinBalance)
	default:
		return new(RoundRobinBalance)
	}
}
