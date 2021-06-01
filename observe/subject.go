package observe

type Subject interface {
	Register(o Observer)
	Deregister(o Observer) error
	NotifyObservers()
}

type Observer interface {
	Update(name, status string)
	GetId() int
}
