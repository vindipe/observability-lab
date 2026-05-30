# Snap example

This directory contains a conceptual `snapcraft.yaml` for packaging the Go app as a snap.

This is not part of the default runnable demo because the main target is a containerized/Kubernetes observability stack. It is included so you can talk about Linux packaging concepts in the technical review.

Build idea:

```bash
cd snap
snapcraft
```

In a real project, the source path and build step would be adjusted depending on whether the snap builds from the local repository, a Git tag, or a released tarball.
