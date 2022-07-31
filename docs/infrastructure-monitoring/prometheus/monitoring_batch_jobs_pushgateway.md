---
id: monitoring_batch_jobs_with_the_pushgateway
title: Monitoring Batch Jobs with the Pushgateway
sidebar_label: Monitoring Batch Jobs with the Pushgateway
sidebar_position: 11
---

Batch jobs and other ephemeral jobs often don't run long enough to be scrapped realiably by Prometheus. For such batch jobs, we can use the Prometheus Pushgateway. 

The batch job can proactively send its metrics to the Pushgateway and which can then hold onto those metrics and expose them to Prometheus for continouous scrapping so we can easily query them. 

![](./img/pushgateway-architecture-overview-723d31f6f844746c2b7a66af6800b5a8.svg#center)