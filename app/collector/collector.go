package collector

import (
	"net/url"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/rca0/kafka-connect-exporter/app/client"
	"github.com/sirupsen/logrus"
)

type connectors []string

type Collector struct {
	client                  client.Client
	URI                     string
	Up                      prometheus.Gauge
	ConnectorsCount         prometheus.Gauge
	IsConnectorRunning      *prometheus.Desc
	AreConnectorTaskRunning *prometheus.Desc
}

var supportedSchema = map[string]bool{
	"http":  true,
	"https": true,
}

func NewCollector(uri, nameSpace string) *Collector {
	parseURI, err := url.Parse(uri)
	if err != nil {
		logrus.Errorf("%v", err)
		os.Exit(1)
	}
	if !supportedSchema[parseURI.Scheme] {
		logrus.Error("scheme not supported")
		os.Exit(1)
	}

	logrus.Infoln("collecting data from: ", uri)

	return &Collector{
		client: client.NewClient(uri),
		URI:    uri,
		Up: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: nameSpace,
			Name:      "up",
			Help:      "was the last scrape of kafka successful?",
		}),
		ConnectorsCount: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: nameSpace,
			Subsystem: "connectors",
			Name:      "count",
			Help:      "number of deployed connectors",
		}),
		IsConnectorRunning: prometheus.NewDesc(
			prometheus.BuildFQName(nameSpace, "connector", "state_running"),
			"is the connector running?",
			[]string{"connector", "consumer_group", "state", "worker_id", "topic", "bucket_name", "source_database_host", "source_database_name"},
			nil,
		),
		AreConnectorTaskRunning: prometheus.NewDesc(
			prometheus.BuildFQName(nameSpace, "connector", "tasks_state"),
			"the state of tasks. 0-failed, 1-running, 2-unassigned, 3-paused",
			[]string{"connector", "consumer_group", "state", "worker_id", "id", "topic", "bucket_name", "source_database_host", "source_database_name"},
			nil,
		),
	}
}

func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	c.Up.Describe(ch)
}
