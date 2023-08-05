FROM golang:1.19.4-buster as build

ENV BIN_FILE "/opt/orderservice/orderservice"
ENV CODE_DIR "/go/src/app"

COPY . ${CODE_DIR}
WORKDIR ${CODE_DIR}

ARG LDFLAGS
RUN CGO_ENABLED=0 go build -ldflags "$LDFLAGS" -o ${BIN_FILE} main.go

FROM scratch

LABEL SERVICE="orderservice"
LABEL MAINTAINERS="astrviktor@gmail.com"

COPY --from=build /opt/orderservice/orderservice /srv/orderservice

WORKDIR /srv
ENTRYPOINT ["/srv/orderservice"]





