package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	_ "embed"
)

//go:embed hosts.json
var hosts []byte

var hostmap = loadHosts()

func loadHosts() map[string]string {
	var hostmap map[string]string
	err := json.Unmarshal(hosts, &hostmap)
	if err != nil {
		log.Fatal(err)
	}

	return hostmap
}

// determines the domain name requested and returns the www redirect
func wwwRedirect(w http.ResponseWriter, r *http.Request) {
	// get the host
	host := r.Host

	log.Println(host)

	// check if the host is www
	rd := hostmap[host]
	if rd != "" && !strings.HasPrefix(host, "www") {
		// redirect to www
		log.Println("redirecting:", rd)
		http.Redirect(w, r, rd, http.StatusPermanentRedirect)
	}

	http.NotFound(w, r)
}

func main() {
	// create the server
	server := http.Server{
		Addr: "0.0.0.0:80",
	}

	// handle the www redirect
	http.HandleFunc("/", wwwRedirect)

	// start the server
	server.ListenAndServe()
}
