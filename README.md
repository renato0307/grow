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

