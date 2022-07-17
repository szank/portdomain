FROM ubuntu

MAINTAINER teamName@90poe.io

COPY portdomain /app/portdomain
ENTRYPOINT ["/app/portdomain"]
