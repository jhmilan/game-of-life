package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/jhmilan/game-of-life/internal/app/gameoflife/entities"
)

func main() {
	generations := make(map[string]int)

	world := entities.NewWorld(10, 10)
	world.SetCellAlive(0, 0)
	world.SetCellAlive(7, 5)
	world.SetCellAlive(7, 6)
	world.SetCellAlive(8, 7)
	world.SetCellAlive(8, 7)
	world.SetCellAlive(9, 8)
	world.SetCellAlive(9, 8)
	world.SetCellAlive(9, 7)
	world.Print()
	ticker := time.NewTicker(100 * time.Millisecond)
	done := make(chan bool)
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		for i := 0; ; i++ {
			select {
			case <-done:
				fmt.Println("time up")
				return
			case <-ticker.C:
				numAlive := world.CountAlive()
				hash := world.StatusHash()
				status := world.GetStatus(generations)
				fmt.Printf("Iteration: %d - Status: %s - Num cells alive: %d - Hash: %s", i, status, numAlive, hash)
				world = world.Evolve()
				world.Print()
				switch status {
				case entities.Extinction:
					fallthrough
				case entities.Static:
					fmt.Println("Stopping")
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
