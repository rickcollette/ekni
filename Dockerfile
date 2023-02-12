#build stage
FROM golang:alpine AS builder
RUN apk add --no-cache git
WORKDIR /go/src/ekni
RUN git clone https://github.com/rickcollette/ekni.git .
RUN go get -d -v ./...
RUN go build -o /go/bin/ekni -v ./...

#final stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
RUN apk add nginx
COPY --from=builder /go/bin/ekni /ekni
RUN adduser -D ekni
RUN chown -R ekni:ekni /ekni
RUN chmod -R 700 /ekni
RUN setcap cap_net_admin,cap_net_bind_service,cap_net_raw+ep /ekni
EXPOSE 8675
EXPOSE 443
EXPOSE 51820
COPY nginx.conf /etc/nginx/nginx.conf
RUN chown -R ekni:ekni /etc/nginx/nginx.conf
RUN chmod -R 700 /etc/nginx/nginx.conf
RUN echo "daemon off;" >> /etc/nginx/nginx.conf
CMD nginx && su ekni -c "/ekni"