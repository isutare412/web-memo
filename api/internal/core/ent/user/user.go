// Code generated by ent, DO NOT EDIT.

package user

import (
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/google/uuid"
)

const (
	// Label holds the string label denoting the user type in the database.
	Label = "user"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldCreateTime holds the string denoting the create_time field in the database.
	FieldCreateTime = "create_time"
	// FieldUpdateTime holds the string denoting the update_time field in the database.
	FieldUpdateTime = "update_time"
	// FieldEmail holds the string denoting the email field in the database.
	FieldEmail = "email"
	// FieldUserName holds the string denoting the user_name field in the database.
	FieldUserName = "user_name"
	// FieldGivenName holds the string denoting the given_name field in the database.
	FieldGivenName = "given_name"
	// FieldFamilyName holds the string denoting the family_name field in the database.
	FieldFamilyName = "family_name"
	// FieldPhotoURL holds the string denoting the photo_url field in the database.
	FieldPhotoURL = "photo_url"
	// EdgeMemos holds the string denoting the memos edge name in mutations.
	EdgeMemos = "memos"
	// Table holds the table name of the user in the database.
	Table = "users"
	// MemosTable is the table that holds the memos relation/edge.
	MemosTable = "memos"
	// MemosInverseTable is the table name for the Memo entity.
	// It exists in this package in order to avoid circular dependency with the "memo" package.
	MemosInverseTable = "memos"
	// MemosColumn is the table column denoting the memos relation/edge.
	MemosColumn = "user_memos"
)

// Columns holds all SQL columns for user fields.
var Columns = []string{
	FieldID,
	FieldCreateTime,
	FieldUpdateTime,
	FieldEmail,
	FieldUserName,
	FieldGivenName,
	FieldFamilyName,
	FieldPhotoURL,
}

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	return false
}

var (
	// DefaultCreateTime holds the default value on creation for the "create_time" field.
	DefaultCreateTime func() time.Time
	// DefaultUpdateTime holds the default value on creation for the "update_time" field.
	DefaultUpdateTime func() time.Time
	// UpdateDefaultUpdateTime holds the default value on update for the "update_time" field.
	UpdateDefaultUpdateTime func() time.Time
	// EmailValidator is a validator for the "email" field. It is called by the builders before save.
	EmailValidator func(string) error
	// UserNameValidator is a validator for the "user_name" field. It is called by the builders before save.
	UserNameValidator func(string) error
	// GivenNameValidator is a validator for the "given_name" field. It is called by the builders before save.
	GivenNameValidator func(string) error
	// FamilyNameValidator is a validator for the "family_name" field. It is called by the builders before save.
	FamilyNameValidator func(string) error
	// PhotoURLValidator is a validator for the "photo_url" field. It is called by the builders before save.
	PhotoURLValidator func(string) error
	// DefaultID holds the default value on creation for the "id" field.
	DefaultID func() uuid.UUID
)

// OrderOption defines the ordering options for the User queries.
type OrderOption func(*sql.Selector)

// ByID orders the results by the id field.
func ByID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldID, opts...).ToFunc()
}

// ByCreateTime orders the results by the create_time field.
func ByCreateTime(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldCreateTime, opts...).ToFunc()
}

// ByUpdateTime orders the results by the update_time field.
func ByUpdateTime(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldUpdateTime, opts...).ToFunc()
}

// ByEmail orders the results by the email field.
func ByEmail(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldEmail, opts...).ToFunc()
}

// ByUserName orders the results by the user_name field.
func ByUserName(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldUserName, opts...).ToFunc()
}

// ByGivenName orders the results by the given_name field.
func ByGivenName(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldGivenName, opts...).ToFunc()
}

// ByFamilyName orders the results by the family_name field.
func ByFamilyName(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldFamilyName, opts...).ToFunc()
}

// ByPhotoURL orders the results by the photo_url field.
func ByPhotoURL(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldPhotoURL, opts...).ToFunc()
}

// ByMemosCount orders the results by memos count.
func ByMemosCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newMemosStep(), opts...)
	}
}

// ByMemos orders the results by memos terms.
func ByMemos(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newMemosStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}
func newMemosStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(MemosInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, MemosTable, MemosColumn),
	)
}