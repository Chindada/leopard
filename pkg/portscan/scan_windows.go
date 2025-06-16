//go:build windows

// Package portscan package portscan
package portscan

import (
	"strconv"
	"strings"

	"github.com/chindada/leopard/pkg/command"
)

func (p *PortScan) Scan() (ExcludePortArr, error) {
	args := []string{
		"interface", "ipv4", "show", "excludedportrange", "protocol=tcp",
	}
	cmd := command.NewCMD("netsh", args...)
	result, err := command.RunAndParse(cmd)
	if err != nil {
		return []ExcludePort{}, err
	}
	for _, v := range result {
		split := strings.Fields(v)
		var startPort, endPort int64
		if count := len(split); count == 2 {
			for i := 0; i < count; i++ {
				n, err := strconv.ParseInt(split[i], 10, 64)
				if err != nil {
					break
				}

				switch i {
				case 0:
					startPort = n
				case 1:
					endPort = n
				}
			}
		}
		if startPort != 0 && endPort != 0 {
			p.ExcludePorts = append(p.ExcludePorts, ExcludePort{
				StartPort: startPort,
				EndPort:   endPort,
			})
		}
	}
	return p.ExcludePorts, nil
}
