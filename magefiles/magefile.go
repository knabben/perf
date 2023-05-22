//go:build mage

package main

import (
	_ "embed"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"log"
	"os"
	"path/filepath"
)

const (
	exporterPath = "/usr/local/bin/node_exporter"
	exporterName = "node_exporter-1.5.0.linux-amd64"
	exporterURL  = "https://github.com/prometheus/node_exporter/releases/download/v1.5.0/" + exporterName + ".tar.gz"
	exporterTAR  = "node.tar.gz"
	binaryFolder = "/usr/local/bin/"

	systemdPath = "/etc/systemd/system/"
)

type (
	Prom mg.Namespace
	Node mg.Namespace
)

var (
	logger *log.Logger

	//go:embed templates/node_exporter.service
	nodeExporterSvc     string
	nodeExporterSvcPath = filepath.Join(systemdPath, "node_exporter.service")

	//go:embed templates/prometheus.yaml
	prometheusCR string

	//go:embed templates/grafana.yaml
	grafanaSpec string
)

func init() {
	logger = log.New(os.Stdout, "", log.Ldate)
}

func withCmd(cmd string, args []string) error {
	logger.Printf("%s %s", cmd, args)
	return sh.Run(cmd, args...)
}

func runOrFatal(cmd string, args []string) {
	if err := withCmd(cmd, args); err != nil {
		logger.Fatal(err)
	}
}

func createWriteFile(path string, content []byte) error {
	if _, err := os.Stat(path); err != nil {
		file, err := os.Create(path)
		if err != nil {
			logger.Fatal(err)
		}
		defer file.Close()
		_, err = file.Write(content)
		return err
	}
	return nil
}
