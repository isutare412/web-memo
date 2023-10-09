// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
	"github.com/isutare412/web-memo/api/internal/core/ent/memo"
	"github.com/isutare412/web-memo/api/internal/core/ent/predicate"
	"github.com/isutare412/web-memo/api/internal/core/ent/tag"
	"github.com/isutare412/web-memo/api/internal/core/ent/user"
)

// MemoUpdate is the builder for updating Memo entities.
type MemoUpdate struct {
	config
	hooks    []Hook
	mutation *MemoMutation
}

// Where appends a list predicates to the MemoUpdate builder.
func (mu *MemoUpdate) Where(ps ...predicate.Memo) *MemoUpdate {
	mu.mutation.Where(ps...)
	return mu
}

// SetUpdateTime sets the "update_time" field.
func (mu *MemoUpdate) SetUpdateTime(t time.Time) *MemoUpdate {
	mu.mutation.SetUpdateTime(t)
	return mu
}

// SetOwnerID sets the "owner_id" field.
func (mu *MemoUpdate) SetOwnerID(u uuid.UUID) *MemoUpdate {
	mu.mutation.SetOwnerID(u)
	return mu
}

// SetTitle sets the "title" field.
func (mu *MemoUpdate) SetTitle(s string) *MemoUpdate {
	mu.mutation.SetTitle(s)
	return mu
}

// SetContent sets the "content" field.
func (mu *MemoUpdate) SetContent(s string) *MemoUpdate {
	mu.mutation.SetContent(s)
	return mu
}

// SetOwner sets the "owner" edge to the User entity.
func (mu *MemoUpdate) SetOwner(u *User) *MemoUpdate {
	return mu.SetOwnerID(u.ID)
}

// AddTagIDs adds the "tags" edge to the Tag entity by IDs.
func (mu *MemoUpdate) AddTagIDs(ids ...int) *MemoUpdate {
	mu.mutation.AddTagIDs(ids...)
	return mu
}

// AddTags adds the "tags" edges to the Tag entity.
func (mu *MemoUpdate) AddTags(t ...*Tag) *MemoUpdate {
	ids := make([]int, len(t))
	for i := range t {
		ids[i] = t[i].ID
	}
	return mu.AddTagIDs(ids...)
}

// Mutation returns the MemoMutation object of the builder.
func (mu *MemoUpdate) Mutation() *MemoMutation {
	return mu.mutation
}

// ClearOwner clears the "owner" edge to the User entity.
func (mu *MemoUpdate) ClearOwner() *MemoUpdate {
	mu.mutation.ClearOwner()
	return mu
}

// ClearTags clears all "tags" edges to the Tag entity.
func (mu *MemoUpdate) ClearTags() *MemoUpdate {
	mu.mutation.ClearTags()
	return mu
}

// RemoveTagIDs removes the "tags" edge to Tag entities by IDs.
func (mu *MemoUpdate) RemoveTagIDs(ids ...int) *MemoUpdate {
	mu.mutation.RemoveTagIDs(ids...)
	return mu
}

// RemoveTags removes "tags" edges to Tag entities.
func (mu *MemoUpdate) RemoveTags(t ...*Tag) *MemoUpdate {
	ids := make([]int, len(t))
	for i := range t {
		ids[i] = t[i].ID
	}
	return mu.RemoveTagIDs(ids...)
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (mu *MemoUpdate) Save(ctx context.Context) (int, error) {
	mu.defaults()
	return withHooks(ctx, mu.sqlSave, mu.mutation, mu.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (mu *MemoUpdate) SaveX(ctx context.Context) int {
	affected, err := mu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (mu *MemoUpdate) Exec(ctx context.Context) error {
	_, err := mu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (mu *MemoUpdate) ExecX(ctx context.Context) {
	if err := mu.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (mu *MemoUpdate) defaults() {
	if _, ok := mu.mutation.UpdateTime(); !ok {
		v := memo.UpdateDefaultUpdateTime()
		mu.mutation.SetUpdateTime(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (mu *MemoUpdate) check() error {
	if v, ok := mu.mutation.Title(); ok {
		if err := memo.TitleValidator(v); err != nil {
			return &ValidationError{Name: "title", err: fmt.Errorf(`ent: validator failed for field "Memo.title": %w`, err)}
		}
	}
	if v, ok := mu.mutation.Content(); ok {
		if err := memo.ContentValidator(v); err != nil {
			return &ValidationError{Name: "content", err: fmt.Errorf(`ent: validator failed for field "Memo.content": %w`, err)}
		}
	}
	if _, ok := mu.mutation.OwnerID(); mu.mutation.OwnerCleared() && !ok {
		return errors.New(`ent: clearing a required unique edge "Memo.owner"`)
	}
	return nil
}

func (mu *MemoUpdate) sqlSave(ctx context.Context) (n int, err error) {
	if err := mu.check(); err != nil {
		return n, err
	}
	_spec := sqlgraph.NewUpdateSpec(memo.Table, memo.Columns, sqlgraph.NewFieldSpec(memo.FieldID, field.TypeUUID))
	if ps := mu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := mu.mutation.UpdateTime(); ok {
		_spec.SetField(memo.FieldUpdateTime, field.TypeTime, value)
	}
	if value, ok := mu.mutation.Title(); ok {
		_spec.SetField(memo.FieldTitle, field.TypeString, value)
	}
	if value, ok := mu.mutation.Content(); ok {
		_spec.SetField(memo.FieldContent, field.TypeString, value)
	}
	if mu.mutation.OwnerCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   memo.OwnerTable,
			Columns: []string{memo.OwnerColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(user.FieldID, field.TypeUUID),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := mu.mutation.OwnerIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   memo.OwnerTable,
			Columns: []string{memo.OwnerColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(user.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if mu.mutation.TagsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   memo.TagsTable,
			Columns: memo.TagsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(tag.FieldID, field.TypeInt),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := mu.mutation.RemovedTagsIDs(); len(nodes) > 0 && !mu.mutation.TagsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   memo.TagsTable,
			Columns: memo.TagsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(tag.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := mu.mutation.TagsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   memo.TagsTable,
			Columns: memo.TagsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(tag.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, mu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{memo.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	mu.mutation.done = true
	return n, nil
}

// MemoUpdateOne is the builder for updating a single Memo entity.
type MemoUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *MemoMutation
}

// SetUpdateTime sets the "update_time" field.
func (muo *MemoUpdateOne) SetUpdateTime(t time.Time) *MemoUpdateOne {
	muo.mutation.SetUpdateTime(t)
	return muo
}

// SetOwnerID sets the "owner_id" field.
func (muo *MemoUpdateOne) SetOwnerID(u uuid.UUID) *MemoUpdateOne {
	muo.mutation.SetOwnerID(u)
	return muo
}

// SetTitle sets the "title" field.
func (muo *MemoUpdateOne) SetTitle(s string) *MemoUpdateOne {
	muo.mutation.SetTitle(s)
	return muo
}

// SetContent sets the "content" field.
func (muo *MemoUpdateOne) SetContent(s string) *MemoUpdateOne {
	muo.mutation.SetContent(s)
	return muo
}

// SetOwner sets the "owner" edge to the User entity.
func (muo *MemoUpdateOne) SetOwner(u *User) *MemoUpdateOne {
	return muo.SetOwnerID(u.ID)
}

// AddTagIDs adds the "tags" edge to the Tag entity by IDs.
func (muo *MemoUpdateOne) AddTagIDs(ids ...int) *MemoUpdateOne {
	muo.mutation.AddTagIDs(ids...)
	return muo
}

// AddTags adds the "tags" edges to the Tag entity.
func (muo *MemoUpdateOne) AddTags(t ...*Tag) *MemoUpdateOne {
	ids := make([]int, len(t))
	for i := range t {
		ids[i] = t[i].ID
	}
	return muo.AddTagIDs(ids...)
}

// Mutation returns the MemoMutation object of the builder.
func (muo *MemoUpdateOne) Mutation() *MemoMutation {
	return muo.mutation
}

// ClearOwner clears the "owner" edge to the User entity.
func (muo *MemoUpdateOne) ClearOwner() *MemoUpdateOne {
	muo.mutation.ClearOwner()
	return muo
}

// ClearTags clears all "tags" edges to the Tag entity.
func (muo *MemoUpdateOne) ClearTags() *MemoUpdateOne {
	muo.mutation.ClearTags()
	return muo
}

// RemoveTagIDs removes the "tags" edge to Tag entities by IDs.
func (muo *MemoUpdateOne) RemoveTagIDs(ids ...int) *MemoUpdateOne {
	muo.mutation.RemoveTagIDs(ids...)
	return muo
}

// RemoveTags removes "tags" edges to Tag entities.
func (muo *MemoUpdateOne) RemoveTags(t ...*Tag) *MemoUpdateOne {
	ids := make([]int, len(t))
	for i := range t {
		ids[i] = t[i].ID
	}
	return muo.RemoveTagIDs(ids...)
}

// Where appends a list predicates to the MemoUpdate builder.
func (muo *MemoUpdateOne) Where(ps ...predicate.Memo) *MemoUpdateOne {
	muo.mutation.Where(ps...)
	return muo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (muo *MemoUpdateOne) Select(field string, fields ...string) *MemoUpdateOne {
	muo.fields = append([]string{field}, fields...)
	return muo
}

// Save executes the query and returns the updated Memo entity.
func (muo *MemoUpdateOne) Save(ctx context.Context) (*Memo, error) {
	muo.defaults()
	return withHooks(ctx, muo.sqlSave, muo.mutation, muo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (muo *MemoUpdateOne) SaveX(ctx context.Context) *Memo {
	node, err := muo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (muo *MemoUpdateOne) Exec(ctx context.Context) error {
	_, err := muo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (muo *MemoUpdateOne) ExecX(ctx context.Context) {
	if err := muo.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (muo *MemoUpdateOne) defaults() {
	if _, ok := muo.mutation.UpdateTime(); !ok {
		v := memo.UpdateDefaultUpdateTime()
		muo.mutation.SetUpdateTime(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (muo *MemoUpdateOne) check() error {
	if v, ok := muo.mutation.Title(); ok {
		if err := memo.TitleValidator(v); err != nil {
			return &ValidationError{Name: "title", err: fmt.Errorf(`ent: validator failed for field "Memo.title": %w`, err)}
		}
	}
	if v, ok := muo.mutation.Content(); ok {
		if err := memo.ContentValidator(v); err != nil {
			return &ValidationError{Name: "content", err: fmt.Errorf(`ent: validator failed for field "Memo.content": %w`, err)}
		}
	}
	if _, ok := muo.mutation.OwnerID(); muo.mutation.OwnerCleared() && !ok {
		return errors.New(`ent: clearing a required unique edge "Memo.owner"`)
	}
	return nil
}

func (muo *MemoUpdateOne) sqlSave(ctx context.Context) (_node *Memo, err error) {
	if err := muo.check(); err != nil {
		return _node, err
	}
	_spec := sqlgraph.NewUpdateSpec(memo.Table, memo.Columns, sqlgraph.NewFieldSpec(memo.FieldID, field.TypeUUID))
	id, ok := muo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "Memo.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := muo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, memo.FieldID)
		for _, f := range fields {
			if !memo.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != memo.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := muo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := muo.mutation.UpdateTime(); ok {
		_spec.SetField(memo.FieldUpdateTime, field.TypeTime, value)
	}
	if value, ok := muo.mutation.Title(); ok {
		_spec.SetField(memo.FieldTitle, field.TypeString, value)
	}
	if value, ok := muo.mutation.Content(); ok {
		_spec.SetField(memo.FieldContent, field.TypeString, value)
	}
	if muo.mutation.OwnerCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   memo.OwnerTable,
			Columns: []string{memo.OwnerColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(user.FieldID, field.TypeUUID),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := muo.mutation.OwnerIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   memo.OwnerTable,
			Columns: []string{memo.OwnerColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(user.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if muo.mutation.TagsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   memo.TagsTable,
			Columns: memo.TagsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(tag.FieldID, field.TypeInt),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := muo.mutation.RemovedTagsIDs(); len(nodes) > 0 && !muo.mutation.TagsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   memo.TagsTable,
			Columns: memo.TagsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(tag.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := muo.mutation.TagsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   memo.TagsTable,
			Columns: memo.TagsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(tag.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_node = &Memo{config: muo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, muo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{memo.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	muo.mutation.done = true
	return _node, nil
}
