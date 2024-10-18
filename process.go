package main

import (
	"log"
)

func (a *Ay) Process() {

	tagSet := NewIntSet()

	for t := range a.Tags {

		if t.Antena == 0 {
			/*
				Antena 0 não existe
			*/

			continue
		}

		tagSet.Insert(t.Epc)
		log.Println(tagSet.Count)
	}
}
