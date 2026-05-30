# Cassandra as an advanced stateful scenario

The runnable demo uses Redis to keep local startup fast and stable.

Cassandra would be valuable in an advanced profile because it is a distributed, stateful database with rich operational signals.

## What the app would do

- connect to Cassandra with a Go driver;
- create/use an `orders` keyspace/table;
- execute a query in `/api/slow`;
- record query latency as a histogram;
- record connection/query errors;
- create a trace span around each Cassandra operation;
- log query failures with trace IDs.

## Metrics to observe

- read latency;
- write latency;
- coordinator latency;
- compaction throughput/backlog;
- dropped messages;
- pending tasks;
- disk usage;
- node availability;
- JVM heap and GC;
- client-side timeout/error rate.

## Why not default?

Cassandra is heavier than Redis and usually needs more memory, more startup time and more careful tuning. For an technical review demo, a stable base stack is more valuable than a fragile heavyweight dependency.
