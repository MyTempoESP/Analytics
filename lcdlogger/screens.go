package lcdlogger

import (
	"fmt"

	"github.com/MyTempoesp/flick"
)

const (
	SCREEN_TAGS = iota
	SCREEN_ADDR
	SCREEN_WIFI
	SCREEN_STAT
)

type IPOctets [4]int

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

func (display *SerialDisplay) ScreenAddr(nome, commVerif int, ping int64, ip IPOctets, leitorOk int) {

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
