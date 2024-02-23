package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
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
		field.UUID("owner_id", uuid.UUID{}),
		field.String("title").
			NotEmpty().
			MaxLen(512),
		field.String("content").
			MaxLen(20_000),
		field.Bool("is_published").
			Default(false),
		field.Time("create_time").
			Default(time.Now).
			Immutable(),
		field.Time("update_time").
			Default(time.Now),
	}
}

// Edges of the Memo.
func (Memo) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("owner", User.Type).
			Ref("memos").
			Required().
			Unique().
			Field("owner_id"),
		edge.To("tags", Tag.Type),
		edge.From("subscribers", User.Type).
			Ref("subscribing_memos").
			Through("subscriptions", Subscription.Type),
	}
}
