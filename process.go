package main

import (
	"log"
)

func (a *Ay) AtualizarTagsUnicas(tagsUnicas int) {
	_, err := a.db.Exec(`UPDATE equipamento SET tags_unicas = ? WHERE id = 1`, tagsUnicas)

	if err != nil {
		log.Println("(AtualizarTagsUnicas)", err)
	}
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
