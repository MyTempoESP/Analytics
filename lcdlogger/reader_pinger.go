package lcdlogger

import (
	"github.com/prometheus-community/pro-bing"
	"log"
	"os"
	"sync/atomic"
	"time"
)

func NewSimplePinger(ip string) (p *probing.Pinger, err error) {

	p, err = probing.NewPinger(ip)

	p.SetPrivileged(true)

	if err != nil {

		return
	}

	//p.Size = *size
	p.Interval = 2 * time.Second
	p.Timeout = p.Interval * 2
	//p.TTL = *ttl
	//p.InterfaceName = *iface
	//p.SetPrivileged(*privileged)
	//p.SetTrafficClass(uint8(*tclass))

	return
}

type ReaderPinger struct {
	pinger *probing.Pinger
	ip     string

	Octets [4]int
	State  atomic.Bool
	Ping   int32
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
		atomic.StoreInt32(&r.Ping, int32(pkt.Rtt))
	}

	p.Run()

	r.pinger = p

	return
}
