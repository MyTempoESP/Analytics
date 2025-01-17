package lcdlogger

import (
	"log"
	"os"
	"sync/atomic"
	"time"

	"github.com/prometheus-community/pro-bing"
)

func NewSimplePinger(ip string) (p *probing.Pinger, err error) {

	p, err = probing.NewPinger(ip)

	//p.SetPrivileged(true)

	if err != nil {

		return
	}

	p.Count = 0xFFFE // basically all we need
	//p.Size = *size
	p.Interval = 4 * time.Second
	//p.Timeout = p.Interval
	//p.TTL = *ttl
	//p.InterfaceName = *iface
	//p.SetPrivileged(*privileged)
	//p.SetTrafficClass(uint8(*tclass))

	return
}

type ReaderPinger struct {
	Pinger *probing.Pinger
	ip     string

	Octets [4]int
	State  atomic.Bool
	Ping   atomic.Int64
}

func NewReaderPinger() (r ReaderPinger, err error) {

	r.ip = os.Getenv("READER_IP")
	r.Octets = IPIfy(r.ip)

	p, err := NewSimplePinger(r.ip)

	if err != nil {

		return
	}

	p.OnSend = func(pkt *probing.Packet) {

		log.Printf("IP Addr: %s\n", pkt.IPAddr)

		r.State.Store(false)
	}

	p.OnRecv = func(pkt *probing.Packet) {

		log.Printf("IP Addr: %s receive, RTT: %v\n", pkt.IPAddr, pkt.Rtt)

		r.State.Store(true)
		r.Ping.Store(pkt.Rtt.Milliseconds())
	}

	p.Run()

	r.Pinger = p

	return
}
