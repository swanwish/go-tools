package main

import (
	"flag"
	"fmt"
	"net/http"
	"strings"

	"./utils"
)

const (
	DefaultPort = 8080
)

var (
	port     int64
	rootPath string
)

func init() {
	flag.Int64Var(&port, "port", DefaultPort, "The port to listen")
	flag.StringVar(&rootPath, "path", ".", "The path of the md files")
	flag.Parse()
}

func main() {
	if port == 0 {
		port = DefaultPort
	}

	localIps, err := utils.GetLocalIPAddrs()
	if err != nil {
		fmt.Println("Failed to get local ip addresses.")
		return
	}

	fmt.Printf("Service listen on port \x1b[31;1m%d\x1b[0m and server ip addresses are \x1b[31;1m%s\x1b[0m\n", port, strings.Join(localIps, ", "))

	httpAddr := fmt.Sprintf(":%d", port)
	//	log.Printf("The root path is %s\n", rootPath)
	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir(rootPath))))
	if err := http.ListenAndServe(httpAddr, nil); err != nil {
		fmt.Printf("http.ListendAndServer() failed with %s\n", err)
	}

	fmt.Printf("Exited\n")
}
