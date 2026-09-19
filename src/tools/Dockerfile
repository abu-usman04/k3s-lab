FROM alpine:3.22

RUN apk add --no-cache \
    bash \
    curl \
    ca-certificates \
    tar

ARG HELM_VERSION=v3.21.4

RUN curl -fsSL \
    https://get.helm.sh/helm-${HELM_VERSION}-linux-amd64.tar.gz \
    -o /tmp/helm.tar.gz \
    && tar -xzf /tmp/helm.tar.gz -C /tmp \
    && mv /tmp/linux-amd64/helm /usr/local/bin/helm \
    && chmod +x /usr/local/bin/helm \
    && rm -rf /tmp/*

ENTRYPOINT ["sleep", "infinity"]