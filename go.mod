module github.com/ductone/protoc-gen-pgdb

go 1.25

require (
	github.com/clipperhouse/jargon v1.0.9
	github.com/clipperhouse/uax29 v1.16.0
	github.com/davecgh/go-spew v1.1.1
	github.com/doug-martin/goqu/v9 v9.19.0
	github.com/gaissmai/extnetip v1.3.1
	github.com/jackc/pgx/v5 v5.7.6
	github.com/json-iterator/go v1.1.12
	github.com/lyft/protoc-gen-star/v2 v2.0.4
	github.com/pquerna/protoc-gen-dynamo v0.9.0
	github.com/segmentio/ksuid v1.0.4
	github.com/stretchr/testify v1.8.1
	golang.org/x/exp v0.0.0-20251209150349-8475f28825e9
	google.golang.org/protobuf v1.36.10
)

require (
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/kljensen/snowball v0.10.0 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/lib/pq v1.10.2 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rogpeppe/go-internal v1.9.0 // indirect
	github.com/spf13/afero v1.15.0 // indirect
	golang.org/x/crypto v0.46.0 // indirect
	golang.org/x/mod v0.31.0 // indirect
	golang.org/x/net v0.48.0 // indirect
	golang.org/x/sync v0.19.0 // indirect
	golang.org/x/text v0.32.0 // indirect
	golang.org/x/tools v0.40.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

// https://github.com/lyft/protoc-gen-star/pull/132
replace github.com/lyft/protoc-gen-star/v2 => github.com/pquerna/protoc-gen-star/v2 v2.0.0-20250415201647-653a078eb414

// Upstream PR: https://github.com/clipperhouse/jargon/pull/18
replace github.com/clipperhouse/jargon => github.com/pquerna/jargon v0.0.0-20250830190803-0dfe4377f37f
