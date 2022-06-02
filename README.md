# go-otl-demo
Just a demo application with OpenTelemetry Instrumentation and Jaeger Exporter

# How to run
## Setup Jaeger
Run the docker command bellow:
```sh
docker run -d --name jaeger \                
-e COLLECTOR_ZIPKIN_HTTP_PORT=9411 \
  -p 5775:5775/udp \
  -p 6831:6831/udp \
  -p 6832:6832/udp \
  -p 5778:5778 \
  -p 16686:16686 \
  -p 14268:14268 \
  -p 9411:9411 \
  jaegertracing/all-in-one:1.6
```
## Setup Go application
Run the code bellow with environment variable:

```sh
APPLICATION_NAME=application-name \
APPLICATION_PORT=9000 \
JAEGER_HOST=http://localhost:14268/api/traces \
JSON_PLACEHOLDER_API_HOST=https://jsonplaceholder.typicode.com \
go run cmd/main.go
```