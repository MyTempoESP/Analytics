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

func (display *SerialDisplay) S(r1, r2, r3, r4 string, l1, v1, l2, v2, l3, v3, l4, v4 int64) {

	display.Forth.Send(
		fmt.Sprintf(r1+" "+r2+" "+r3+" "+r4, l1, v1, l2, v2, l3, v3, l4, v4),
	)
}

func (display *SerialDisplay) ScreenTags(data DisplayInfo) {

	nome_equip := int64(data.nome_equip)
	comm_verif := int64(data.comm_verif)
	tags_unica := data.tags_unica
	tags_total := data.tags_total

	display.S(
		"%d lbl %d num",
		"%d lbl %d num",
		"%d lbl %d num",
		"%d lbl %d val",

		0, nome_equip,
		flick.REGIST, tags_total,
		flick.UNICAS, tags_unica,
		flick.COMUNICANDO, comm_verif,
	)
}

func (display *SerialDisplay) ScreenAddr(data DisplayInfo) {

	nome_equip := int64(data.nome_equip)
	comm_verif := int64(data.comm_verif)
	addr_equip := int64(data.addr_equip)
	read_verif := int64(data.read_verif)

	display.S(
		"%d lbl %d num",
		"%d lbl %d num",
		"%d lbl %d num",
		"%d lbl %d val",

		flick.PORTAL, nome_equip,
		flick.IP, addr_equip,
		flick.LEITOR, read_verif,
		flick.COMUNICANDO, comm_verif,
	)
}
