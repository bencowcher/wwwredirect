package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"golang.org/x/exp/slices"
)

func isWhiteListed(host string) bool {
	env := os.Getenv("WHITELISTED_DOMAINS")
	domains := strings.Split(env, ",")

	return slices.Contains(domains, host)
}

// determines the domain name requested and returns the www redirect
func wwwRedirect(w http.ResponseWriter, r *http.Request) {
	// get the host
	host := r.Host

	log.Println(host)
	// check if the host is www
	if isWhiteListed(host) && !strings.HasPrefix(host, "www") {
		// redirect to www
		log.Println("redirecting:", "https://www."+host)
		http.Redirect(w, r, "https://www."+host, http.StatusPermanentRedirect)
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
