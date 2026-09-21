#!/usr/bin/env bash
set -eu

kubectl apply -f /workspace/workloads/topology-demo/namespace.yaml
kubectl apply -f /workspace/workloads/topology-demo/backend.yaml
kubectl apply -f /workspace/workloads/topology-demo/api.yaml
kubectl apply -f /workspace/workloads/topology-demo/frontend.yaml
kubectl -n demo rollout status deployment/backend --timeout=2m
kubectl -n demo rollout status deployment/api --timeout=2m
kubectl -n demo rollout status deployment/frontend --timeout=2m
