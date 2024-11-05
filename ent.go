package xo

import (
	"database/sql"
	"database/sql/driver"
	"fmt"

	"entgo.io/ent/schema/field"
	"github.com/samber/lo"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

var _ field.TypeValueScanner[*any] = (*ProtoValueScanner[any])(nil)

/*
ProtoValueScanner is a field.ValueScanner that implements the ent.ValueScanner interface as helper for
working with protobuf messages. It is used to scan and convert protobuf messages to and from the database.

	func (SomeTable) Fields() []ent.Field {
		return []ent.Field{
			field.
			String("payload").
			ValueScanner(utils.ProtoValueScanner[somepb.YourMessage]{}).
			GoType(&somepb.YourMessage{}).
			SchemaType(map[string]string{
				dialect.Postgres: "jsonb",
				dialect.MySQL:    "json",
				dialect.SQLite:   "json",
			}),
		}
	}
*/
type ProtoValueScanner[T any] struct {
}

func (s ProtoValueScanner[T]) v(data *T) (driver.Value, error) {
	if data == nil {
		return sql.NullString{}, nil
	}

	pbMessage, ok := any(data).(proto.Message)
	pbMessage = lo.Must(pbMessage, ok)

	bytes, err := protojson.Marshal(pbMessage)
	if err != nil {
		return nil, err
	}

	return &sql.NullString{String: string(bytes), Valid: true}, nil
}
func (s ProtoValueScanner[T]) s(sqlData *sql.NullString) (*T, error) {
	if sqlData == nil {
		return nil, nil
	}
	if !sqlData.Valid {
		return nil, nil
	}

	var data T

	pbMessage, ok := any(&data).(proto.Message)
	pbMessage = lo.Must(pbMessage, ok)

	err := protojson.Unmarshal([]byte(sqlData.String), pbMessage)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

// Value returns the driver.Valuer for the GoType.
func (s ProtoValueScanner[T]) Value(data *T) (driver.Value, error) {
	return s.v(data)
}

// ScanValue returns a new ValueScanner that functions as an
// intermediate result between database value and GoType value.
// For example, sql.NullString or sql.NullInt.
func (s ProtoValueScanner[T]) ScanValue() field.ValueScanner {
	return new(sql.NullString)
}

// FromValue returns the field instance from the ScanValue
// above after the database value was scanned.
func (s ProtoValueScanner[T]) FromValue(value driver.Value) (vt *T, err error) {
	switch v := value.(type) {
	case *sql.NullString:
		return s.s(v)
	case *T:
		return v, nil
	case *any:
		return s.s(FromPtrAny[*sql.NullString](v))
	case any:
		vFromAny, _ := v.(*sql.NullString)
		return s.s(vFromAny)
	}

	str, ok := value.(*sql.NullString)
	if !ok {
		return vt, fmt.Errorf("unexpected input for FromValue: %T", value)
	}

	return s.s(str)
}
