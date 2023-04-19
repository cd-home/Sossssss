ARG v
FROM golang:${v} AS Builder
LABEL stage=Build
WORKDIR /app/build/
COPY . .
ENV GOOS=linux CGO_ENABLED=0 GOARCH=arm64 GOPROXY=https://goproxy.cn,direct
RUN go build -gcflags="-N -l" -o app  main.go

FROM alpine:3.14
LABEL stage=Run
WORKDIR /app/release/
COPY --from=Builder /app/build/app /app/release/
EXPOSE 8999
ENTRYPOINT ["./app"]
# docker build -t my-app -f Dockerfile --build-arg v=1.18 .