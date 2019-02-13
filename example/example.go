package main

import (
	"fmt"
	"github.com/hlfstr/logit"
	"os"
	"time"
)

//Globally Declare your Logs
var ThisLog *logit.Logger
var ThatLog *logit.Logger

func main() {
	//Create an "Error" variable
	var err error

	//Create a logfile named "this" in the relative directory
	thisFile, err := logit.OpenFile("thisfile.log")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	ThisLog, err = logit.StartLogger(thisFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	//Create a terminal Logger
	ThatLog, err = logit.StartLogger(logit.TermLog())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	//Create 10 Goroutines
	for i := 0; i < 10; i++ {
		go test(i)
	}

	//Give 15 seconds for them to finish
	time.Sleep(time.Second * 15)

	//Close the logs
	ThisLog.Quit()
	ThatLog.Quit()
}

func test(i int) {
	//Temp string
	var s string
	for j := 0; j < 100; j++ {
		//Append string to end of "this" log
		s = fmt.Sprintf("This: I am number %d, look at me!", i)
		ThisLog.Log(s)

		//Append string to end of "this" log
		s = fmt.Sprintf("That: I am number %d, look at me!", i)
		ThatLog.Log(s)

		//Make the goroutine wait .25 of a second
		time.Sleep(time.Millisecond * 250)
	}
}
