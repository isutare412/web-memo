package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
)

// Collaboration holds the schema definition for the Collaboration entity.
type Collaboration struct {
	ent.Schema
}

// Fields of the Collaboration.
func (Collaboration) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("user_id", uuid.UUID{}),
		field.UUID("memo_id", uuid.UUID{}),
		field.Bool("approved").
			Default(false),
		field.Time("create_time").
			Default(time.Now).
			Immutable(),
		field.Time("update_time").
			Default(time.Now),
	}
}

// Indexes of the Collaboration.
func (Collaboration) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("memo_id", "user_id").
			Unique(),
		index.Fields("user_id"),
	}
}

// Edges of the Collaboration.
func (Collaboration) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("collaborator", User.Type).
			Required().
			Unique().
			Annotations(entsql.OnDelete(entsql.Cascade)).
			Field("user_id"),
		edge.To("memo", Memo.Type).
			Required().
			Unique().
			Annotations(entsql.OnDelete(entsql.Cascade)).
			Field("memo_id"),
	}
}
