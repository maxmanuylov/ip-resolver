FROM alpine
COPY ip-resolver /opt/bin/ip-resolver
ENTRYPOINT ["/opt/bin/ip-resolver"]
