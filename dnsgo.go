package main

import (
	"github.com/grt1st/dnsgo/server"
	"time"
	"flag"
	"fmt"
	"os"
)

const versionNumber = "1.0.0#20180623"

func main() {
	// lookup
	// localhost
	// log
	version := flag.Bool("version", false, "Show program's version number and exit")
	host := flag.String("host", "localhost", "Address to bind")
	query := flag.Bool("query", true, "Whether to send dns request")
	logFilename := flag.String("log", "", "Filename of log file")
	help := flag.Bool("h", false, "Show usage")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage:\n  %s [Options]\n\nOptions\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

	if *help {
		flag.Usage()
		return
	}
	if *version {
		fmt.Println(versionNumber)
		return
	}

	server := &server.Server{
		Host:     *host,
		Port:     53,
		RTimeout: 5 * time.Second,
		WTimeout: 5 * time.Second,
	}

	server.Run(*query, *logFilename)
}
