//go:build mage

package main

import (
	"fmt"
	"github.com/magefile/mage/sh"
	"log"
	"os"
	"path"
)

const (
	exporterPath = "/usr/local/bin/node_exporter"
	exporterName = "node_exporter-1.5.0.linux-amd64"
	exporterURL  = "https://github.com/prometheus/node_exporter/releases/download/v1.5.0/" + exporterName + ".tar.gz"
	exporterTAR  = "node.tar.gz"
	binaryFolder = "/usr/local/bin"
)

func init() {
	log.SetOutput(os.Stdout)
}

// Install node exporter as systemd service
func Install() error {
	if err := withCmd("which", []string{exporterPath}); err != nil {
		// No exporter found, downloading
		args := []string{
			"-L",
			exporterURL,
			"-s",
			"-o",
			exporterTAR,
		}
		if err := withCmd("curl", args); err != nil {
			return err
		}

		// Decompress downloaded tar file
		args = []string{
			"zxvf",
			exporterTAR,
			"-C",
			binaryFolder,
			path.Join(exporterName, "node_exporter"),
			"--strip-components=1",
		}
		if err := withCmd("tar", args); err != nil {
			return err
		}

		// Remove compressed file
		if err := sh.Run("rm", exporterPath); err != nil {
			return err
		}

		return nil
	}

	log.Println(fmt.Sprintf("%s was found. Skipping rule.", exporterPath))
	return nil
}

func LocalInstall() {

}

func withCmd(cmd string, args []string) error {
	log.Printf("Running: %s %s", cmd, args)
	return sh.Run(cmd, args...)
}
