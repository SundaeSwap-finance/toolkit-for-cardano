FROM golang:1.16 as golang

ADD . /work
WORKDIR /work

RUN CGO_ENABLED=0 GOOS=linux go build -a -tags netgo -ldflags '-w' -o toolkit-for-cardano


FROM node as node

ADD ui /work
WORKDIR /work

RUN yarn install && \
  yarn local:clean && \
  yarn local:build


FROM inputoutput/cardano-node:1.29.0

EXPOSE 3200

ENV PORT                     3200
ENV ASSETS                   /opt/toolkit-for-cardano/assets
ENV CARDANO_CLI              /usr/local/bin/cardano-cli
ENV CARDANO_NODE_SOCKET_PATH /ipc/node.sock

COPY --from=golang /work/toolkit-for-cardano /opt/toolkit-for-cardano/bin/toolkit-for-cardano
COPY --from=node   /work/dist                /opt/toolkit-for-cardano/assets

ENTRYPOINT [ "/opt/toolkit-for-cardano/bin/toolkit-for-cardano" ]
