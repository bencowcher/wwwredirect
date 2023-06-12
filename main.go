package main

import (
	"log"
	"net/http"
)

// determines the domain name requested and returns the www redirect
func wwwRedirect(w http.ResponseWriter, r *http.Request) {
	// get the host
	host := r.Host
	log.Println(host)
	// check if the host is www
	if host[:4] != "www" {
		// redirect to www
		log.Println("redirecting:", "https://www."+host)
		http.Redirect(w, r, "https://www."+host, http.StatusPermanentRedirect)
	}
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
