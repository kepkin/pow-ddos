FROM golang:1.18-alpine as build-stage
WORKDIR /opt
RUN apk add git

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build ./cmd/example-client
RUN CGO_ENABLED=0 GOOS=linux go build ./cmd/example-server
RUN CGO_ENABLED=0 GOOS=linux go build ./cmd/guide-server


FROM alpine:3.15

WORKDIR /opt/app
RUN apk add libc6-compat

COPY --from=build-stage /opt/example-client .
COPY --from=build-stage /opt/example-server .
COPY --from=build-stage /opt/guide-server .