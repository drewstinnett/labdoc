FROM alpine:3.14
COPY labdoc /
ENTRYPOINT ["/labdoc"]
