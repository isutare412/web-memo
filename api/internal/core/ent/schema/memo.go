package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
	"github.com/isutare412/web-memo/api/internal/core/enum"
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
			MaxLen(60_000),
		field.Enum("publish_state").
			GoType(enum.PublishState("")).
			Default(string(enum.PublishStatePrivate)),
		field.Int("version").
			Default(0),
		field.Bool("is_embedded").
			Default(false),
		field.Time("create_time").
			Default(time.Now).
			Immutable(),
		field.Time("update_time").
			Default(time.Now),
	}
}

// Indexes of the Memo.
func (Memo) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("owner_id"),
		index.Fields("is_embedded"),
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
		edge.From("collaborators", User.Type).
			Ref("collaborating_memos").
			Through("collaborations", Collaboration.Type),
	}
}
