package main

import (
	"crypto/tls"
	"github.com/elastic/elastic-transport-go/v8/elastictransport"
	"github.com/elastic/go-elasticsearch/v8"
	"log"
	"net/http"
	"time"
)

const addr = "http://10.211.55.18:9200"

func main() {
	cfg := elasticsearch.Config{
		Addresses: []string{
			addr,
		},
		//Username: "foo",
		//Password: "bar",
		Transport: &http.Transport{
			MaxIdleConnsPerHost:   10,
			ResponseHeaderTimeout: time.Second,
			TLSClientConfig: &tls.Config{
				MinVersion: tls.VersionTLS12,
			},
		},
	}

	es, _ := elasticsearch.NewClient(cfg)
	log.Print(es.Transport.(*elastictransport.Client).URLs())
}
