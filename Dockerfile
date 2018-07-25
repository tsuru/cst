FROM golang:1.10-alpine3.8 AS build

ARG repository="github.com/tsuru/cst"

COPY ./ "${GOPATH}/src/${repository}/"

RUN apk add --update make && \
    CSTBIN=/tmp/cst make -C "${GOPATH}/src/${repository}" build

FROM alpine:3.8

COPY --from=build /tmp/cst /usr/local/bin/cst

USER nobody

EXPOSE 8443

ENTRYPOINT ["cst"]
