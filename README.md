## What is Metrics?
Metrics exposes deafult Go runtime  and couple of custom [Prometheus](https://prometheus.io/docs/guides/go-application/)
with HTTP authorization via token.

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
METRICS_TOKEN_FILE=./tokens.csv ./app
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
## Examples
Get all metrics with token:
```
curl -H "Authorization: Bearer 90d64460d14870c08c81352a05dedd3465940a7" http://localhost:8080/metrics
```
Get custom metrics (Counter/Summary):
```
curl -s -H "Authorization: Bearer 90d64460d14870c08c81352a05dedd3465940a7" http://localhost:8080/metrics | grep '^my'
my_own_private_metrics_example_counter 1
my_own_private_metrics_example_summary{service="random_summary_value",quantile="0.5"} 0.5238203060500009
my_own_private_metrics_example_summary{service="random_summary_value",quantile="0.9"} 0.9002751472643116
my_own_private_metrics_example_summary{service="random_summary_value",quantile="0.99"} 0.9891320920320221
my_own_private_metrics_example_summary_sum{service="random_summary_value"} 268.91727211548044
my_own_private_metrics_example_summary_count{service="random_summary_value"} 550
```