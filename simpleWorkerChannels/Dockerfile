FROM golang:buster AS builder
RUN apt-get install git
WORKDIR /go/src/github.com/dalmarcogd/test-go-worker-pool/simpleWorkerChannels
ADD ./go* ./
RUN go mod download
ADD . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o app .

FROM alpine
WORKDIR /app
COPY --from=builder /go/src/github.com/dalmarcogd/test-go-worker-pool/simpleWorkerChannels/app /app/
EXPOSE 8001
ENTRYPOINT ./app