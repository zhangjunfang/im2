package tcp

import (
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

func Dial(p, ip, port string) {
	conn, err := net.Dial(p, fmt.Sprintf("%s:%s", ip, port))
	if err != nil {
		fmt.Println(err)
		return
	}
	c := Newconnection(conn)
	defer c.Close()
	c.Write([]byte("Test"))
	c.Write([]byte("Test"))

	r, size, err := c.Read()
	if err != nil {
		fmt.Println(err, size)
		return
	}
	_, err = io.Copy(os.Stdout, r)
	if err != nil && err != io.EOF {
		fmt.Println(err)
	}
}
func Listener(proto, addr string) {
	lis, err := net.Listen(proto, addr)
	if err != nil {
		panic("Listen port error:" + err.Error())
		return
	}
	defer lis.Close()
	for {
		conn, err := lis.Accept()
		if err != nil {
			time.Sleep(1e7)
			continue
		}
		go handler(conn)
	}
}
func handler(conn net.Conn) {
	c := Newconnection(conn)
	msgchan := make(chan struct{})
	defer c.Close()
	go func(ch chan struct{}) {
		<-msgchan
		f, _ := os.Open("tcp_test.go")
		defer f.Close()
		info, _ := f.Stat()
		c.Writer(int(info.Size()), f)
		c.Close()
	}(msgchan)
	for {
		r, size, err := c.Read()
		if err != nil {
			fmt.Println(err)
			return
		}
		n, err := io.Copy(os.Stdout, r)
		if err != nil || n != int64(size) {
			if err == io.EOF {
				continue
			}
			fmt.Println("读取数据失败:", err)
			return
		}
		time.Sleep(2e9)
		msgchan <- struct{}{}
	}
}
