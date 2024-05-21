// Cameron, a web fuzzer
package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"sort"
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
	var host string = ""
	var targetIP string = ""
	var wg sync.WaitGroup
	tokens := make(chan struct{}, *r)
	var startTimer time.Time
	var scanResults = map[string]int{}
	wordlistFile = *l
	wordlist := []string{
		"/graphql",
		// "/v1/graphql",
		// "/v2/graphql",
		// "/v3/graphql",
		// "/graphiql",
		// "/v1/graphiql",
		// "/v2/graphiql",
		// "/v3/graphiql",
		// "/playground",
		// "/v1/playground",
		// "/v2/playground",
		// "/v3/playground",
	}

	if isVerbose {
		startTimer = time.Now()
	}

	// Check for empty argument list and then validate target input
	if len(os.Args) <= 1 {
		prHeader()
		os.Exit(0)
	}
	targetCheck(&host, &targetIP)

	// Client
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Fuzz scan
	for _, targetWord := range wordlist {
		wg.Add(1)
		go fuzz(wordlist, host, targetWord, &wg, &tokens, client, scanResults)
	}
	wg.Wait()

	printResults(scanResults)

	// Time program execution
	stopTimer := time.Now()
	if isVerbose {
		duration := stopTimer.Sub(startTimer)
		fmt.Println("")
		fmt.Println("Scan duration: ", duration)
	}

}

//
func fuzz(wordlist []string, target string, targetWord string, wg *sync.WaitGroup, tokens *chan struct{}, client *http.Client, scanResults map[string]int) {
	defer wg.Done()
	*tokens <- struct{}{}

	// Fuzzing
	targetCombined := fmt.Sprintf("%s%s", target, targetWord)
	fmt.Println("Scanning", targetCombined)
	resp, err := client.Get(targetCombined)
	if err != nil {
		fmt.Printf("Error fetching URL %s: %s", targetCombined, err)
		return
	}
	defer resp.Body.Close()

	scanResults[targetCombined] = resp.StatusCode
	time.Sleep(1 * time.Second)
	<-*tokens
}

// Pretty print results in table
func printResults(scanResults map[string]int) {

	// Sort map by string and print
	fmt.Printf("%-30v %v\n", "HOST", "HTTP Status")
	keys := make([]string, 0)
	for k := range scanResults {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		if scanResults[k] > 0 {
			fmt.Printf("%-30s %-4d \n", k, scanResults[k])
		}
	}
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

// Check if user input for target is valid IP or URI
func targetCheck(host *string, targetIP *string) {

	// Check for valid IP in input
	checkIP := net.ParseIP(*t)
	if checkIP != nil {
		*host = *t
		*targetIP = *t
		return
	}

	// Check for valid URI in input
	_, err := url.ParseRequestURI(*t)
	if err == nil {
		tempHost := fmt.Sprintf("%s%s", "http://", *t)
		*host = tempHost
		*targetIP = getIP(*host)
		return
	}

	// Check for if input is string localhost
	if *t == "localhost" {
		tempHost := fmt.Sprintf("%s%s", "http://", *t)
		*host = tempHost
		*targetIP = getIP(*host)
		return
	}

	// Add http prefix to check isURI again
	tempHost := fmt.Sprintf("%s%s", "http://", *t)
	_, err2 := url.ParseRequestURI(tempHost)
	if err2 == nil {
		*host = tempHost
		*targetIP = getIP(*host)
		return
	}

	// Exit program since no valid input
	prHeader()
	fmt.Println("Error: No valid IP or URI given")
	fmt.Println("Error on input target candidate: ", *t)
	os.Exit(0)

}

// Return IPv4 from URL
func getIP(host string) string {
	ips, _ := net.LookupIP(host)
	var tempIP string
	for _, ip := range ips {
		if ipv4 := ip.To4(); ipv4 != nil {
			tempIP = fmt.Sprintf("%v", ipv4)

		}
	}
	return tempIP

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
