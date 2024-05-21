// Cameron, a web fuzzer
package main

import (
	"flag"
	"fmt"
	"os"
	"sync"
	"time"
)

// Flags
var t = flag.String("t", "localhost", "set target IP/URL")
var l = flag.String("l", "localhost", "input wordlist")
var v = flag.Bool("v", false, "enable verbose output")
var r = flag.Int("r", 5, "set requests per second")

var maxRequests int
var isVerbose bool
var wordlist string

//
func main() {
	target := *t
	var wg sync.WaitGroup
	tokens := make(chan struct{}, *r)
	var startTimer time.Time

	if isVerbose {
		startTimer = time.Now()
	}

	// Check for empty argument list
	if len(os.Args) <= 1 {
		prHeader()
		os.Exit(0)
	}

	// Fuzz scan
	for i := 1; i <= 10; i++ {
		wg.Add(1)
		go fuzz(wordlist, target, &wg, &tokens)
	}
	wg.Wait()

	// Time program execution
	stopTimer := time.Now()
	if isVerbose {
		duration := stopTimer.Sub(startTimer)
		fmt.Println("")
		fmt.Println("Scan duration: ", duration)
	}

}

//
func fuzz(wordlist string, target string, wg *sync.WaitGroup, tokens *chan struct{}) {
	defer wg.Done()
	*tokens <- struct{}{}
	fmt.Println("Scanning")
	time.Sleep(1 * time.Second)
	<-*tokens
}

// Set initial values from flags and other values
func init() {
	flag.Parse()

	if *l != "" {
		wordlist = *l
	}

	if *r > 0 {
		maxRequests = *r
	} else {
		// Default on negative input
		maxRequests = 5
	}

	if *v {
		isVerbose = true
		fmt.Println("Cameron is in a talkative mood right now")
	} else {
		isVerbose = false
	}

	if isVerbose {
		fmt.Println("Requests per second: ", maxRequests)
	}
}

// Print header when no arguments in CLI or on error
func prHeader() {
	fmt.Println("Cameron, a web fuzzer by BenPapple")
	fmt.Println("")
	// ANSI Shadow
	fmt.Println(" ██████╗ █████╗ ███╗   ███╗███████╗██████╗  ██████╗ ███╗   ██╗")
	fmt.Println("██╔════╝██╔══██╗████╗ ████║██╔════╝██╔══██╗██╔═══██╗████╗  ██║")
	fmt.Println("██║     ███████║██╔████╔██║█████╗  ██████╔╝██║   ██║██╔██╗ ██║")
	fmt.Println("██║     ██╔══██║██║╚██╔╝██║██╔══╝  ██╔══██╗██║   ██║██║╚██╗██║")
	fmt.Println("╚██████╗██║  ██║██║ ╚═╝ ██║███████╗██║  ██║╚██████╔╝██║ ╚████║")
	fmt.Println(" ╚═════╝╚═╝  ╚═╝╚═╝     ╚═╝╚══════╝╚═╝  ╚═╝ ╚═════╝ ╚═╝  ╚═══╝")
	fmt.Println("")
	fmt.Println("Use -h for help")
	fmt.Println("Example use case: cameron -t 127.0.0.1")
	fmt.Println("Example use case: cameron -t localhost")
	fmt.Println("Example use case: cameron -t URL")
	fmt.Println("")
}
