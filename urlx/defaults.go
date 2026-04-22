package urlx

import "strings"

var gDefaultPorts = map[string]uint16{
	"http":              80,
	"https":             443,
	"ssh":               22,
	"dns":               53,
	"mysql":             3306,
	"postgres":          5432,
	"pg":                5432,
	"oracle":            1521,
	"mssql":             1433,
	"tidb":              4000,
	"redis":             6379,
	"memcached":         11211,
	"mongodb":           27017,
	"etcd":              2379,
	"zookeeper":         2181,
	"amqp":              5672,
	"kafka":             9092,
	"pulsar":            6650,
	"mqtt":              1883,
	"rocketmq":          9876,
	"prometheus":        9090,
	"grafana":           3000,
	"jaeger":            16686,
	"elasticsearch":     9200,
	"es":                9200,
	"elasticsearch-api": 9300,
	"es-api":            9300,
	"clickhouse":        8123,
	"clickhouse-api":    9000,
	"minio":             9000,
}

func DefaultPort(scheme string) (port uint16, ok bool) {
	port, ok = gDefaultPorts[strings.ToLower(scheme)]
	return
}
