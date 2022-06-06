module github.com/kamva/hexa-echo

go 1.13

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/google/uuid v1.1.2
	github.com/gorilla/sessions v1.2.1
	github.com/kamva/gutil v0.0.0-20210827084201-35b6a3421580
	github.com/kamva/hexa v0.0.0-20220326081636-83739571eaff
	github.com/kamva/tracer v0.0.0-20201115122932-ea39052d56cd
	github.com/labstack/echo/v4 v4.7.2
	github.com/labstack/gommon v0.3.1
	github.com/magiconair/properties v1.8.1
	github.com/mailru/easyjson v0.7.7
	github.com/sethvargo/go-limiter v0.7.2
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
