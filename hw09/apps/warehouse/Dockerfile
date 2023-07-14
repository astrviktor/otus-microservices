FROM golang:1.19.4-buster as build

ENV BIN_FILE "/opt/warehouse/warehouse"
ENV CODE_DIR "/go/src/app"

COPY . ${CODE_DIR}
WORKDIR ${CODE_DIR}

ARG LDFLAGS
RUN CGO_ENABLED=0 go build -ldflags "$LDFLAGS" -o ${BIN_FILE} main.go

FROM scratch

LABEL SERVICE="warehouse"
LABEL MAINTAINERS="astrviktor@gmail.com"

COPY --from=build /opt/warehouse/warehouse /srv/warehouse

WORKDIR /srv
ENTRYPOINT ["/srv/warehouse"]





