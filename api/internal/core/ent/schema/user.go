package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
	"github.com/google/uuid"

	"github.com/isutare412/web-memo/api/internal/core/enum"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New),
		field.String("email").
			NotEmpty().
			MaxLen(256).
			Unique(),
		field.String("user_name").
			NotEmpty().
			MaxLen(800),
		field.String("given_name").
			Optional().
			NotEmpty().
			MaxLen(800),
		field.String("family_name").
			Optional().
			NotEmpty().
			MaxLen(800),
		field.String("photo_url").
			Optional().
			NotEmpty(),
		field.Enum("type").
			GoType(enum.UserType("")).
			Default(string(enum.UserTypeClient)),
	}
}

// Mixin of the User.
func (User) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.Time{},
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("memos", Memo.Type).
			Annotations(entsql.OnDelete(entsql.Cascade)),
	}
}
