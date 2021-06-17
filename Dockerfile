FROM golang:1.16.5-alpine3.13 as builder
RUN apk add --no-cache git openssh gcc musl-dev
WORKDIR /src
COPY go.mod .
COPY go.sum .

# Get dependancies - will also be cached if we won't change mod/sum
RUN go mod download

# Build
COPY . .
RUN go build -o server .

FROM philipslabs/siderite:v0.7.1 AS siderite

FROM alpine:latest
LABEL maintainer="andy.lo-a-foe@philips.com"
RUN apk add --no-cache git openssh openssl bash postgresql-client
WORKDIR /app
COPY --from=siderite /app/siderite /app/siderite
COPY --from=builder /src/server /app

CMD ["/app/siderite","function"]
