//go:build !windows

// Package portscan package portscan
package portscan

func (p *PortScan) Scan() (ExcludePortArr, error) {
	return p.ExcludePorts, nil
}
