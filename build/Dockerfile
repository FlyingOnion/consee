ARG maintainer="FlyingOnion"
ARG version="0.25.12-alpha1"
ARG build_date="2025-12-05"
ARG url="https://github.com/FlyingOnion/consee"

FROM oven/bun:1.3-alpine as frontend
LABEL maintainer=${maintainer} \
      version=${version} \
      build_date=${build_date} \
      url=${url}

# For Chinese friends:
# Uncomment following commands to switch registry to npmmirror
#
# RUN cat <<EOF > ~/.bunfig.toml
# [install]
# registry = "https://registry.npmmirror.com/"
# EOF
WORKDIR /frontend
COPY ./frontend .
RUN bun install && bun run build



FROM golang:alpine as backend
LABEL maintainer=${maintainer} \
      version=${version} \
      build_date=${build_date} \
      url=${url}

# For Chinese friends:
# Uncomment the following line to use goproxy.cn as the GOPROXY
#
# RUN go env -w GOPROXY=https://goproxy.cn,direct
WORKDIR /go/src/github.com/FlyingOnion/consee/backend
COPY ./backend .
COPY --from=frontend /frontend/dist ./public
RUN go mod tidy && go build -o consee


# Run consee
FROM alpine:latest
LABEL maintainer=${maintainer} \
      version=${version} \
      build_date=${build_date} \
      url=${url}

# For Chinese friends:
# Uncomment the following line to switch alpine repo to aliyun
#
# RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories

RUN apk add --no-cache tzdata ca-certificates
WORKDIR /opt/consee
VOLUME [/opt/consee/config, /opt/consee/logs]
EXPOSE 3668

COPY --from=backend /go/src/github.com/FlyingOnion/consee/backend/consee .

RUN chmod 755 consee

RUN addgroup -S consee && adduser -S -G consee consee
USER consee

CMD ["./consee"]