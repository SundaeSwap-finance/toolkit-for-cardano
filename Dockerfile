FROM golang:1.16 as golang

ADD . /work
WORKDIR /work

RUN CGO_ENABLED=0 GOOS=linux go build -a -tags netgo -ldflags '-w' -o cardano-toolkit


FROM scratch

ENV PORT 80
EXPOSE 80

COPY --from=golang /work/cardano-toolkit /opt/cardano-toolkit/bin/cardano-toolkit

ENTRYPOINT [ "/opt/cardano-toolkit/bin/cardano-toolkit" ]
