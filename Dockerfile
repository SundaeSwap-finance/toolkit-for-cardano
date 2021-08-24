FROM golang:1.16 as golang

ADD . /work
WORKDIR /work

RUN CGO_ENABLED=0 GOOS=linux go build -a -tags netgo -ldflags '-w' -o cardano-toolkit


FROM node as node

ADD ui /work
WORKDIR /work

RUN yarn install && yarn local:build


FROM scratch

EXPOSE 8888

ENV PORT   8888
ENV ASSETS /opt/cardano-toolkit/assets

COPY --from=golang /work/cardano-toolkit /opt/cardano-toolkit/bin/cardano-toolkit
COPY --from=node   /work/dist            /opt/cardano-toolkit/assets

ENTRYPOINT [ "/opt/cardano-toolkit/bin/cardano-toolkit" ]
