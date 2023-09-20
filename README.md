# Tracer CLI
Simple tool for sending sample traces to Jaeger or Zipkin

## Run
```bash
tracer-cli trace jaeger -n "sample-trace" -u jaeger:6831
```
```bash
tracer-cli trace zipkin -n "sample-trace" -c http://jaeger:9411/api/v2/spans -r https://www.google.com
```

## Build
```bash
go build
# OR
docker build -t tracer-cli .
```