# Grow

Plant monitoring system.

Uses Raspberry PIs, NATS, and Prometheus.

## How

Collects information about plants and soil moisture.

Publishs the collected data into a NATS cluster.

Processes the information and displays it with Prometheus/Grafana.

## Structure

* `charts` - helm charts from this repository (more info [here](https://renato0307.github.io/grow/))
* `k8s` - setups a k8s cluster to run NATS and Prometheus
* `monitor-ghm` - monitor running on raspberry pi with Pimonori Grow HAT Mini (GHM)
* `ingestion-service` - service to ingest readings from NATS jetstream

## Useful NATS commands

|What|Command|
|----|-------|
|Stream summary|`nats --server nats://192.168.1.131:4222 stream report`|
|Delete stream |`nats --server nats://192.168.1.131:4222 stream rm PlantReadings`|
|List stream consumers|`nats --server nats://192.168.1.131:4222 consumer ls PlantReadings`|
|Delete stream consumer|`nats --server nats://192.168.1.131:4222 consumer rm PlantReadings PlantReadingsIngestion`|

## TODO

### Milestone 1 - home plants

1. ~~Send data to NATS server/jetstreams~~
1. ~~Consume jetstream messages and publish to prometheus~~
1. Dinamic configuration of the sensors
1. ~~Show data on Grafana~~
1. Calibrate sensors
1. Alarms for plants with low soil moisture
1. Loadbalancer and external IP for NATS
1. NATS security (TLS, auth, etc.)
1. Mobile/Slack notifications
1. IaC for K8s cluster (FluxCD)

### Milestone 2 - remote monitoring

1. Acquire Raspberry Pico or ESP32, basic sensors, batery and GSM module
1. Basic monitor
1. Configure NATS to support MQTT
1. Comunication MQTT via GSM
1. Deep sleep
1. Measure consumptions and batery life
1. Batery life optimizations

## Build

### Helm charts

The [chart releaser action](https://github.com/helm/chart-releaser-action) is used
to release helm charts.

Check the workflow file at `.github/workflows/helm-release.yaml`.

### Docker images

Each service has its own versioning.

We use [GitVersion](https://gitversion.net/) to automatically generate
new versions using [SemVer](https://semver.org/).

Each service contains a `GitVersion.yml` file and a release action.

Example for the ingestion service:

* `ingestion-service/GitVersion.yml`
* `.github/workflows/ingestion-service-release.yaml`