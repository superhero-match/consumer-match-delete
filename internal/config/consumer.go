package config

// Consumer holds the configuration values for the Kafka consumer.
type Consumer struct {
	Brokers []string `env:"KAFKA_BROKERS" default:"[192.168.178.17:9092]"`
	Topic   string   `env:"KAFKA_DELETE_MATCH_MATCH" default:"delete.match.match"`
	GroupID string   `env:"KAFKA_CONSUMER_MATCH_DELETE_GROUP_ID" default:"consumer-match-delete-group"`
}
