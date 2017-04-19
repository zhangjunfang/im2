package tcp

import (
	"errors"
	"io"
	"net"
	"sync"
	"time"
)

//支持的最大消息长度
const maxLength int = 1<<32 - 1 // 4294967295

var (
	rHeadBytes = [4]byte{0, 0, 0, 0}
	wHeadBytes = [4]byte{0, 0, 0, 0}
	errMsgRead = errors.New("Message read length error")
	errHeadLen = errors.New("Message head length error")
	errMsgLen  = errors.New("Message length is no longer in normal range")
)
var connPool sync.Pool

//从对象池中获取一个对象,不存在则申明
func Newconnection(conn net.Conn) Conn {
	c := connPool.Get()
	if cnt, ok := c.(*connection); ok {
		cnt.rwc = conn
		return cnt
	}
	return &connection{rlen: 0, rwc: conn}
}

//定义一个结构体,用来封装Conn.
type connection struct {
	rlen  int        //消息长度
	rwc   net.Conn   //原始的网络链接
	rlock sync.Mutex //Conn读锁
	wlock sync.Mutex //Conn写锁
}

type Conn interface {
	Read() (r string, err error)
	Write(p []byte) (n int, err error)
	Writer(size int, r io.Reader) (n int, err error)
	RemoteAddr() net.Addr
	LocalAddr() net.Addr
	SetDeadline(t time.Time) error
	SetReadDeadline(t time.Time) error
	SetWriteDeadline(t time.Time) error
	Close() (err error)
}

func (c *connection) Read() (r string, err error) {
	b := make([]byte, 64)
	var n int = 0
	for {
		c.rlock.Lock()
		temp, err := c.rwc.Read(b)
		if temp == 0 && err == nil {
			n = temp
			break
		}
		if err != nil {
			c.Close()
			c.rlock.Unlock()
			b = make([]byte, 64)
			n = 0
			continue
		}
		if "^" != b[n-1] { //是否读到数据结束符
			n = n + temp
			b = append(b, b[:n]...)
			c.rlock.Unlock()
			continue
		}
		r = string(b[:n])
		c.rlock.Unlock()
	}
	return r, n, err
}
func (c *connection) Write(p []byte) (n int, err error) {
	c.wlock.Lock()
	n, err = c.rwc.Write(p)
	c.wlock.Unlock()
	return
}
func (c *connection) Writer(size int, r io.Reader) (n int, err error) {
	b := make([]byte, size)
	c.wlock.Lock()
	r.Read(b)
	n, err = c.rwc.Write(b)
	c.wlock.Unlock()
	return
}
func (c *connection) RemoteAddr() net.Addr {
	return c.rwc.RemoteAddr()
}
func (c *connection) LocalAddr() net.Addr {
	return c.rwc.LocalAddr()
}
func (c *connection) SetDeadline(t time.Time) error {
	return c.rwc.SetDeadline(t)
}
func (c *connection) SetReadDeadline(t time.Time) error {
	return c.rwc.SetReadDeadline(t)
}
func (c *connection) SetWriteDeadline(t time.Time) error {
	return c.rwc.SetWriteDeadline(t)
}
func (c *connection) Close() (err error) {
	c.wlock.Lock()
	err = c.rwc.Close()
	c.rlen = 0
	connPool.Put(c)
	c.wlock.Unlock()
	return
}
