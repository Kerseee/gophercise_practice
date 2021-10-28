package main

import (
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"os"
	"time"
	
)

// problemFile is a flag for csv file containing problems.
var (
	problemFile = flag.String(
		"csv", 
		"problems.csv", 
		"a csv file in the format of 'question,answer'",
	)
	timeLimit = flag.Int("limit", 30, "the time limit for the quiz in seconds")
	shuffle = flag.Bool("shuffle", false, "set true to shuffle the quiz")
)

// score record the number of correct answer from user.
var score int

type Problem struct{
	id 		int
	ques	string
	ans 	string
}

type Quiz []Problem

// check return true if ans equals to p.ans.
func (p Problem) check(ans string) bool {
	return p.ans == ans
}

// exit print the msg and exit
func exit(msg string){
	fmt.Println(msg)
	os.Exit(1)
}

// problemReader read the csv file and return the quiz.
func problemReader(csvFile io.Reader) (Quiz, error) {
	// Read problems from csv.
	problems, err := csv.NewReader(csvFile).ReadAll()
	if err != nil {
		return Quiz{}, err
	}
	if len(problems) <= 1 {
		return Quiz{}, errors.New("No problem in the file.")
	}
	// Store problems
	problemList := make(Quiz, 0)
	// Drop the header
	for i, q := range problems[1:] {
		problemList = append(problemList, Problem{id: i+1, ques: q[0], ans: q[1]})
	}
	return problemList, nil
}

// Start the quiz
func (q Quiz) Start(done chan bool) {
	// Quiz start
	for _, p := range q {
		fmt.Printf(`Problem #%v: %v = `, p.id, p.ques)
		var ans string
		_, err := fmt.Scanln(&ans)
		if err != nil {
			fmt.Println("\tInput error!")
		}
		if p.check(ans) {
			score ++
		}
	}
	done <- true
}

// Swap swap two problem.
func (q Quiz) Swap(i, j int) {
	q[i], q[j] = q[j], q[i]
}

// Shuffle shuffle the whole quiz.
func (q Quiz) Shuffle(){
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(q), q.Swap)
}

func main() {
	// Parse the flag.
	flag.Parse()

	// Open the problem csv file.
	file, err := os.Open(*problemFile)
	if err != nil {
		exit(fmt.Sprintf(`File %v is not exist.`, *problemFile))
	}
	defer file.Close()

	// Read problems
	quiz, err := problemReader(file)
	if err != nil {
		exit(err.Error())
	}

	if *shuffle {
		quiz.Shuffle()
	}

	// Start quiz
	quizEnd := make(chan bool)
	go quiz.Start(quizEnd)

	// Start timer
	timer := time.NewTimer(time.Second * time.Duration(*timeLimit))
	
	// Move on either when time is out or all quiz are done.
	select {
	case _ = <- timer.C:
		fmt.Println("\nTime Out!")
	case _ = <- quizEnd:
		fmt.Println("All quiz are done in time!")
	}	
	
	// Print score
	fmt.Printf(`You score %v out of %v!`, score, len(quiz))
}