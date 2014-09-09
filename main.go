package main

import (
	"flag"
	"fmt"
	"log"
	"os/exec"
	"runtime"
	"strings"
	"sync"
)

var procs = flag.Int("procs", 1, "number of processors to use")
var jobs = flag.Int("jobs", 4, "number of jobs to start concurrently")

func worker(cmdChan chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	for cmdString := range cmdChan {
		parsed := strings.SplitAfterN(cmdString, " ", 2)
		cmdName := parsed[0]
		cmdArgs := ""
		if len(parsed) == 2 {
			cmdArgs = parsed[1]
		}
		cmd := exec.Command(cmdName, cmdArgs)
		/*
			if err := cmd.Start(); err != nil {
				log.Fatalf("cmd.Start " + cmdString + " %v")
			}
			//... Wait()
			//... Output()
		*/
		if err := cmd.Run(); err != nil {
			log.Println(1, cmdString, err)
		} else {
			log.Println(0, cmdString, err)
		}
	}
}

func main() {
	flag.Parse()
	fmt.Println("procs", *procs)
	fmt.Println("jobs", *jobs)

	runtime.GOMAXPROCS(*procs)

	cmdChan := make(chan string)
	wg := new(sync.WaitGroup)

	for i := 0; i < *jobs; i++ {
		wg.Add(1)
		go worker(cmdChan, wg)
	}

	for _, cmd := range flag.Args() {
		cmdChan <- cmd
	}
	close(cmdChan)
	wg.Wait()
}
