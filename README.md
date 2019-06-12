# kubectl dfi

[![Build Status](https://travis-ci.org/makocchi-git/kubectl-dfi.svg?branch=master)](https://travis-ci.org/makocchi-git/kubectl-dfi)
[![Maintainability](https://api.codeclimate.com/v1/badges/b92591d00becc95b11ca/maintainability)](https://codeclimate.com/github/makocchi-git/kubectl-dfi/maintainability)
[![Go Report Card](https://goreportcard.com/badge/github.com/makocchi-git/kubectl-dfi)](https://goreportcard.com/report/github.com/makocchi-git/kubectl-dfi)
[![kubectl plugin](https://img.shields.io/badge/kubectl-plugin-blue.svg)](https://github.com/topics/kubectl-plugin)
[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat)](LICENSE)

Print disk usage of container image on Kubernetes node(s) like a linux "df" command.  

```shell
$ kubectl dfi
NAME                               IMAGE USED   ALLOCATABLE   CAPACITY     %USED
node1-default-pool-500decb4-5q58   1982531K     47093746K     101241290K   1%
node2-default-pool-500decb4-7wpk   1891326K     47093746K     101241290K   1%
node3-default-pool-500decb4-9dd4   1982531K     47093746K     101241290K   1%
```

And list images on Kubernetes node(s).

```shell
$ kubectl dfi --list node1-default-pool-500decb4-5q58
node1-default-pool-500decb4-5q58   286572K      k8s.gcr.io/node-problem-detector:v0.4.1
node1-default-pool-500decb4-5q58   223242K      gcr.io/stackdriver-agents/stackdriver-logging-agent:0.6-1.6.0-1
node1-default-pool-500decb4-5q58   135716K      k8s.gcr.io/fluentd-elasticsearch:v2.0.4
node1-default-pool-500decb4-5q58   103488K      k8s.gcr.io/fluentd-gcp-scaler:0.5
node1-default-pool-500decb4-5q58   102992K      k8s.gcr.io/kube-proxy:v1.11.8-gke.6
node1-default-pool-500decb4-5q58   102319K      k8s.gcr.io/kubernetes-dashboard-amd64:v1.8.3
...
```

## Install

```shell
$ make
$ mv _output/kubectl-dfi /usr/local/bin/.

# Happy dfi time!
$ kubectl dfi
```

## Usage

```shell
# Show image usage of Kubernetes nodes.
kubectl dfi

# Using label selector.
kubectl dfi -l key=value

# Use image count with image disk usage.
kubectl dfi --count

# Print raw(bytes) usage.
kubectl dfi --bytes --without-unit

# Using binary prefix unit (GiB, MiB, etc)
kubectl dfi -g -B

# List images on nodes.
kubectl dfi --list
```

## Notice

`IMAGE USED` is simply sum up of container image size reported by kubelet.  
In fact, node disk might be not used so much by container images because of cache by layered filesystem.

## License

This software is released under the MIT License.
