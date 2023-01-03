package main

import (
	"github.com/ductone/protoc-gen-pgdb/internal/pgdb"
	pgs "github.com/lyft/protoc-gen-star"
	pgsgo "github.com/lyft/protoc-gen-star/lang/go"
)

func main() {
	pgs.Init(pgs.DebugEnv("DEBUG_PG_PGDB")).
		RegisterModule(pgdb.New()).
		RegisterPostProcessor(pgsgo.GoImports()).
		RegisterPostProcessor(pgsgo.GoFmt()).
		Render()
}
