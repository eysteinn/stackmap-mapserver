FROM golang:1.19.1-alpine AS build_base

ENV CGO_ENABLED=1
ENV GO111MODULE=on
RUN apk add --no-cache git  git gcc g++

# Set the Current Working Directory inside the container
WORKDIR /src

# We want to populate the module cache based on the go.{mod,sum} files.
COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

# Build the Go app
RUN go build -o ./out/app ./main.go


# Start fresh from a smaller image
FROM alpine:3.12
RUN apk add ca-certificates

WORKDIR /app

COPY --from=build_base /src/out/app /app/config_runner
#COPY --from=build_base /src/data /app/data
COPY ./data /app/data
#RUN mkdir -p ./mapfiles
RUN ln -s /mapfiles ./mapfiles
#RUN ln -s ./mapfiles /mapfiles
COPY ./data/rasters.map /mapfiles/rasters.map
RUN chmod +x config_runner

# This container exposes port 8080 to the outside world
EXPOSE 3000

# Run the binary program produced by `go install`
#ENTRYPOINT ./config_runner
ENTRYPOINT [ "/bin/sh", "-c", "cp -R /app/data/rasters.map /mapfiles/ && ./config_runner" ]
