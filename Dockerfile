# https://levelup.gitconnected.com/complete-guide-to-create-docker-container-for-your-golang-application-80f3fb59a15e
FROM golang:1.15-alpine AS builder

ARG PORT=8085
EXPOSE $PORT
# https://docs.docker.com/engine/reference/builder/#expose

# Set necessary environmet variables needed for our image
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Move to working directory /build
WORKDIR /build

# Copy and download dependency using go mod
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy the code into the container
COPY . .

# Build the application
RUN go build -o mastro .


# multistage build - we only copy the result (binary) into a fresh scratch image which is super light
FROM scratch

# copy binary
COPY --from=builder /build/mastro ./

# set default vars
ENV MASTRO_CONFIGPATH=cfg.yml

# set config.yaml using wget or local copy
COPY conf $MASTRO_CONFIGPATH 

# Command to run when starting the container
ENTRYPOINT ["./mastro"]


