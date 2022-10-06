---
id: patterns
title: Patterns
sidebar_label: Patterns
sidebar_position: 21
---

## Structural Patterns

### Init Container

#### Problem

If you have one or more containers in a Pod that have prerequisites before starting up then you can use _init containers_ which allow separation of initializing activities from the main application duties.

#### Solution

_Init containers_ in Kubernetes are part of the Pod definition, and they separate all containers in a Pod into two groups: _init containers_ and application containters. All init containers are executed in a sequence, one by one, and all of them have to terminate succesfully before the application containers are started up. 

:::warning
If an init container fails, the whole Pod is restarted (unless it is marked with `RestartNever`), causing all init containers to run again. Thus to prevent any side effects, making init containers idempotent is a good practice. Thus, to prevent any side effects, making init containers idempotent is a good practice.
:::


In relation with health-checking and resource-handling semantics:

+ There is no readiness check for init containers because all init containers must terminate succesfully before the Pod startup processes can continue with application containers.

+ Init containers affects Pod resource requirements, scheduling, autoscaling, and quota management. The effective Pod-level request and limit value become the highest value of the following two:

    + The highest init container request/limit value
    + The sum of all application container values for request/limit

:::info Keep a Pod running

For debugging the outcome of init containers, it helps if the command of the application container is replaced temporarily with a dummy sleep command so that you have time to examine the situation. This trick is particularly useful if your init container fails to start up and your application fails to start because the configuration is missing or broken. The following command within the Pod
declaration gives you an hour to debug the volumes mounted by
entering the Pod with `kubectl exec -it <Pod> sh`:

```bash
command:
- /bin/sh
- "-c"
- "sleep 3600"
```

:::


##### Example 1. Shows an init container that copies data into an empty volume

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: www
  labels:
    app: www
spec:
  initContainers:
  - name: download
    image: axeclbr/git
    # Clone an HTML page to be served
    command:
    - git
    - clone
    - https://github.com/mdn/beginner-html-site-scripted
    - /var/lib/data
    # Shared volume with main container
    volumeMounts:
    - mountPath: /var/lib/data
      name: source
  containers:
  # Simple static HTTP server for serving these pages
  - name: run
    image: docker.io/centos/httpd
    ports:
    - containerPort: 80
    # Shared volume with main container
    volumeMounts:
    - mountPath: /var/www/html
      name: source
  volumes:
  - emptyDir: {}
    name: source
```