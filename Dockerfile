FROM alpine:latest
COPY ./prober-demo /prober-demo
ENTRYPOINT ["/prober-demo"]