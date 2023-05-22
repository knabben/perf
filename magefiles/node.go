//go:build mage

package main

import (
	"fmt"
	"path/filepath"
)

// Download node exporter and save on binary folder
func (Node) Download() error {
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
func (Node) Install() error {
	for _, f := range []map[string]string{
		{
			"path":    nodeExporterSvcPath,
			"content": nodeExporterSvc,
		},
	} {
		err := createWriteFile(f["path"], []byte(f["content"]))
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
func (Node) Clean() {
	runOrFatal("systemctl", []string{"disable", "node_exporter"})
	runOrFatal("systemctl", []string{"stop", "node_exporter"})
	runOrFatal("rm", []string{nodeExporterSvcPath})
	runOrFatal("systemctl", []string{"daemon-reload"})
}
