# execute from project root

# build inside an own docker
FROM golang:alpine AS build

RUN apk update \
    && apk --no-cache add make git

COPY . /go/src/github.com/IceflowRE/irc-discord-relay

RUN cd /go/src/github.com/IceflowRE/irc-discord-relay \
    && make

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
    && mkdir /idr \
    && mkdir /idr/config \
    && chown iceflower: -R /idr

WORKDIR /idr
COPY --from=build /go/src/github.com/IceflowRE/irc-discord-relay/build/irc-discord-relay .

USER iceflower
ENTRYPOINT ["./irc-discord-relay", "-c", "./config/config.json"]
