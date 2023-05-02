# The image is always native to our current machine
FROM --platform=$BUILDPLATFORM golang:1.20-alpine AS build

WORKDIR /app

ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn,direct

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./

ARG TARGETOS TARGETARCH
RUN GOOS=$TARGETOS GOARCH=$TARGETARCH go build -o /grafana-feishu


FROM alpine

WORKDIR /
COPY --from=build /grafana-feishu /grafana-feishu

EXPOSE 2387
USER nonroot:nonroot

ENTRYPOINT ["/grafana-feishu"]
