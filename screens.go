package main

import (
	"fmt"
)

const (
	SCREEN_TAGS = iota
	SCREEN_ADDR
	SCREEN_WIFI
	SCREEN_STAT
)

func (display *SerialDisplay) ScreenTags(data DisplayInfo) {

	nome_equip := data.nome_equip
	comm_verif := data.comm_verif
	tags_unica := fmt.Sprintf("UNICAS   %d", data.tags_unica)
	tags_total := fmt.Sprintf("REGIST.  %d", data.tags_total)

	display.forth.Run(
		fmt.Sprintf("%s %s %s %s SCX", // Call SCX with screen data
			display.forth.GetBytes(nome_equip),
			display.forth.GetBytes(tags_unica),
			display.forth.GetBytes(tags_total),
			display.forth.GetBytes(comm_verif),
		),
	)
}

func (display *SerialDisplay) ScreenAddr(data DisplayInfo) {

	nome_equip := data.nome_equip
	comm_verif := data.comm_verif
	addr_equip := data.addr_equip
	read_verif := data.read_verif

	display.forth.Run(
		fmt.Sprintf("%s %s %s %s SCX", // Call SCX with screen data
			display.forth.GetBytes(nome_equip),
			display.forth.GetBytes(addr_equip),
			display.forth.GetBytes(read_verif),
			display.forth.GetBytes(comm_verif),
		),
	)
}
