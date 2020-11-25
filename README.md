<p align="center">
    <img alt="Toad Logo" src="https://raw.githubusercontent.com/clivern/Toad/master/assets/img/gopher.png?v=0.2.4" width="150" />
    <h3 align="center">Toad</h3>
    <p align="center">Containerized Application for Testing Purposes</p>
    <p align="center">
        <a href="https://travis-ci.com/Clivern/Toad"><img src="https://travis-ci.com/Clivern/Toad.svg?branch=master"></a>
        <a href="https://github.com/Clivern/Toad/releases"><img src="https://img.shields.io/badge/Version-0.2.4-red.svg"></a>
        <a href="https://goreportcard.com/report/github.com/Clivern/Toad"><img src="https://goreportcard.com/badge/github.com/clivern/Toad?v=0.2.4"></a>
        <a href="https://hub.docker.com/r/clivern/toad"><img src="https://img.shields.io/badge/Docker-Latest-green"></a>
        <a href="https://github.com/Clivern/Toad/blob/master/LICENSE"><img src="https://img.shields.io/badge/LICENSE-MIT-orange.svg"></a>
    </p>
</p>

## Documentation

### Usage

Get [the latest binary.](https://github.com/Clivern/Toad/releases)

```zsh
$ curl -sL https://github.com/Clivern/Toad/releases/download/x.x.x/Toad_x.x.x_OS_x86_64.tar.gz | tar xz
```

Run Toad:

```zsh
$ ./Toad
```

Check the release:

```zsh
$ ./Toad --get release
```

Health check:

```zsh
$ ./Toad --get health
```

To run on stateful mode:

```zsh
# If stateful, it will use redis by default
$ export IS_STATEFUL=true

# Redis configs
$ export REDIS_HOST=localhost
$ export REDIS_PORT=6379
$ export REDIS_PASSWORD=
```

Test it:

```zsh
# Home
$ curl http://127.0.0.1:8080/ -v

# Health check
$ curl http://127.0.0.1:8080/_health -v

# Get State (Only Stateful Mode)
$ curl http://127.0.0.1:8080/do/get -v

# Change State (Only Stateful Mode)
$ curl http://127.0.0.1:8080/do/change -v

# Reset State (Only Stateful Mode)
$ curl http://127.0.0.1:8080/do/reset -v

# Enable Host (Only Stateful Mode)
$ curl http://127.0.0.1:8080/do/host_up -v

# Disable Host (Only Stateful Mode)
$ curl http://127.0.0.1:8080/do/host_down -v

# Enable All Hosts (Only Stateful Mode)
$ curl http://127.0.0.1:8080/do/all_up -v

# Disable  All Hosts (Only Stateful Mode)
$ curl http://127.0.0.1:8080/do/all_down -v
```

### Deploy on Kubernetes

```zsh
$ kubectl apply -f deployment/k8s/configs.yaml --record

$ kubectl get configmaps configs -o yaml

$ kubectl apply -f deployment/k8s/redis.yaml --record

$ kubectl apply -f deployment/k8s/toad.yaml --record

$ kubectl get deployments -o wide
$ kubectl get svc -o wide

$ kubectl get pods -o wide
$ kubectl logs $PodName
```

### Deployment Strategies

Ramped: release a new version on a rolling update fashion, one after the other.

```zsh
$ kubectl apply -f deployment/k8s/ramped/toad_ramped_strategy.yaml --record
```

Recreate: terminate the old version and release the new one.

```zsh
$ kubectl apply -f deployment/k8s/recreate/toad_recreate_strategy.yaml --record
```

Blue/Green: release a new version alongside the old version then switch traffic.

```zsh
$ kubectl apply -f deployment/k8s/blue_green/toad_blue_green_strategy.yaml --record

# Create service to load balance the old version (blue)
$ kubectl apply -f deployment/k8s/blue_green/switch_to_blue.yaml --record

# Create service to load balance the new version (green)
$ kubectl apply -f deployment/k8s/blue_green/switch_to_green.yaml --record
```

Canary: release a new version to a subset of users, then proceed to a full rollout.

```zsh
$ ~
```


### Kubernetes Ingresses

It allow you to flexibly route traffic from outside your Kubernetes cluster to Services inside of your cluster.

```zsh
$ kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v0.41.2/deploy/static/provider/cloud/deploy.yaml
$ kubectl get pods -n ingress-nginx \
  -l app.kubernetes.io/name=ingress-nginx --watch
$ kubectl get svc --namespace=ingress-nginx

$ kubectl apply -f deployment/k8s/ingress/toad.yaml --record

$ kubectl get ingress
$ kubectl describe ingress toad-ing1
$ kubectl describe ingress toad-ing2

$ curl http://example.com/toad1/
$ curl http://example.com/toad2/

$ curl http://example.com/toad1/_ready
$ curl http://example.com/toad2/_ready

$ curl http://toad1.example.com/
$ curl http://toad2.example.com/
```


## Versioning

For transparency into our release cycle and in striving to maintain backward compatibility, Toad is maintained under the [Semantic Versioning guidelines](https://semver.org/) and release process is predictable and business-friendly.

See the [Releases section of our GitHub project](https://github.com/clivern/toad/releases) for changelogs for each release version of Toad. It contains summaries of the most noteworthy changes made in each release.


## Bug tracker

If you have any suggestions, bug reports, or annoyances please report them to our issue tracker at https://github.com/clivern/toad/issues


## Security Issues

If you discover a security vulnerability within Toad, please send an email to [hello@clivern.com](mailto:hello@clivern.com)


## Contributing

We are an open source, community-driven project so please feel free to join us. see the [contributing guidelines](CONTRIBUTING.md) for more details.


## License

© 2020, clivern. Released under [MIT License](https://opensource.org/licenses/mit-license.php).

**Toad** is authored and maintained by [@clivern](http://github.com/clivern).
