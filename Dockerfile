FROM golang:1.16 as golang

ADD . /work
WORKDIR /work

RUN CGO_ENABLED=0 GOOS=linux go build -a -tags netgo -ldflags '-w' -o cardano-toolkit


FROM node as node

ADD ui /work
WORKDIR /work

RUN yarn install && \
  yarn local:clean && \
  yarn local:build


FROM inputoutput/cardano-node:1.29.0-rc3

EXPOSE 3200

ENV PORT                     3200
ENV ASSETS                   /opt/cardano-toolkit/assets
ENV CARDANO_CLI              /usr/local/bin/cardano-cli
ENV CARDANO_NODE_SOCKET_PATH /ipc/node.sock

COPY --from=golang /work/cardano-toolkit /opt/cardano-toolkit/bin/cardano-toolkit
COPY --from=node   /work/dist            /opt/cardano-toolkit/assets

ENTRYPOINT [ "/opt/cardano-toolkit/bin/cardano-toolkit" ]
