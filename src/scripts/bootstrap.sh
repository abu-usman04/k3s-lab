#!/usr/bin/env bash
set -eu

KUBECONFIG="${KUBECONFIG:-/root/.kube/kubeconfig.yaml}"
export KUBECONFIG

until test -s "$KUBECONFIG"; do sleep 1; done
until curl -kfsS https://127.0.0.1:6443/readyz >/dev/null; do sleep 2; done

helm repo add cilium https://helm.cilium.io >/dev/null
helm repo update >/dev/null
helm upgrade --install cilium cilium/cilium \
  --version 1.20.2 \
  --namespace kube-system \
  --values /workspace/cilium/values.yaml \
  --wait \
  --timeout 10m

kubectl -n kube-system rollout status daemonset/cilium --timeout=10m
kubectl -n kube-system rollout status daemonset/cilium-envoy --timeout=10m
kubectl -n kube-system rollout status deployment/cilium-operator --timeout=10m
kubectl -n kube-system rollout status deployment/hubble-relay --timeout=10m
kubectl -n kube-system rollout status deployment/hubble-ui --timeout=10m
kubectl -n kube-system rollout status deployment/coredns --timeout=10m
cilium status --wait --wait-duration 10m
hubble status -P -n kube-system
