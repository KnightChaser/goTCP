FROM golang:alpine AS gobuilder
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64
WORKDIR /build
COPY go.mod go.sum tcpserver.go ./
RUN go mod download
RUN go build -o tcpserver tcpserver.go

FROM scratch
COPY --from=gobuilder /build .
ENTRYPOINT [ "tcpserver" ]