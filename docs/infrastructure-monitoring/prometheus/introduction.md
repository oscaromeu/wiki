---
id: introduction
title: Introduction
sidebar_label: Introduction
sidebar_position: 1
---

## Configuring Prometheus

Save the following Prometheus configuration as a file named `prometheus.yml` This configures Prometheus to scrape metrics from itself every 5 seconds.

```bash
global:
  scrape_interval: 5s

scrape_configs:
  - job_name: 'prometheus'
    static_configs:
      - targets: ['localhost:9090']
```

+ `global`: configures certain global settings and default values.
+ `scrape_configs`: specifies which targets Prometheus should scrape. Targets may be statically configured via the `static_configs` parameter or dynamically discovered using one of the supported service-discovery mechanisms. 

Bind-mount the `prometheus.yml` from the host by running

```bash
docker run \
    -p 9090:9090 \
    -v /path/to/prometheus.yml:/etc/prometheus/prometheus.yml \
    prom/prometheus
```