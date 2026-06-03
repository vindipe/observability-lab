# OpenStack context

This repository does not deploy OpenStack locally. That would make the lab too heavy and distract from the observability goal.

## What OpenStack is

OpenStack is an IaaS layer. It manages compute, storage and networking resources and exposes cloud-like APIs for virtual machines, volumes, images, security groups and networks.

## OpenStack vs Kubernetes

| OpenStack | Kubernetes |
|---|---|
| Infrastructure-as-a-Service | Container orchestration platform |
| VMs, networks, volumes, images | Pods, Deployments, Services, ConfigMaps |
| Nova, Neutron, Cinder, Glance, Keystone | kube-apiserver, scheduler, kubelet, controllers |
| Often runs below Kubernetes | Often runs workloads above IaaS/bare metal |

A realistic enterprise architecture can be:

```text
bare metal -> OpenStack -> Kubernetes -> applications -> observability stack
```

## Observability of OpenStack

Signals to care about:

- Nova scheduler/API latency and errors;
- Neutron agent health and network errors;
- Cinder volume attach/detach failures;
- Keystone authentication latency/errors;
- RabbitMQ queue depth if used by control plane services;
- database health for control plane state;
- hypervisor resource saturation;
- VM lifecycle events;
- storage and network throughput.