/*
  Copyright (C) 2019 - 2020 MWSOFT
  This program is free software: you can redistribute it and/or modify
  it under the terms of the GNU General Public License as published by
  the Free Software Foundation, either version 3 of the License, or
  (at your option) any later version.
  This program is distributed in the hope that it will be useful,
  but WITHOUT ANY WARRANTY; without even the implied warranty of
  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
  GNU General Public License for more details.
  You should have received a copy of the GNU General Public License
  along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/
package config

// Consumer holds the configuration values for the Kafka consumer.
type Consumer struct {
	Brokers []string `env:"KAFKA_BROKERS" default:"[192.168.178.17:9092]"`
	Topic   string   `env:"KAFKA_DELETE_MATCH_MATCH" default:"delete.match.match"`
	GroupID string   `env:"KAFKA_CONSUMER_MATCH_DELETE_GROUP_ID" default:"consumer-match-delete-group"`
}
