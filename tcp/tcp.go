package tcp

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

//定义一个结构体,用来封装Conn.
type Connection struct {
	//rlen  int        //消息长度
	Rwc   net.Conn   //原始的网络链接
	Rlock sync.Mutex //Conn读锁
	Wlock sync.Mutex //Conn写锁
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

func (c Connection) Read() (r string, err error) {
	b := make([]byte, 64)
	var n int = 0
	for {
		c.Rlock.Lock()
		temp, err := c.Rwc.Read(b)
		if temp == 0 && err == nil {
			n = temp
			break
		}
		if err != nil {
			c.Close()
			c.Rlock.Unlock()
			b = make([]byte, 64)
			n = 0
			continue
		}
		fmt.Println("read:nnn:::", n)
		if '^' != b[n-1] { //是否读到数据结束符
			n = n + temp
			b = append(b, b[:n]...)
			c.Rlock.Unlock()
			continue
		}
		r = string(b[:n])
		c.Rlock.Unlock()
		fmt.Println("read:", r)
	}
	return r, err
}
func (c Connection) Write(p []byte) (n int, err error) {
	c.Wlock.Lock()
	n, err = c.Rwc.Write(p)
	c.Wlock.Unlock()
	fmt.Println("write:", string(p))
	return
}
func (c *Connection) Writer(size int, r io.Reader) (n int, err error) {
	b := make([]byte, size)
	c.Wlock.Lock()
	r.Read(b)
	n, err = c.Rwc.Write(b)
	c.Wlock.Unlock()
	return
}
func (c *Connection) RemoteAddr() net.Addr {
	return c.Rwc.RemoteAddr()
}
func (c *Connection) LocalAddr() net.Addr {
	return c.Rwc.LocalAddr()
}
func (c *Connection) SetDeadline(t time.Time) error {
	return c.Rwc.SetDeadline(t)
}
func (c *Connection) SetReadDeadline(t time.Time) error {
	return c.Rwc.SetReadDeadline(t)
}
func (c *Connection) SetWriteDeadline(t time.Time) error {
	return c.Rwc.SetWriteDeadline(t)
}
func (c *Connection) Close() (err error) {
	c.Wlock.Lock()
	err = c.Rwc.Close()
	c.Wlock.Unlock()
	return
}
