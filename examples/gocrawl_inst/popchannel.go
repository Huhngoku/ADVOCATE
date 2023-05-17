package main

import (
	"time"

	"github.com/ErikKassubek/deadlockDetectorGo/src/dedego"
)

type popChannel struct {
	dedego.Chan[[]*URLContext]
}

func newPopChannel() popChannel {

	return popChannel{dedego.NewChan[[]*URLContext](int(1))}
}

func (pc popChannel) stack(cmd ...*URLContext) {
	toStack := cmd
	for {
		{
			dedego.PreSelect(false, pc.GetIdPre(false), pc.GetIdPre(true))
			selectCaseDedego_4 := dedego.BuildMessage(toStack)
			switch dedegoFetchOrder[2] {
			case 0:
				select {

				case pc.GetChan() <- selectCaseDedego_4:
					pc.Post(false, selectCaseDedego_4)
					return
				case <-time.After(2 * time.Second):
					select {
					case pc.GetChan() <- selectCaseDedego_4:
						pc.Post(false, selectCaseDedego_4)
						return
					case selectCaseDedego_5 := <-pc.GetChan():
						pc.Post(true, selectCaseDedego_5)
						old := selectCaseDedego_5.GetInfo()

						toStack = append(old, toStack...)
					}
				}
			case 1:
				select {
				case selectCaseDedego_5 := <-pc.GetChan():
					pc.Post(true, selectCaseDedego_5)
					old := selectCaseDedego_5.GetInfo()

					toStack = append(old, toStack...)
				case <-time.After(2 * time.Second):
					select {
					case pc.GetChan() <- selectCaseDedego_4:
						pc.Post(false, selectCaseDedego_4)
						return
					case selectCaseDedego_5 := <-pc.GetChan():
						pc.Post(true, selectCaseDedego_5)
						old := selectCaseDedego_5.GetInfo()

						toStack = append(old, toStack...)
					}
				}
			default:
				select {
				case pc.GetChan() <- selectCaseDedego_4:
					pc.Post(false, selectCaseDedego_4)
					return
				case selectCaseDedego_5 := <-pc.GetChan():
					pc.Post(true, selectCaseDedego_5)
					old := selectCaseDedego_5.GetInfo()

					toStack = append(old, toStack...)
				}
			}
		}
	}
}
