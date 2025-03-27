module replicator-chess

go 1.23.5

replace github.com/josephchapman/datasources/cmn => ../../cmn

require github.com/josephchapman/datasources/cmn v0.0.0-20250327020326-24fe48d58f1d // direct

require (
	github.com/apapsch/go-jsonmerge/v2 v2.0.0 // indirect
	github.com/google/uuid v1.3.1 // indirect
	github.com/influxdata/influxdb-client-go/v2 v2.14.0 // indirect
	github.com/influxdata/line-protocol v0.0.0-20200327222509-2487e7298839 // indirect
	github.com/oapi-codegen/runtime v1.0.0 // indirect
	golang.org/x/net v0.23.0 // indirect
)
