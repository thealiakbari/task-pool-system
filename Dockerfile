FROM golang:1.25.1 as builder
WORKDIR /go/app
COPY . .
RUN go mod download
RUN go build -o ./src/build ./cmd/executor

FROM golang:1.25.1
WORKDIR /root/
COPY config/config.yml ./config/config.yml
COPY ./cmd/migration/scripts ./cmd/migration/scripts
COPY --from=builder /go/app/src/build .
EXPOSE 1212
ENTRYPOINT ["./build"]
