FROM golang:1.21.5-alpine3.19 AS gobuilder
ENV GO111MODULE=on \
    CGO_ENABLED=0
WORKDIR /build
COPY go.mod go.sum tcpclient.go ./
RUN go mod download
RUN go build -o tcpclient tcpclient.go

FROM scratch
COPY --from=gobuilder /build .