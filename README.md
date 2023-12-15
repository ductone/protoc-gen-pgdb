# protoc-gen-pgdb


## interfaces

## Quick Example of adding a `MessageOption`

In one case we wanted to add a partitioned option to a message.

Like so,
```
message Pasta {
  option (pgdb.v1.msg).partitioned = true;
  ...
}
```

To do this we added a new field to the `MessageOptions` message in `pgdb.proto` like so,
```
message MessageOptions {
  bool disabled = 1;
  message Index {
    ...
  }
 
  ...

  // if this message is used then we create a partitioned table and partition by
  // tenant_id.
  This is new --> bool partitioned = 5;
}
```

### Some context
```
message Foo {
    ...
}
```
Foo is a `descriptor` and it has a `DescriptorFieldOption` struct in `pgdb/v1/descriptor.go` that is used to add options to the `Foo` message.

Field options are the same thing but are `RecordOptions` in `pgdb/v1/record_options.go` which also has a struct like so 

```
type RecordOption struct {
	Prefix   string
	IsNested bool
	Nulled   bool
}
```

### Going to back to our example 
I added a new option to the descriptor struct like so,
```
type DescriptorFieldOption struct {
	Prefix        string
	IsNested      bool
	IsPartitioned bool
}
```

Then I added a new method `IsPartitioned` to the `DescriptorFieldOption` interface like so,
```
type Descriptor interface {
	TableName() string

	Fields(opts ...DescriptorFieldOptionFunc) []*Column

    ...

	IsPartitioned() bool
}

```

But we want all of our messages to have this option so we need to add it to the descriptor template (`internal/pgdb/templates/descriptor.tmpl`) like so,
```
func (d *{{.ReceiverType}}) IsPartitioned() bool {
    return {{.IsPartitioned}}
}
```

If you search for where this template is used you will find it in `internal/pgdb/pgdb_descriptor.go` in the `renderDescriptor` method.

Next I add the newly added message option to the context
```
type descriptorTemplateContext struct {
	...
	IsPartitioned   bool
}
```

And then we pass it into the template creation like so,
```
    // This gets the message options for a message
    fext := pgdb_v1.MessageOptions{}
	_, err := m.Extension(pgdb_v1.E_Msg, &fext)
	if err != nil {
		panic(err)
	}

    ...

    c := &descriptorTemplateContext{
		...
		IsPartitioned:   fext.Partitioned,
	}

	return templates["descriptor.tmpl"].Execute(w, c)
```

Hence in the end my original Pasta message should include this in the generated pb file
```
    func (m *Pasta) IsPartitioned() bool {
        return true
    }
```

Hope this helps. It should be very similar for FieldOptions and IndexOptions I imagine.

