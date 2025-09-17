package main

import "log"

var conf *Config

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	conf = initConfig()

	initMiner()
}
