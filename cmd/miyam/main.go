package main

import (
	"flag"
	"fmt"
	"net/url"
	"os"
	"time"

	"github.com/lamg/miyam"
)

func main() {
	var info bool
	var itag string
	var timeout time.Duration
	var attempts uint
	flag.BoolVar(&info, "i", false,
		"Get information about available formats")
	flag.DurationVar(&timeout, "t", 5*time.Second,
		"Timeout for client")
	flag.StringVar(&itag, "f", "", "Selected video format tag")
	flag.UintVar(&attempts, "a", 1,
		"Number of attempts for getting the whole file")
	flag.Parse()
	proxy := os.Getenv("http_proxy")
	if proxy == "" {
		proxy = os.Getenv("https_proxy")
	}
	ur := flag.Arg(0)
	my := miyam.NewMiyam(proxy, timeout)
	var e error
	for attempts != 0 {
		if info {
			e = my.InfoTerm(ur)
		} else {
			e = my.DownloadTerm(ur, itag)
		}
		if e == nil {
			attempts = 0
		} else {
			attempts = attempts - 1
			report(e)
		}
	}
}

func report(e error) {
	if e != nil {
		_, ok := e.(*url.Error)
		if ok {
			fmt.Fprintln(os.Stderr, "No network connection")
		} else {
			fmt.Fprintln(os.Stderr, e.Error())
		}
	}
}
