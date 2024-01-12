FROM golang:1.20-buster AS gobuilder
ENV GO111MODULE=on \
    CGO_ENABLED=0
WORKDIR /build
COPY go.mod tcpserver.go ./
RUN go mod download
RUN go build -o tcpserver tcpserver.go

FROM scratch
COPY --from=gobuilder /build .
ENTRYPOINT [ "./tcpserver" ]