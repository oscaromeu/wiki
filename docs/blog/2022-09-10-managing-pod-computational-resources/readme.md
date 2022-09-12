---
id: _pod_computational_resources_
title: Managing pod's computational resources
sidebar_label: Managing pod's computational resources
sidebar_position: 1
hide_table_of_contents: true
---

Setting both how much a pod is expected to consume and the maximum amount it's allowed to consume is vital part of any pod definition. Setting these two parameters makes sure that a pod takes only its fair share of the resources provided by the Kubernetes cluster and also affects how pods are scheduled across the cluster. 

:::info This post covers
+ Requesting CPU, memory, and other computational resources for containers
+ Setting a hard limit for CPU and memory
+ Understanding Quality of Service guarantees for pods
+ Setting default, min, and max resources for pods in a namespace
+ Limiting the total amount of resources available in a namespace
:::

## Request resources for a pod's containers

When we create a pod, we can optionally specify how much resources a container needs  (known as __requests__) and a hard limit on what it may consume (known as __limits__). 
Kubernetes has support for many different types of resources including CPU, memory, storage, network bandwidth, and the use of special devices such as graphical processing units (aka GPUs). We'll focus on the most commonly specified resources types: CPU and memory.
### Creating a pod with resource requests and limits
We will start our journey with an ordinary pod manifest spec:

```yaml
---
apiVersion: v1
kind: Pod
metadata:
  name: requests-pod
spec:
  containers:
  - image: busybox
    command: ["dd", "if=/dev/zero", "of=/dev/null"]
    name: main
    resources:
      requests:
        cpu: 200m
        memory: 10Mi
      limits:
        memory: 10Mi
```


The pod manifest specification is simple, it only has one container executing the command:

```bash
$ dd if=/dev/zero of=/dev/null
```

Our humble container requires one-fifth of a CPU core to run properly and also 10 mebibytes of memory.
## Requests for planning

First let's see how `requests` are used for scheduling workloads to nodes, later we will see that we can constrain our container on what it may consume, that is, setting a hard limit known as `limits`. By specifying resource requests we are telling to Kubernetes the minimum amount of resources our pod needs. Resource requests are used by the Scheduler to assign a node to a pod. Each node has a certain amount of CPU and memory to allocate pods on nodes. 

By default Pods can consume all the available capacity on a node. This can be an issue because nodes typically run quite a few system daemons that power the OS and Kubernetes itself. So it is the job of kubernetes administrator to set aside resources for these system daemons configuring `Node Allocatable` based on the workload density on each node in the cluster.

![Requests are for planning!](./img/managing_computational_resources_1.png#center "Requests are for planning")

:::info

Kubelet reports the resources that each node has to the API server which are available to the scheduler through the Node resource. We can see that by using the `kubectl describe` command.

```bash
$ kubectl describe nodes
Name:               lbmwk1
...
Capacity:
  cpu:                4
  ephemeral-storage:  61275608Ki
  hugepages-1Gi:      0
  hugepages-2Mi:      0
  hugepages-32Mi:     0
  hugepages-64Ki:     0
  memory:             3880868Ki
  pods:               110
Allocatable:
  cpu:                4
  ephemeral-storage:  59608911416
  hugepages-1Gi:      0
  hugepages-2Mi:      0
  hugepages-32Mi:     0
  hugepages-64Ki:     0
  memory:             3880868Ki
  pods:               110
...
```

The output shows two sets of amounts related to the available resources on the node:
the node’s capacity and allocatable resources. The capacity represents the total resources
of a node, which may not all be available to pods. Certain resources may be reserved
for Kubernetes and/or system components. The Scheduler bases its decisions only on
the allocatable resource amounts.
:::

When assigning a pod the scheduler will only consider nodes with enough unallocated resources to meet the pod’s resources requirements. The scheduler does not care about the resources that are being used at the exact time of scheduling but instead the sum of resources requested by the existing pods deployed on the node. Note that scheduling another pod based on actual resource usage would break the guarantee given to the already deployed pods. 

#### Example #1

Let's imagine that there are three pods (_PODs A_, _B_ and _C_) deployed on a node which they've requested 80% of the node's CPU and 60% of the node's memory. Pod D, shown at the bottom right of the figure, cannot be scheduled onto the node because the unallocated CPU is roughly a 20% and the pod requests is 25% which is more. Note that the three pods are using 70% of the CPU and this makes no differentece.

![The journey of a limit!](./img/managing_computational_resources_sched_1.png#center "The journey of a limit")


#### Example #2

### CPU requests affect CPU time sharing

The CPU requests are generally fulfilled at the Kubernetes scheduler level but they also determine how the remaining unused CPU time is distributed between pods as we will see later in more detail.

## QoS class in kubernetes

In addition to provide resource isolation, resource request and limits determine the pod QoS class. That is, every pod has a QoS class based on the request and limit values that we set. There are three types which are summarized in the following table

| **QoS Class** 	| **Condition**                                                                                                                           	| **Priority (Lower is better)** 	|
|---------------	|-----------------------------------------------------------------------------------------------------------------------------------------	|--------------------------------	|
| Guaranteed    	| Limits and optionally requests (not equal to 0) are set for all resources across all containers and they are equal                      	| 1                              	|
| Burstable     	| Requests and optionally limits are set (not equal to 0) for one or more resources across one or more containers, and they are not equal 	| 2                              	|
| BestEffort    	| Requests and limits are not set for all of the resources, across all containers                                                         	| 3                              	|


## Limits by level

Now we are going to talk about limits. The main thing to know about limits is that limits are for enforcing the rules. This way we say this pod can't go over this limit of resources. Let's take a look of what happens when a limit comes in to Kubernetes. 

### The Journey of a Pod Limit

Let's see how resource limit value is being used in Kubernetes. Once pod spec is registered to Kubernetes, kube scheduler fetches the new pod specs and then assigns a node to the Pod you want to create, but the limit value is not directly used at this moment yet. Kubelet on each node, runs the sync process to fetch the latest information of the Pod that are assigned. Kubelet, sees the limit value from pods spec then converts the CPU cores value to CFS period and quota milliseconds. Then kubelet calls container runtime interface to create actual containers on the Linux host. 

![The journey of a limit!](./img/managing_computational_resources_limits_1.png#center "The journey of a limit")

Once container runtimes get the entry from Kubelet it executes the container creation by calling cgroup's. Now let's deep dive into the three parts here: cgroup's, CFS period and quota and container runtimes.


## Container primitives

When we talk about containers, it is actually implemented with Linux kernel features that is called cgroup's and namespaces: __namespaces__ are used to isolate process on Linux hosts and __cgroups__ is used to limit resources. This time we talk about resource managemement in Kubernetes so we are not going to talk about namespaces.


### Control groups or cgroup's


Control groups often referred as cgroups is a mechanism to limit the resources,
such as memory, CPU, and network input/ouput that a group of processes can use. Kubelet and the underlying container runtime need to interface with cgroups to enforce resource management for pods and containers which includes cpu/memory requests and limits for containerized workloads. Note that this is needed to prevent starvation from one pod to another. (IMPROVE)

:::info
When we refer to control groups, we refer tp the unified version or simply version 2. The cgroup version depends on the Linux distribution and the default cgroup version configured on the OS. To check which cgroup version our distribution uses, we can run 
```bash 
stat -fc %T /sys/fs/cgroup/
``` 
the output for version will be `cgroup2fs` and for version 1 will be `tmpfs`.
:::

The CPU request value that we set in Kubernetes will be converted to CPU weight in cgroup world and the CPU limit will be converted from CPU core numbers to the CPU time value of CFS and stored in cgroups as CFS period and quota in `cpu.max`. Those cgroup values can be seen for our running pods on each node by seeing files under `/sys/fs/cgroup` s' directory. Memory limit is also stored in cgroups. In cgroup v1, and prior to this feature, the container runtime never took into account and effectively ignored `spec.containers[].resources.requests["memory"]`. This is unlike CPU, in which the container runtime consider both requests and limits. Furthermore, memory actually can't be compressed in cgroup v1. Because there is no way to throttle memory usage, if a container goes past its memory limit it will be terminated by the kernel with an OOM (Out of Memory) kill. In contrast, CPU is considered a "compressible" resource. If your app starts hitting your CPU limits, Kubernetes starts throttling your container, giving your app potentially worse performance. However, it won't be terminated. That is what "compressible" means.


| Kubernetes manifest parameter                 	| Cgroup V1               	| Cgroup V2     	| Comment                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                 	|
|-----------------------------------------------	|-------------------------	|---------------	|---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------	|
| `spec.containers[].resources.requests.cpu`    	| `cpu.shares`            	| `cpu.weight`  	| `cpu.shares` is replaced with `cpu.weight` and operates on the standard scale defined by CGROUP_WEIGHT_MIN/DFL/MAX (1, 100, 10000). The weight is scaled to scheduler weight so that 100 maps to 1024 and the ratio relationship is preserved - if weight is W and its scaled value is S, W / 100 == S / 1024.  While the mapped range is a bit smaller than the orignal scheduler weight range, the dead zones on both sides are relatively small and covers wider range than the nice value mappings. 	|
| `spec.containers[].resources.limits.cpu`      	| `cpu.max`               	|               	| `cpu.cfs_quota_us` and `cpu.cfs_period_us` are replaced by `cpu.max` which contains both quota and period.                                                                                                                                                                                                                                                                                                                                                                                              	|
| `spec.containers[].resources.requests.memory` 	| NA                      	| `memory.min`  	| `memory.min` specifies a minimum amount of memory the cgroup must always retain, i.e., memory that can never be reclaimed by the system. If the cgroup's memory usage reaches this low limit and can’t be increased, the system OOM killer will be invoked.                                                                                                                                                                                                                                             	|
| `spec.containers[].resources.limits.memory`   	| `memory.limit_in_bytes` 	| `memory.high` 	| `memory.high` is the memory usage throttle limit. This is the main mechanism to control a cgroup's memory use. If a cgroup's memory use goes over the high boundary specified here, the cgroup’s processes are throttled and put under heavy reclaim pressure. The default is max, meaning there is no limit.              


## Memory requests and memory limit

A Pod's memory request represents a guaranteed lower bound for the memory resource. However, they can consume additional memory if it's available on the node. This is problematic because the Pod uses memory that the scheduler can assign to other workloads or tenants. When a new Pod is scheduled onto the same node, the Pods may fight over the memory. To honor the memory requests of both Pods, the Pod consuming memory above its request is terminated.

In order to control the amount of memory that tenants can consume, we must include memory limits on the workloads, which enforce an upper bound on the amount of memory available to a given workload. If the workload attempts to con‐
sume memory above the limit, the workload is terminated. This is because memory is a noncompressible resource. There is no way to throttle memory, and thus the process must be terminated when the node’s memory is under contention. The following snippet shows a container that was out-of-memory killed (OOMKilled).

A common question we encounter in the field is whether one should allow tenants to set memory limits higher than requests. In other words, whether nodes should be oversubscribed on memory. This question boils down to a trade-off between node density and stability. When you oversubscribe your nodes, you increase node density but decrease workload stability. As we’ve seen, workloads that consume memory above their requests get terminated when memory comes under contention. In most cases, we encourage platform teams to avoid oversubscribing nodes, as they typically consider stability more important than tightly packing nodes. This is especially the case in clusters hosting production workloads.
Now that we’ve covered memory requests and limits, let’s shift our discussion to CPU. In contrast to memory, CPU is a compressible resource. You can throttle processes when CPU is under contention. For this reason, CPU requests and limits are some what more complex than memory requests and limits.

## CPU requests and CPU shares

CPU requests and limits are specified using CPU units. In most cases, 1 CPU unit is equivalent to 1 CPU core. Requests and limits can be fractional (e.g., 0.5 CPU) and they can be expressed using millis by adding an m suffix. 1 CPU unit equals 1000m CPU.

When containers within a Pod specify CPU requests, the scheduler finds a node with enough capacity to place the Pod. Once placed, the kubelet converts the requested CPU units into cgroup CPU shares. CPU shares is a mechanism in the Linux kernel that grants CPU time to cgroups (i.e, the processes within the cgroup). The following are critical aspects of CPU shares to keep in mind:

+ __CPU shares are relative__. 1000 CPU shares does not mean 1 CPU core or 1000 CPU cores. Instead, the CPU capacity is proportionally divided among all cgroups according to their relative shares. For example, consider two processes in
different cgroups. If process 1 (P1) has 2000 shares, and process 2 (P2) has 1000 shares, P1 will get twice the CPU time as P2. 

+ __CPU shares come into effect only when the CPU is under contention__. If the CPU is not fully utilized, processes are not throttled and can consume additional CPU cycles. Following the preceding example, P1 will get twice the CPU time as P2
only when the CPU is 100% busy 


+ Control of CPU scheduling is handled by the cgroup v2 __cpu controller__. Fair share scheduling is handled through a weight-based CPU time distribution model (altough you can also impose usage limits). When the CPU controller is enabled for a cgroup, a `cpu.weight` file appears in all children (with a default of 100). When distributing CPU time to the children, all of their `cpu.weight` values are summed up and then each active child gets CPU in proportion to their weight relative to the total. This means that if all `cpu.weight` files have all the same value, all children will get equal shares of the CPU time. The actual `cpu.weight` values only matter if they're different; if they're all the same the value is arbitrary.

![cpu-shares!](./img/cpu-shares.png#center "CPU shares")

CPU shares (CPU requests) provide the CPU resource isolation necessary to run different tenants on the same node. As long as tenants declare CPU requests, the CPU capacity is shared according to those requests. Consequently, tenants are unable to starve other tenants from getting CPU time. In other words, Pods are guaranteed to get the amount of CPU they requested, they may or may not get additional CPU time (depending on the other jobs running.). 

Excess CPU resources will be distributed based on the the amount of CPU requested. For example, suppose container A requests for 600m CPUs, and container B requests for 300m CPUs. Suppose that both containers are trying to use as much as CPU as they can. Then the extra 100m CPUs will be distributed to A and B in a 2:1 ratio. Pods will be throttled if they exceed their limit but if limit is unspecified, then the Pods can use excess CPU when available.

:::warning
Disable your CPU limits! If you give all your k8s pods accurate CPU requests, then no one can throttle them because CPU is reserved for them if they need it. This has nothing to do with limits.
:::

CPU limits work differently. They set an upper bound on the CPU time that each container can use. Kubernetes leverages the bandwidth control feature of the Completely Fair Scheduler (CFS) to implement CPU limits. CFS bandwidth control uses
time periods to limit CPU consumption. Each container gets a quota within a configurable period. The quota determines how much CPU time can be consumed in every period. If the container exhausts the quota, the container is throttled for the rest of the period.

### Best Practices for Kubernetes Limits and Requests

![](./img/best_practices_tim_hockin.png#center)

:::info
The development of the Linux cgroup subsystem started in 2006 at Google, led primarly by Rohit Seth and Paul Menage. The cgroup functionality was merged into the Linux kernel mainline in kernel version 2.6.24, which was released in January 2008. Afterwards this is called cgroups version 1. 
:::

:::info
The development and maintenance of cgroups was then taken by Tejun Heo who redesigned and rewrote cgroups. This rewrite is now called version 2.
:::

:::info
Cgroup v2 focuses on simplicity: `/sys/fs/cgroup/x/foo` (x is the controller for example cpu or memory) in v1 are now unified as `/sys/fs/cgroup/foo` , and a process can no longer join different groups for different controllers. If the process joins foo ( `/sys/fs/cgroup/foo` ), all controllers enabled for `foo` will take the control of the process.
:::

## CFS Quota? Period?

By default, Kubernetes sets the period to 100 ms. A container with a limit of 0.5 CPUs gets 50 ms of CPU time every 100 ms, as depicted in Figure 12-4. A container with a limit of 3 CPUs gets 300 ms of CPU time in every 100 millisecond period, effectively allowing the container to consume up to 3 CPUs every 100 ms.

Due to the nature of CPU limits, they can sometimes result in surprising behavior or unexpected throttling. This is usually the case in multithreaded applications that can consume the entire quota at the very beginning of the period. For example, a container with a limit of 1 CPU will get 100 ms of CPU time every 100 ms. Assuming the container has 5 threads using CPU, the container consumes the 100 ms quota in 20 ms and gets throttled for the remaining 80 ms. 

Enforcing CPU limits is useful to minimize the variability of an application’s performance, especially when running multiple replicas across different nodes. This variability in performance stems from the fact that, without CPU limits, replicas can burst and consume idle CPU cycles, which might be available at different times. By setting the CPU limits equal to the CPU requests, you remove the variability as the workloads get precisely the CPU they requested. (Google and IBM published an excellent whitepaper that discusses CFS bandwidth control in more detail.) In a similar vein,
CPU limits play a critical role in performance testing and benchmarking. Without any CPU limits, your benchmarks will produce inconclusive results, as the CPU available to your workloads will vary based on the nodes where they got scheduled
upper bound on CPU cycles is not necessary. When the CPU resources on a node are under contention, the CPU shares mechanism ensures that workloads get their fair share of CPU time, according to their container’s CPU requests. When the CPU is not under contention, the idle CPU cycles are not wasted as workloads opportunisti‐
cally consume them. 

:::info

Another issue with CPU limits is a Linux kernel bug that throttles containers unnecessarily. This has a significant impact on latency-sensitive workloads, such as web services. To avoid this issue, Kubernetes users resorted to different workarounds, including:
+ Removing CPU limits from Pod specifications
+ Disabling enforcement of CPU limits by setting the kubelet flag `--cpu-cfsquota=false`
+ Reducing the CFS period to 5–10ms by setting the kubelet flag `--cpu-cfsquota-period`

Depending on your Linux kernel version, you might not have to implement these workarounds, as the bug has been fixed in version 5.4 of the Linux kernel and backported to versions 4.14.154+, 4.19.84+, and 5.3.9+. If you need to enforce CPU limits, consider upgrading your Linux kernel version to avoid this bug.
:::

+ CFS = "Completely Fair" Scheduler A process scheduler in Linux
+ Container isolation is based on cgroups (a Linux kernel functionality) resource limitation
+ Cgroups uses CFS to implement CPU resource restriction
+ CFS scheduling is based on processing time but not core. 
+ Scheduling period is every 100ms

## CRI Runtime vs OCI Runtime

How container runtime works on Kubernetes

## References

[1] &nbsp; [Cgroup v2 fair share scheduling](https://utcc.utoronto.ca/~cks/space/blog/linux/CgroupV2FairShareScheduling)

[2] &nbsp; [QoS memory resources](https://kubernetes.io/blog/2021/11/26/qos-memory-resources/)

[3] &nbsp; [[PATCH 1/2] sched: Misc preps for cgroup unified hierarchy interface](https://lore.kernel.org/lkml/20160812221742.GA24736@cmpxchg.org/T/)

[4] &nbsp; [Kubernetes memory limits](https://home.robusta.dev/blog/kubernetes-memory-limit/)

[5] &nbsp; [Stop using CPU limits](https://home.robusta.dev/blog/stop-using-cpu-limits/)

[6] &nbsp; [An introduction to control groups (cgroups) v2](https://man7.org/conf/ndctechtown2021/cgroups-v2-part-1-intro-NDC-TechTown-2021-Kerrisk.pdf)

[7] &nbsp; [Linux Scheduler - CFS and Virtual Runtime](https://oakbytes.wordpress.com/2012/07/03/linux-scheduler-cfs-and-virtual-run-time/)

[8] &nbsp; [Setting CPU memory requests and limits](https://learnk8s.io/setting-cpu-memory-limits-requests)

[9] &nbsp; [Resource Requests and Limits Under the Hood](https://kccnceu2021.sched.com/event/iE2K/resource-requests-and-limits-under-the-hood-the-journey-of-a-pod-spec-kohei-ota-hewlett-packard-enterprise-kaslin-fields-google)

[10] &nbsp; [LISA21 - 5 Years of Cgroup v2: The Future of Linux Resource Control](https://www.youtube.com/watch?v=kPMZYoRxtmg)