package main

import "sync"
import "time"

var a string
var b bool
var once sync.Once

func setup() {
	a = "hello, world"
	b = true
}

func doprint() {
	once.Do(setup)
	if b {
		print(a)
	}


}

func twoprint() {
	go doprint()
	go doprint()
}

func main() {
	twoprint()
	time.Sleep(1 * time.Second)
}
