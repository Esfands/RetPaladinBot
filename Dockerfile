# Go image tag
ARG GOLANG_TAG=1.20-bullseye

FROM golang:$GOLANG_TAG as builder

ENV TWITCH_BOT_PREFIX=""
ENV TWITCH_BOT_CHANNEL=""
ENV TWITCH_BOT_CHANNEL_ID=""
ENV TWITCH_BOT_USERNAME=""
ENV TWITCH_BOT_OAUTH=""
ENV TWITCH_HELIX_CLIENT_ID=""
ENV TWITCH_HELIX_CLIENT_SECRET=""

ENV HTTP_ADDRESS=""
ENV HTTP_PORTS_REST=""

ENV API_KEYS_LASTFM=""

# Package path for ldflags
ARG PACKAGE=""
# Version of the REST build
ARG VERSION=""
# Commit hash
ARG COMMIT=""

# Github user required for downloading private Go modules
ARG ACCESS_TOKEN_USR=""
# Github user access token required for downloading private Go modules
ARG ACCESS_TOKEN_PWD=""

# Install required dependencies
RUN apt-get update && apt-get install -y \
    ca-certificates \
    git \
    build-essential \
    gcc

# Create the user and group files that will be used in the running
# container to run the process as an unprivileged user.
RUN mkdir /user && \
    echo 'nobody:x:65534:65534:nobody:/:' > /user/passwd && \
    echo 'nobody:x:65534:' > /user/group

# Create a netrc file using the credentials specified using --build-arg
RUN printf "machine github.com\n\
     login ${ACCESS_TOKEN_USR}\n\
     password ${ACCESS_TOKEN_PWD}\n\
     \n\
     machine api.github.com\n\
     login ${ACCESS_TOKEN_USR}\n\
     password ${ACCESS_TOKEN_PWD}\n"\
     >> /root/.netrc

RUN chmod 600 /root/.netrc

# Disable terminal prompts for Git
ENV GIT_TERMINAL_PROMPT=0

# Set the working directory outside $GOPATH to enable the support for modules.
WORKDIR /src

# Import code from context.
COPY . .

# Import the vendor directory into the Docker build.
COPY ./vendor ./vendor

COPY ./config /app/config

# Build the executable to `/app`. Mark the build as statically linked.
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod=vendor -o /app/server -ldflags="-X 'main.Version=${VERSION}' -X 'main.CommitHash=${COMMIT}'" ./cmd/app/main.go

# Use a more complete base image for the final stage
FROM debian:buster-slim AS final

# Import the user and group files from the first stage.
COPY --from=builder /user/group /user/passwd /etc/

# Import the Certificate-Authority certificates for enabling HTTPS.
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Import the compiled executable from the first stage.
COPY --from=builder /app/server /app/server

COPY --from=builder /app/config /config

# Create a writable tmp directory
RUN mkdir -p /tmp && chmod 1777 /tmp

# Perform any further action as an unprivileged user.
USER nobody:nobody

# Run the compiled binary.
ENTRYPOINT ["/app/server"]