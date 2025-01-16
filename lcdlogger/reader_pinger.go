package lcdlogger

import (
	"github.com/tatsushid/go-fastping"
	"log"
	"net"
	"os"
	"time"
)

func NewSimplePinger(ip string) (p *fastping.Pinger, err error) {

	p = fastping.NewPinger()
	ra, err := net.ResolveIPAddr("ip4:icmp", ip)

	if err != nil {

		return
	}

	p.AddIPAddr(ra)

	p.MaxRTT = time.Second * 2

	err = p.Run()

	if err != nil {

		return
	}

	return
}

type ReaderPinger struct {
	pinger *fastping.Pinger
	ip     string

	Octets [4]int
	State  bool
	Ping   int
}

func NewReaderPinger() (r ReaderPinger, err error) {

	r.ip = os.Getenv("READER_IP")
	r.Octets = IPIfy(r.ip)

	p, err := NewSimplePinger(r.ip)

	if err != nil {

		return
	}

	p.OnRecv = func(addr *net.IPAddr, rtt time.Duration) {
		log.Printf("IP Addr: %s receive, RTT: %v\n", addr.String(), rtt)

		r.State = true
		r.Ping = int(rtt / time.Millisecond)
	}

	p.OnIdle = func() { r.State = false }

	r.pinger = p

	return
}
