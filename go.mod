module main.go

go 1.16

replace github.com/byuoitav/central-via-alert-service => /home/creeder/go/src/github.com/byuoitav/central-via-alert-service

require (
	github.com/byuoitav/central-via-alert-service v0.0.0-00010101000000-000000000000
	github.com/byuoitav/common v0.0.0-20200521193927-1fdf4e0a4271
	github.com/byuoitav/kramer-driver v0.1.16
	github.com/go-kivik/couchdb/v3 v3.2.8
	github.com/go-kivik/kivik/v3 v3.2.3
	github.com/labstack/echo v3.3.10+incompatible
	github.com/spf13/pflag v1.0.5
	go.uber.org/zap v1.18.1
)
