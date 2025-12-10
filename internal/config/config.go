package config

import (
	"log"
	"os"
	"strconv"
)

func GetAppPort() int {
	const defaultPort = 8080

	portStr := os.Getenv("APP_PORT")
	if portStr == "" {
		return defaultPort
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Printf("Invalid app port %q. Service will use default port: %d", portStr, defaultPort)
		return defaultPort
	}
	return port
}
