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

	display.Forth.Run(
		fmt.Sprintf("%s %s %s %s SCX", // Call SCX with screen data
			display.Forth.GetBytes(nome_equip),
			display.Forth.GetBytes(tags_unica),
			display.Forth.GetBytes(tags_total),
			display.Forth.GetBytes(comm_verif),
		),
	)
}

func (display *SerialDisplay) ScreenAddr(data DisplayInfo) {

	nome_equip := data.nome_equip
	comm_verif := data.comm_verif
	addr_equip := data.addr_equip
	read_verif := data.read_verif

	display.Forth.Run(
		fmt.Sprintf("%s %s %s %s SCX", // Call SCX with screen data
			display.Forth.GetBytes(nome_equip),
			display.Forth.GetBytes(addr_equip),
			display.Forth.GetBytes(read_verif),
			display.Forth.GetBytes(comm_verif),
		),
	)
}
