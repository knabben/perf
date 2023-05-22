//go:build mage

package main

import (
	"fmt"
	"github.com/magefile/mage/sh"
)

const (
	tempCR      = "/tmp/.prom_cr"
	tempT       = "/tmp/target.yaml"
	tempGrafana = "/tmp/.prom_grafana"

	prometheusBundle = "https://raw.githubusercontent.com/prometheus-operator/prometheus-operator/main/bundle.yaml"
)

var template = `- job_name: "catalogue"
  static_configs:
    - targets: ["%s"]
`

// Install Prometheus Operator and initial CR
func (Prom) Install(targets string) error {
	// Install Prometheus Operator from bundle
	runOrFatal("kubectl", []string{"create", "-f", prometheusBundle})

	// Create a new target file for secret
	err := createWriteFile(tempT, []byte(fmt.Sprintf(template, targets)))
	if err != nil {
		return err
	}
	runOrFatal("kubectl", []string{
		"create",
		"secret",
		"generic",
		"additional-scrape-configs",
		"--from-file=" + tempT,
	})

	// Create a new Prometheus CR
	err = createWriteFile(tempCR, []byte(prometheusCR))
	if err != nil {
		return err
	}
	runOrFatal("kubectl", []string{"create", "-f", tempCR})

	// Create the Grafana spec
	err = createWriteFile(tempGrafana, []byte(grafanaSpec))
	if err != nil {
		return err
	}
	runOrFatal("kubectl", []string{"create", "-f", tempGrafana})

	return nil
}

// Remove Prometheus assets
func (Prom) Clean() error {
	sh.Run("kubectl", "delete", "-f", tempCR)
	sh.Run("kubectl", "delete", "-f", prometheusBundle)
	sh.Run("kubectl", "delete", "-f", tempGrafana)
	sh.Run("kubectl", "delete", "secret", "additional-scrape-configs")
	runOrFatal("rm", []string{tempCR, tempT, tempGrafana})
	return nil
}
