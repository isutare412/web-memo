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
	"github.com/isutare412/web-memo/api/internal/core/ent/subscription"
	"github.com/isutare412/web-memo/api/internal/core/ent/user"
)

// SubscriptionCreate is the builder for creating a Subscription entity.
type SubscriptionCreate struct {
	config
	mutation *SubscriptionMutation
	hooks    []Hook
	conflict []sql.ConflictOption
}

// SetUserID sets the "user_id" field.
func (sc *SubscriptionCreate) SetUserID(u uuid.UUID) *SubscriptionCreate {
	sc.mutation.SetUserID(u)
	return sc
}

// SetMemoID sets the "memo_id" field.
func (sc *SubscriptionCreate) SetMemoID(u uuid.UUID) *SubscriptionCreate {
	sc.mutation.SetMemoID(u)
	return sc
}

// SetCreateTime sets the "create_time" field.
func (sc *SubscriptionCreate) SetCreateTime(t time.Time) *SubscriptionCreate {
	sc.mutation.SetCreateTime(t)
	return sc
}

// SetNillableCreateTime sets the "create_time" field if the given value is not nil.
func (sc *SubscriptionCreate) SetNillableCreateTime(t *time.Time) *SubscriptionCreate {
	if t != nil {
		sc.SetCreateTime(*t)
	}
	return sc
}

// SetSubscriberID sets the "subscriber" edge to the User entity by ID.
func (sc *SubscriptionCreate) SetSubscriberID(id uuid.UUID) *SubscriptionCreate {
	sc.mutation.SetSubscriberID(id)
	return sc
}

// SetSubscriber sets the "subscriber" edge to the User entity.
func (sc *SubscriptionCreate) SetSubscriber(u *User) *SubscriptionCreate {
	return sc.SetSubscriberID(u.ID)
}

// SetMemo sets the "memo" edge to the Memo entity.
func (sc *SubscriptionCreate) SetMemo(m *Memo) *SubscriptionCreate {
	return sc.SetMemoID(m.ID)
}

// Mutation returns the SubscriptionMutation object of the builder.
func (sc *SubscriptionCreate) Mutation() *SubscriptionMutation {
	return sc.mutation
}

// Save creates the Subscription in the database.
func (sc *SubscriptionCreate) Save(ctx context.Context) (*Subscription, error) {
	sc.defaults()
	return withHooks(ctx, sc.sqlSave, sc.mutation, sc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (sc *SubscriptionCreate) SaveX(ctx context.Context) *Subscription {
	v, err := sc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (sc *SubscriptionCreate) Exec(ctx context.Context) error {
	_, err := sc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (sc *SubscriptionCreate) ExecX(ctx context.Context) {
	if err := sc.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (sc *SubscriptionCreate) defaults() {
	if _, ok := sc.mutation.CreateTime(); !ok {
		v := subscription.DefaultCreateTime()
		sc.mutation.SetCreateTime(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (sc *SubscriptionCreate) check() error {
	if _, ok := sc.mutation.UserID(); !ok {
		return &ValidationError{Name: "user_id", err: errors.New(`ent: missing required field "Subscription.user_id"`)}
	}
	if _, ok := sc.mutation.MemoID(); !ok {
		return &ValidationError{Name: "memo_id", err: errors.New(`ent: missing required field "Subscription.memo_id"`)}
	}
	if _, ok := sc.mutation.CreateTime(); !ok {
		return &ValidationError{Name: "create_time", err: errors.New(`ent: missing required field "Subscription.create_time"`)}
	}
	if _, ok := sc.mutation.SubscriberID(); !ok {
		return &ValidationError{Name: "subscriber", err: errors.New(`ent: missing required edge "Subscription.subscriber"`)}
	}
	if _, ok := sc.mutation.MemoID(); !ok {
		return &ValidationError{Name: "memo", err: errors.New(`ent: missing required edge "Subscription.memo"`)}
	}
	return nil
}

func (sc *SubscriptionCreate) sqlSave(ctx context.Context) (*Subscription, error) {
	if err := sc.check(); err != nil {
		return nil, err
	}
	_node, _spec := sc.createSpec()
	if err := sqlgraph.CreateNode(ctx, sc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	_node.ID = int(id)
	sc.mutation.id = &_node.ID
	sc.mutation.done = true
	return _node, nil
}

func (sc *SubscriptionCreate) createSpec() (*Subscription, *sqlgraph.CreateSpec) {
	var (
		_node = &Subscription{config: sc.config}
		_spec = sqlgraph.NewCreateSpec(subscription.Table, sqlgraph.NewFieldSpec(subscription.FieldID, field.TypeInt))
	)
	_spec.OnConflict = sc.conflict
	if value, ok := sc.mutation.CreateTime(); ok {
		_spec.SetField(subscription.FieldCreateTime, field.TypeTime, value)
		_node.CreateTime = value
	}
	if nodes := sc.mutation.SubscriberIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   subscription.SubscriberTable,
			Columns: []string{subscription.SubscriberColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(user.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_node.UserID = nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := sc.mutation.MemoIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   subscription.MemoTable,
			Columns: []string{subscription.MemoColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(memo.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_node.MemoID = nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	return _node, _spec
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.Subscription.Create().
//		SetUserID(v).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.SubscriptionUpsert) {
//			SetUserID(v+v).
//		}).
//		Exec(ctx)
func (sc *SubscriptionCreate) OnConflict(opts ...sql.ConflictOption) *SubscriptionUpsertOne {
	sc.conflict = opts
	return &SubscriptionUpsertOne{
		create: sc,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.Subscription.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (sc *SubscriptionCreate) OnConflictColumns(columns ...string) *SubscriptionUpsertOne {
	sc.conflict = append(sc.conflict, sql.ConflictColumns(columns...))
	return &SubscriptionUpsertOne{
		create: sc,
	}
}

type (
	// SubscriptionUpsertOne is the builder for "upsert"-ing
	//  one Subscription node.
	SubscriptionUpsertOne struct {
		create *SubscriptionCreate
	}

	// SubscriptionUpsert is the "OnConflict" setter.
	SubscriptionUpsert struct {
		*sql.UpdateSet
	}
)

// SetUserID sets the "user_id" field.
func (u *SubscriptionUpsert) SetUserID(v uuid.UUID) *SubscriptionUpsert {
	u.Set(subscription.FieldUserID, v)
	return u
}

// UpdateUserID sets the "user_id" field to the value that was provided on create.
func (u *SubscriptionUpsert) UpdateUserID() *SubscriptionUpsert {
	u.SetExcluded(subscription.FieldUserID)
	return u
}

// SetMemoID sets the "memo_id" field.
func (u *SubscriptionUpsert) SetMemoID(v uuid.UUID) *SubscriptionUpsert {
	u.Set(subscription.FieldMemoID, v)
	return u
}

// UpdateMemoID sets the "memo_id" field to the value that was provided on create.
func (u *SubscriptionUpsert) UpdateMemoID() *SubscriptionUpsert {
	u.SetExcluded(subscription.FieldMemoID)
	return u
}

// UpdateNewValues updates the mutable fields using the new values that were set on create.
// Using this option is equivalent to using:
//
//	client.Subscription.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//		).
//		Exec(ctx)
func (u *SubscriptionUpsertOne) UpdateNewValues() *SubscriptionUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(s *sql.UpdateSet) {
		if _, exists := u.create.mutation.CreateTime(); exists {
			s.SetIgnore(subscription.FieldCreateTime)
		}
	}))
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.Subscription.Create().
//	    OnConflict(sql.ResolveWithIgnore()).
//	    Exec(ctx)
func (u *SubscriptionUpsertOne) Ignore() *SubscriptionUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *SubscriptionUpsertOne) DoNothing() *SubscriptionUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the SubscriptionCreate.OnConflict
// documentation for more info.
func (u *SubscriptionUpsertOne) Update(set func(*SubscriptionUpsert)) *SubscriptionUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&SubscriptionUpsert{UpdateSet: update})
	}))
	return u
}

// SetUserID sets the "user_id" field.
func (u *SubscriptionUpsertOne) SetUserID(v uuid.UUID) *SubscriptionUpsertOne {
	return u.Update(func(s *SubscriptionUpsert) {
		s.SetUserID(v)
	})
}

// UpdateUserID sets the "user_id" field to the value that was provided on create.
func (u *SubscriptionUpsertOne) UpdateUserID() *SubscriptionUpsertOne {
	return u.Update(func(s *SubscriptionUpsert) {
		s.UpdateUserID()
	})
}

// SetMemoID sets the "memo_id" field.
func (u *SubscriptionUpsertOne) SetMemoID(v uuid.UUID) *SubscriptionUpsertOne {
	return u.Update(func(s *SubscriptionUpsert) {
		s.SetMemoID(v)
	})
}

// UpdateMemoID sets the "memo_id" field to the value that was provided on create.
func (u *SubscriptionUpsertOne) UpdateMemoID() *SubscriptionUpsertOne {
	return u.Update(func(s *SubscriptionUpsert) {
		s.UpdateMemoID()
	})
}

// Exec executes the query.
func (u *SubscriptionUpsertOne) Exec(ctx context.Context) error {
	if len(u.create.conflict) == 0 {
		return errors.New("ent: missing options for SubscriptionCreate.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *SubscriptionUpsertOne) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}

// Exec executes the UPSERT query and returns the inserted/updated ID.
func (u *SubscriptionUpsertOne) ID(ctx context.Context) (id int, err error) {
	node, err := u.create.Save(ctx)
	if err != nil {
		return id, err
	}
	return node.ID, nil
}

// IDX is like ID, but panics if an error occurs.
func (u *SubscriptionUpsertOne) IDX(ctx context.Context) int {
	id, err := u.ID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// SubscriptionCreateBulk is the builder for creating many Subscription entities in bulk.
type SubscriptionCreateBulk struct {
	config
	err      error
	builders []*SubscriptionCreate
	conflict []sql.ConflictOption
}

// Save creates the Subscription entities in the database.
func (scb *SubscriptionCreateBulk) Save(ctx context.Context) ([]*Subscription, error) {
	if scb.err != nil {
		return nil, scb.err
	}
	specs := make([]*sqlgraph.CreateSpec, len(scb.builders))
	nodes := make([]*Subscription, len(scb.builders))
	mutators := make([]Mutator, len(scb.builders))
	for i := range scb.builders {
		func(i int, root context.Context) {
			builder := scb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*SubscriptionMutation)
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
					_, err = mutators[i+1].Mutate(root, scb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					spec.OnConflict = scb.conflict
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, scb.driver, spec); err != nil {
						if sqlgraph.IsConstraintError(err) {
							err = &ConstraintError{msg: err.Error(), wrap: err}
						}
					}
				}
				if err != nil {
					return nil, err
				}
				mutation.id = &nodes[i].ID
				if specs[i].ID.Value != nil {
					id := specs[i].ID.Value.(int64)
					nodes[i].ID = int(id)
				}
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
		if _, err := mutators[0].Mutate(ctx, scb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (scb *SubscriptionCreateBulk) SaveX(ctx context.Context) []*Subscription {
	v, err := scb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (scb *SubscriptionCreateBulk) Exec(ctx context.Context) error {
	_, err := scb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (scb *SubscriptionCreateBulk) ExecX(ctx context.Context) {
	if err := scb.Exec(ctx); err != nil {
		panic(err)
	}
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.Subscription.CreateBulk(builders...).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.SubscriptionUpsert) {
//			SetUserID(v+v).
//		}).
//		Exec(ctx)
func (scb *SubscriptionCreateBulk) OnConflict(opts ...sql.ConflictOption) *SubscriptionUpsertBulk {
	scb.conflict = opts
	return &SubscriptionUpsertBulk{
		create: scb,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.Subscription.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (scb *SubscriptionCreateBulk) OnConflictColumns(columns ...string) *SubscriptionUpsertBulk {
	scb.conflict = append(scb.conflict, sql.ConflictColumns(columns...))
	return &SubscriptionUpsertBulk{
		create: scb,
	}
}

// SubscriptionUpsertBulk is the builder for "upsert"-ing
// a bulk of Subscription nodes.
type SubscriptionUpsertBulk struct {
	create *SubscriptionCreateBulk
}

// UpdateNewValues updates the mutable fields using the new values that
// were set on create. Using this option is equivalent to using:
//
//	client.Subscription.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//		).
//		Exec(ctx)
func (u *SubscriptionUpsertBulk) UpdateNewValues() *SubscriptionUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(s *sql.UpdateSet) {
		for _, b := range u.create.builders {
			if _, exists := b.mutation.CreateTime(); exists {
				s.SetIgnore(subscription.FieldCreateTime)
			}
		}
	}))
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.Subscription.Create().
//		OnConflict(sql.ResolveWithIgnore()).
//		Exec(ctx)
func (u *SubscriptionUpsertBulk) Ignore() *SubscriptionUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *SubscriptionUpsertBulk) DoNothing() *SubscriptionUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the SubscriptionCreateBulk.OnConflict
// documentation for more info.
func (u *SubscriptionUpsertBulk) Update(set func(*SubscriptionUpsert)) *SubscriptionUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&SubscriptionUpsert{UpdateSet: update})
	}))
	return u
}

// SetUserID sets the "user_id" field.
func (u *SubscriptionUpsertBulk) SetUserID(v uuid.UUID) *SubscriptionUpsertBulk {
	return u.Update(func(s *SubscriptionUpsert) {
		s.SetUserID(v)
	})
}

// UpdateUserID sets the "user_id" field to the value that was provided on create.
func (u *SubscriptionUpsertBulk) UpdateUserID() *SubscriptionUpsertBulk {
	return u.Update(func(s *SubscriptionUpsert) {
		s.UpdateUserID()
	})
}

// SetMemoID sets the "memo_id" field.
func (u *SubscriptionUpsertBulk) SetMemoID(v uuid.UUID) *SubscriptionUpsertBulk {
	return u.Update(func(s *SubscriptionUpsert) {
		s.SetMemoID(v)
	})
}

// UpdateMemoID sets the "memo_id" field to the value that was provided on create.
func (u *SubscriptionUpsertBulk) UpdateMemoID() *SubscriptionUpsertBulk {
	return u.Update(func(s *SubscriptionUpsert) {
		s.UpdateMemoID()
	})
}

// Exec executes the query.
func (u *SubscriptionUpsertBulk) Exec(ctx context.Context) error {
	if u.create.err != nil {
		return u.create.err
	}
	for i, b := range u.create.builders {
		if len(b.conflict) != 0 {
			return fmt.Errorf("ent: OnConflict was set for builder %d. Set it on the SubscriptionCreateBulk instead", i)
		}
	}
	if len(u.create.conflict) == 0 {
		return errors.New("ent: missing options for SubscriptionCreateBulk.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *SubscriptionUpsertBulk) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}
