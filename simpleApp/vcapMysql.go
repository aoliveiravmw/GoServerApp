package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type MySQLCredentials struct {
	Hostname string `json:"hostname"`
	JDBCUrl  string `json:"jdbcUrl"`
	Name     string `json:"name"`
	Password string `json:"password"`
	Port     int    `json:"port"`
	URI      string `json:"uri"`
	Username string `json:"username"`
}

type MySQLService struct {
	BindingName    *string          `json:"binding_name"`
	Credentials    MySQLCredentials `json:"credentials"`
	InstanceName   string           `json:"instance_name"`
	Label          string           `json:"label"`
	Name           string           `json:"name"`
	Plan           string           `json:"plan"`
	Provider       *string          `json:"provider"`
	SyslogDrainURL *string          `json:"syslog_drain_url"`
	Tags           []string         `json:"tags"`
	VolumeMounts   []interface{}    `json:"volume_mounts"`
}

func vcapSqlService(cfMysql *MySQLService) error {
	varvalue := os.Getenv("VCAP_SERVICES")
	if varvalue == "" {
		return fmt.Errorf("Environment variable %s is not set.\n", "VCAP_SERVICES")
	}

	var jsonData map[string][]MySQLService
	err := json.Unmarshal([]byte(varvalue), &jsonData)
	if err != nil {
		return fmt.Errorf("Failed to parse JSON value of %s: %s\n", "VCAP_SERVICES", err)
	}

	services, exists := jsonData["p.mysql"]
	if !exists || len(services) == 0 {
		fmt.Printf("No 'p.mysql' services found in the JSON.\n")
		return nil
	}
	*cfMysql = services[0]

	/* For Debugging
	   fmt.Printf("Binding Name: %v\n", service.BindingName)
	   fmt.Printf("Instance Name: %s\n", service.InstanceName)
	   fmt.Printf("Label: %s\n", service.Label)
	   fmt.Printf("Name: %s\n", service.Name)
	   fmt.Printf("Plan: %s\n", service.Plan)
	   fmt.Printf("Provider: %v\n", service.Provider)
	   fmt.Printf("Syslog Drain URL: %v\n", service.SyslogDrainURL)
	   fmt.Printf("Tags: %v\n", service.Tags)

	   fmt.Println("Credentials:")
	   fmt.Printf("  Hostname: %s\n", service.Credentials.Hostname)
	   fmt.Printf("  JDBC URL: %s\n", service.Credentials.JDBCUrl)
	   fmt.Printf("  Name: %s\n", service.Credentials.Name)
	   fmt.Printf("  Password: %s\n", service.Credentials.Password)
	   fmt.Printf("  Port: %d\n", service.Credentials.Port)
	   fmt.Printf("  URI: %s\n", service.Credentials.URI)
	   fmt.Printf("  Username: %s\n", service.Credentials.Username)

	   fmt.Printf("Volume Mounts: %v\n", service.VolumeMounts)
	*/
	return nil
}
