FROM golang:1.20

WORKDIR /go/src/github.com/ThreeDotsLabs/monolith-microservice-shop
COPY . .

RUN go mod download
RUN go install github.com/cespare/reflex@latest
