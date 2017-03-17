package gochat

import (
	"bufio"
	"io"
)

type client struct {
	*bufio.Reader
	*bufio.Writer
	wc chan string
}

func StartClient(msgCh chan<-string, cn io.ReadWriteCloser, quit chan struct{}) (chan<- string, chan struct{}) {
	c := new(client)
	c.Reader = bufio.NewReader(cn)
	c.Writer = bufio.NewWriter(cn)
	c.wc = make(chan string)
	done := make(chan struct{})
	// setup the reader
	go func() {
		scanner := bufio.NewScanner(c.Reader)
		for scanner.Scan(){
			logger.Println(scanner.Text())
			msgCh <- scanner.Text()
		}
		done <- struct{}{}
	}()
	// setup the writer
	c.writerMonitor()
	go func() {
		select {
		case <-quit:
			cn.Close()
		case <-done:
		}
	}()
	return c.wc, done
}

func (c *client)writerMonitor() {
	go func() {
		for s := range c.wc {
			c.WriteString(s + "\n")
			c.Flush()
		}
	}()
}