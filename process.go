package main

import (
	"log"
	"sync/atomic"
	"time"
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

		antennas [4]int64
	)

	tagSet := NewIntSet()

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

			//log.Println(tagSet.Count)
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

				antennas[0] = 0
				antennas[1] = 0
				antennas[2] = 0
				antennas[3] = 0
			}

			a.AtualizarAntenas(&antennas)
		}
	}()

	select {}
}
