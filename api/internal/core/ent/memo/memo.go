// Code generated by ent, DO NOT EDIT.

package memo

import (
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/google/uuid"
)

const (
	// Label holds the string label denoting the memo type in the database.
	Label = "memo"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldOwnerID holds the string denoting the owner_id field in the database.
	FieldOwnerID = "owner_id"
	// FieldTitle holds the string denoting the title field in the database.
	FieldTitle = "title"
	// FieldContent holds the string denoting the content field in the database.
	FieldContent = "content"
	// FieldIsPublished holds the string denoting the is_published field in the database.
	FieldIsPublished = "is_published"
	// FieldVersion holds the string denoting the version field in the database.
	FieldVersion = "version"
	// FieldCreateTime holds the string denoting the create_time field in the database.
	FieldCreateTime = "create_time"
	// FieldUpdateTime holds the string denoting the update_time field in the database.
	FieldUpdateTime = "update_time"
	// EdgeOwner holds the string denoting the owner edge name in mutations.
	EdgeOwner = "owner"
	// EdgeTags holds the string denoting the tags edge name in mutations.
	EdgeTags = "tags"
	// EdgeSubscribers holds the string denoting the subscribers edge name in mutations.
	EdgeSubscribers = "subscribers"
	// EdgeCollaborators holds the string denoting the collaborators edge name in mutations.
	EdgeCollaborators = "collaborators"
	// EdgeSubscriptions holds the string denoting the subscriptions edge name in mutations.
	EdgeSubscriptions = "subscriptions"
	// EdgeCollaborations holds the string denoting the collaborations edge name in mutations.
	EdgeCollaborations = "collaborations"
	// Table holds the table name of the memo in the database.
	Table = "memos"
	// OwnerTable is the table that holds the owner relation/edge.
	OwnerTable = "memos"
	// OwnerInverseTable is the table name for the User entity.
	// It exists in this package in order to avoid circular dependency with the "user" package.
	OwnerInverseTable = "users"
	// OwnerColumn is the table column denoting the owner relation/edge.
	OwnerColumn = "owner_id"
	// TagsTable is the table that holds the tags relation/edge. The primary key declared below.
	TagsTable = "memo_tags"
	// TagsInverseTable is the table name for the Tag entity.
	// It exists in this package in order to avoid circular dependency with the "tag" package.
	TagsInverseTable = "tags"
	// SubscribersTable is the table that holds the subscribers relation/edge. The primary key declared below.
	SubscribersTable = "subscriptions"
	// SubscribersInverseTable is the table name for the User entity.
	// It exists in this package in order to avoid circular dependency with the "user" package.
	SubscribersInverseTable = "users"
	// CollaboratorsTable is the table that holds the collaborators relation/edge. The primary key declared below.
	CollaboratorsTable = "collaborations"
	// CollaboratorsInverseTable is the table name for the User entity.
	// It exists in this package in order to avoid circular dependency with the "user" package.
	CollaboratorsInverseTable = "users"
	// SubscriptionsTable is the table that holds the subscriptions relation/edge.
	SubscriptionsTable = "subscriptions"
	// SubscriptionsInverseTable is the table name for the Subscription entity.
	// It exists in this package in order to avoid circular dependency with the "subscription" package.
	SubscriptionsInverseTable = "subscriptions"
	// SubscriptionsColumn is the table column denoting the subscriptions relation/edge.
	SubscriptionsColumn = "memo_id"
	// CollaborationsTable is the table that holds the collaborations relation/edge.
	CollaborationsTable = "collaborations"
	// CollaborationsInverseTable is the table name for the Collaboration entity.
	// It exists in this package in order to avoid circular dependency with the "collaboration" package.
	CollaborationsInverseTable = "collaborations"
	// CollaborationsColumn is the table column denoting the collaborations relation/edge.
	CollaborationsColumn = "memo_id"
)

// Columns holds all SQL columns for memo fields.
var Columns = []string{
	FieldID,
	FieldOwnerID,
	FieldTitle,
	FieldContent,
	FieldIsPublished,
	FieldVersion,
	FieldCreateTime,
	FieldUpdateTime,
}

var (
	// TagsPrimaryKey and TagsColumn2 are the table columns denoting the
	// primary key for the tags relation (M2M).
	TagsPrimaryKey = []string{"memo_id", "tag_id"}
	// SubscribersPrimaryKey and SubscribersColumn2 are the table columns denoting the
	// primary key for the subscribers relation (M2M).
	SubscribersPrimaryKey = []string{"user_id", "memo_id"}
	// CollaboratorsPrimaryKey and CollaboratorsColumn2 are the table columns denoting the
	// primary key for the collaborators relation (M2M).
	CollaboratorsPrimaryKey = []string{"user_id", "memo_id"}
)

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
	// TitleValidator is a validator for the "title" field. It is called by the builders before save.
	TitleValidator func(string) error
	// ContentValidator is a validator for the "content" field. It is called by the builders before save.
	ContentValidator func(string) error
	// DefaultIsPublished holds the default value on creation for the "is_published" field.
	DefaultIsPublished bool
	// DefaultVersion holds the default value on creation for the "version" field.
	DefaultVersion int
	// DefaultCreateTime holds the default value on creation for the "create_time" field.
	DefaultCreateTime func() time.Time
	// DefaultUpdateTime holds the default value on creation for the "update_time" field.
	DefaultUpdateTime func() time.Time
	// DefaultID holds the default value on creation for the "id" field.
	DefaultID func() uuid.UUID
)

// OrderOption defines the ordering options for the Memo queries.
type OrderOption func(*sql.Selector)

// ByID orders the results by the id field.
func ByID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldID, opts...).ToFunc()
}

// ByOwnerID orders the results by the owner_id field.
func ByOwnerID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldOwnerID, opts...).ToFunc()
}

// ByTitle orders the results by the title field.
func ByTitle(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldTitle, opts...).ToFunc()
}

// ByContent orders the results by the content field.
func ByContent(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldContent, opts...).ToFunc()
}

// ByIsPublished orders the results by the is_published field.
func ByIsPublished(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldIsPublished, opts...).ToFunc()
}

// ByVersion orders the results by the version field.
func ByVersion(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldVersion, opts...).ToFunc()
}

// ByCreateTime orders the results by the create_time field.
func ByCreateTime(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldCreateTime, opts...).ToFunc()
}

// ByUpdateTime orders the results by the update_time field.
func ByUpdateTime(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldUpdateTime, opts...).ToFunc()
}

// ByOwnerField orders the results by owner field.
func ByOwnerField(field string, opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newOwnerStep(), sql.OrderByField(field, opts...))
	}
}

// ByTagsCount orders the results by tags count.
func ByTagsCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newTagsStep(), opts...)
	}
}

// ByTags orders the results by tags terms.
func ByTags(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newTagsStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}

// BySubscribersCount orders the results by subscribers count.
func BySubscribersCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newSubscribersStep(), opts...)
	}
}

// BySubscribers orders the results by subscribers terms.
func BySubscribers(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newSubscribersStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}

// ByCollaboratorsCount orders the results by collaborators count.
func ByCollaboratorsCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newCollaboratorsStep(), opts...)
	}
}

// ByCollaborators orders the results by collaborators terms.
func ByCollaborators(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newCollaboratorsStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}

// BySubscriptionsCount orders the results by subscriptions count.
func BySubscriptionsCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newSubscriptionsStep(), opts...)
	}
}

// BySubscriptions orders the results by subscriptions terms.
func BySubscriptions(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newSubscriptionsStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}

// ByCollaborationsCount orders the results by collaborations count.
func ByCollaborationsCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newCollaborationsStep(), opts...)
	}
}

// ByCollaborations orders the results by collaborations terms.
func ByCollaborations(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newCollaborationsStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}
func newOwnerStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(OwnerInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2O, true, OwnerTable, OwnerColumn),
	)
}
func newTagsStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(TagsInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2M, false, TagsTable, TagsPrimaryKey...),
	)
}
func newSubscribersStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(SubscribersInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2M, true, SubscribersTable, SubscribersPrimaryKey...),
	)
}
func newCollaboratorsStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(CollaboratorsInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2M, true, CollaboratorsTable, CollaboratorsPrimaryKey...),
	)
}
func newSubscriptionsStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(SubscriptionsInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.O2M, true, SubscriptionsTable, SubscriptionsColumn),
	)
}
func newCollaborationsStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(CollaborationsInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.O2M, true, CollaborationsTable, CollaborationsColumn),
	)
}
