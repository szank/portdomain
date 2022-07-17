FROM ubuntu

MAINTAINER teamName@90poe.io

COPY portdomain-linux /app/portdomain
ENTRYPOINT ["/app/portdomain"]
