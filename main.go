package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	collector "github.com/rca0/kafka-connect-exporter/app/collector"
	"github.com/rca0/kafka-connect-exporter/misc"
	"github.com/sirupsen/logrus"
)

const (
	nameSpace  = "kafka_connect"
	versionUrl = "https://github.com/rca0/kafka-connect-exporter"
)

var (
	showVersion   = flag.Bool("version", false, "show version and exit")
	listenAddress = flag.String("listen-address", ":8000", "Address on which to expose metrics")
	metricsPath   = flag.String("telemetry-path", "/metrics", "Path under which to expose metrics")
	scrapeURI     = flag.String("scrape-uri", "http://127.0.0.1:8000,http://127.0.0.1:8003,http://127.0.0.1:8080", "URI on which to scrape Kafka connect.")
	gitRef        string
	gitSHA        string
)

func main() {
	logrus.Info("Starting kafka_connect_exporter")

	if gitRef != "" {
		logrus.Debug("found version details built in, will use that...")
		misc.SetGitRevisionDetails(gitRef, gitSHA)
	} else {
		misc.DetectVersionDuringRuntime()
	}

	flag.Parse()

	if *showVersion {
		fmt.Print("kafka_connect_exporter\n url: %s\n", versionUrl)
		fmt.Print("version %s -- git SHA %s -- git ref %s", misc.GetVersion(), gitSHA, gitRef)
		os.Exit(2)
	}

	uris := strings.Split(*scrapeURI, ",")

	prometheus.Unregister(collectors.NewGoCollector())
	prometheus.Unregister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))

	registeredURIs := make(map[string]bool)

	for _, u := range uris {
		co := collector.NewCollector(u, nameSpace)
		if _, ok := registeredURIs[u]; !ok {
			if err := tryRegister(co); err != nil {
				logrus.Errorf("error entering metric for URI: %s\n, error %v", u, err)
				continue
			}
			registeredURIs[u] = true
		}
	}

	http.Handle(*metricsPath, promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("OK"))
		if err != nil {
			logrus.Errorf("internal server error: %v", err)
			os.Exit(1)
		}
	})

	log.Fatal(http.ListenAndServe(*listenAddress, nil))
}

func tryRegister(metric prometheus.Collector) error {
	defer func() {
		if r := recover(); r != nil {
			err := fmt.Errorf("%v", r)
			logrus.Infof(fmt.Sprintf("%+v\n", err))
		}
	}()
	prometheus.MustRegister(metric)
	return nil
}
