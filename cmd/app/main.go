package main

import (
	"fmt"
	"time"

	"workerpool/internal/pool"
)

func main() {
	p := pool.New(8)
	defer p.Shutdown()

	for i := 0; i < 3; i++ {
		p.AddWorker()
	}

	for i := 1; i <= 5; i++ {
		p.Submit(fmt.Sprintf("task #%d", i))
	}
	time.Sleep(time.Second)

	p.AddWorker()
	p.RemoveWorker()

	for i := 6; i <= 10; i++ {
		p.Submit(fmt.Sprintf("task #%d", i))
	}

	time.Sleep(2 * time.Second)
}
