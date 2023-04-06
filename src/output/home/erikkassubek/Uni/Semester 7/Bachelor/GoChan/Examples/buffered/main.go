package main

import (
	"time"
	"github.com/ErikKassubek/GoChan/goChan"
)

func main() {
	goChan.Init(20)
	defer goChan.RunAnalyzer()
	defer time.Sleep(time.Millisecond)
	c := goChan.NewChan[int](int(1))
	d := goChan.NewChan[int](int(0))
	func() {
		GoChanRoutineIndex := goChan.SpawnPre()
		go func() {
			goChan.SpawnPost(GoChanRoutineIndex)
			{
				c.Send(1)
				c.Send(1)
				d.Send(1)
			}
		}()
	}()
	d.Receive()
	c.Receive()
	c.Receive()

	time.Sleep(3 * time.Second)
}

var goChanFetchOrder = make(map[int]int)
