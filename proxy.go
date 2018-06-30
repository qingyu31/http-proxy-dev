package main

import (
	"bufio"
	"net"
	"net/http"
	"net/http/httputil"
	"strconv"
	"time"
)

const (
	LogRequestMessage  = 0x1
	LogResponseMessage = 0x2
)

//Proxy is a local reverse proxy that prints out http message.
type Proxy struct {
}

func (this *Proxy) ListenAndServe() error {
	ln, err := net.Listen("tcp", ":"+strconv.Itoa(*port))
	if err != nil {
		return err
	}
	var tempDelay time.Duration
	for {
		conn, e := ln.Accept()
		if e != nil {
			if ne, ok := e.(net.Error); ok && ne.Temporary() {
				if tempDelay == 0 {
					tempDelay = 5 * time.Millisecond
				} else {
					tempDelay *= 2
				}
				if max := 1 * time.Second; tempDelay > max {
					tempDelay = max
				}
				logger.Printf("Proxy: Accept error: %v; retrying in %v", e, tempDelay)
				time.Sleep(tempDelay)
				continue
			}
			return e
		}
		tempDelay = 0
		go this.serveConn(conn)
	}
}

func (this *Proxy) serveConn(conn net.Conn) {
	defer conn.Close()
	rd := bufio.NewReader(conn)
	req, err := http.ReadRequest(rd)
	if err != nil {
		logger.Printf("serveConn: %v", err)
		return
	}
	reqm, _ := httputil.DumpRequest(req, *logMessageSwitch&LogRequestMessage > 0)
	logger.Println(string(reqm))
	nreq, err := http.NewRequest(req.Method, *towards, req.Body)
	if err != nil {
		logger.Printf("serveConn: %v", err)
		return
	}
	for k, v := range req.Header {
		nreq.Header[k] = v
	}
	client := new(http.Client)
	res, err := client.Do(nreq)
	if err != nil {
		logger.Printf("serveConn: %v", err)
		return
	}
	defer res.Body.Close()
	resm, _ := httputil.DumpResponse(res, *logMessageSwitch&LogResponseMessage > 0)
	logger.Println(string(resm))
	err = res.Write(conn)
	if err != nil {
		logger.Printf("serveConn: %v", err)
	}
}
