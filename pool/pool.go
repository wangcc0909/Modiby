package pool

import "errors"

var (
	//ErrClosed 连接池已经关闭的error
	ErrClosed = errors.New("pool is closed")
)

//Pool 基本方法
type Pool interface {
	Get() (interface{}, error)
	Put(interface{}) error
	Close(interface{}) error
	Release()
	Len() int
}

/*type ConnRes interface {
	Close() error
}

type Factory func() (ConnRes, error)

type Conn struct {
	conn ConnRes
	time time.Time
}

func (c Conn) Close() error {

	return nil
}

//连接池
type ConnPool struct {
	mu          sync.Mutex
	conns       chan *Conn
	factory     Factory
	closed      bool
	connTimeOut time.Duration
}

func NewConnPool(factory Factory, cap int, connTimeOut time.Duration) (*ConnPool, error) {
	if cap <= 0 {
		return nil,errors.New("cap不能小于0")
	}
	if connTimeOut <= 0 {
		return nil, errors.New("connTimeOut不能小于0")
	}
	cp := &ConnPool{
		mu: sync.Mutex{},
		conns: make(chan *Conn,cap),
		factory: factory,
		closed: false,
		connTimeOut: connTimeOut,
	}
	for i := 0; i < cap; i++ {
		connRes,err := cp.factory()
		if err != nil {
			cp.Close()
			return nil, errors.New("factory出错")
		}
		cp.conns <- &Conn{connRes, time.Now()}
	}
	return cp, nil
}

//获取连接资源
func (cp *ConnPool) Get() (ConnRes, error) {
	if cp.closed {
		return nil,errors.New("连接池已关闭")
	}
	for {
		select {
		case connRes,ok := <-cp.conns:
			if !ok {
				return nil, errors.New("连接池已关闭")
			}
			if time.Now().Sub(connRes.time) > cp.connTimeOut {
				connRes.conn.Close()
				continue
			}
			return connRes, nil
		default:
			connRes,err := cp.factory()
			if err != nil {
				return nil, err
			}
			return connRes, nil
		}
	}
}

func (cp *ConnPool) Put(conn ConnRes) error {
	if cp.closed {
		return errors.New("连接池已关闭")
	}
	select {
	case cp.conns <- &Conn{conn: conn,time: time.Now()}:
		return nil
	default:
		conn.Close()
		return errors.New("连接池已满")
	}
}

//关闭连接池
func (cp *ConnPool) Close() {
	if cp.closed {
		return
	}
	cp.mu.Lock()
	defer cp.mu.Unlock()
	cp.closed = true
	close(cp.conns)
	for conn := range cp.conns {
		conn.conn.Close()
	}
}*/
