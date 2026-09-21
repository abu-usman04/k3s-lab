#!/usr/bin/env bash
set -eu

KUBECONFIG="${KUBECONFIG:-/root/.kube/kubeconfig.yaml}"
export KUBECONFIG

test "$(kubectl get nodes --no-headers | wc -l | tr -d ' ')" = 1
test "$(kubectl get nodes -o jsonpath='{.items[0].metadata.name}')" = k3s
kubectl wait --for=condition=Ready node/k3s --timeout=2m
kubectl -n kube-system rollout status daemonset/cilium --timeout=2m
kubectl -n kube-system rollout status daemonset/cilium-envoy --timeout=2m
kubectl -n kube-system rollout status deployment/cilium-operator --timeout=2m
kubectl -n kube-system rollout status deployment/coredns --timeout=2m
kubectl -n kube-system rollout status deployment/hubble-relay --timeout=2m
kubectl -n kube-system rollout status deployment/hubble-ui --timeout=2m
cilium status --wait --wait-duration 2m
hubble status -P -n kube-system
