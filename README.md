## What is Metrics?
Metrics exposes deafult Go runtime  and couple of custom [Prometheus](https://prometheus.io/docs/guides/go-application/) metrics.

## Dependencies

Build dependencies:
* golang 1.17
* [go-guardian](https://pkg.go.dev/github.com/shaj13/go-guardian/v2)
* [envconfig](https://github.com/kelseyhightower/envconfig)

## Build
To build binary run:
```  
go build -o app main.go
```
To build docker image ([gcr.io/distroless/base](https://github.com/GoogleContainerTools/distroless) is used for base image):
```
docker build -t metrics:0.0.1 .
```
## Run
To run just simply execute built binary
```  
./app
```
To run in docker:
```
docker run -ti -p 8080:8080 -v "$(pwd)"/tokens.csv:/tokens.csv -e METRICS_TOKEN_FILE=/tokens.csv --rm metrics:0.0.1
```
## Run tests
Couple of unit tests are implemented for `metrics` module
```
cd metrics/
go test
```
## Usage
Metrics settings are controlled via  **METRICS_*** environment variables
```
METRICS_HTTP_ADDR - HTTP <addr>:<port> to bind (default:":8080", mandatory: "false") 
METRICS_HANDLER_PATH - HTTP GET handler path (default:"/metrics", mandatory: "false")  
METRICS_TOKEN_FILE - Path to CSV token file (no defaults, mandatory: "true")
```
`METRICS_TOKEN_FILE` is defined in [CSV format](https://pkg.go.dev/github.com/shaj13/go-guardian/v2/auth/strategies/token#NewStaticFromFile). Sample file is provided in `tokens.csv`
