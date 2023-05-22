## Perf testing

DO NOT RUN in production. A bunch of scripts to install Prometheus node_exporter
on a systemd based test environment + Automatic helm installation of 
Prometheus operator on a Kind cluster.

### Mage compile

```
mage --compile ./mage
```

### Prometheus and Grafana

#### Install Prometheus operator

On a working cluster install the Prometheus Operator, the default CR and Grafana

```
./mage installMonitoring
```

#### Adding static targets

To add new targets create a new secret

```
./merge targetMonitoring --targets="xyz:9100,zvh:9100"
```

### Node Worker

#### Download node_exporter

This will install node_exporter on /usr/local/bin

```
sudo ./mage downloadExporter
```

#### Install the node_exporter systemd service

Creates a systemd listening on 9100 for all NICs.

```
sudo ./mage installExporter
```

#### Delete node_exporter systemd service

```
sudo ./mage cleanExporter
```