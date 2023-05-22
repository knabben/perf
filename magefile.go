//go:build mage

package main

import (
	_ "embed"
	"fmt"
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

var (
	logger *log.Logger

	//go:embed templates/node_exporter.service
	nodeExporterSvc     string
	nodeExporterSvcPath = filepath.Join(systemdPath, "node_exporter.service")
)

func init() {
	logger = log.New(os.Stdout, "", log.Ldate)
}

// Download node exporter and save on binary folder
func Download() error {
	if err := withCmd("which", []string{exporterPath}); err != nil {
		// No exporter found, downloading
		runOrFatal("curl", []string{
			"-L",
			exporterURL,
			"-s",
			"-o",
			exporterTAR,
		})

		// Decompress downloaded tar file
		runOrFatal("tar", []string{
			"zxvf",
			exporterTAR,
			filepath.Join(exporterName, "node_exporter"),
			"--strip-components",
			"1",
		})

		// Move to binary folder
		runOrFatal("mv", []string{"node_exporter", binaryFolder})

		// Remove compressed file
		runOrFatal("rm", []string{exporterTAR})

		return nil
	}

	logger.Println(fmt.Sprintf("%s was found. Skipping rule.", exporterPath))
	return nil
}

// Install the systemd
func Install() error {
	for _, f := range []map[string]string{
		{
			"path":    nodeExporterSvcPath,
			"content": nodeExporterSvc,
		},
	} {
		err := createWriteFile(f)
		if err != nil {
			logger.Fatal(err)
		}

		for _, n := range []string{"start", "enable"} {
			runOrFatal("systemctl", []string{n, "node_exporter"})
		}
	}
	return nil
}

// Clean removes the service from systemd
func Clean() {
	runOrFatal("systemctl", []string{"disable", "node_exporter"})
	runOrFatal("systemctl", []string{"stop", "node_exporter"})
	runOrFatal("rm", []string{nodeExporterSvcPath})
	runOrFatal("systemctl", []string{"daemon-reload"})
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

func createWriteFile(f map[string]string) error {
	if _, err := os.Stat(f["path"]); err != nil {
		file, err := os.Create(f["path"])
		if err != nil {
			logger.Fatal(err)
		}
		defer file.Close()
		_, err = file.Write([]byte(f["content"]))
		return err
	}
	return nil
}
