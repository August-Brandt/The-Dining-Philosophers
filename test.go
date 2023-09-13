package main

import (
	"fmt"
	"time"
)

func philosopher(id int, ate *[]int, sending1, sending2 chan bool, receiving1, receiving2 chan string) {
	eating := false
	for {
		// We can always wait for fork1 because it will always be put back if currently looked at
		fork1 := <-sending1
		// We use a select statement to create a timeout case that ensures we dont get stuck waiting for fork2 while locking fork1
		select {
		case fork2 := <-sending2: // In case we can look at fork2 we can then proceed to look at the state of the forks
			if eating {
				eating = false
				fmt.Println("Philo", id, ": Now thinking")
				(*ate)[id]++
				receiving1 <- "put back"
				receiving2 <- "put back"
			} else {
				if fork1 && fork2 { // We only takes forks if both are available
					receiving1 <- "take"
					receiving2 <- "take"
					eating = true
					fmt.Println("Philo", id, ": Now eating")
				} else { // If we can't take both forks then we tell the forks to repost their value on the sending channel
					receiving1 <- "not take"
					receiving2 <- "not take"
				}
			}

		case <-time.After(100 * time.Millisecond):
			receiving1 <- "not take" // If we can't access fork2 then don't take fork1
		}

		//time.Sleep(time.Millisecond * 100)
	}
}

func fork(id string, sending chan bool, receiving chan string) {
	available := true
	sending <- available
	for {
		receive := <-receiving // take the reponse from the philosopher about what action on the fork
		// Update availability based on reponse from philosopher and
		// send out availability so that other pilosophers can look at the fork
		if receive == "take" {
			available = false
			sending <- available
		} else if receive == "put back" {
			available = true
			sending <- available
		} else if receive == "not take" {
			sending <- available
		}
		// time.Sleep(time.Millisecond * 100)
	}
}

func main() {
	fmt.Println("Starting")

	// Create sending and receiving channels for all forks
	var ch0S = make(chan bool)
	var ch0R = make(chan string)
	var ch1S = make(chan bool)
	var ch1R = make(chan string)
	var ch2S = make(chan bool)
	var ch2R = make(chan string)
	var ch3S = make(chan bool)
	var ch3R = make(chan string)
	var ch4S = make(chan bool)
	var ch4R = make(chan string)

	// Create forks
	go fork("0", ch0S, ch0R)
	go fork("1", ch1S, ch1R)
	go fork("2", ch2S, ch2R)
	go fork("3", ch3S, ch3R)
	go fork("4", ch4S, ch4R)

	// Create philosophers
	var ate = make([]int, 5)
	go philosopher(0, &ate, ch0S, ch4S, ch0R, ch4R)
	go philosopher(1, &ate, ch1S, ch0S, ch1R, ch0R)
	go philosopher(2, &ate, ch2S, ch1S, ch2R, ch1R)
	go philosopher(3, &ate, ch3S, ch2S, ch3R, ch2R)
	go philosopher(4, &ate, ch4S, ch3S, ch4R, ch3R)

	// Check if all philosophers have eaten at least 3 times
	run := true
	for run {
		run = false
		for _, timeAte := range ate {
			if timeAte < 3 {
				run = true
			}
		}
	}
	fmt.Println(ate)
	fmt.Println("Done")
}
