package main

import (
	"github.com/ductone/protoc-gen-pgdb/internal/pgdb"
	pgs "github.com/lyft/protoc-gen-star/v2"
	pgsgo "github.com/lyft/protoc-gen-star/v2/lang/go"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
)

func main() {
	minEdition := int32(descriptorpb.Edition_EDITION_PROTO2)
	maxEdition := int32(descriptorpb.Edition_EDITION_2023)
	features := uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL | pluginpb.CodeGeneratorResponse_FEATURE_SUPPORTS_EDITIONS)
	pgs.Init(
		pgs.DebugEnv("DEBUG_PG_PGDB"),
		pgs.SupportedFeatures(&features),
		pgs.MinimumEdition(&minEdition),
		pgs.MaximumEdition(&maxEdition),
	).
		RegisterModule(pgdb.New()).
		RegisterPostProcessor(pgsgo.GoImports()).
		RegisterPostProcessor(pgsgo.GoFmt()).
		Render()
}
