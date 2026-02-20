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

// Subscription holds the schema definition for the Subscription entity.
type Subscription struct {
	ent.Schema
}

// Fields of the Subscription.
func (Subscription) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("user_id", uuid.UUID{}),
		field.UUID("memo_id", uuid.UUID{}),
		field.Bool("approved").
			Default(true),
		field.Time("create_time").
			Default(time.Now).
			Immutable(),
	}
}

// Indexes of the Subscription.
func (Subscription) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("memo_id", "user_id").
			Unique(),
		index.Fields("user_id"),
	}
}

// Edges of the Subscription.
func (Subscription) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("subscriber", User.Type).
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
