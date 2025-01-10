package main

import (
	"fmt"
	"log"
	"time"

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
	Screen int

	forth flick.SerialForth
}

func NewSerialDisplay() (display SerialDisplay, err error) {

	f, err := flick.NewSerialForth()

	if err != nil {

		log.Printf("Erro ao iniciar a comunicação com o arduino: %v\n", err)

		return
	}

	display.forth = f

	return
}

func (display *SerialDisplay) Start() (info chan<- DisplayInfo) {

	infoExchange := make(chan DisplayInfo)

	display.Info = infoExchange

	go func() {

		display.forth.Run(": DRW 0 m $ d a ;")
		display.forth.Run(": SCX 3 FOR I DRW NXT 0 DRW ;")

		data := <-display.Info

		for {
			switch display.Screen {
			case SCREEN_TAGS:
				display.ScreenTags(data)
			case SCREEN_ADDR:
				display.ScreenAddr(data)
			}

			fmt.Println(display.forth.Query("6 IN 1 = ."))

			time.Sleep(500 * time.Millisecond)
		}
	}()

	info = infoExchange

	return
}
