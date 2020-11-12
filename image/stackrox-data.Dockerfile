ARG DOCS_VERSION
FROM stackrox/docs:embed-$DOCS_VERSION AS docs

FROM alpine:3.11
ARG ALPINE_MIRROR=sjc.edge.kernel.org

RUN mkdir /stackrox-data

RUN echo http://$ALPINE_MIRROR/alpine/v3.11/main > /etc/apk/repositories; \
    echo http://$ALPINE_MIRROR/alpine/v3.11/community >> /etc/apk/repositories

RUN apk update && \
    apk add --no-cache \
        openssl zip \
        && \
    apk --purge del apk-tools \
    ;

COPY --from=docs /docs/public /stackrox-data/product-docs
# Basic sanity check: are the docs in the right place?
RUN ls /stackrox-data/product-docs/index.html

RUN mkdir -p /stackrox-data/cve/k8s && \
    wget -O /stackrox-data/cve/k8s/checksum "https://definitions.stackrox.io/cve/k8s/checksum" && \
    wget -O /stackrox-data/cve/k8s/cve-list.json "https://definitions.stackrox.io/cve/k8s/cve-list.json" && \
    mkdir -p /stackrox-data/cve/istio && \
    wget -O /stackrox-data/cve/istio/checksum "https://definitions.stackrox.io/cve/istio/checksum" && \
    wget -O /stackrox-data/cve/istio/cve-list.json "https://definitions.stackrox.io/cve/istio/cve-list.json"

RUN mkdir -p /stackrox-data/external-networks && \
    wget -O /stackrox-data/external-networks/checksum "https://definitions.stackrox.io/external-networks/latest/checksum" && \
    wget -O /stackrox-data/external-networks/networks "https://definitions.stackrox.io/external-networks/latest/networks"

RUN zip -jr /stackrox-data/external-networks/external-networks.zip /stackrox-data/external-networks

COPY ./policies/files /stackrox-data/policies/files
COPY ./docs/api/v1/swagger.json /stackrox-data/docs/api/v1/swagger.json

COPY ./keys /tmp/keys

RUN set -eo pipefail; \
	( cd /stackrox-data ; tar -czf - * ; ) | \
    openssl enc -aes-256-cbc \
        -K "$(hexdump -e '32/1 "%02x"' </tmp/keys/data-key)" \
        -iv "$(hexdump -e '16/1 "%02x"' </tmp/keys/data-iv)" \
        -out /stackrox-data.tgze
