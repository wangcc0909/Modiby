package pool

import (
	"errors"
	"fmt"
	"log"
	"sync"
	"time"
)

var (
	//MaxActiveConnReached 连接池超限
	ErrMaxActiveConnReached = errors.New("MaxActiveConnReached")
)

//config 连接池相关配置
type Config struct {
	//连接池中拥有的最小连接数
	InitialCap int
	//最大并发存活连接数
	MaxCap int
	//最大空闲连接数
	MaxIdle int
	//生成连接的方法
	Factory func() (interface{}, error)
	//关闭连接的方法
	Close func(interface{}) error
	//检查连接是否有效的方法
	Ping func(interface{}) error
	//连接最大空闲时间，超过该时间则失效
	IdleTimeout time.Duration
}

type idleConn struct {
	conn interface{}
	t    time.Time
}

type connReq struct {
	idleConn *idleConn
}

//channelPool 存放连接信息
type channelPool struct {
	mu                       sync.Mutex
	conns                    chan *idleConn
	factory                  func() (interface{}, error)
	close                    func(interface{}) error
	ping                     func(interface{}) error
	idleTimeout, waitTimeout time.Duration
	maxActive                int
	openingConns             int
	connReqs                 []chan connReq
}

func NewChannelPool(poolConfig *Config) (Pool, error) {
	if !(poolConfig.InitialCap <= poolConfig.MaxIdle && poolConfig.MaxCap >= poolConfig.MaxIdle && poolConfig.InitialCap >= 0) {
		return nil, errors.New("invalid capacity settings")
	}
	if poolConfig.Factory == nil {
		return nil, errors.New("invalid factory func settings")
	}
	if poolConfig.Close == nil {
		return nil, errors.New("invalid close func settings")
	}
	c := &channelPool{
		conns:        make(chan *idleConn, poolConfig.MaxCap),
		factory:      poolConfig.Factory,
		close:        poolConfig.Close,
		idleTimeout:  poolConfig.IdleTimeout,
		maxActive:    poolConfig.MaxCap,
		openingConns: poolConfig.InitialCap,
	}
	if poolConfig.Ping != nil {
		c.ping = poolConfig.Ping
	}
	for i := 0; i < poolConfig.InitialCap; i++ {
		conn, err := c.factory()
		if err != nil {
			c.Release()
			return nil, fmt.Errorf("factory is not able to fill the pool: %s", err)
		}
		c.conns <- &idleConn{conn: conn, t: time.Now()}
	}
	return c, nil
}

func (c *channelPool) getConns() chan *idleConn {
	c.mu.Lock()
	defer c.mu.Unlock()
	conns := c.conns
	return conns
}

//从pool获取一个连接
func (c *channelPool) Get() (interface{}, error) {
	conns := c.getConns()
	if conns == nil {
		return nil, ErrClosed
	}
	for {
		select {
		case wrapConn := <-conns:
			if wrapConn == nil {
				return nil, ErrClosed
			}
			//判断时候超时，超时则丢弃
			if timeout := c.idleTimeout; timeout > 0 {
				if wrapConn.t.Add(timeout).Before(time.Now()) {
					//丢弃并关闭该连接
					c.Close(wrapConn.conn)
					continue
				}
			}
			//判断是否失效，失效则丢弃 如果用户没有设定ping方法，就不检查
			if c.ping != nil {
				if err := c.Ping(wrapConn.conn); err != nil {
					c.Close(wrapConn.conn)
					continue
				}
			}
			return wrapConn.conn, nil
		default:
			c.mu.Lock()
			log.Printf("openConn %v %v \n", c.openingConns, c.maxActive)
			if c.openingConns >= c.maxActive {
				req := make(chan connReq, 1)
				c.connReqs = append(c.connReqs, req)
				c.mu.Unlock()
				ret, ok := <-req
				if !ok {
					return nil, ErrMaxActiveConnReached
				}
				if timeout := c.idleTimeout; timeout > 0 {
					if ret.idleConn.t.Add(timeout).Before(time.Now()) {
						//丢弃并关闭该连接
						c.Close(ret.idleConn.conn)
						continue
					}
				}
				return ret.idleConn.conn, nil
			}
			if c.factory == nil {
				c.mu.Unlock()
				return nil, ErrClosed
			}
			conn, err := c.factory()
			if err != nil {
				c.mu.Unlock()
				return nil, err
			}
			c.openingConns++
			c.mu.Unlock()
			return conn, nil
		}
	}
}

func (c *channelPool) Put(conn interface{}) error {
	if conn == nil {
		return errors.New("connection is nil.rejecting")
	}
	c.mu.Lock()
	if c.conns == nil {
		c.mu.Unlock()
		return c.Close(conn)
	}
	if l := len(c.connReqs); l > 0 {
		req := c.connReqs[0]
		copy(c.connReqs, c.connReqs[1:])
		c.connReqs = c.connReqs[:l-1]
		req <- connReq{
			idleConn: &idleConn{conn: conn, t: time.Now()},
		}
		c.mu.Unlock()
		return nil
	} else {
		select {
		case c.conns <- &idleConn{conn: conn, t: time.Now()}:
			c.mu.Unlock()
			return nil
		default:
			c.mu.Unlock()
			//连接池已满 直接关闭该连接
			return c.Close(conn)
		}
	}
}

func (c *channelPool) Close(conn interface{}) error {
	if conn == nil {
		return errors.New("connection is nil.rejecting")
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.close == nil {
		return nil
	}
	c.openingConns--
	return c.close(conn)
}

func (c *channelPool) Ping(conn interface{}) error {
	if conn == nil {
		return errors.New("connection is nil, rejecting")
	}
	return c.ping(conn)
}

func (c *channelPool) Release() {
	c.mu.Lock()
	conns := c.conns
	c.conns = nil
	c.factory = nil
	c.ping = nil
	closeFunc := c.close
	c.close = nil
	c.mu.Unlock()

	if conns == nil {
		return
	}
	close(conns)
	for conn := range conns {
		closeFunc(conn)
	}
}

func (c *channelPool) Len() int {
	return len(c.getConns())
}
