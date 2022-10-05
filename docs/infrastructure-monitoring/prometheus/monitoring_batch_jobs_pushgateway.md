---
id: monitoring_batch_jobs_with_the_pushgateway
title: Monitoring Batch Jobs with the Pushgateway
sidebar_label: Monitoring Batch Jobs with the Pushgateway
sidebar_position: 11
---

Batch jobs and other ephemeral jobs often don't run long enough to be scrapped realiably by Prometheus. For such batch jobs, we can use the Prometheus Pushgateway. 

The batch job can proactively send its metrics to the Pushgateway and which can then hold onto those metrics and expose them to Prometheus for continouous scrapping so we can easily query them. 

![](./img/pushgateway-architecture-overview-723d31f6f844746c2b7a66af6800b5a8.svg#center)

The Pushgateway is a core Prometheus component that runs as a separate server and allows other jobs to push labeled groups of metrics to it over HTTP using Prometheus' text-based exposition format. The Pushgateway caches the last received sample value for each metric and exposes all pushed metrics to Prometheus for regular scraping.

+ Pushgateway does not act as an aggregator or event counter. It only stores the latest value of each metric that was pushed to it. Prometheus will then periodically scrape the last value.

+ Multiple jobs can push metrics to the Pushgateway without interfering with each other by defining grouping labels tipically the job label. A pushing job can also choose wheter it wants to replace an entire group of metrics or only update metrics that are included in the latest push to that group (`PUT` vs. `POST` HTTP methods).