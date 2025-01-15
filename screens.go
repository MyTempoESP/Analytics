package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/MyTempoesp/flick"
)

const (
	SCREEN_TAGS = iota
	SCREEN_ADDR
	SCREEN_WIFI
	SCREEN_STAT
)

type ForthNumber struct {
	Value     int64
	Magnitude int // 1, 10, 100, 1000 (10^Magnitude)
}

func IPIfy(ip string) (out [4]int, err error) {

	parts := strings.Split(ip, ".")

	if len(parts) != 4 {

		err = fmt.Errorf("invalid IP address format: %s", ip)

		return
	}

	var result []int
	var num int

	for _, part := range parts {

		num, err = strconv.Atoi(part)

		if err != nil {

			err = fmt.Errorf("invalid number in IP address: %s", part)

			return
		}

		if num < 0 || num > 255 {

			fmt.Errorf("IP address octet out of range: %d", num)

			return
		}

		result = append(result, num)
	}

	return
}

func ToForthNumber(n int64) (f ForthNumber) {

	if n < 1000 {

		f.Value = n
		f.Magnitude = 0

		return
	}

	if n < 1_000_000 {

		f.Value = n / 1000
		f.Magnitude = 3

		return
	}

	f.Value = n / 1_000_000
	f.Magnitude = 6

	return
}

func (display *SerialDisplay) ScreenTags(nome, commVerif int, tags, tagsUnicas int64) {

	display.Forth.Send(
		fmt.Sprintf(
			"%d lbl %d num"+
				" %d lbl %d num"+
				" %d lbl %d num"+
				" %d lbl %d val",

			flick.PORTAL, nome,
			flick.REGIST, tags,
			flick.UNICAS, tagsUnicas,
			flick.COMUNICANDO, commVerif,
		),
	)
}

func (display *SerialDisplay) ScreenAddr(nome, commVerif int, ip [4]int, leitorOk int) {

	display.Forth.Send(
		fmt.Sprintf(
			"%d lbl %d num"+
				" %d lbl %d %d %d %d ip"+
				" %d lbl %d val"+
				" %d lbl %d val",

			flick.PORTAL, nome,
			flick.IP, ip[3], ip[2], ip[1], ip[0],
			flick.LEITOR, leitorOk,
			flick.COMUNICANDO, commVerif,
		),
	)
}

func (display *SerialDisplay) ScreenWifi(nome, commVerif, wifiVerif, LTE4GVerif int) {

	display.Forth.Send(
		fmt.Sprintf(
			"%d lbl %d num"+
				" %d lbl %d val"+
				" %d lbl %d val"+
				" %d lbl %d val",

			flick.PORTAL, nome,
			flick.WIFI, wifiVerif,
			flick.LTE4G, LTE4GVerif,
			flick.COMUNICANDO, commVerif,
		),
	)
}

func (display *SerialDisplay) ScreenStat(nome, commVerif int, a1, a2, a3, a4 ForthNumber) {

	display.Forth.Send(
		fmt.Sprintf(
			"%d lbl %d num"+
				" %d %d"+ // A4 Val+Mag
				" %d %d"+ // A3 Val+Mag
				" %d %d"+ // A2 Val+Mag
				" %d %d atn"+ // A1 Val+Mag then display
				" %d lbl %d val",

			flick.PORTAL, nome,
			a4.Value, a4.Magnitude,
			a3.Value, a3.Magnitude,
			a2.Value, a2.Magnitude,
			a1.Value, a1.Magnitude,
			flick.COMUNICANDO, commVerif,
		),
	)
}
