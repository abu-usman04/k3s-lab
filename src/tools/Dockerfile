FROM alpine:3.22

ARG HELM_VERSION=v3.21.4
ARG KUBECTL_VERSION=v1.35.8
ARG CILIUM_CLI_VERSION=v0.18.5
ARG HUBBLE_CLI_VERSION=v1.20.2

RUN apk add --no-cache bash ca-certificates curl jq tar

RUN set -eux; \
    curl -fsSL "https://get.helm.sh/helm-${HELM_VERSION}-linux-amd64.tar.gz" -o /tmp/helm.tgz; \
    tar -xzf /tmp/helm.tgz -C /tmp; \
    install /tmp/linux-amd64/helm /usr/local/bin/helm; \
    curl -fsSL "https://dl.k8s.io/release/${KUBECTL_VERSION}/bin/linux/amd64/kubectl" -o /usr/local/bin/kubectl; \
    chmod +x /usr/local/bin/kubectl; \
    curl -fsSL "https://github.com/cilium/cilium-cli/releases/download/${CILIUM_CLI_VERSION}/cilium-linux-amd64.tar.gz" | tar -xz -C /usr/local/bin; \
    curl -fsSL "https://github.com/cilium/hubble/releases/download/${HUBBLE_CLI_VERSION}/hubble-linux-amd64.tar.gz" | tar -xz -C /usr/local/bin; \
    rm -rf /tmp/*

WORKDIR /workspace
ENTRYPOINT ["sleep", "infinity"]
