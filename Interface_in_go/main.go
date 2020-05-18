package main

import "fmt"

type Engine interface {
	Start()
	Stop()
}

type CarEngine struct {
}

func (c CarEngine) Start() {
	fmt.Println("My car engine is started")
}

func (c CarEngine) Stop() {
	fmt.Println("My car engine is stoped")
}

type TrainEngine struct {
}

func (t TrainEngine) Start() {
	fmt.Println("My train engine is started")
}

func (t TrainEngine) Stop() {
	fmt.Println("My train engine is stoped")
}

type FanEngine struct {
}

func (t FanEngine) Start() {
	fmt.Println("My fan engine is started")
}

func (t FanEngine) Stop() {
	fmt.Println("My fan engine is stoped")
}

func Starting(e Engine) {
	e.Start()
}

func Stoping(e Engine) {
	e.Stop()
}

func main() {
	carEngine := CarEngine{}
	trainEngine := TrainEngine{}
	FanEngine := FanEngine{}

	engines := []Engine{carEngine, trainEngine, FanEngine}

	for _, engine := range engines {
		Starting(engine)
		Stoping(engine)
	}
}
