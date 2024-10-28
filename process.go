package main

import (
	"log"
	"sync/atomic"
	"time"
)

func (a *Ay) AtualizarTagsUnicas(tagsUnicas int64) {
	_, err := a.db.Exec(`UPDATE equipamento SET tags_unicas = ? WHERE id = 1`, tagsUnicas)

	if err != nil {
		log.Println("(AtualizarTagsUnicas)", err)
	}
}

func (a *Ay) AtualizarTags(tags int64) (ok bool) {

	var totalAnterior int64 = 0

	ok = true

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

		_, err = a.db.Exec(`REPLACE INTO stats (tags_total) VALUES (1) WHERE id = 1`)

		ok = false

		return
	}

	_, err = a.db.Exec(`REPLACE INTO stats (tags_total) VALUES (?) WHERE id = 1`, tags)

	if err != nil {
		log.Println("(AtualizarTags)", err)
	}

	return
}

func (a *Ay) Process() {

	var tags int64 /* shared */

	tagSet := NewIntSet()

	go func() {
		for t := range a.Tags {

			if t.Antena == 0 {
				/*
					Antena 0 não exist
				*/

				continue
			}

			atomic.AddInt64(&tags, 1)

			if tagSet.Insert(t.Epc) {

				a.AtualizarTagsUnicas(tagSet.Count())
			}

			//log.Println(tagSet.Count)
		}
	}()

	go func() {

		atualizaContagem := time.NewTicker(30 * time.Second)

		for {
			<-atualizaContagem.C

			/*
				true  => prosseguir
				false => resetar
			*/
			if !a.AtualizarTags(atomic.LoadInt64(&tags)) {
				tagSet.Clear()
				atomic.StoreInt64(&tags, 0)
			}
		}
	}()

	select {}
}
