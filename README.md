# ekni
## WireGuard Configuration Server

-----

Enki is a god in Sumerian mythology who was worshipped as the god of wisdom, magic, and water. He was considered the patron deity of the city of Eridu and was said to have been the creator of humans and the one who brought civilization to the world. Enki was often depicted with the symbols of a goat and a fish, representing his association with both land and water. In some legends, he was also described as a trickster and a cunning deity who used his intelligence and wisdom to manipulate situations to his advantage. Despite his mischievous nature, Enki was revered as a powerful deity who had a significant impact on the lives of the ancient Sumerians.

-----

# Dockerfile Documentation

This Dockerfile builds an image for the ekni application. The image is based on the Alpine Linux image and includes the ekni binary, nginx, and a custom nginx configuration file. The image is optimized for deployment in a production environment.

The Dockerfile uses a multi-stage build process, with two stages: the build stage and the final stage.
Build Stage

The build stage uses the golang:alpine image and installs the necessary packages for building the ekni binary. This stage sets the working directory to /go/src/app, copies the source code into the image, and runs go get to retrieve the required dependencies. Finally, it runs go build to compile the ekni binary and places it in the /go/bin directory.

# Build Stage

The build stage uses the golang:alpine image and installs the necessary packages for building the ekni binary. This stage sets the working directory to /go/src/app, copies the source code into the image, and runs go get to retrieve the required dependencies. Finally, it runs go build to compile the ekni binary and places it in the /go/bin directory.

```
FROM golang:alpine AS builder
RUN apk add --no-cache git
WORKDIR /go/src/app
COPY . .
RUN go get -d -v ./...
RUN go build -o /go/bin/app -v ./...
```
# Final Stage

The final stage uses the Alpine Linux image and installs the necessary packages for running the ekni application and nginx. This stage copies the ekni binary from the build stage into the image and sets it as the entrypoint for the image. The final stage also exposes the necessary ports (TCP 8675, TCP 443, and UDP 51820) and includes a label with the name and version of the ekni application.

```
FROM alpine:latest
RUN apk --no-cache add ca-certificates
RUN apk add nginx
COPY --from=builder /go/bin/app /app
ENTRYPOINT /app
LABEL Name=ekni Version=0.0.1
EXPOSE 8675
EXPOSE 443
EXPOSE 51820
```
