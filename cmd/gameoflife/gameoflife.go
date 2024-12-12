package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/jhmilan/game-of-life/internal/app/gameoflife/entities"
)

func main() {
	world := entities.NewWorld(
		10,
		10,
		[][]int{
			{0, 0},
			{6, 5},
			{7, 4},
			{7, 5},
			{8, 5},
			{8, 6},
		})
	world.Print()
	ticker := time.NewTicker(5 * time.Millisecond)
	done := make(chan bool)
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		for {
			select {
			case <-done:
				fmt.Println("time up")
				return
			case <-ticker.C:
				numAlive := world.CountAlive()
				generationStep := world.GetGenerationStep()
				fmt.Printf("Generation: %d - Num cells alive: %d\n", generationStep, numAlive)
				world.Step()
				world.Print()

				status := world.GetStatus()
				switch status {
				case entities.Extinction:
					fallthrough
				case entities.Static:
					fmt.Printf("Stopping: %s\n", status)
					return
				}
			}
		}
	}()

	go func() {
		// stop after a maximum of 30 seconds
		time.Sleep(30 * time.Second)
		done <- true
	}()

	wg.Wait()
	ticker.Stop()
	fmt.Println("Ticker stopped")
}
