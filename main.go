package main

import (
	_ "github.com/lib/pq"
	"flag"
	"fmt"
	"math/rand"
	"time"
)


var maxTime int
var shuffle bool
var terminal bool

func init() {
	flag.IntVar(&maxTime, "time", 10,"set duration of quiz")
	flag.BoolVar(&shuffle, "rand", false,"shuffle questions")
	flag.BoolVar(&terminal, "cli", false,"run in terminal")
}

func main() {
	flag.Parse()

	var probs []problem
	var score int

	db := connectToDB()

	if shuffle {
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(probs), func(i, j int) {probs[i], probs[j] = probs[j], probs[i]})
	}

	if terminal {
		fmt.Println("Press enter to start quiz, timelimit is set to ", maxTime)
		fmt.Scanln()

		cn := make(chan bool)
		go startTimer(cn)
		go runQuiz(probs, &score, cn)

		for <-cn{
			fmt.Println("Score: ", score, "/", len(probs))
			return
		}
	} else {
		initiateServer(db)
	}
}