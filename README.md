## Perf testing

DO NOT RUN in production. A bunch of scripts to install Prometheus node_exporter
on a systemd based test environment + Automatic helm installation of 
Prometheus operator on a Kind cluster.

### Mage compile

```
mage --compile ./mage
```

### Worker Node Exporter

#### Download node_exporter

This will install node_exporter on /usr/local/bin

```
sudo ./mage download
```

#### Install node_exporter

Creates a systemd listening on 9100 for all NICs.

```
sudo ./mage install
```