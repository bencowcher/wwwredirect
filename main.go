package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"os"

	_ "embed"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

//go:embed hosts.json
var hosts []byte

var hostmap = map[string]string{}

func loadEmbeddedHosts() map[string]string {
	var hostmap map[string]string
	err := json.Unmarshal(hosts, &hostmap)
	if err != nil {
		log.Fatal(err)
	}

	return hostmap
}

func loadKubernetesHosts(ns string) map[string]string {
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	confmap, err := clientset.CoreV1().ConfigMaps(ns).Get(context.Background(), "hosts", metav1.GetOptions{})
	if err != nil {
		panic(err.Error())
	}
	return confmap.Data
}

// determines the domain name requested and returns the www redirect
func wwwRedirect(w http.ResponseWriter, r *http.Request) {
	// get the host
	host := r.Host

	log.Println(host)

	rd := hostmap[host]
	if rd != "" {
		// redirect to www
		loc, err := url.JoinPath(rd, r.URL.Path)
		if err != nil {
			log.Println("error joining path:", err)
		} else {
			rd = loc
		}
		log.Println("redirecting:", rd)
		http.Redirect(w, r, rd, http.StatusPermanentRedirect)
		return
	}

	http.NotFound(w, r)
}

func main() {
	src := os.Getenv("APP_HOST_SOURCE") // embedded, kubernetes
	if src == "" {
		src = "embedded"
	}

	if src == "kubernetes" {
		hostmap = loadKubernetesHosts(os.Getenv("APP_KUBERNETES_NAMESPACE"))
	} else {
		hostmap = loadEmbeddedHosts()
	}

	// create the server
	server := http.Server{
		Addr: "0.0.0.0:80",
	}

	// handle the www redirect
	http.HandleFunc("/", wwwRedirect)

	// start the server
	server.ListenAndServe()
}
