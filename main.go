//A http reverse proxy for developer to debug.

package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
)

//logpath is the file path for log.
var logpath = flag.String("l", "/dev/stdout", "path of log file")

//port is the port which server listened.
var port = flag.Int("p", 80, "port which server listened")

//towards is the url where server proxy to.
var towards = flag.String("t", "http://127.0.0.1:8080", "url which proxy redirect to")
var towardsUrl *url.URL

//logMessageSwitch is switch if server print full http request or response.
//0 presents neither  request nor response.
//1 presents only request.
//2 presents only response.
//3 presents both request and reponse.
var logMessageSwitch = flag.Int("m", 1, "switch if print full http message")

//logger is std log.
var logger *log.Logger

func main() {
	flag.Parse()
	//if towards is not validable, stop
	var err error
	towardsUrl, err = url.Parse(*towards)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s is illegal url:%v\n", *towards, err)
		return
	}
	//if logpath is not available, use stdout instead.
	file, err := os.OpenFile(*logpath, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0664)
	if err == nil {
		defer file.Close()
		logger = log.New(file, "", log.Lshortfile)
	} else {
		fmt.Fprintf(os.Stderr, "open %s failed:%v\n", *logpath, err)
		logger = log.New(os.Stdout, "", log.Lshortfile)
	}
	srv := new(Proxy)
	logger.Fatal(srv.ListenAndServe())
}
