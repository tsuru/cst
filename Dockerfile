FROM golang:1.10-alpine3.7

ARG repository="github.com/tsuru/cst"

COPY ./ "${GOPATH}/src/${repository}/"

RUN apk update && \
    apk upgrade && \
    apk add --virtual .build-deps git make && \
    cd "${GOPATH}/src/${repository}" && \
    CSTBIN="/usr/local/bin/cst" make get-dev-deps build && \
    rm -rf "${GOPATH}/pkg" "${GOPATH}/src" ${GOPATH}/bin/* && \
    ln -s /usr/local/bin/cst ${GOPATH}/bin/ && \
    apk del .build-deps && \
    rm /var/cache/apk/* && \
    cd /

USER nobody

ENTRYPOINT ["cst"]
