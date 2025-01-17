package main

import (
	"log"
	"sync/atomic"
	"time"

	"analytics/intSet"
	"analytics/lcdlogger"
	"github.com/MyTempoesp/flick"
)

const (
	QUERY_ATUALIZAR_TAGS_TOTAL = `
	INSERT INTO stats (
	    id,
	    tags_total,
	    tags_unicas
	)
	VALUES(
		1,
		0,
		0
	)
	ON DUPLICATE KEY
	UPDATE
		tags_total = ?
	`

	QUERY_ATUALIZAR_TAGS_UNICAS = `
	INSERT INTO stats (
	    id,
	    tags_unicas,
	    tags_total
	)
	VALUES(
		1,
		0,
		0
	)
	ON DUPLICATE KEY
	UPDATE
		tags_unicas = ?
	`

	QUERY_ATUALIZAR_ANTENAS = `
	INSERT INTO stats (
	    id,
	    tags_unicas,
	    tags_total,
	    Antena1,
	    Antena2,
	    Antena3,
	    Antena4
	)
	VALUES(
		1,
		0,
		0,
		0,
		0,
		0,
		0
	)
	ON DUPLICATE KEY
	UPDATE
		Antena1 = ?,
		Antena2 = ?,
		Antena3 = ?,
		Antena4 = ?
	`
)

func (a *Ay) AtualizarTagsUnicas(tagsUnicas int64) {

	_, err := a.db.Exec(QUERY_ATUALIZAR_TAGS_UNICAS, tagsUnicas)

	if err != nil {
		log.Println("(AtualizarTagsUnicas)", err)
	}
}

func (a *Ay) AtualizarAntenas(antenas *[4]atomic.Int64) {

	_, err := a.db.Exec(QUERY_ATUALIZAR_ANTENAS,
		antenas[0].Load(),
		antenas[1].Load(),
		antenas[2].Load(),
		antenas[3].Load(),
	)

	if err != nil {
		log.Println("(AtualizarTagsUnicas)", err)
	}
}

func (a *Ay) AtualizarTagsTotal(tags int64) {

	_, err := a.db.Exec(QUERY_ATUALIZAR_TAGS_TOTAL, tags)

	if err != nil {

		log.Println("(AtualizarTags)", err)
	}
}

func (a *Ay) AtualizarTags(tags int64) (ok bool) {

	var totalAnterior int64 = 0

	ok = true

	if tags <= 0 {

		return
	}

	res, err := a.db.Query("SELECT tags_total FROM stats")

	if err != nil {

		log.Println("(AtualizarTags)", err)

		return
	}

	defer res.Close()

	if res.Next() {

		err = res.Scan(&totalAnterior)

		if err != nil {

			log.Println("(AtualizarTags)", err)

			return
		}
	}

	if totalAnterior == 0 {

		a.AtualizarTagsTotal(1)

		ok = false

		return
	}

	a.AtualizarTagsTotal(tags)

	return
}

func (a *Ay) Process() {

	var (
		tags     atomic.Int64
		antennas [4]atomic.Int64
	)

	tagSet := intSet.New()

	go func() {

		for t := range a.Tags {

			if t.Antena == 0 {
				/*
					Antena 0 não exist
				*/

				continue
			}

			antennas[(t.Antena-1)%4].Add(1)

			tags.Add(1)

			if tagSet.Insert(t.Epc) {

				a.AtualizarTagsUnicas(tagSet.Count())
			}
		}
	}()

	go func() {

		atualizaContagem := time.NewTicker(2 * time.Second)

		for {
			<-atualizaContagem.C

			/*
				true  => prosseguir
				false => resetar
			*/
			log.Println("Atualizando tags")

			if !a.AtualizarTags(tags.Load()) {

				tagSet.Clear()

				tags.Store(0)

				for i := range 4 {

					antennas[i].Store(0)
				}
			}

			a.AtualizarAntenas(&antennas)
		}
	}()

	var readerState atomic.Bool
	var readerPing atomic.Int64
	var readerIP [4]int

	display, displayErr := lcdlogger.NewSerialDisplay()

	/* > Monitoring can be skipped if NewSerialDisplay() errors out, disabling the routine in Line 221 */
	if displayErr != nil {

		goto skip_monitoring
	}

	go func() {

		reader, readerErr := lcdlogger.NewReaderPinger()

		if readerErr != nil {

			return
		}

		readerIP = reader.Octets

		for {
			<-time.After(reader.Pinger.Interval)
			readerState.Store(reader.State.Load())
			readerPing.Store(reader.Ping.Load())

			log.Println(reader.State.Load())
			log.Println(reader.Ping.Load())
		}
	}()

	go func() {

		const NUM_EQUIP = 701

		for {

			comm_verif := flick.WEB

			switch display.Screen {
			case lcdlogger.SCREEN_TAGS:
				display.ScreenTags(
					NUM_EQUIP,
					comm_verif,
					/* Tags */ tags.Load(),
					/* Atletas */ tagSet.Count(),
				)
			case lcdlogger.SCREEN_ADDR:

				ip := readerIP
				leitor := flick.OK

				if !readerState.Load() {

					ip = [4]int{0, 0, 0, 0}
					leitor = flick.DESLIGAD
				}

				display.ScreenAddr(
					NUM_EQUIP,
					readerPing.Load(),
					/* IP */ ip,
					/* Leitor */ leitor,
				)
			case lcdlogger.SCREEN_WIFI:
				display.ScreenWifi(
					NUM_EQUIP,
					comm_verif,
					/* WIFI */ flick.CONECTAD,
					/* 4G */ flick.DESLIGAD,
				)
			case lcdlogger.SCREEN_STAT:
				display.ScreenStat(
					NUM_EQUIP,
					comm_verif,
					lcdlogger.ToForthNumber(antennas[0].Load()),
					lcdlogger.ToForthNumber(antennas[1].Load()),
					lcdlogger.ToForthNumber(antennas[2].Load()),
					lcdlogger.ToForthNumber(antennas[3].Load()),
				)
			}

			display.SwitchScreens()

			time.Sleep(100 * time.Millisecond)
		}
	}()

skip_monitoring:
	select {}
}
