package lcdlogger

import (
	"log"

	"github.com/MyTempoesp/flick"
)

type DisplayInfo struct {
	tags_unica, tags_total int64
	addr_equip, read_verif int
	nome_equip, comm_verif int
	wifi_verif, lt4g_verif int
}

type SerialDisplay struct {
	Info   <-chan DisplayInfo
	Forth  *flick.MyTempo_Forth
	Screen int

	switchButtonToggled bool
}

func NewSerialDisplay() (display SerialDisplay, err error) {

	f, err := flick.NewForth("/dev/ttyUSB0")

	if err != nil {

		log.Printf("Erro ao iniciar a comunicação com o arduino: %v\n", err)

		return
	}

	f.Start()

	display.Forth = &f

	f.Send("VAR bac")
	f.Send("VAR bst")
	f.Send(": btn 7 IN 0 = ;")
	f.Send(": chb bac @ NOT IF bst @ btn DUP ROT SWP NOT AND bac ! bst ! THN ;")
	f.Send("10 0 TMI chb 1 TME")

	return
}

func (display *SerialDisplay) SwitchScreens() {

	// TODO: onrelease actions

	res, err := display.Forth.Send("bac @ .")
	defer display.Forth.Send("0 bac !")

	if err != nil {

		return
	}

	if res[0] == '-' && !display.switchButtonToggled {

		display.Screen++
		display.Screen %= 4

		display.switchButtonToggled = true
	}

	if res[0] == '0' && display.switchButtonToggled {

		display.switchButtonToggled = false
	}
}
