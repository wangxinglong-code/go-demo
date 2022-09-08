FROM ""

#go build
MAINTAINER ""

ENV PORT 3003
EXPOSE $PORT

COPY bin/go-demo /
COPY ./conf /conf

ENTRYPOINT ["/go-demo"]
