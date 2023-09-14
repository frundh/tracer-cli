FROM golang:1.19 AS build
ENV GOPATH=/go
COPY . /src/
WORKDIR /src/
RUN go install

FROM ubuntu
COPY --from=build /go/bin/tracer-cli /usr/local/bin/tracer-cli
ENTRYPOINT [ "tracer-cli" ]
CMD [ "--help" ]