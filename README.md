# K3s Playground

A small local Kubernetes playground for testing Kubernetes networking, Cilium, Hubble, and related infrastructure tools.

Everything runs locally with Docker Compose.

## What's Included

- **K3s** — lightweight Kubernetes cluster
- **Cilium** — Kubernetes networking and network policies
- **Hubble** — network flow observability
- **Helm** — package management for Kubernetes
- **Headlamp** — optional Kubernetes dashboard

## Start

```bash
docker compose up -d --build
```

## Check the cluster:

```bash
docker compose exec k3s kubectl get nodes
docker compose exec k3s kubectl get pods -A
```

## Stop or completely remove

```bash
docker compose down
docker compose down -v
```
