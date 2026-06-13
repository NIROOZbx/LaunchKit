FROM golang:1.26.4-alpine
WORKDIR /launchkit

RUN go install github.com/air-verse/air@latest

COPY shared/ ./shared/
COPY core/go.mod core/go.sum ./core/

WORKDIR /launchkit/core
RUN go mod download

CMD ["air", "-c", ".air.toml"]