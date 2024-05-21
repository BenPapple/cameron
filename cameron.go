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
var wordlistFile string

//
func main() {
	target := *t
	var wg sync.WaitGroup
	tokens := make(chan struct{}, *r)
	var startTimer time.Time
	wordlist := []string{
		"/graphql",
		"/v1/graphql",
		"/v2/graphql",
		"/v3/graphql",
		"/graphiql",
		"/v1/graphiql",
		"/v2/graphiql",
		"/v3/graphiql",
		"/playground",
		"/v1/playground",
		"/v2/playground",
		"/v3/playground",
	}

	if isVerbose {
		startTimer = time.Now()
	}

	// Check for empty argument list
	if len(os.Args) <= 1 {
		prHeader()
		os.Exit(0)
	}

	// Fuzz scan
	for _, targetWord := range wordlist {
		wg.Add(1)
		go fuzz(wordlist, target, targetWord, &wg, &tokens)
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
func fuzz(wordlist []string, target string, targetWord string, wg *sync.WaitGroup, tokens *chan struct{}) {
	defer wg.Done()
	*tokens <- struct{}{}
	targetCombined := fmt.Sprintf("%s%s", target, targetWord)
	fmt.Println("Scanning", targetCombined)
	time.Sleep(1 * time.Second)
	<-*tokens
}

// Set initial values from flags and other values
func init() {
	flag.Parse()

	if *l != "" {
		wordlistFile = *l
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
