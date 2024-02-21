// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"time"

	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
	"github.com/isutare412/web-memo/api/internal/core/ent/memo"
	"github.com/isutare412/web-memo/api/internal/core/ent/subscription"
	"github.com/isutare412/web-memo/api/internal/core/ent/tag"
	"github.com/isutare412/web-memo/api/internal/core/ent/user"
)

// MemoCreate is the builder for creating a Memo entity.
type MemoCreate struct {
	config
	mutation *MemoMutation
	hooks    []Hook
	conflict []sql.ConflictOption
}

// SetCreateTime sets the "create_time" field.
func (mc *MemoCreate) SetCreateTime(t time.Time) *MemoCreate {
	mc.mutation.SetCreateTime(t)
	return mc
}

// SetNillableCreateTime sets the "create_time" field if the given value is not nil.
func (mc *MemoCreate) SetNillableCreateTime(t *time.Time) *MemoCreate {
	if t != nil {
		mc.SetCreateTime(*t)
	}
	return mc
}

// SetUpdateTime sets the "update_time" field.
func (mc *MemoCreate) SetUpdateTime(t time.Time) *MemoCreate {
	mc.mutation.SetUpdateTime(t)
	return mc
}

// SetNillableUpdateTime sets the "update_time" field if the given value is not nil.
func (mc *MemoCreate) SetNillableUpdateTime(t *time.Time) *MemoCreate {
	if t != nil {
		mc.SetUpdateTime(*t)
	}
	return mc
}

// SetOwnerID sets the "owner_id" field.
func (mc *MemoCreate) SetOwnerID(u uuid.UUID) *MemoCreate {
	mc.mutation.SetOwnerID(u)
	return mc
}

// SetTitle sets the "title" field.
func (mc *MemoCreate) SetTitle(s string) *MemoCreate {
	mc.mutation.SetTitle(s)
	return mc
}

// SetContent sets the "content" field.
func (mc *MemoCreate) SetContent(s string) *MemoCreate {
	mc.mutation.SetContent(s)
	return mc
}

// SetIsPublished sets the "is_published" field.
func (mc *MemoCreate) SetIsPublished(b bool) *MemoCreate {
	mc.mutation.SetIsPublished(b)
	return mc
}

// SetNillableIsPublished sets the "is_published" field if the given value is not nil.
func (mc *MemoCreate) SetNillableIsPublished(b *bool) *MemoCreate {
	if b != nil {
		mc.SetIsPublished(*b)
	}
	return mc
}

// SetID sets the "id" field.
func (mc *MemoCreate) SetID(u uuid.UUID) *MemoCreate {
	mc.mutation.SetID(u)
	return mc
}

// SetNillableID sets the "id" field if the given value is not nil.
func (mc *MemoCreate) SetNillableID(u *uuid.UUID) *MemoCreate {
	if u != nil {
		mc.SetID(*u)
	}
	return mc
}

// SetOwner sets the "owner" edge to the User entity.
func (mc *MemoCreate) SetOwner(u *User) *MemoCreate {
	return mc.SetOwnerID(u.ID)
}

// AddTagIDs adds the "tags" edge to the Tag entity by IDs.
func (mc *MemoCreate) AddTagIDs(ids ...int) *MemoCreate {
	mc.mutation.AddTagIDs(ids...)
	return mc
}

// AddTags adds the "tags" edges to the Tag entity.
func (mc *MemoCreate) AddTags(t ...*Tag) *MemoCreate {
	ids := make([]int, len(t))
	for i := range t {
		ids[i] = t[i].ID
	}
	return mc.AddTagIDs(ids...)
}

// AddSubscriberIDs adds the "subscribers" edge to the User entity by IDs.
func (mc *MemoCreate) AddSubscriberIDs(ids ...uuid.UUID) *MemoCreate {
	mc.mutation.AddSubscriberIDs(ids...)
	return mc
}

// AddSubscribers adds the "subscribers" edges to the User entity.
func (mc *MemoCreate) AddSubscribers(u ...*User) *MemoCreate {
	ids := make([]uuid.UUID, len(u))
	for i := range u {
		ids[i] = u[i].ID
	}
	return mc.AddSubscriberIDs(ids...)
}

// AddSubscriptionIDs adds the "subscriptions" edge to the Subscription entity by IDs.
func (mc *MemoCreate) AddSubscriptionIDs(ids ...int) *MemoCreate {
	mc.mutation.AddSubscriptionIDs(ids...)
	return mc
}

// AddSubscriptions adds the "subscriptions" edges to the Subscription entity.
func (mc *MemoCreate) AddSubscriptions(s ...*Subscription) *MemoCreate {
	ids := make([]int, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return mc.AddSubscriptionIDs(ids...)
}

// Mutation returns the MemoMutation object of the builder.
func (mc *MemoCreate) Mutation() *MemoMutation {
	return mc.mutation
}

// Save creates the Memo in the database.
func (mc *MemoCreate) Save(ctx context.Context) (*Memo, error) {
	mc.defaults()
	return withHooks(ctx, mc.sqlSave, mc.mutation, mc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (mc *MemoCreate) SaveX(ctx context.Context) *Memo {
	v, err := mc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (mc *MemoCreate) Exec(ctx context.Context) error {
	_, err := mc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (mc *MemoCreate) ExecX(ctx context.Context) {
	if err := mc.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (mc *MemoCreate) defaults() {
	if _, ok := mc.mutation.CreateTime(); !ok {
		v := memo.DefaultCreateTime()
		mc.mutation.SetCreateTime(v)
	}
	if _, ok := mc.mutation.UpdateTime(); !ok {
		v := memo.DefaultUpdateTime()
		mc.mutation.SetUpdateTime(v)
	}
	if _, ok := mc.mutation.IsPublished(); !ok {
		v := memo.DefaultIsPublished
		mc.mutation.SetIsPublished(v)
	}
	if _, ok := mc.mutation.ID(); !ok {
		v := memo.DefaultID()
		mc.mutation.SetID(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (mc *MemoCreate) check() error {
	if _, ok := mc.mutation.CreateTime(); !ok {
		return &ValidationError{Name: "create_time", err: errors.New(`ent: missing required field "Memo.create_time"`)}
	}
	if _, ok := mc.mutation.UpdateTime(); !ok {
		return &ValidationError{Name: "update_time", err: errors.New(`ent: missing required field "Memo.update_time"`)}
	}
	if _, ok := mc.mutation.OwnerID(); !ok {
		return &ValidationError{Name: "owner_id", err: errors.New(`ent: missing required field "Memo.owner_id"`)}
	}
	if _, ok := mc.mutation.Title(); !ok {
		return &ValidationError{Name: "title", err: errors.New(`ent: missing required field "Memo.title"`)}
	}
	if v, ok := mc.mutation.Title(); ok {
		if err := memo.TitleValidator(v); err != nil {
			return &ValidationError{Name: "title", err: fmt.Errorf(`ent: validator failed for field "Memo.title": %w`, err)}
		}
	}
	if _, ok := mc.mutation.Content(); !ok {
		return &ValidationError{Name: "content", err: errors.New(`ent: missing required field "Memo.content"`)}
	}
	if v, ok := mc.mutation.Content(); ok {
		if err := memo.ContentValidator(v); err != nil {
			return &ValidationError{Name: "content", err: fmt.Errorf(`ent: validator failed for field "Memo.content": %w`, err)}
		}
	}
	if _, ok := mc.mutation.IsPublished(); !ok {
		return &ValidationError{Name: "is_published", err: errors.New(`ent: missing required field "Memo.is_published"`)}
	}
	if _, ok := mc.mutation.OwnerID(); !ok {
		return &ValidationError{Name: "owner", err: errors.New(`ent: missing required edge "Memo.owner"`)}
	}
	return nil
}

func (mc *MemoCreate) sqlSave(ctx context.Context) (*Memo, error) {
	if err := mc.check(); err != nil {
		return nil, err
	}
	_node, _spec := mc.createSpec()
	if err := sqlgraph.CreateNode(ctx, mc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	if _spec.ID.Value != nil {
		if id, ok := _spec.ID.Value.(*uuid.UUID); ok {
			_node.ID = *id
		} else if err := _node.ID.Scan(_spec.ID.Value); err != nil {
			return nil, err
		}
	}
	mc.mutation.id = &_node.ID
	mc.mutation.done = true
	return _node, nil
}

func (mc *MemoCreate) createSpec() (*Memo, *sqlgraph.CreateSpec) {
	var (
		_node = &Memo{config: mc.config}
		_spec = sqlgraph.NewCreateSpec(memo.Table, sqlgraph.NewFieldSpec(memo.FieldID, field.TypeUUID))
	)
	_spec.OnConflict = mc.conflict
	if id, ok := mc.mutation.ID(); ok {
		_node.ID = id
		_spec.ID.Value = &id
	}
	if value, ok := mc.mutation.CreateTime(); ok {
		_spec.SetField(memo.FieldCreateTime, field.TypeTime, value)
		_node.CreateTime = value
	}
	if value, ok := mc.mutation.UpdateTime(); ok {
		_spec.SetField(memo.FieldUpdateTime, field.TypeTime, value)
		_node.UpdateTime = value
	}
	if value, ok := mc.mutation.Title(); ok {
		_spec.SetField(memo.FieldTitle, field.TypeString, value)
		_node.Title = value
	}
	if value, ok := mc.mutation.Content(); ok {
		_spec.SetField(memo.FieldContent, field.TypeString, value)
		_node.Content = value
	}
	if value, ok := mc.mutation.IsPublished(); ok {
		_spec.SetField(memo.FieldIsPublished, field.TypeBool, value)
		_node.IsPublished = value
	}
	if nodes := mc.mutation.OwnerIDs(); len(nodes) > 0 {
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
		_node.OwnerID = nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := mc.mutation.TagsIDs(); len(nodes) > 0 {
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
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := mc.mutation.SubscribersIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   memo.SubscribersTable,
			Columns: memo.SubscribersPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(user.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		createE := &SubscriptionCreate{config: mc.config, mutation: newSubscriptionMutation(mc.config, OpCreate)}
		createE.defaults()
		_, specE := createE.createSpec()
		edge.Target.Fields = specE.Fields
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := mc.mutation.SubscriptionsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   memo.SubscriptionsTable,
			Columns: []string{memo.SubscriptionsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(subscription.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	return _node, _spec
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.Memo.Create().
//		SetCreateTime(v).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.MemoUpsert) {
//			SetCreateTime(v+v).
//		}).
//		Exec(ctx)
func (mc *MemoCreate) OnConflict(opts ...sql.ConflictOption) *MemoUpsertOne {
	mc.conflict = opts
	return &MemoUpsertOne{
		create: mc,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.Memo.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (mc *MemoCreate) OnConflictColumns(columns ...string) *MemoUpsertOne {
	mc.conflict = append(mc.conflict, sql.ConflictColumns(columns...))
	return &MemoUpsertOne{
		create: mc,
	}
}

type (
	// MemoUpsertOne is the builder for "upsert"-ing
	//  one Memo node.
	MemoUpsertOne struct {
		create *MemoCreate
	}

	// MemoUpsert is the "OnConflict" setter.
	MemoUpsert struct {
		*sql.UpdateSet
	}
)

// SetUpdateTime sets the "update_time" field.
func (u *MemoUpsert) SetUpdateTime(v time.Time) *MemoUpsert {
	u.Set(memo.FieldUpdateTime, v)
	return u
}

// UpdateUpdateTime sets the "update_time" field to the value that was provided on create.
func (u *MemoUpsert) UpdateUpdateTime() *MemoUpsert {
	u.SetExcluded(memo.FieldUpdateTime)
	return u
}

// SetOwnerID sets the "owner_id" field.
func (u *MemoUpsert) SetOwnerID(v uuid.UUID) *MemoUpsert {
	u.Set(memo.FieldOwnerID, v)
	return u
}

// UpdateOwnerID sets the "owner_id" field to the value that was provided on create.
func (u *MemoUpsert) UpdateOwnerID() *MemoUpsert {
	u.SetExcluded(memo.FieldOwnerID)
	return u
}

// SetTitle sets the "title" field.
func (u *MemoUpsert) SetTitle(v string) *MemoUpsert {
	u.Set(memo.FieldTitle, v)
	return u
}

// UpdateTitle sets the "title" field to the value that was provided on create.
func (u *MemoUpsert) UpdateTitle() *MemoUpsert {
	u.SetExcluded(memo.FieldTitle)
	return u
}

// SetContent sets the "content" field.
func (u *MemoUpsert) SetContent(v string) *MemoUpsert {
	u.Set(memo.FieldContent, v)
	return u
}

// UpdateContent sets the "content" field to the value that was provided on create.
func (u *MemoUpsert) UpdateContent() *MemoUpsert {
	u.SetExcluded(memo.FieldContent)
	return u
}

// SetIsPublished sets the "is_published" field.
func (u *MemoUpsert) SetIsPublished(v bool) *MemoUpsert {
	u.Set(memo.FieldIsPublished, v)
	return u
}

// UpdateIsPublished sets the "is_published" field to the value that was provided on create.
func (u *MemoUpsert) UpdateIsPublished() *MemoUpsert {
	u.SetExcluded(memo.FieldIsPublished)
	return u
}

// UpdateNewValues updates the mutable fields using the new values that were set on create except the ID field.
// Using this option is equivalent to using:
//
//	client.Memo.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//			sql.ResolveWith(func(u *sql.UpdateSet) {
//				u.SetIgnore(memo.FieldID)
//			}),
//		).
//		Exec(ctx)
func (u *MemoUpsertOne) UpdateNewValues() *MemoUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(s *sql.UpdateSet) {
		if _, exists := u.create.mutation.ID(); exists {
			s.SetIgnore(memo.FieldID)
		}
		if _, exists := u.create.mutation.CreateTime(); exists {
			s.SetIgnore(memo.FieldCreateTime)
		}
	}))
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.Memo.Create().
//	    OnConflict(sql.ResolveWithIgnore()).
//	    Exec(ctx)
func (u *MemoUpsertOne) Ignore() *MemoUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *MemoUpsertOne) DoNothing() *MemoUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the MemoCreate.OnConflict
// documentation for more info.
func (u *MemoUpsertOne) Update(set func(*MemoUpsert)) *MemoUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&MemoUpsert{UpdateSet: update})
	}))
	return u
}

// SetUpdateTime sets the "update_time" field.
func (u *MemoUpsertOne) SetUpdateTime(v time.Time) *MemoUpsertOne {
	return u.Update(func(s *MemoUpsert) {
		s.SetUpdateTime(v)
	})
}

// UpdateUpdateTime sets the "update_time" field to the value that was provided on create.
func (u *MemoUpsertOne) UpdateUpdateTime() *MemoUpsertOne {
	return u.Update(func(s *MemoUpsert) {
		s.UpdateUpdateTime()
	})
}

// SetOwnerID sets the "owner_id" field.
func (u *MemoUpsertOne) SetOwnerID(v uuid.UUID) *MemoUpsertOne {
	return u.Update(func(s *MemoUpsert) {
		s.SetOwnerID(v)
	})
}

// UpdateOwnerID sets the "owner_id" field to the value that was provided on create.
func (u *MemoUpsertOne) UpdateOwnerID() *MemoUpsertOne {
	return u.Update(func(s *MemoUpsert) {
		s.UpdateOwnerID()
	})
}

// SetTitle sets the "title" field.
func (u *MemoUpsertOne) SetTitle(v string) *MemoUpsertOne {
	return u.Update(func(s *MemoUpsert) {
		s.SetTitle(v)
	})
}

// UpdateTitle sets the "title" field to the value that was provided on create.
func (u *MemoUpsertOne) UpdateTitle() *MemoUpsertOne {
	return u.Update(func(s *MemoUpsert) {
		s.UpdateTitle()
	})
}

// SetContent sets the "content" field.
func (u *MemoUpsertOne) SetContent(v string) *MemoUpsertOne {
	return u.Update(func(s *MemoUpsert) {
		s.SetContent(v)
	})
}

// UpdateContent sets the "content" field to the value that was provided on create.
func (u *MemoUpsertOne) UpdateContent() *MemoUpsertOne {
	return u.Update(func(s *MemoUpsert) {
		s.UpdateContent()
	})
}

// SetIsPublished sets the "is_published" field.
func (u *MemoUpsertOne) SetIsPublished(v bool) *MemoUpsertOne {
	return u.Update(func(s *MemoUpsert) {
		s.SetIsPublished(v)
	})
}

// UpdateIsPublished sets the "is_published" field to the value that was provided on create.
func (u *MemoUpsertOne) UpdateIsPublished() *MemoUpsertOne {
	return u.Update(func(s *MemoUpsert) {
		s.UpdateIsPublished()
	})
}

// Exec executes the query.
func (u *MemoUpsertOne) Exec(ctx context.Context) error {
	if len(u.create.conflict) == 0 {
		return errors.New("ent: missing options for MemoCreate.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *MemoUpsertOne) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}

// Exec executes the UPSERT query and returns the inserted/updated ID.
func (u *MemoUpsertOne) ID(ctx context.Context) (id uuid.UUID, err error) {
	if u.create.driver.Dialect() == dialect.MySQL {
		// In case of "ON CONFLICT", there is no way to get back non-numeric ID
		// fields from the database since MySQL does not support the RETURNING clause.
		return id, errors.New("ent: MemoUpsertOne.ID is not supported by MySQL driver. Use MemoUpsertOne.Exec instead")
	}
	node, err := u.create.Save(ctx)
	if err != nil {
		return id, err
	}
	return node.ID, nil
}

// IDX is like ID, but panics if an error occurs.
func (u *MemoUpsertOne) IDX(ctx context.Context) uuid.UUID {
	id, err := u.ID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// MemoCreateBulk is the builder for creating many Memo entities in bulk.
type MemoCreateBulk struct {
	config
	err      error
	builders []*MemoCreate
	conflict []sql.ConflictOption
}

// Save creates the Memo entities in the database.
func (mcb *MemoCreateBulk) Save(ctx context.Context) ([]*Memo, error) {
	if mcb.err != nil {
		return nil, mcb.err
	}
	specs := make([]*sqlgraph.CreateSpec, len(mcb.builders))
	nodes := make([]*Memo, len(mcb.builders))
	mutators := make([]Mutator, len(mcb.builders))
	for i := range mcb.builders {
		func(i int, root context.Context) {
			builder := mcb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*MemoMutation)
				if !ok {
					return nil, fmt.Errorf("unexpected mutation type %T", m)
				}
				if err := builder.check(); err != nil {
					return nil, err
				}
				builder.mutation = mutation
				var err error
				nodes[i], specs[i] = builder.createSpec()
				if i < len(mutators)-1 {
					_, err = mutators[i+1].Mutate(root, mcb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					spec.OnConflict = mcb.conflict
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, mcb.driver, spec); err != nil {
						if sqlgraph.IsConstraintError(err) {
							err = &ConstraintError{msg: err.Error(), wrap: err}
						}
					}
				}
				if err != nil {
					return nil, err
				}
				mutation.id = &nodes[i].ID
				mutation.done = true
				return nodes[i], nil
			})
			for i := len(builder.hooks) - 1; i >= 0; i-- {
				mut = builder.hooks[i](mut)
			}
			mutators[i] = mut
		}(i, ctx)
	}
	if len(mutators) > 0 {
		if _, err := mutators[0].Mutate(ctx, mcb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (mcb *MemoCreateBulk) SaveX(ctx context.Context) []*Memo {
	v, err := mcb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (mcb *MemoCreateBulk) Exec(ctx context.Context) error {
	_, err := mcb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (mcb *MemoCreateBulk) ExecX(ctx context.Context) {
	if err := mcb.Exec(ctx); err != nil {
		panic(err)
	}
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.Memo.CreateBulk(builders...).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.MemoUpsert) {
//			SetCreateTime(v+v).
//		}).
//		Exec(ctx)
func (mcb *MemoCreateBulk) OnConflict(opts ...sql.ConflictOption) *MemoUpsertBulk {
	mcb.conflict = opts
	return &MemoUpsertBulk{
		create: mcb,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.Memo.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (mcb *MemoCreateBulk) OnConflictColumns(columns ...string) *MemoUpsertBulk {
	mcb.conflict = append(mcb.conflict, sql.ConflictColumns(columns...))
	return &MemoUpsertBulk{
		create: mcb,
	}
}

// MemoUpsertBulk is the builder for "upsert"-ing
// a bulk of Memo nodes.
type MemoUpsertBulk struct {
	create *MemoCreateBulk
}

// UpdateNewValues updates the mutable fields using the new values that
// were set on create. Using this option is equivalent to using:
//
//	client.Memo.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//			sql.ResolveWith(func(u *sql.UpdateSet) {
//				u.SetIgnore(memo.FieldID)
//			}),
//		).
//		Exec(ctx)
func (u *MemoUpsertBulk) UpdateNewValues() *MemoUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(s *sql.UpdateSet) {
		for _, b := range u.create.builders {
			if _, exists := b.mutation.ID(); exists {
				s.SetIgnore(memo.FieldID)
			}
			if _, exists := b.mutation.CreateTime(); exists {
				s.SetIgnore(memo.FieldCreateTime)
			}
		}
	}))
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.Memo.Create().
//		OnConflict(sql.ResolveWithIgnore()).
//		Exec(ctx)
func (u *MemoUpsertBulk) Ignore() *MemoUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *MemoUpsertBulk) DoNothing() *MemoUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the MemoCreateBulk.OnConflict
// documentation for more info.
func (u *MemoUpsertBulk) Update(set func(*MemoUpsert)) *MemoUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&MemoUpsert{UpdateSet: update})
	}))
	return u
}

// SetUpdateTime sets the "update_time" field.
func (u *MemoUpsertBulk) SetUpdateTime(v time.Time) *MemoUpsertBulk {
	return u.Update(func(s *MemoUpsert) {
		s.SetUpdateTime(v)
	})
}

// UpdateUpdateTime sets the "update_time" field to the value that was provided on create.
func (u *MemoUpsertBulk) UpdateUpdateTime() *MemoUpsertBulk {
	return u.Update(func(s *MemoUpsert) {
		s.UpdateUpdateTime()
	})
}

// SetOwnerID sets the "owner_id" field.
func (u *MemoUpsertBulk) SetOwnerID(v uuid.UUID) *MemoUpsertBulk {
	return u.Update(func(s *MemoUpsert) {
		s.SetOwnerID(v)
	})
}

// UpdateOwnerID sets the "owner_id" field to the value that was provided on create.
func (u *MemoUpsertBulk) UpdateOwnerID() *MemoUpsertBulk {
	return u.Update(func(s *MemoUpsert) {
		s.UpdateOwnerID()
	})
}

// SetTitle sets the "title" field.
func (u *MemoUpsertBulk) SetTitle(v string) *MemoUpsertBulk {
	return u.Update(func(s *MemoUpsert) {
		s.SetTitle(v)
	})
}

// UpdateTitle sets the "title" field to the value that was provided on create.
func (u *MemoUpsertBulk) UpdateTitle() *MemoUpsertBulk {
	return u.Update(func(s *MemoUpsert) {
		s.UpdateTitle()
	})
}

// SetContent sets the "content" field.
func (u *MemoUpsertBulk) SetContent(v string) *MemoUpsertBulk {
	return u.Update(func(s *MemoUpsert) {
		s.SetContent(v)
	})
}

// UpdateContent sets the "content" field to the value that was provided on create.
func (u *MemoUpsertBulk) UpdateContent() *MemoUpsertBulk {
	return u.Update(func(s *MemoUpsert) {
		s.UpdateContent()
	})
}

// SetIsPublished sets the "is_published" field.
func (u *MemoUpsertBulk) SetIsPublished(v bool) *MemoUpsertBulk {
	return u.Update(func(s *MemoUpsert) {
		s.SetIsPublished(v)
	})
}

// UpdateIsPublished sets the "is_published" field to the value that was provided on create.
func (u *MemoUpsertBulk) UpdateIsPublished() *MemoUpsertBulk {
	return u.Update(func(s *MemoUpsert) {
		s.UpdateIsPublished()
	})
}

// Exec executes the query.
func (u *MemoUpsertBulk) Exec(ctx context.Context) error {
	if u.create.err != nil {
		return u.create.err
	}
	for i, b := range u.create.builders {
		if len(b.conflict) != 0 {
			return fmt.Errorf("ent: OnConflict was set for builder %d. Set it on the MemoCreateBulk instead", i)
		}
	}
	if len(u.create.conflict) == 0 {
		return errors.New("ent: missing options for MemoCreateBulk.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *MemoUpsertBulk) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}
