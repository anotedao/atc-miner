package main

import (
	"log"
	"time"
)

type Miner struct {
	height uint64
}

func (m *Miner) start() {
	for {
		nh := getHeight()

		if m.height < nh {
			m.height = getHeight()
			log.Println(m.height)

			if m.height%10 == 0 {
				m.mine()
				log.Println("New block!")
			}
		}

		time.Sleep(time.Second * MinerTick)
	}
}

func (m *Miner) mine() {
	callMine()
}

func initMiner() {
	m := &Miner{}
	m.start()
}
