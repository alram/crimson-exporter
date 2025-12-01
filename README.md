# No support

I'm not in a postion where I'm able to maintain or debug issues anymore, feel free to fork if you've a use for this project

# crimson-exporter

Quick and dirty exporter that runs dump_metrics on all OSD asok to gather metrics for crimson-osd testing.
It doesn't record any cache metric cause that's a lot but if some are useful create an issue or preferably PR to add them.

This isn't meant for any prod use just to gather metrics while doing perf testing and I'll most definitely not have the time to maintain it.

It exports metrics on http://0.0.0.0:9092/metrics

# how to use

```
$ go build
$ ./crimson-exporterd
```

# TODO

- Add latency bucket metrics
- Docker build
- Cleanup code
- env vars for port, ceph bins, ceph dirs, etc.
- There almost definetely bugs since there's a lot of data, there could be some metrics overwritting md particularly for the ones that have more labels than osd, shard
- Slow even with a thread per OSD (querying dump_metrics itself take >1s), on a host with 24 OSDs:
```
real	0m4.979s
user	0m0.005s
sys	0m0.010s
```
