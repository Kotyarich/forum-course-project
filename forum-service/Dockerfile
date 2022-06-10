FROM golang:1.16.10-alpine as builder
COPY go.mod go.sum /go/src/
WORKDIR /go/src/
RUN go mod download
COPY . /go/src/
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o build/ ./...

FROM alpine

COPY --from=builder /go/src/build /usr/bin/forum-service
EXPOSE 5000 5000

ENTRYPOINT ["/usr/bin/forum-service/api"]
