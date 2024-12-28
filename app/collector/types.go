package collector

import "strings"

type status struct {
	Name      string    `json:"name"`
	Connector connector `json:"connector"`
	Tasks     []task    `json:"tasks"`
}

type connector struct {
	State    string `json:"state"`
	WorkerId string `json:"worker_id"`
}

type task struct {
	State    string  `json:"state"`
	Id       float64 `json:"id"`
	WorkerId string  `json:"worker_id"`
}

type config struct {
	ConnectorClass string `json:"connector.class"`
}

type sinkConfig struct {
	ConnectorClass string `json:"connector.class"`
	Name           string `json:"name"`
	S3BucketName   string `json:"s3.bucket_name"`
	Topic          string `json:"topics"`
}

type sourceConfig struct {
	Name             string `json:"name"`
	DatabaseHostname string `json:"database.hostname"`
	DatabaseName     string `json:"database.names"`
	Topic            string `json:"transforms.Reroute.topic.replacement"`
}

type sourceJdbcConfig struct {
	Name  string `json:"name"`
	Topic string `json:"topic.prefix"`
}

func (s status) ConsumerGroup() string {
	// ignore source
	if strings.Contains(s.Name, "source") {
		return ""
	}
	// add connect- into prefix name
	return "connect-" + s.Name
}

func (s sinkConfig) BucketName() string {
	// ignore source
	if strings.Contains(s.Name, "source") {
		return ""
	}
	return string(s.S3BucketName)
}

func (s sourceConfig) DBName() string {
	// ignore sinkConfigs: if source not contains in name, return empty
	if !strings.Contains(s.Name, "source") {
		return ""
	}
	return string(s.DatabaseName)
}

func (s sourceConfig) DBHost() string {
	// ignore sinkConfigs: if source not contains in name, return empty
	if !strings.Contains(s.Name, "source") {
		return ""
	}
	return string(s.DatabaseHostname)
}

func (j sourceJdbcConfig) TopicName() string {
	return string(j.Topic)
}
