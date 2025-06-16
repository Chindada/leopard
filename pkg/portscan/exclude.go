package portscan

import "strconv"

type ExcludePort struct {
	StartPort int64
	EndPort   int64
}

type ExcludePortArr []ExcludePort

func (e ExcludePortArr) IsExcluded(port string) bool {
	for _, v := range e {
		n, err := strconv.ParseInt(port, 10, 64)
		if err != nil {
			break
		}
		if v.StartPort <= n && n <= v.EndPort {
			return true
		}
	}
	return false
}

type PortScan struct {
	ExcludePorts []ExcludePort
}

func NewPortScan() *PortScan {
	return &PortScan{}
}
