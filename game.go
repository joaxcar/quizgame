package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)


func runQuiz(probs []problem, score *int, cn chan bool) {
	scanner := bufio.NewScanner(os.Stdin)
	for _,p := range probs {
		fmt.Println("Fråga: ", p.Question)
		scanner.Scan()
		ans := strings.TrimSpace(scanner.Text())

		if ans == p.Answer {
			*score++
		}
	}
	cn <- true
}

func startTimer(cn chan bool) {
	duration := time.Second * time.Duration(maxTime)
	time.Sleep(duration);
	cn <- true
}

