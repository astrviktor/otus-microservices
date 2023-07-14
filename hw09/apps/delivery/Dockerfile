FROM golang:1.19.4-buster as build

ENV BIN_FILE "/opt/delivery/delivery"
ENV CODE_DIR "/go/src/app"

COPY . ${CODE_DIR}
WORKDIR ${CODE_DIR}

ARG LDFLAGS
RUN CGO_ENABLED=0 go build -ldflags "$LDFLAGS" -o ${BIN_FILE} main.go

FROM scratch

LABEL SERVICE="delivery"
LABEL MAINTAINERS="astrviktor@gmail.com"

COPY --from=build /opt/delivery/delivery /srv/delivery

WORKDIR /srv
ENTRYPOINT ["/srv/delivery"]





