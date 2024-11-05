package xo

import (
	"database/sql"
	"testing"

	"github.com/nekomeowww/xo/protobufs/testpb"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/encoding/protojson"
)

func TestProtoValueScanner(t *testing.T) {
	t.Run("NULL", func(t *testing.T) {
		scanner := ProtoValueScanner[testpb.TestMessage]{}
		require.NotNil(t, scanner)

		val, err := scanner.Value(nil)
		require.NoError(t, err)
		assert.Equal(t, sql.NullString{}, val)

		value, err := val.(sql.NullString).Value()
		require.NoError(t, err)
		require.Nil(t, value)

		pb, err := scanner.FromValue(&sql.NullString{String: "", Valid: false})
		require.NoError(t, err)
		require.Nil(t, pb)

		pb, err = scanner.FromValue(lo.ToPtr(val.(sql.NullString)))
		require.NoError(t, err)
		require.Nil(t, pb)
	})

	t.Run("NonNULL", func(t *testing.T) {
		original := &testpb.TestMessage{
			Property_1: "Hello, World!",
			Property_2: "John Doe",
			OneofField: &testpb.TestMessage_PossibleOne{
				PossibleOne: &testpb.PossibleOne{
					Property_1: "Hello, World!",
					Property_2: "John Doe",
				},
			},
		}

		scanner := ProtoValueScanner[testpb.TestMessage]{}
		require.NotNil(t, scanner)

		val, err := scanner.Value(original)
		require.NoError(t, err)

		str, ok := val.(*sql.NullString)
		require.True(t, ok)
		require.NotEmpty(t, str)

		bytes, err := protojson.Marshal(original)
		require.NoError(t, err)
		assert.Equal(t, string(bytes), str.String)

		pb, err := scanner.FromValue(&sql.NullString{String: str.String, Valid: true})
		require.NoError(t, err)
		require.NotNil(t, pb)

		assert.Equal(t, original, pb)
	})
}
