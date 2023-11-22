FROM alpine
ARG K6_BINARY

RUN apk add --no-cache ca-certificates && \
    adduser -D -u 1000 -g 1000 k6
COPY ${K6_BINARY} /usr/bin/k6

USER 1000
WORKDIR /home/k6
ENTRYPOINT ["k6"]
