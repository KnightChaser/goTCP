FROM golang:1.20-buster AS gobuilder
ENV GO111MODULE=on \
    CGO_ENABLED=0
WORKDIR /build
COPY go.mod tcpclient.go ./
RUN go mod download
RUN go build -o tcpclient tcpclient.go

FROM scratch
COPY --from=gobuilder /build .
ARG USERNAME
ENV USERNAME=$USERNAME
ENTRYPOINT [ "./tcpclient", "-username", "$USERNAME" ]
