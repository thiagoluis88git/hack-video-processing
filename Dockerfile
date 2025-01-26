FROM golang:1.22 AS build-stage

WORKDIR /go/src
ENV PATH="/go/src:${PATH}"

# Install Certificate
RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates

COPY . ./

RUN go mod download
RUN go mod tidy

ENV CGO_ENABLED 1
ENV GOOS=linux

RUN \
  --mount=target=. \
  --mount=target=/root/.cache,type=cache \
  CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
  go build \
  -ldflags "-s -d -w" \
  -o /VideoProcessing cmd/api/main.go

FROM scratch

WORKDIR /app

# Copy Certificate
COPY --from=build-stage /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=build-stage /VideoProcessing /VideoProcessing
# COPY --from=build-stage /go/src/docs/ /docs/

EXPOSE 3210 3211

ENTRYPOINT ["/VideoProcessing"]
