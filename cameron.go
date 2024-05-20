// A web fuzzer
package main

import (
	"flag"
	"fmt"
	"sync"
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
	prHeader()
	target := *t
	var wg sync.WaitGroup

	// Fuzz scan
	for i := 1; i <= 100; i++ {
		wg.Add(1)
		go fuzz(wordlist, target, &wg)
	}
	wg.Wait()

}

//
func fuzz(wordlist string, target string, wg *sync.WaitGroup) {
	defer wg.Done()
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
		fmt.Println("Gordo is in a talkative mood right now")
	} else {
		isVerbose = false
	}

	if isVerbose {
		fmt.Println("Requests per second: ", maxRequests)
	}
}

// Print header when no arguments in CLI or on error
func prHeader() {
	var Reset = "\033[0m"
	var White = "\033[97m"
	fmt.Println("Cameron, a web fuzzer by BenPapple")
	fmt.Println("")
	// ANSI Shadow
	fmt.Println(White + " ██████╗ █████╗ ███╗   ███╗███████╗██████╗  ██████╗ ███╗   ██╗")
	fmt.Println("██╔════╝██╔══██╗████╗ ████║██╔════╝██╔══██╗██╔═══██╗████╗  ██║")
	fmt.Println("██║     ███████║██╔████╔██║█████╗  ██████╔╝██║   ██║██╔██╗ ██║")
	fmt.Println("██║     ██╔══██║██║╚██╔╝██║██╔══╝  ██╔══██╗██║   ██║██║╚██╗██║")
	fmt.Println("╚██████╗██║  ██║██║ ╚═╝ ██║███████╗██║  ██║╚██████╔╝██║ ╚████║")
	fmt.Println(" ╚═════╝╚═╝  ╚═╝╚═╝     ╚═╝╚══════╝╚═╝  ╚═╝ ╚═════╝ ╚═╝  ╚═══╝" + Reset)
	fmt.Println("")
	fmt.Println("Use -h for help")
	fmt.Println("Example use case: cameron -t 127.0.0.1")
	fmt.Println("Example use case: cameron -t localhost")
	fmt.Println("Example use case: cameron -t URL")
	fmt.Println("")
}
