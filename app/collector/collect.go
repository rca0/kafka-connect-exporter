package collector

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	var (
		topic                 string
		isRunning             float64 = 0
		connectorConfig       config
		connectorSinkConfig   sinkConfig
		connectorJdbc         sourceJdbcConfig
		connectorSourceConfig sourceConfig
		connectorStatus       status
	)

	c.Up.Set(0)

	resp, err := c.client.Get("/connectors")
	if err != nil {
		if os.IsTimeout(err) {
			logrus.Errorf("timeout occurred while fetching: %v", err)
		} else {
			logrus.Errorf("can't scrape kafka connect: %v", err)
		}
		return
	}
	defer func() {
		err = resp.Body.Close()
		if err != nil {
			logrus.Errorf("can't close connection to kafka connect: %v", err)
		}
	}()

	output, err := io.ReadAll(resp.Body)
	if err != nil {
		logrus.Errorf("can't scrape kafka connect: %v", err)
		return
	}

	var connectorsList connectors
	if err := json.Unmarshal(output, &connectorsList); err != nil {
		logrus.Errorf("error when json unmarshal, can't scrape kafka connect: %v", err)
		return
	}

	c.ConnectorsCount.Set(float64(len(connectorsList)))

	ch <- c.Up
	ch <- c.ConnectorsCount

	for _, connector := range connectorsList {
		// connectors  /config path
		logrus.Infof("fetch connenctor: %s", connector)

		connectorsConfigResponse, err := c.client.Get("/connectors" + connector + "/config")
		if err != nil {
			if os.IsTimeout(err) {
				logrus.Errorf("timeout occurred while fetching: %v", err)
				continue
			} else {
				logrus.Errorf("can't scrape /config for: %v", err)
				continue
			}
		}

		connectorsConfigOutput, err := io.ReadAll(connectorsConfigResponse.Body)
		if err != nil {
			logrus.Errorf("can't read Body for: %v", err)
			continue
		}

		if err := json.Unmarshal(connectorsConfigOutput, &connectorConfig); err != nil {
			logrus.Errorf("error while json unmarshal, can't decode response for: %v", err)
			continue
		}

		if err := json.Unmarshal(connectorsConfigOutput, &connectorSinkConfig); err != nil {
			logrus.Errorf("error while json unmarshal, can't decode response for: %v", err)
			continue
		}

		if err := json.Unmarshal(connectorsConfigOutput, &connectorSourceConfig); err != nil {
			logrus.Errorf("error while json unmarshal, can't decode response for: %v", err)
			continue
		}

		if err := json.Unmarshal(connectorsConfigOutput, &connectorJdbc); err != nil {
			logrus.Errorf("error while json unmarshal, can't decode response for: %v", err)
			continue
		}

		err = connectorsConfigResponse.Body.Close()
		if err != nil {
			logrus.Errorf("can't close connection to connector: %v", err)
		}

		// connectors /status path

		connectorsStatusResponse, err := c.client.Get("/connectors" + connector + "/status")
		if err != nil {
			if os.IsTimeout(err) {
				logrus.Errorf("timeout occurred while fetching: %v", err)
				continue
			} else {
				logrus.Errorf("can't scrape /status for: %v", err)
				continue
			}
		}

		connectorsStatusOutput, err := io.ReadAll(connectorsStatusResponse.Body)
		if err != nil {
			logrus.Errorf("can't read Body for: %v", err)
			continue
		}

		if err := json.Unmarshal(connectorsStatusOutput, &connectorStatus); err != nil {
			logrus.Errorf("error while json unmarshal, can't decode response for: %v", err)
			continue
		}

		if strings.ToLower(connectorStatus.Connector.State) == "running" {
			isRunning = 1
		}

		if !strings.Contains(connectorJdbc.Name, "sink") || !strings.Contains(connectorJdbc.Name, "source") {
			topic = connectorJdbc.TopicName()
		}

		if strings.Contains(connectorConfig.ConnectorClass, "S3Sink") {
			topic = connectorSinkConfig.Topic
		}

		if strings.Contains(connectorSourceConfig.Name, "source") {
			topic = connectorSourceConfig.Topic
		}

		ch <- prometheus.MustNewConstMetric(
			c.IsConnectorRunning, prometheus.GaugeValue, isRunning,
			connectorStatus.Name, connectorStatus.ConsumerGroup(),
			strings.ToLower(connectorStatus.Connector.State),
			connectorStatus.Connector.WorkerId,
			topic, connectorSinkConfig.BucketName(), connectorSourceConfig.DBHost(), connectorSourceConfig.DBName(),
		)

		for _, connectorTask := range connectorStatus.Tasks {

			var state float64
			switch taskState := strings.ToLower(connectorTask.State); taskState {
			case "running":
				state = 1
			case "unassigned":
				state = 2
			case "paused":
				state = 3
			default:
				state = 0
			}

			ch <- prometheus.MustNewConstMetric(
				c.AreConnectorTaskRunning, prometheus.GaugeValue, state,
				connectorStatus.Name, connectorStatus.ConsumerGroup(),
				strings.ToLower(connectorTask.State),
				connectorTask.WorkerId, fmt.Sprintf("%d", int(connectorTask.Id)),
				topic, connectorSinkConfig.BucketName(), connectorSourceConfig.DBHost(), connectorSourceConfig.DBName(),
			)
		}

		err = connectorsStatusResponse.Body.Close()
		if err != nil {
			logrus.Errorf("cannot close connection to connector: %v", err)
		}
	}
}
