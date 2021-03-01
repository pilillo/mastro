# https://github.com/moby/moby/issues/37345
ARG ARTIFACT=mastro-crawlers

# https://levelup.gitconnected.com/complete-guide-to-create-docker-container-for-your-golang-application-80f3fb59a15e
FROM golang:1.15-alpine AS builder
ARG ARTIFACT
ARG PORT=8085
EXPOSE $PORT
# https://docs.docker.com/engine/reference/builder/#expose

# Set necessary environmet variables needed for our image
ENV GO111MODULE=on \
    CGO_ENABLED=1 \
    GOOS=linux \
    GOARCH=amd64

# Move to working directory /build
WORKDIR /build

RUN apk add --no-cache git make build-base krb5-pkinit krb5-dev krb5

# Copy and download dependency using go mod
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy the code into the container
COPY . .

# Build the application
RUN go build -tags=kerberos -o ${ARTIFACT} .

# multistage build - we only copy the result (binary) into a fresh image which is super light
FROM alpine:3.12.4
ARG ARTIFACT
ENV ARTIFACT=${ARTIFACT}

# copy binary
COPY --from=builder /build/${ARTIFACT} ./

# set default vars
ENV MASTRO_CONFIG=/conf/catalogue/

# set config.yaml using wget or local copy
COPY conf $MASTRO_CONFIG

ENV GIN_MODE=release

# Command to run when starting the container
ENTRYPOINT ["sh", "-c", "./${ARTIFACT}"]