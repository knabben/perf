//go:build mage

package main

// Install Prometheus Operator and initial CR
func (Prom) Install() {
	// Install Prometheus Operator from bundle
	runOrFatal("kubectl", []string{"create", "-f", prometheusBundle})

	// Create a new CR
	tempCR := "/tmp/prom_cr"
	createWriteFile(tempCR, []byte(prometheusCR))
	runOrFatal("kubectl", []string{"create", "-f", tempCR})
}
