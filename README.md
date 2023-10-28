# Grow

Plant monitoring system.

Uses Raspberry PIs, NATS, and Prometheus.

## How

Collects information about plants and soil moisture.

Publishs the collected data into a NATS cluster.

Processes the information and displays it with Prometheus/Grafana.

## Structure

* `k8s` - setups a k8s cluster to run NATS and Prometheus
* `raspi-grow-hat-mini` - monitor running on raspberry pi zero with Pimonori Grow HAT Mini (GHM)

## TODO

### Milestone 1 - home plants

1. ~~Send data to NATS server/jetstreams~~
1. Consume jetstream messages and publish to prometheus
1. Dinamic configuration of the sensors
1. Show data on Grafana
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