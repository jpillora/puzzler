package x

import (
	"net"
	"time"
)

func Online(domain string) bool {
	d := net.Dialer{Timeout: 3 * time.Second}
	conn, err := d.Dial("tcp", domain+":80")
	if err == nil {
		conn.Close()
		return true
	}
	return false
}
