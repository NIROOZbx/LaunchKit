FROM golang:1.26.4-alpine
WORKDIR /launchkit

RUN go install github.com/air-verse/air@latest

COPY shared/ ./shared/
COPY db/ ./db/
COPY gateway/go.mod gateway/go.sum ./gateway/

WORKDIR /launchkit/gateway
RUN go mod download

CMD ["air", "-c", ".air.toml"]