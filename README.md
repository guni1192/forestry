# forestry

![Unit Test](https://github.com/guni1192/forestry/workflows/Unit%20Test/badge.svg)
![Lint](https://github.com/guni1192/forestry/workflows/Lint/badge.svg)

Distributed Logging Platform

**Note**: This repository is created for demonstration.
We not assumption for production environment.

![](figs/architecture.png)

## Getting Started (local)

**Requirements**

- docker
- helm
- minikube(docker driver)
- kubectl

```console
minikube start
eval $(minikube docker-env)
kubectl get node -o wide
```

### Setup Grafana Loki

```console
helm repo add grafana https://grafana.github.io/helm-charts
helm upgrade --install loki grafana/loki-stack --set grafana.enabled=true
kubectl get secret loki-grafana -o jsonpath="{.data.admin-password}" | base64 --decode ; echo
```

### Build and Deploy forestry

```console
make docker-build
make deploy
```

### Port-Forward

```console
kubectl port-forward service/loki-grafana 3000:80 --address 0.0.0.0
```

Grafana
- username: admin
- password: `kubectl get secret loki-grafana ...`
