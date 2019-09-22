# execute from project root

# build inside an own docker
FROM golang:alpine AS build

# without musl-dev gcc it wont work in arm64v8, even if it does on amd64
RUN apk update \
    && apk --no-cache add make git gcc musl-dev

COPY ./ /build

RUN cd /build \
    && go build -x -o irc-discord-relay

###################################
# Set the base image
FROM alpine:latest

# bot version
ARG VERSION

# informations
LABEL author="Iceflower S" \
      maintainer="Iceflower S" \
      email="iceflower@gmx.de" \
      version="$VERSION" \
      description="IRC Discord Relay"

# discord needs certificates
RUN apk update && apk add --no-cache ca-certificates

# add server user and set permissions
RUN adduser -S -D iceflower \
    && mkdir /app \
    && mkdir /app/config \
    && chown iceflower: -R /app

WORKDIR /app
COPY --from=build /build/irc-discord-relay .

USER iceflower
ENTRYPOINT ["./irc-discord-relay", "-c", "./config/config.json"]
