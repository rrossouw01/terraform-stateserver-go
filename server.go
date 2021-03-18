package main

import (
	"flag"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

var (
	flagDataPath string
	flagCertFile string
	flagKeyFile  string

	flagListenAddress string
)

func init() {
	flag.StringVar(&flagCertFile, "certfile", "", "path to the certs file for https")
	flag.StringVar(&flagKeyFile, "keyfile", "", "path to the keyfile for https")
	flag.StringVar(&flagDataPath, "data_path", os.TempDir(), "path to the data storage directory")

	flag.StringVar(&flagListenAddress, "listen_address", "0.0.0.0:8080", "address:port to bind listener on")
}

func main() {
	if !flag.Parsed() {
		flag.Parse()
	}

	router := http.NewServeMux()
	router.HandleFunc("/", requestHandler)

	if flagCertFile != "" && flagKeyFile != "" {
          log.Fatal(http.ListenAndServeTLS(flagListenAddress, flagCertFile, flagKeyFile, router))
	} else {
          log.Fatal(http.ListenAndServe(flagListenAddress, router))
	}
}

func requestHandler(res http.ResponseWriter, req *http.Request) {
	if req.URL.Path == "/" {
		http.NotFound(res, req)
		return
	}

	stateStorageFile := filepath.Join(flagDataPath, req.URL.Path)
	stateStorageDir := filepath.Dir(stateStorageFile)

	switch req.Method {
	case "GET":
		log.Printf("received GET for state file: %s\n", stateStorageFile)
		fh, err := os.Open(stateStorageFile)
		if err != nil {
			log.Printf("during GET cannot open file: %s\n", err)
			goto not_found
		}
		defer fh.Close()

		res.WriteHeader(200)

		io.Copy(res, fh)

		return
	case "POST":
		log.Printf("received POST for state file: %s\n", stateStorageFile)
		if err := os.MkdirAll(stateStorageDir, 0750); err != nil && !os.IsExist(err) {
			log.Printf("during POST cannot create parent directories: %s\n", err)
			goto not_found
		}

		fh, err := os.OpenFile(stateStorageFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
		if err != nil {
			log.Printf("cannot open file: %s\n", err)
			goto not_found
		}
		defer fh.Close()

		if _, err := io.Copy(fh, req.Body); err != nil {
			log.Printf("failed to stream data into statefile: %s\n", err)
			goto not_found
		}

		res.WriteHeader(200)

		return
	case "DELETE":
		log.Printf("received DELETE for state file: %s\n", stateStorageFile)
		if err := os.RemoveAll(stateStorageFile); err != nil {
			log.Printf("cannot remove file: %s\n", err)
			goto not_found
		}

		res.WriteHeader(200)

		return
	}

not_found:
	http.NotFound(res, req)
}
