// Code generated by ent, DO NOT EDIT.

package migrate

import (
	"entgo.io/ent/dialect/sql/schema"
	"entgo.io/ent/schema/field"
)

var (
	// CollaborationsColumns holds the columns for the "collaborations" table.
	CollaborationsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "approved", Type: field.TypeBool, Default: false},
		{Name: "create_time", Type: field.TypeTime},
		{Name: "update_time", Type: field.TypeTime},
		{Name: "user_id", Type: field.TypeUUID},
		{Name: "memo_id", Type: field.TypeUUID},
	}
	// CollaborationsTable holds the schema information for the "collaborations" table.
	CollaborationsTable = &schema.Table{
		Name:       "collaborations",
		Columns:    CollaborationsColumns,
		PrimaryKey: []*schema.Column{CollaborationsColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "collaborations_users_collaborator",
				Columns:    []*schema.Column{CollaborationsColumns[4]},
				RefColumns: []*schema.Column{UsersColumns[0]},
				OnDelete:   schema.Cascade,
			},
			{
				Symbol:     "collaborations_memos_memo",
				Columns:    []*schema.Column{CollaborationsColumns[5]},
				RefColumns: []*schema.Column{MemosColumns[0]},
				OnDelete:   schema.Cascade,
			},
		},
		Indexes: []*schema.Index{
			{
				Name:    "collaboration_memo_id_user_id",
				Unique:  true,
				Columns: []*schema.Column{CollaborationsColumns[5], CollaborationsColumns[4]},
			},
			{
				Name:    "collaboration_user_id",
				Unique:  false,
				Columns: []*schema.Column{CollaborationsColumns[4]},
			},
		},
	}
	// MemosColumns holds the columns for the "memos" table.
	MemosColumns = []*schema.Column{
		{Name: "id", Type: field.TypeUUID},
		{Name: "title", Type: field.TypeString, Size: 512},
		{Name: "content", Type: field.TypeString, Size: 20000},
		{Name: "is_published", Type: field.TypeBool, Default: false},
		{Name: "version", Type: field.TypeInt, Default: 0},
		{Name: "create_time", Type: field.TypeTime},
		{Name: "update_time", Type: field.TypeTime},
		{Name: "owner_id", Type: field.TypeUUID},
	}
	// MemosTable holds the schema information for the "memos" table.
	MemosTable = &schema.Table{
		Name:       "memos",
		Columns:    MemosColumns,
		PrimaryKey: []*schema.Column{MemosColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "memos_users_memos",
				Columns:    []*schema.Column{MemosColumns[7]},
				RefColumns: []*schema.Column{UsersColumns[0]},
				OnDelete:   schema.Cascade,
			},
		},
		Indexes: []*schema.Index{
			{
				Name:    "memo_owner_id",
				Unique:  false,
				Columns: []*schema.Column{MemosColumns[7]},
			},
		},
	}
	// SubscriptionsColumns holds the columns for the "subscriptions" table.
	SubscriptionsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "create_time", Type: field.TypeTime},
		{Name: "user_id", Type: field.TypeUUID},
		{Name: "memo_id", Type: field.TypeUUID},
	}
	// SubscriptionsTable holds the schema information for the "subscriptions" table.
	SubscriptionsTable = &schema.Table{
		Name:       "subscriptions",
		Columns:    SubscriptionsColumns,
		PrimaryKey: []*schema.Column{SubscriptionsColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "subscriptions_users_subscriber",
				Columns:    []*schema.Column{SubscriptionsColumns[2]},
				RefColumns: []*schema.Column{UsersColumns[0]},
				OnDelete:   schema.Cascade,
			},
			{
				Symbol:     "subscriptions_memos_memo",
				Columns:    []*schema.Column{SubscriptionsColumns[3]},
				RefColumns: []*schema.Column{MemosColumns[0]},
				OnDelete:   schema.Cascade,
			},
		},
		Indexes: []*schema.Index{
			{
				Name:    "subscription_memo_id_user_id",
				Unique:  true,
				Columns: []*schema.Column{SubscriptionsColumns[3], SubscriptionsColumns[2]},
			},
			{
				Name:    "subscription_user_id",
				Unique:  false,
				Columns: []*schema.Column{SubscriptionsColumns[2]},
			},
		},
	}
	// TagsColumns holds the columns for the "tags" table.
	TagsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "create_time", Type: field.TypeTime},
		{Name: "update_time", Type: field.TypeTime},
		{Name: "name", Type: field.TypeString, Unique: true, Size: 80},
	}
	// TagsTable holds the schema information for the "tags" table.
	TagsTable = &schema.Table{
		Name:       "tags",
		Columns:    TagsColumns,
		PrimaryKey: []*schema.Column{TagsColumns[0]},
	}
	// UsersColumns holds the columns for the "users" table.
	UsersColumns = []*schema.Column{
		{Name: "id", Type: field.TypeUUID},
		{Name: "create_time", Type: field.TypeTime},
		{Name: "update_time", Type: field.TypeTime},
		{Name: "email", Type: field.TypeString, Unique: true, Size: 256},
		{Name: "user_name", Type: field.TypeString, Size: 800},
		{Name: "given_name", Type: field.TypeString, Nullable: true, Size: 800},
		{Name: "family_name", Type: field.TypeString, Nullable: true, Size: 800},
		{Name: "photo_url", Type: field.TypeString, Nullable: true},
		{Name: "type", Type: field.TypeEnum, Enums: []string{"client", "operator"}, Default: "client"},
	}
	// UsersTable holds the schema information for the "users" table.
	UsersTable = &schema.Table{
		Name:       "users",
		Columns:    UsersColumns,
		PrimaryKey: []*schema.Column{UsersColumns[0]},
	}
	// MemoTagsColumns holds the columns for the "memo_tags" table.
	MemoTagsColumns = []*schema.Column{
		{Name: "memo_id", Type: field.TypeUUID},
		{Name: "tag_id", Type: field.TypeInt},
	}
	// MemoTagsTable holds the schema information for the "memo_tags" table.
	MemoTagsTable = &schema.Table{
		Name:       "memo_tags",
		Columns:    MemoTagsColumns,
		PrimaryKey: []*schema.Column{MemoTagsColumns[0], MemoTagsColumns[1]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "memo_tags_memo_id",
				Columns:    []*schema.Column{MemoTagsColumns[0]},
				RefColumns: []*schema.Column{MemosColumns[0]},
				OnDelete:   schema.Cascade,
			},
			{
				Symbol:     "memo_tags_tag_id",
				Columns:    []*schema.Column{MemoTagsColumns[1]},
				RefColumns: []*schema.Column{TagsColumns[0]},
				OnDelete:   schema.Cascade,
			},
		},
	}
	// Tables holds all the tables in the schema.
	Tables = []*schema.Table{
		CollaborationsTable,
		MemosTable,
		SubscriptionsTable,
		TagsTable,
		UsersTable,
		MemoTagsTable,
	}
)

func init() {
	CollaborationsTable.ForeignKeys[0].RefTable = UsersTable
	CollaborationsTable.ForeignKeys[1].RefTable = MemosTable
	MemosTable.ForeignKeys[0].RefTable = UsersTable
	SubscriptionsTable.ForeignKeys[0].RefTable = UsersTable
	SubscriptionsTable.ForeignKeys[1].RefTable = MemosTable
	MemoTagsTable.ForeignKeys[0].RefTable = MemosTable
	MemoTagsTable.ForeignKeys[1].RefTable = TagsTable
}
