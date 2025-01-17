FROM golang:1.23-alpine as go-builder

WORKDIR /app
COPY go.mod go.sum ./
RUN apk add --update --no-cache git
RUN go mod download && go mod verify

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o application ./cmd/app

FROM scratch
# FROM debian:bookworm as debian-builder
COPY --from=go-builder /app/application /app/application
# RUN set -ex \
#     && sed -i -- 's/Types: deb/Types: deb deb-src/g' /etc/apt/sources.list.d/debian.sources \
#     && apt-get update \
#     && apt-get install -y --no-install-recommends \
#                build-essential \
#                cdbs \
#                devscripts \
#                equivs \
#                fakeroot \
#     && apt-get clean \
#     && rm -rf /tmp/* /var/tmp/*

CMD ["/app/application"]
