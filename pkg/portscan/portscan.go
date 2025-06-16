// Package portscan package portscan
package portscan

import (
	"crypto/rand"
	"math/big"
	"net"
	"strconv"
	"time"
)

const (
	_target = "127.0.0.1"
)

func (p *PortScan) GetRandomPort() string {
	randomPort, _ := rand.Int(rand.Reader, big.NewInt((30000)))
	if !p.GetPortIsUsed(randomPort.String()) {
		port := randomPort.Int64() + 30000
		portStr := strconv.FormatInt(port, 10)
		exclude, err := p.Scan()
		if err != nil {
			return ""
		}
		if exclude.IsExcluded(portStr) {
			return p.GetRandomPort()
		}
		return strconv.FormatInt(port, 10)
	}
	return p.GetRandomPort()
}

func (p *PortScan) GetPortIsUsed(port string) bool {
	var conn net.Conn
	var err error
	defer func() {
		if conn != nil {
			if err = conn.Close(); err != nil {
				return
			}
		}
	}()
	conn, err = net.DialTimeout("tcp", net.JoinHostPort(_target, port), 500*time.Millisecond)
	if err != nil {
		return false
	}
	if conn != nil {
		return true
	}
	return false
}
