/*
  Copyright (C) 2019 - 2022 MWSOFT
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

// Health holds configuration for health server.
type Health struct {
	Port             string `env:"CONSUMER_MATCH_DELETE_HEALTH_SERVER_PORT" yaml:"port" default:":8030"`
	ShutdownEndpoint string `env:"CONSUMER_MATCH_DELETE_HEALTH_SERVER_SHUTDOWN_ENDPOINT" yaml:"shutdown_endpoint" default:"/api/v1/consumer_match_delete_health/shutdown"`
	ContentType      string `env:"CONSUMER_MATCH_DELETE_HEALTH_SERVER_CONTENT_TYPE" yaml:"content_type" default:"application/json"`
}
