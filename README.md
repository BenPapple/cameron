# cameron
Cameron, a web fuzzer in Go.

Example scan on a web server:

![Scan with gordo on metasploitable2](https://github.com/BenPapple/cameron/blob/main/pics/cameron1.png?raw=true)

# Flags
-fc "filter status code"

-l "input wordlist"

-mc "match status code"

-r "set requests per second"

-t "set target IP/URL"

-v "enable verbose output"

# Use case
Input at least a URL with FUZZ keyword and a wordlist.

go run cameron.go -l ~/yourwordlists.txt -t URL/FUZZ

go run cameron.go -l ~/yourwordlists.txt -t 127.0.0.1/FUZZ


