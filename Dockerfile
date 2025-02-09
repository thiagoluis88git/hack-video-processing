FROM golang:1.22 AS build-stage

WORKDIR /go/src
ENV PATH="/go/src:${PATH}"

# Install Certificate
RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates
RUN apt-get update && apt install ffmpeg -y && \
    ffmpeg -version && \
    ffprobe -version

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

EXPOSE 3210 3211

ENTRYPOINT ["/VideoProcessing"]
