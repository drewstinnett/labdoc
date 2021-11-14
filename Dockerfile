FROM alpine:3.14
# hadolint ignore=DL3018
RUN apk add --no-cache git
COPY labdoc /
ENTRYPOINT ["/labdoc"]
