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

func (a *Ay) AtualizarAntenas(antenas *[4]int64) {

	_, err := a.db.Exec(QUERY_ATUALIZAR_ANTENAS,
		atomic.LoadInt64(&antenas[0]),
		atomic.LoadInt64(&antenas[1]),
		atomic.LoadInt64(&antenas[2]),
		atomic.LoadInt64(&antenas[3]),
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
		tags int64 /* shared */

		antennas [4]int64 /* shared */
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

			atomic.AddInt64(&antennas[(t.Antena-1)%4], 1)

			atomic.AddInt64(&tags, 1)

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

			if !a.AtualizarTags(atomic.LoadInt64(&tags)) {

				tagSet.Clear()

				atomic.StoreInt64(&tags, 0)

				for i := range 4 {
					atomic.StoreInt64(&antennas[i], 0)
				}
			}

			a.AtualizarAntenas(&antennas)
		}
	}()

	display, displayErr := lcdlogger.NewSerialDisplay()
	reader, readerErr := lcdlogger.NewReaderPinger()

	/* > Monitoring can be skipped if NewSerialDisplay() errors out, disabling the routine in Line 221 */
	if displayErr != nil || readerErr != nil {

		goto skip_monitoring
	}

	go func() {

		const NUM_EQUIP = 701

		for {

			comm_verif := flick.WEB

			switch display.Screen {
			case lcdlogger.SCREEN_TAGS:
				display.ScreenTags(
					NUM_EQUIP,
					comm_verif,
					/* Tags */ atomic.LoadInt64(&tags),
					/* Atletas */ tagSet.Count(),
				)
			case lcdlogger.SCREEN_ADDR:

				ip := reader.Octets
				leitor := flick.OK

				if !reader.State {

					ip = [4]int{0, 0, 0, 0}
					leitor = flick.DESLIGAD
				}

				display.ScreenAddr(
					NUM_EQUIP,
					reader.Ping,
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
					lcdlogger.ToForthNumber(atomic.LoadInt64(&antennas[0])),
					lcdlogger.ToForthNumber(atomic.LoadInt64(&antennas[1])),
					lcdlogger.ToForthNumber(atomic.LoadInt64(&antennas[2])),
					lcdlogger.ToForthNumber(atomic.LoadInt64(&antennas[3])),
				)
			}

			display.SwitchScreens()

			time.Sleep(100 * time.Millisecond)
		}
	}()

skip_monitoring:
	select {}
}
