// A web fuzzer
package main

import "fmt"

//
func main() {
	prHeader()
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
	//fmt.Println("Use -h for help")
	//fmt.Println("Example use case: cameron -t 127.0.0.1")
	//fmt.Println("Example use case: cameron -t localhost")
	//fmt.Println("Example use case: cameron -t URL")
	fmt.Println("")
}
