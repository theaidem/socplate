FROM alpine:3.10

RUN apk add -U --no-cache curl && \
    curl https://github.com/pressly/goose/releases/download/v2.6.0/goose-linux64 -L --output /bin/goose && \
    apk del curl && \
    chmod +x /bin/goose

CMD tail -f /dev/null 