# Minimal K3s, Cilium, and Hubble playground

This repository runs a single-node K3s cluster inside Docker Compose with Cilium and Hubble. It is a small learning environment for observing service traffic and trying the Hubble Observer API.

## Start and verify

```sh
cd src
docker compose up -d --build
docker compose exec tools bash /workspace/scripts/verify.sh
```

The generated kubeconfig is written to `src/.kube` and is not committed. Cluster state is kept in Docker volumes. Do not use `docker compose down -v` unless the cluster should be reset.

## Minimal demo

```sh
docker compose exec tools bash /workspace/workloads/topology-demo/apply.sh
docker compose exec tools hubble observe -P --namespace demo
docker compose exec tools hubble observe -P --namespace demo --verdict DROPPED
```

The frontend calls one API route every five seconds. The API proxies that route to the backend's one response endpoint.

## Hubble access

The Relay and UI are internal services. Use local forwarding when needed:

```sh
docker compose exec tools cilium hubble port-forward
kubectl -n kube-system port-forward service/hubble-ui 8080:80
```

The small Go experiment is under `src/experiments/hubble-client` and covers `ServerStatus`, `GetNodes`, `GetNamespaces`, and bounded `GetFlows` calls.
