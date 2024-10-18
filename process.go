package main

import (
	"log"
)

func (a *Ay) AtualizarTagsUnicas(tagsUnicas int) {
	a.db.Exec(`UPDATE equipamento SET tags_unicas = ? WHERE id = 1`, tagsUnicas)
}

func (a *Ay) Process() {

	tagSet := NewIntSet()

	for t := range a.Tags {

		if t.Antena == 0 {
			/*
				Antena 0 não exist
			*/

			continue
		}

		if tagSet.Insert(t.Epc) {
			a.AtualizarTagsUnicas(tagSet.Count)
		}

		//log.Println(tagSet.Count)
	}
}
