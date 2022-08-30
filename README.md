# feelty

`feeltty` quantifies the latency experience of a TTY.

# Install

```shell script
go install github.com/coder/feeltty@master
```

# Basic Usage

```shell script
$ feeltty ssh coder.c
connect              1278.56467ms
input sample size    32
input mean           98.824ms
input stddev         152.745ms
```

or local:

```shell script
$ feeltty bash
connect              41.907167ms
input sample size    32
input mean           0s
input stddev         0s
```
