package main

import (
	"analyzer/reader"
)

func main() {
	file_path := "./dedego.log"
	reader.Read(file_path)
}
