module github.com/kamva/hexa-echo

go 1.13

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/google/uuid v1.1.2
	github.com/kamva/gutil v0.0.0-20210827084201-35b6a3421580
	github.com/kamva/hexa v0.0.0-20211128175703-59125a2fe5ec
	github.com/kamva/tracer v0.0.0-20201115122932-ea39052d56cd
	github.com/labstack/echo/v4 v4.1.17
	github.com/labstack/gommon v0.3.0
	github.com/magiconair/properties v1.8.1
	github.com/mattn/go-colorable v0.1.8 // indirect
	github.com/stretchr/testify v1.7.0
	go.opentelemetry.io/otel v1.2.0
	go.opentelemetry.io/otel/exporters/jaeger v1.0.0-RC3
	go.opentelemetry.io/otel/exporters/prometheus v0.25.0
	go.opentelemetry.io/otel/metric v0.25.0
	go.opentelemetry.io/otel/sdk v1.2.0
	go.opentelemetry.io/otel/sdk/export/metric v0.25.0
	go.opentelemetry.io/otel/sdk/metric v0.25.0
	go.opentelemetry.io/otel/trace v1.2.0
)
