# Snap packaging

This lab uses containers for the runnable demo because the target architecture is cloud-native. Snap packaging is included as a conceptual Cloud-Native extension.

## What Snap is

A snap is a Linux application package with metadata, confinement, interfaces, channels and automatic refresh mechanics. It is useful for distributing applications to Ubuntu machines, IoT devices and workstations.

## Snap vs container image

| Snap | Container image |
|---|---|
| Packages an app for Linux systems | Packages a filesystem/process for container runtimes |
| Uses snapd, channels and confinement | Uses Docker/containerd/Kubernetes |
| Good for CLIs, agents, desktop apps, IoT | Good for cloud-native services |
| Runs as host-managed app/service | Runs as container workload |

## Example

See `snap/snapcraft.yaml` for a minimal conceptual package of the Go app.

## Technical Review framing

> I would not force Snap into the main demo because the workload is intentionally containerized. But Snap is relevant in the cloud-native platform ecosystem for distributing agents, CLIs or host-level services, whereas containers are natural for this Kubernetes observability lab.
