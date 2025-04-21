protoc-gen-pgdb

A Protocol Buffers code generator plugin that generates PostgreSQL database code from protobuf definitions. This tool allows you to define your database schema, indexes, and queries using protobuf annotations, and then generates the necessary Go code to interact with PostgreSQL databases.

## Overview

protoc-gen-pgdb bridges the gap between Protocol Buffers and PostgreSQL by automatically generating type-safe database access code. It supports:

- Mapping protobuf messages to PostgreSQL tables
- Creating and managing indexes
- Support for various PostgreSQL data types
- Query building

## Installation

### Prerequisites

- Go 1.23 or higher
- Protocol Buffers compiler
- [Buf](https://buf.build/) (recommended for code generation)

### Installing the plugin

```bash
go install github.com/ductone/protoc-gen-pgdb@latest
```

Or build from source:

```bash
git clone https://github.com/ductone/protoc-gen-pgdb.git
cd protoc-gen-pgdb
make build
```

## Interfaces

### Message Options

The `pgdb.v1.msg` extension provides options for configuring how a protobuf message maps to a PostgreSQL table:

- `disabled`: Disables code generation for this message
- `indexes`: Defines database indexes
- `tenant_id_field`: Specifies the field to use for multi-tenancy
- `nested_only`: Indicates that this message should only be used as a nested type
- `partitioned`: Enables table partitioning
- `partitioned_by_created_at`: Partitions by the created_at timestamp
- `partitioned_by_date_range`: Specifies the date range for partitioning
- `stats`: Configures PostgreSQL statistics collection

### Field Options

The `pgdb.v1.options` extension provides options for configuring how a protobuf field maps to a PostgreSQL column:

- `full_text_type`: Configures full-text search type (NONE, EXACT, ENGLISH)
- `full_text_weight`: Sets the weight for full-text search (A, B, C, D)
- `message_behavior`: Configures how message fields are stored (JSONB, etc.)

## Examples

See the `example` directory for complete examples of how to use protoc-gen-pgdb.

## Development

### Building

```bash
make build
```

### Running tests

```bash
make test
```

### Generating example code

```bash
make example
```

## Dependencies

- [protoc-gen-star](https://github.com/lyft/protoc-gen-star): Framework for building Protobuf code generators
- [pgx](https://github.com/jackc/pgx): PostgreSQL driver and toolkit for Go
- [goqu](https://github.com/doug-martin/goqu): SQL builder for Go

## License

See the [LICENSE](LICENSE) file for details.