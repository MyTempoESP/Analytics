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

func (a *Ay) AtualizarTags(tags int64) {
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

			a.AtualizarTags(atomic.LoadInt64(&tags))
		}
	}()

	select {}
}
