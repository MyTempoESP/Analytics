package main

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
				" %d lbl %d num"+
				" %d lbl %d val",

			flick.PORTAL, nome,
			flick.IP, ip[0], ip[1], ip[2], ip[3],
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
