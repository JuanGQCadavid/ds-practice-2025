FROM golang:alpine3.21 as builder

ENV GOOS="linux"
ENV CGO_ENABLED="0"
ARG GOARCH="arm64"

ARG CMD="grpc"

WORKDIR /app

COPY cmd/${CMD} cmd/${CMD}
COPY internal internal

COPY go.mod go.mod
COPY go.sum go.sum

RUN go build -o main ./cmd/${CMD}/main.go

FROM alpine:3.21 as prod
EXPOSE 8000
COPY --from=builder /app/main /bin/
ENTRYPOINT  ["/bin/main"]