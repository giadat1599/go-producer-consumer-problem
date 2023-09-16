package main

import (
	"container/list"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/fatih/color"
)

const MAX_BUFFER_SIZE = 10
const LOOP_TIMES = 100

var (
	buffer   *list.List
	wg       *sync.WaitGroup
	mutex    *sync.Mutex
	cond_var *sync.Cond
)

func Producer() {
	defer wg.Done()

	for i := 0; i < LOOP_TIMES; i++ {
		time.Sleep(500 * time.Millisecond)

		rand_val := rand.Intn(100)

		cond_var.L.Lock()

		for buffer.Len() >= MAX_BUFFER_SIZE {
			color.Green("NO MORE SPACE FOR PRODUCER -> PRODUCER IS GOING TO SLEEP....")
			cond_var.Wait()
		}

		buffer.PushBack(rand_val)
		color.Green("Producer push %d to the buffer", rand_val)

		cond_var.Signal()
		cond_var.L.Unlock()
	}
}

func Consumer() {
	defer wg.Done()

	for i := 0; i < LOOP_TIMES; i++ {
		time.Sleep(1 * time.Second)

		cond_var.L.Lock()

		for buffer.Len() == 0 {
			color.Red("THERE ARE NO ITEMS FOR CONSUMER -> CONSUMER IS GOING TO SLEEP....")
			cond_var.Wait()
		}

		n := buffer.Front()
		buffer.Remove(n)
		color.Red("Consumer get %d from the buffer", n.Value)

		cond_var.Signal()
		cond_var.L.Unlock()

	}
}

func main() {
	buffer = list.New()
	mutex = new(sync.Mutex)
	wg = &sync.WaitGroup{}
	cond_var = sync.NewCond(mutex)

	wg.Add(2)

	go Producer()
	go Consumer()

	wg.Wait()
	fmt.Println("Done.......")
}
