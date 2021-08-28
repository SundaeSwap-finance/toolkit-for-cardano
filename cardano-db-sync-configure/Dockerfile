FROM ubuntu:20.04

RUN apt-get update && apt-get install -y vim jq curl wget
RUN curl -L https://github.com/mikefarah/yq/releases/download/v4.12.1/yq_linux_amd64 -o /usr/local/bin/yq && \
  chmod +x /usr/local/bin/yq

ADD configure.sh                /usr/sbin/configure.sh
ADD cardano-db-sync-config.json /etc/cardano-db-sync/cardano-db-sync-config.json

ENTRYPOINT ["/usr/sbin/configure.sh"]
