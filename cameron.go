// Cameron, a web fuzzer
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
	"unicode"
)

// Flags
var t = flag.String("t", "localhost", "set target IP/URL")
var l = flag.String("l", "localhost", "input wordlist")
var v = flag.Bool("v", false, "enable verbose output")
var r = flag.Int("r", 5, "set requests per second")
var fc = flag.String("fc", "", "filter status code")
var mc = flag.String("mc", "", "match status code")

var maxRequests int
var isVerbose bool

var wordlistFile string

// Start fuzzing
func main() {
	var host string = ""
	var wg sync.WaitGroup
	tokens := make(chan struct{}, *r)
	var startTimer time.Time
	var scanResults sync.Map
	wordlistFile = *l
	filterCode := *fc
	matchCode := *mc
	var progressCounter uint64

	if isVerbose {
		startTimer = time.Now()
	}

	// Check for empty argument list and then validate target input
	if len(os.Args) <= 1 {
		prHeader()
		os.Exit(0)
	}
	targetCheck(&host)

	// Read file
	wlFile := getFile(wordlistFile)

	// Client
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Progress tests
	if !isVerbose {
		go progressBar(wlFile, &progressCounter)
	}

	// Fuzz scan
	for _, targetWord := range wlFile {
		wg.Add(1)
		go fuzz(wlFile, host, targetWord, &wg, &tokens, client, &scanResults, &progressCounter)
	}

	wg.Wait()

	printResults(scanResults, host, filterCode, matchCode)

	// Time program execution
	stopTimer := time.Now()
	if isVerbose {
		duration := stopTimer.Sub(startTimer)
		fmt.Println("")
		fmt.Println("Scan duration: ", duration)
	}

}

// Fuzz a URL with words from wordlist
func fuzz(wordlist []string, target string, targetWord string, wg *sync.WaitGroup, tokens *chan struct{}, client *http.Client, scanResults *sync.Map, progressCounter *uint64) {
	defer wg.Done()
	*tokens <- struct{}{}

	// Fuzzing
	targetCombined := replaceFUZZ(target, targetWord)
	resp, err := client.Get(targetCombined)
	if err != nil {
		fmt.Printf("\nError fetching URL %s: %s", targetCombined, err)
		atomic.AddUint64(progressCounter, 1)
		<-*tokens
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("\nError reading response body:", err)
		atomic.AddUint64(progressCounter, 1)
		<-*tokens
		return
	}

	// Get the size of the response body in bytes
	responseSize := len(body)

	// Count the number of lines in the response body
	lineCount := countLines(string(body))

	// Count the number of words in the response body
	wordCount := countWords(string(body))

	// Store results from response body
	respData := [4]int{resp.StatusCode, responseSize, wordCount, lineCount}
	scanResults.Store(targetCombined, respData)
	atomic.AddUint64(progressCounter, 1)
	time.Sleep(1 * time.Second)
	<-*tokens
}

// Progress bar
func progressBar(wlFile []string, progressCounter *uint64) {
	var count uint64
	count = 1

	for int(count) <= len(wlFile) {
		count = atomic.LoadUint64(progressCounter)
		fmt.Printf("\033[2J\033[0;0HProgress: %d %s %d %s", count, "of", len(wlFile), "targets done.")
		if int(count) == len(wlFile) {
			fmt.Print("\n\n")
			time.Sleep(1000 * time.Millisecond)
		} else {
			time.Sleep(40 * time.Millisecond)
		}
	}

}

// Read wordlist from file
func getFile(wordlistFile string) []string {
	list := []string{}

	file, err := os.Open(wordlistFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		list = append(list, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return list
}

// Count the lines in a string
func countLines(s string) int {
	n := strings.Count(s, "\n")
	if !strings.HasSuffix(s, "\n") {
		n++
	}
	return n
}

// Count the words in a string
func countWords(text string) int {
	words := strings.FieldsFunc(text, func(c rune) bool {
		return !unicode.IsLetter(c) && !unicode.IsNumber(c)
	})
	return len(words)
}

// Replace FUZZ with word from wordlist
func replaceFUZZ(host string, fuzz string) string {
	out := strings.Replace(host, "FUZZ", fuzz, 1)
	return out
}

// Pretty print results in table
func printResults(scanResults sync.Map, host string, filterCode string, matchCode string) {

	tempMap := map[string][4]int{}
	scanResults.Range(func(key, value interface{}) bool {
		tempMap[fmt.Sprint(key)] = value.([4]int)
		return true
	})

	// Sort map by string and print
	fmt.Printf("%-20v %-6v  %-4v  %-5v  %-5v\n", "HOST", "Status", "Size", "Words", "Lines")
	keys := make([]string, 0)
	for k := range tempMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {

		trimmedHost := strings.TrimSuffix(host, "/FUZZ")
		trimmedHost = strings.TrimPrefix(k, trimmedHost)

		// not in filter and in match
		if !strings.Contains(filterCode, strconv.Itoa(tempMap[k][0])) && strings.Contains(matchCode, strconv.Itoa(tempMap[k][0])) {
			fmt.Printf("%-20s %-6d  %-4d  %-5d  %-5d \n", trimmedHost, tempMap[k][0], tempMap[k][1], tempMap[k][2], tempMap[k][3])
		}
		// empty and in match
		if filterCode == "" && strings.Contains(matchCode, strconv.Itoa(tempMap[k][0])) {
			fmt.Printf("%-20s %-6d  %-4d  %-5d  %-5d \n", trimmedHost, tempMap[k][0], tempMap[k][1], tempMap[k][2], tempMap[k][3])
		}
		// in filter and empty
		if strings.Contains(filterCode, strconv.Itoa(tempMap[k][0])) && matchCode == "" {
			fmt.Printf("%-20s %-6d  %-4d  %-5d  %-5d \n", trimmedHost, tempMap[k][0], tempMap[k][1], tempMap[k][2], tempMap[k][3])
		}
		// both empty
		if filterCode == "" && matchCode == "" {
			fmt.Printf("%-20s %-6d  %-4d  %-5d  %-5d \n", trimmedHost, tempMap[k][0], tempMap[k][1], tempMap[k][2], tempMap[k][3])
		}
		// both the same
		if strings.Contains(filterCode, strconv.Itoa(tempMap[k][0])) && strings.Contains(matchCode, strconv.Itoa(tempMap[k][0])) && matchCode != "" {
			fmt.Println("Error: Can't match and filter the same code.")
			os.Exit(0)
		}

	}
}

// Check if user input for target is valid IP or URI
func targetCheck(host *string) {

	// Check for FUZZ keyword
	if !strings.Contains(*t, "FUZZ") {
		fmt.Println("Error: Input is missing FUZZ keyword.")
		os.Exit(0)
	}

	// Check for valid IP in input
	checkIP := net.ParseIP(*t)
	if checkIP != nil {
		*host = *t
		fmt.Println("db: Check IP")
		return
	}

	// Check for valid URI in input
	_, err := url.ParseRequestURI(*t)
	if err == nil {
		*host = *t
		fmt.Println("db: URI")
		return
	}

	// Check for if input is string localhost
	if *t == "localhost" {
		tempHost := fmt.Sprintf("%s%s", "http://", *t)
		*host = tempHost
		fmt.Println("db: localhost")
		return
	}

	// Add http prefix to check isURI again
	tempHost := fmt.Sprintf("%s%s", "http://", *t)
	_, err2 := url.ParseRequestURI(tempHost)
	if err2 == nil {
		*host = tempHost
		fmt.Println("db: add http then check URI")
		return
	}

	// Exit program since no valid input
	prHeader()
	fmt.Println("Error: No valid IP or URI given")
	fmt.Println("Error on input target candidate: ", *t)
	os.Exit(0)

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
