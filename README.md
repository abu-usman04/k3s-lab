# K3s Lab

Local single-node Kubernetes cluster (K3s) with the Headlamp dashboard, managed entirely through Docker Compose. No host-side kubectl required.

## Requirements

- Docker with Compose v2

## Quick Start

```bash
docker compose up -d
```

The first start takes 1–2 minutes while images are pulled.

Verify:

```bash
docker exec headlamp-server kubectl get nodes
docker exec headlamp-server kubectl get pods -A
```

## Dashboard Access

1. Generate a login token:

   ```bash
   docker exec headlamp-server kubectl create token headlamp-admin -n kube-system --duration 24h
   ```

2. Open http://localhost:8080 and paste the token.

Tokens expire after 24 hours. Rerun the command to generate a new one.

## Components

| Component  | Location                | Purpose                                  |
|------------|-------------------------|------------------------------------------|
| k3s-server | Docker container        | Kubernetes control plane and node        |
| headlamp   | Pod in `kube-system`    | Web dashboard, exposed on port 8080      |
| cpu-burner | Deployment in `demo`    | Sample CPU-bound workload for monitoring |

## Working with the Sample Workload

```bash
docker exec headlamp-server kubectl top pods -n demo
docker exec headlamp-server kubectl scale deployment cpu-burner -n demo --replicas=5
docker exec headlamp-server kubectl delete pod -n demo -l app=cpu-burner
```

Each `cpu-burner` pod is limited to 200m CPU. Deleted pods are recreated automatically by the Deployment.

## Lifecycle

```bash
docker compose down       # Stop; cluster state is preserved
docker compose up -d      # Start again with existing state
docker compose down -v    # Stop and delete all cluster data
```

Cluster state persists in the named Docker volume `k3s-data`.
