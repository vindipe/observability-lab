# Optional eBPF / Parca extension

This lab does not start Parca by default because eBPF-based profilers often require host privileges, kernel support and environment-specific configuration. That can be fragile in WSL and Docker Desktop.

## Why eBPF is useful

OpenTelemetry traces answer: "Which request was slow and which spans did it cross?"

Continuous profiling answers: "Where did the process spend CPU time over time?"

With eBPF, profilers can sample stack traces from user space and kernel space with low overhead, often without changing application code.

## How Parca would fit

```text
Go app process -> Parca Agent/eBPF -> Parca Server -> Grafana/profile UI
```

Possible production metrics/questions:

- Which Go functions dominate CPU time?
- Did an upgrade increase CPU cost?
- Is JSON encoding, Redis client code, locking or GC dominating runtime?
- Does a latency spike correlate with CPU saturation?

## Why optional here

- WSL may not expose all kernel features needed for eBPF profiling.
- Docker Desktop may require privileged containers and extra capabilities.
- A brittle privileged setup is worse than a stable technical review demo.

Technical Review framing:

> I know where profiling fits: it complements metrics, logs and traces. I did not enable Parca by default because local eBPF support can be environment-dependent, but I documented how it would extend the stack.
