package actions

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
)

// DownloadWireGuardClientConfig allows the user to download a client configuration for WireGuard
func DownloadWireGuardClientConfig(c buffalo.Context) error {
	clientName := c.Param("client_name")
	userIP := c.Request().Header.Get("X-Forwarded-For")
	if userIP == "" {
		userIP = c.Request().RemoteAddr
	}
	// Retrieve the client configuration from the SQLite3 database
	db, err := pop.Connect("development")
	if err != nil {
		return c.Error(http.StatusInternalServerError, err)
	}
	defer db.Close()

	client := Client{}
	err = db.Where("name = ?", clientName).First(&client)
	if err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	userIPStr := fmt.Sprintf("%s", userIP)

	// Check if the user's IP address matches the IP address associated with the client
	var updatedConfig []byte
	if userIPStr != client.IP {
		// Read the client configuration from the file
		file, err := ioutil.ReadFile(fmt.Sprintf("%s.conf", clientName))
		if err != nil {
			return c.Error(http.StatusInternalServerError, err)
		}

		// Update the client's IP address in the configuration file
		updatedConfig = []byte(strings.Replace(string(file), client.IP, userIPStr, 1))

		// Write the updated configuration to disk
		err = ioutil.WriteFile(fmt.Sprintf("%s.conf", clientName), updatedConfig, 0644)
		if err != nil {
			return c.Error(http.StatusInternalServerError, err)
		}
	} else {
		// Read the client configuration from the file
		file, err := ioutil.ReadFile(fmt.Sprintf("%s.conf", clientName))
		if err != nil {
			return c.Error(http.StatusInternalServerError, err)
		}
		updatedConfig = file
	}

	// Return the client configuration to the user as a download
	c.Response().Header().Set("Content-Type", "application/octet-stream")
	c.Response().Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s.conf", clientName))
	c.Response().Write(updatedConfig)
	return nil
}
