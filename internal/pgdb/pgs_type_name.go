package pgdb

import (
	"fmt"
	"os"
	"regexp"

	pgs "github.com/lyft/protoc-gen-star"
	pgsgo "github.com/lyft/protoc-gen-star/lang/go"
)

func (c *importTracker) Type(f pgs.Field) pgsgo.TypeName {
	ft := f.Type()

	var t pgsgo.TypeName
	switch {
	case ft.IsMap():
		key := scalarType(ft.Key().ProtoType())
		return pgsgo.TypeName(fmt.Sprintf("map[%s]%s", key, c.elType(ft)))
	case ft.IsRepeated():
		return pgsgo.TypeName(fmt.Sprintf("[]%s", c.elType(ft)))
	case ft.IsEmbed():
		return c.importableTypeName(f, ft.Embed()).Pointer()
	case ft.IsEnum():
		fmt.Fprintf(os.Stderr, "🦐 %s\n", f.Name().String())
		t = c.importableTypeName(f, ft.Enum())
	default:
		t = scalarType(ft.ProtoType())
	}

	if f.HasPresence() {
		return t.Pointer()
	}

	return t
}

func (ix *importTracker) importableTypeName(f pgs.Entity, e pgs.Entity) pgsgo.TypeName {
	t := pgsgo.TypeName(ix.ctx.Name(e))

	fmt.Fprintf(os.Stderr, "🌮 %s %s %s\n", f.Name().String(), ix.ctx.ImportPath(e).String(), ix.ctx.ImportPath(f).String())
	if ix.ctx.ImportPath(e) == ix.ctx.ImportPath(f) {
		fmt.Fprintf(os.Stderr, " %s 1\n", t)
		return t
	}
	pkgName := ix.ctx.PackageName(e)
	importName := ix.ctx.ImportPath(e)
	matched, err := regexp.MatchString(`^v(\d)+$`, pkgName.String())
	if err != nil {
		panic(err)
	}

	if matched {
		niceName := importName.Dir().Base()
		pkgName = pgs.Name(fmt.Sprintf("%s_%s", niceName, pkgName.String()))
	}

	for {
		to, ok := ix.typeMapper[pkgName]
		if !ok {
			ix.typeMapper[pkgName] = importName
			break
		}
		if to == importName {
			break
		}
		pkgName += "x"
	}
	fmt.Fprintf(os.Stderr, " %s 2\n", pgsgo.TypeName(fmt.Sprintf("%s.%s", pkgName, t)))
	return pgsgo.TypeName(fmt.Sprintf("%s.%s", pkgName, t))
}

func (ix *importTracker) importableTypeName2(f pgs.Entity, e pgs.Entity) pgsgo.TypeName {
	t := pgsgo.TypeName(ix.ctx.Name(e))

	fmt.Fprintf(os.Stderr, "🌮 %s %s %s\n", f.Name().String(), ix.ctx.ImportPath(e).String(), ix.ctx.ImportPath(f).String())
	if ix.ctx.ImportPath(e) == ix.ctx.ImportPath(f) {
		fmt.Fprintf(os.Stderr, " %s 1\n", t)
		return ""
	}
	pkgName := ix.ctx.PackageName(e)
	importName := ix.ctx.ImportPath(e)
	matched, err := regexp.MatchString(`^v(\d)+$`, pkgName.String())
	if err != nil {
		panic(err)
	}

	if matched {
		niceName := importName.Dir().Base()
		pkgName = pgs.Name(fmt.Sprintf("%s_%s", niceName, pkgName.String()))
	}

	for {
		to, ok := ix.typeMapper[pkgName]
		if !ok {
			ix.typeMapper[pkgName] = importName
			break
		}
		if to == importName {
			break
		}
		pkgName += "x"
	}
	fmt.Fprintf(os.Stderr, " %s 2\n", pgsgo.TypeName(fmt.Sprintf("%s.%s", pkgName, t)))
	return pgsgo.TypeName(pkgName)
}

func (c *importTracker) elType(ft pgs.FieldType) pgsgo.TypeName {
	el := ft.Element()
	switch {
	case el.IsEnum():
		return c.importableTypeName(ft.Field(), el.Enum())
	case el.IsEmbed():
		return c.importableTypeName(ft.Field(), el.Embed()).Pointer()
	default:
		return scalarType(el.ProtoType())
	}
}

func scalarType(t pgs.ProtoType) pgsgo.TypeName {
	switch t {
	case pgs.DoubleT:
		return "float64"
	case pgs.FloatT:
		return "float32"
	case pgs.Int64T, pgs.SFixed64, pgs.SInt64:
		return "int64"
	case pgs.UInt64T, pgs.Fixed64T:
		return "uint64"
	case pgs.Int32T, pgs.SFixed32, pgs.SInt32:
		return "int32"
	case pgs.UInt32T, pgs.Fixed32T:
		return "uint32"
	case pgs.BoolT:
		return "bool"
	case pgs.StringT:
		return "string"
	case pgs.BytesT:
		return "[]byte"
	default:
		panic("unreachable: invalid scalar type")
	}
}
