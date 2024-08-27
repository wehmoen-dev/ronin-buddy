FROM golang:1.22 AS builder

ENV GO111MODULE=on \
  CGO_ENABLED=0 \
  GOOS=linux \
  GOARCH=amd64

#RUN apt update && apt install -y  xz-utils && \
#  wget https://github.com/upx/upx/releases/download/v4.2.4/upx-4.2.4-amd64_linux.tar.xz -O /tmp/upx.tar.xz && \
#  tar -xf /tmp/upx.tar.xz -C /tmp && mv /tmp/upx-4.2.4-amd64_linux/upx /bin/upx

WORKDIR /src
COPY . .

RUN go build \
  -ldflags "-s -w -extldflags '-static'" \
  -o /bin/app \
  . \
  && strip /bin/app
#  && upx -q -9 /bin/app

RUN echo "nobody:x:65534:65534:Nobody:/:" > /etc_passwd



FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc_passwd /etc/passwd
COPY --from=builder --chown=65534:0 /bin/app /app

USER nobody
ENTRYPOINT ["/app"]