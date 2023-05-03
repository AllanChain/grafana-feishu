# The image is always native to our current machine
FROM --platform=$BUILDPLATFORM golang:1.20-alpine AS build

WORKDIR /app

ENV GO111MODULE=on
ENV GOPROXY=$GOPROXY

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./

ARG TARGETOS TARGETARCH
RUN GOOS=$TARGETOS GOARCH=$TARGETARCH go build -o /grafana-feishu


FROM alpine

WORKDIR /
COPY --from=build /grafana-feishu /grafana-feishu

EXPOSE 2387
RUN addgroup -S feishu && adduser -S feishu -G feishu
USER feishu

ENTRYPOINT ["/grafana-feishu"]
