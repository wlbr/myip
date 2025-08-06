package main

import (
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	//geoip2 "github.com/oschwald/geoip2-golang"

	"github.com/wlbr/myip/gotils"
)

//go:generate templify -o myip.go myip.tpl

// Entferne die Definitionen der Funktionen getIP, rawSubHandler, completeProtocol, completePath, fullTemplateSubHandler, handler (Zeilen 267-383)

func main() {
	mode := "http"
	llevel, _ := gotils.LogLevelString(LogLevel)
	logger = gotils.NewLogger(LogFile, llevel)
	logger.Info("Starting up")

	d := make(chan bool, 10)
	go checkDownload(GeoIpUrl, GEOIPFILENAME, d)

	runtime.GOMAXPROCS(runtime.NumCPU()) // use all CPU cores
	n := runtime.NumGoroutine() + 1      // initial number of Goroutines

	// install signal handler
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM)
	// Spawn request handler

	// go func() {
	// 	err := fcgi.Serve(nil, http.HandlerFunc(handler))
	// 	if err != nil {//
	// 		c <- syscall.SIGTERM
	// 		mode = "http"
	// 		//panic(err)

	// 	} else {
	// 		logger.Info("Spawing fcgi handler.")
	// 		mode = "fcgi"
	// 	}
	// }()

	if mode == "fcgi" {
		// catch signal
		_ = <-d

		// give pending requests in fcgi.Serve() some time to enter the request handler
		time.Sleep(time.Millisecond * 100)

		// wait at most 3 seconds for request handlers to finish
		// incase we are downloading the GeoIPFile that may take a while
		for i := 0; i < 30; i++ {
			if runtime.NumGoroutine() <= n {
				return
			}
			time.Sleep(time.Millisecond * 100)
		}

		// catch finished downloader signal
		_ = <-c
	} else {
		logger.Info("Not in fcgi mode so not spawing handler. Will start http server on port %s.", Port)
		http.HandleFunc("/", handler)
		errsrv := http.ListenAndServe(":"+Port, nil)
		if errsrv != nil {
			logger.Fatal("Error starting http server on port %s: %s", Port, errsrv.Error())
		}
	}
}
