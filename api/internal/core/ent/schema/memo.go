package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
	"github.com/google/uuid"
)

// Memo holds the schema definition for the Memo entity.
type Memo struct {
	ent.Schema
}

// Fields of the Memo.
func (Memo) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New),
		field.String("title").
			MaxLen(512),
		field.String("content").
			MaxLen(20_000),
	}
}

// Mixin of the Memo.
func (Memo) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.Time{},
	}
}

// Edges of the Memo.
func (Memo) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("owner", User.Type).
			Ref("memos").
			Required().
			Unique(),
		edge.To("tags", Tag.Type),
	}
}
