### BUILDER
FROM golang:1.13-alpine as builder

WORKDIR /go/src/cf-terraforming

COPY . .

ENV CGO_ENABLED=0
ENV GOOS=linux
RUN go get -u ./...

### IMAGE
FROM alpine
COPY --from=builder /go/bin/cf-terraforming /bin/
ENTRYPOINT ["/bin/cf-terraforming"]
