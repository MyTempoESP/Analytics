package main

import (
	"log"

	"github.com/MyTempoesp/flick"
)

type DisplayInfo struct {
	tags_unica, tags_total int64
	nome_equip, comm_verif string
	addr_equip, read_verif string
	wifi_verif, lt4g_verif string
	nome_prova, nome_local string
}

type SerialDisplay struct {
	Info   <-chan DisplayInfo
	Forth  flick.SerialForth
	Screen int

	switchButtonToggled bool
}

func NewSerialDisplay() (display SerialDisplay, err error) {

	f, err := flick.NewSerialForth("/dev/ttyUSB0")

	if err != nil {

		log.Printf("Erro ao iniciar a comunicação com o arduino: %v\n", err)

		return
	}

	display.Forth = f

	return
}

func (display *SerialDisplay) SwitchScreens() {

	res, err := display.Forth.Query(".")

	if err != nil {

		return
	}

	if res[0] == '0' && !display.switchButtonToggled {

		display.Screen++
		display.Screen %= 2

		display.switchButtonToggled = true 
	}

	if res[0] == '1' && display.switchButtonToggled {

		display.switchButtonToggled = false
	}
}

