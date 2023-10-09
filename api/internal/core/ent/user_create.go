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
	"github.com/isutare412/web-memo/api/internal/core/ent/user"
	"github.com/isutare412/web-memo/api/internal/core/enum"
)

// UserCreate is the builder for creating a User entity.
type UserCreate struct {
	config
	mutation *UserMutation
	hooks    []Hook
	conflict []sql.ConflictOption
}

// SetCreateTime sets the "create_time" field.
func (uc *UserCreate) SetCreateTime(t time.Time) *UserCreate {
	uc.mutation.SetCreateTime(t)
	return uc
}

// SetNillableCreateTime sets the "create_time" field if the given value is not nil.
func (uc *UserCreate) SetNillableCreateTime(t *time.Time) *UserCreate {
	if t != nil {
		uc.SetCreateTime(*t)
	}
	return uc
}

// SetUpdateTime sets the "update_time" field.
func (uc *UserCreate) SetUpdateTime(t time.Time) *UserCreate {
	uc.mutation.SetUpdateTime(t)
	return uc
}

// SetNillableUpdateTime sets the "update_time" field if the given value is not nil.
func (uc *UserCreate) SetNillableUpdateTime(t *time.Time) *UserCreate {
	if t != nil {
		uc.SetUpdateTime(*t)
	}
	return uc
}

// SetEmail sets the "email" field.
func (uc *UserCreate) SetEmail(s string) *UserCreate {
	uc.mutation.SetEmail(s)
	return uc
}

// SetUserName sets the "user_name" field.
func (uc *UserCreate) SetUserName(s string) *UserCreate {
	uc.mutation.SetUserName(s)
	return uc
}

// SetGivenName sets the "given_name" field.
func (uc *UserCreate) SetGivenName(s string) *UserCreate {
	uc.mutation.SetGivenName(s)
	return uc
}

// SetNillableGivenName sets the "given_name" field if the given value is not nil.
func (uc *UserCreate) SetNillableGivenName(s *string) *UserCreate {
	if s != nil {
		uc.SetGivenName(*s)
	}
	return uc
}

// SetFamilyName sets the "family_name" field.
func (uc *UserCreate) SetFamilyName(s string) *UserCreate {
	uc.mutation.SetFamilyName(s)
	return uc
}

// SetNillableFamilyName sets the "family_name" field if the given value is not nil.
func (uc *UserCreate) SetNillableFamilyName(s *string) *UserCreate {
	if s != nil {
		uc.SetFamilyName(*s)
	}
	return uc
}

// SetPhotoURL sets the "photo_url" field.
func (uc *UserCreate) SetPhotoURL(s string) *UserCreate {
	uc.mutation.SetPhotoURL(s)
	return uc
}

// SetNillablePhotoURL sets the "photo_url" field if the given value is not nil.
func (uc *UserCreate) SetNillablePhotoURL(s *string) *UserCreate {
	if s != nil {
		uc.SetPhotoURL(*s)
	}
	return uc
}

// SetType sets the "type" field.
func (uc *UserCreate) SetType(et enum.UserType) *UserCreate {
	uc.mutation.SetType(et)
	return uc
}

// SetNillableType sets the "type" field if the given value is not nil.
func (uc *UserCreate) SetNillableType(et *enum.UserType) *UserCreate {
	if et != nil {
		uc.SetType(*et)
	}
	return uc
}

// SetID sets the "id" field.
func (uc *UserCreate) SetID(u uuid.UUID) *UserCreate {
	uc.mutation.SetID(u)
	return uc
}

// SetNillableID sets the "id" field if the given value is not nil.
func (uc *UserCreate) SetNillableID(u *uuid.UUID) *UserCreate {
	if u != nil {
		uc.SetID(*u)
	}
	return uc
}

// AddMemoIDs adds the "memos" edge to the Memo entity by IDs.
func (uc *UserCreate) AddMemoIDs(ids ...uuid.UUID) *UserCreate {
	uc.mutation.AddMemoIDs(ids...)
	return uc
}

// AddMemos adds the "memos" edges to the Memo entity.
func (uc *UserCreate) AddMemos(m ...*Memo) *UserCreate {
	ids := make([]uuid.UUID, len(m))
	for i := range m {
		ids[i] = m[i].ID
	}
	return uc.AddMemoIDs(ids...)
}

// Mutation returns the UserMutation object of the builder.
func (uc *UserCreate) Mutation() *UserMutation {
	return uc.mutation
}

// Save creates the User in the database.
func (uc *UserCreate) Save(ctx context.Context) (*User, error) {
	uc.defaults()
	return withHooks(ctx, uc.sqlSave, uc.mutation, uc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (uc *UserCreate) SaveX(ctx context.Context) *User {
	v, err := uc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (uc *UserCreate) Exec(ctx context.Context) error {
	_, err := uc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (uc *UserCreate) ExecX(ctx context.Context) {
	if err := uc.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (uc *UserCreate) defaults() {
	if _, ok := uc.mutation.CreateTime(); !ok {
		v := user.DefaultCreateTime()
		uc.mutation.SetCreateTime(v)
	}
	if _, ok := uc.mutation.UpdateTime(); !ok {
		v := user.DefaultUpdateTime()
		uc.mutation.SetUpdateTime(v)
	}
	if _, ok := uc.mutation.GetType(); !ok {
		v := user.DefaultType
		uc.mutation.SetType(v)
	}
	if _, ok := uc.mutation.ID(); !ok {
		v := user.DefaultID()
		uc.mutation.SetID(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (uc *UserCreate) check() error {
	if _, ok := uc.mutation.CreateTime(); !ok {
		return &ValidationError{Name: "create_time", err: errors.New(`ent: missing required field "User.create_time"`)}
	}
	if _, ok := uc.mutation.UpdateTime(); !ok {
		return &ValidationError{Name: "update_time", err: errors.New(`ent: missing required field "User.update_time"`)}
	}
	if _, ok := uc.mutation.Email(); !ok {
		return &ValidationError{Name: "email", err: errors.New(`ent: missing required field "User.email"`)}
	}
	if v, ok := uc.mutation.Email(); ok {
		if err := user.EmailValidator(v); err != nil {
			return &ValidationError{Name: "email", err: fmt.Errorf(`ent: validator failed for field "User.email": %w`, err)}
		}
	}
	if _, ok := uc.mutation.UserName(); !ok {
		return &ValidationError{Name: "user_name", err: errors.New(`ent: missing required field "User.user_name"`)}
	}
	if v, ok := uc.mutation.UserName(); ok {
		if err := user.UserNameValidator(v); err != nil {
			return &ValidationError{Name: "user_name", err: fmt.Errorf(`ent: validator failed for field "User.user_name": %w`, err)}
		}
	}
	if v, ok := uc.mutation.GivenName(); ok {
		if err := user.GivenNameValidator(v); err != nil {
			return &ValidationError{Name: "given_name", err: fmt.Errorf(`ent: validator failed for field "User.given_name": %w`, err)}
		}
	}
	if v, ok := uc.mutation.FamilyName(); ok {
		if err := user.FamilyNameValidator(v); err != nil {
			return &ValidationError{Name: "family_name", err: fmt.Errorf(`ent: validator failed for field "User.family_name": %w`, err)}
		}
	}
	if v, ok := uc.mutation.PhotoURL(); ok {
		if err := user.PhotoURLValidator(v); err != nil {
			return &ValidationError{Name: "photo_url", err: fmt.Errorf(`ent: validator failed for field "User.photo_url": %w`, err)}
		}
	}
	if _, ok := uc.mutation.GetType(); !ok {
		return &ValidationError{Name: "type", err: errors.New(`ent: missing required field "User.type"`)}
	}
	if v, ok := uc.mutation.GetType(); ok {
		if err := user.TypeValidator(v); err != nil {
			return &ValidationError{Name: "type", err: fmt.Errorf(`ent: validator failed for field "User.type": %w`, err)}
		}
	}
	return nil
}

func (uc *UserCreate) sqlSave(ctx context.Context) (*User, error) {
	if err := uc.check(); err != nil {
		return nil, err
	}
	_node, _spec := uc.createSpec()
	if err := sqlgraph.CreateNode(ctx, uc.driver, _spec); err != nil {
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
	uc.mutation.id = &_node.ID
	uc.mutation.done = true
	return _node, nil
}

func (uc *UserCreate) createSpec() (*User, *sqlgraph.CreateSpec) {
	var (
		_node = &User{config: uc.config}
		_spec = sqlgraph.NewCreateSpec(user.Table, sqlgraph.NewFieldSpec(user.FieldID, field.TypeUUID))
	)
	_spec.OnConflict = uc.conflict
	if id, ok := uc.mutation.ID(); ok {
		_node.ID = id
		_spec.ID.Value = &id
	}
	if value, ok := uc.mutation.CreateTime(); ok {
		_spec.SetField(user.FieldCreateTime, field.TypeTime, value)
		_node.CreateTime = value
	}
	if value, ok := uc.mutation.UpdateTime(); ok {
		_spec.SetField(user.FieldUpdateTime, field.TypeTime, value)
		_node.UpdateTime = value
	}
	if value, ok := uc.mutation.Email(); ok {
		_spec.SetField(user.FieldEmail, field.TypeString, value)
		_node.Email = value
	}
	if value, ok := uc.mutation.UserName(); ok {
		_spec.SetField(user.FieldUserName, field.TypeString, value)
		_node.UserName = value
	}
	if value, ok := uc.mutation.GivenName(); ok {
		_spec.SetField(user.FieldGivenName, field.TypeString, value)
		_node.GivenName = value
	}
	if value, ok := uc.mutation.FamilyName(); ok {
		_spec.SetField(user.FieldFamilyName, field.TypeString, value)
		_node.FamilyName = value
	}
	if value, ok := uc.mutation.PhotoURL(); ok {
		_spec.SetField(user.FieldPhotoURL, field.TypeString, value)
		_node.PhotoURL = value
	}
	if value, ok := uc.mutation.GetType(); ok {
		_spec.SetField(user.FieldType, field.TypeEnum, value)
		_node.Type = value
	}
	if nodes := uc.mutation.MemosIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   user.MemosTable,
			Columns: []string{user.MemosColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(memo.FieldID, field.TypeUUID),
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
//	client.User.Create().
//		SetCreateTime(v).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.UserUpsert) {
//			SetCreateTime(v+v).
//		}).
//		Exec(ctx)
func (uc *UserCreate) OnConflict(opts ...sql.ConflictOption) *UserUpsertOne {
	uc.conflict = opts
	return &UserUpsertOne{
		create: uc,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.User.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (uc *UserCreate) OnConflictColumns(columns ...string) *UserUpsertOne {
	uc.conflict = append(uc.conflict, sql.ConflictColumns(columns...))
	return &UserUpsertOne{
		create: uc,
	}
}

type (
	// UserUpsertOne is the builder for "upsert"-ing
	//  one User node.
	UserUpsertOne struct {
		create *UserCreate
	}

	// UserUpsert is the "OnConflict" setter.
	UserUpsert struct {
		*sql.UpdateSet
	}
)

// SetUpdateTime sets the "update_time" field.
func (u *UserUpsert) SetUpdateTime(v time.Time) *UserUpsert {
	u.Set(user.FieldUpdateTime, v)
	return u
}

// UpdateUpdateTime sets the "update_time" field to the value that was provided on create.
func (u *UserUpsert) UpdateUpdateTime() *UserUpsert {
	u.SetExcluded(user.FieldUpdateTime)
	return u
}

// SetEmail sets the "email" field.
func (u *UserUpsert) SetEmail(v string) *UserUpsert {
	u.Set(user.FieldEmail, v)
	return u
}

// UpdateEmail sets the "email" field to the value that was provided on create.
func (u *UserUpsert) UpdateEmail() *UserUpsert {
	u.SetExcluded(user.FieldEmail)
	return u
}

// SetUserName sets the "user_name" field.
func (u *UserUpsert) SetUserName(v string) *UserUpsert {
	u.Set(user.FieldUserName, v)
	return u
}

// UpdateUserName sets the "user_name" field to the value that was provided on create.
func (u *UserUpsert) UpdateUserName() *UserUpsert {
	u.SetExcluded(user.FieldUserName)
	return u
}

// SetGivenName sets the "given_name" field.
func (u *UserUpsert) SetGivenName(v string) *UserUpsert {
	u.Set(user.FieldGivenName, v)
	return u
}

// UpdateGivenName sets the "given_name" field to the value that was provided on create.
func (u *UserUpsert) UpdateGivenName() *UserUpsert {
	u.SetExcluded(user.FieldGivenName)
	return u
}

// ClearGivenName clears the value of the "given_name" field.
func (u *UserUpsert) ClearGivenName() *UserUpsert {
	u.SetNull(user.FieldGivenName)
	return u
}

// SetFamilyName sets the "family_name" field.
func (u *UserUpsert) SetFamilyName(v string) *UserUpsert {
	u.Set(user.FieldFamilyName, v)
	return u
}

// UpdateFamilyName sets the "family_name" field to the value that was provided on create.
func (u *UserUpsert) UpdateFamilyName() *UserUpsert {
	u.SetExcluded(user.FieldFamilyName)
	return u
}

// ClearFamilyName clears the value of the "family_name" field.
func (u *UserUpsert) ClearFamilyName() *UserUpsert {
	u.SetNull(user.FieldFamilyName)
	return u
}

// SetPhotoURL sets the "photo_url" field.
func (u *UserUpsert) SetPhotoURL(v string) *UserUpsert {
	u.Set(user.FieldPhotoURL, v)
	return u
}

// UpdatePhotoURL sets the "photo_url" field to the value that was provided on create.
func (u *UserUpsert) UpdatePhotoURL() *UserUpsert {
	u.SetExcluded(user.FieldPhotoURL)
	return u
}

// ClearPhotoURL clears the value of the "photo_url" field.
func (u *UserUpsert) ClearPhotoURL() *UserUpsert {
	u.SetNull(user.FieldPhotoURL)
	return u
}

// SetType sets the "type" field.
func (u *UserUpsert) SetType(v enum.UserType) *UserUpsert {
	u.Set(user.FieldType, v)
	return u
}

// UpdateType sets the "type" field to the value that was provided on create.
func (u *UserUpsert) UpdateType() *UserUpsert {
	u.SetExcluded(user.FieldType)
	return u
}

// UpdateNewValues updates the mutable fields using the new values that were set on create except the ID field.
// Using this option is equivalent to using:
//
//	client.User.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//			sql.ResolveWith(func(u *sql.UpdateSet) {
//				u.SetIgnore(user.FieldID)
//			}),
//		).
//		Exec(ctx)
func (u *UserUpsertOne) UpdateNewValues() *UserUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(s *sql.UpdateSet) {
		if _, exists := u.create.mutation.ID(); exists {
			s.SetIgnore(user.FieldID)
		}
		if _, exists := u.create.mutation.CreateTime(); exists {
			s.SetIgnore(user.FieldCreateTime)
		}
	}))
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.User.Create().
//	    OnConflict(sql.ResolveWithIgnore()).
//	    Exec(ctx)
func (u *UserUpsertOne) Ignore() *UserUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *UserUpsertOne) DoNothing() *UserUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the UserCreate.OnConflict
// documentation for more info.
func (u *UserUpsertOne) Update(set func(*UserUpsert)) *UserUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&UserUpsert{UpdateSet: update})
	}))
	return u
}

// SetUpdateTime sets the "update_time" field.
func (u *UserUpsertOne) SetUpdateTime(v time.Time) *UserUpsertOne {
	return u.Update(func(s *UserUpsert) {
		s.SetUpdateTime(v)
	})
}

// UpdateUpdateTime sets the "update_time" field to the value that was provided on create.
func (u *UserUpsertOne) UpdateUpdateTime() *UserUpsertOne {
	return u.Update(func(s *UserUpsert) {
		s.UpdateUpdateTime()
	})
}

// SetEmail sets the "email" field.
func (u *UserUpsertOne) SetEmail(v string) *UserUpsertOne {
	return u.Update(func(s *UserUpsert) {
		s.SetEmail(v)
	})
}

// UpdateEmail sets the "email" field to the value that was provided on create.
func (u *UserUpsertOne) UpdateEmail() *UserUpsertOne {
	return u.Update(func(s *UserUpsert) {
		s.UpdateEmail()
	})
}

// SetUserName sets the "user_name" field.
func (u *UserUpsertOne) SetUserName(v string) *UserUpsertOne {
	return u.Update(func(s *UserUpsert) {
		s.SetUserName(v)
	})
}

// UpdateUserName sets the "user_name" field to the value that was provided on create.
func (u *UserUpsertOne) UpdateUserName() *UserUpsertOne {
	return u.Update(func(s *UserUpsert) {
		s.UpdateUserName()
	})
}

// SetGivenName sets the "given_name" field.
func (u *UserUpsertOne) SetGivenName(v string) *UserUpsertOne {
	return u.Update(func(s *UserUpsert) {
		s.SetGivenName(v)
	})
}

// UpdateGivenName sets the "given_name" field to the value that was provided on create.
func (u *UserUpsertOne) UpdateGivenName() *UserUpsertOne {
	return u.Update(func(s *UserUpsert) {
		s.UpdateGivenName()
	})
}

// ClearGivenName clears the value of the "given_name" field.
func (u *UserUpsertOne) ClearGivenName() *UserUpsertOne {
	return u.Update(func(s *UserUpsert) {
		s.ClearGivenName()
	})
}

// SetFamilyName sets the "family_name" field.
func (u *UserUpsertOne) SetFamilyName(v string) *UserUpsertOne {
	return u.Update(func(s *UserUpsert) {
		s.SetFamilyName(v)
	})
}

// UpdateFamilyName sets the "family_name" field to the value that was provided on create.
func (u *UserUpsertOne) UpdateFamilyName() *UserUpsertOne {
	return u.Update(func(s *UserUpsert) {
		s.UpdateFamilyName()
	})
}

// ClearFamilyName clears the value of the "family_name" field.
func (u *UserUpsertOne) ClearFamilyName() *UserUpsertOne {
	return u.Update(func(s *UserUpsert) {
		s.ClearFamilyName()
	})
}

// SetPhotoURL sets the "photo_url" field.
func (u *UserUpsertOne) SetPhotoURL(v string) *UserUpsertOne {
	return u.Update(func(s *UserUpsert) {
		s.SetPhotoURL(v)
	})
}

// UpdatePhotoURL sets the "photo_url" field to the value that was provided on create.
func (u *UserUpsertOne) UpdatePhotoURL() *UserUpsertOne {
	return u.Update(func(s *UserUpsert) {
		s.UpdatePhotoURL()
	})
}

// ClearPhotoURL clears the value of the "photo_url" field.
func (u *UserUpsertOne) ClearPhotoURL() *UserUpsertOne {
	return u.Update(func(s *UserUpsert) {
		s.ClearPhotoURL()
	})
}

// SetType sets the "type" field.
func (u *UserUpsertOne) SetType(v enum.UserType) *UserUpsertOne {
	return u.Update(func(s *UserUpsert) {
		s.SetType(v)
	})
}

// UpdateType sets the "type" field to the value that was provided on create.
func (u *UserUpsertOne) UpdateType() *UserUpsertOne {
	return u.Update(func(s *UserUpsert) {
		s.UpdateType()
	})
}

// Exec executes the query.
func (u *UserUpsertOne) Exec(ctx context.Context) error {
	if len(u.create.conflict) == 0 {
		return errors.New("ent: missing options for UserCreate.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *UserUpsertOne) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}

// Exec executes the UPSERT query and returns the inserted/updated ID.
func (u *UserUpsertOne) ID(ctx context.Context) (id uuid.UUID, err error) {
	if u.create.driver.Dialect() == dialect.MySQL {
		// In case of "ON CONFLICT", there is no way to get back non-numeric ID
		// fields from the database since MySQL does not support the RETURNING clause.
		return id, errors.New("ent: UserUpsertOne.ID is not supported by MySQL driver. Use UserUpsertOne.Exec instead")
	}
	node, err := u.create.Save(ctx)
	if err != nil {
		return id, err
	}
	return node.ID, nil
}

// IDX is like ID, but panics if an error occurs.
func (u *UserUpsertOne) IDX(ctx context.Context) uuid.UUID {
	id, err := u.ID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// UserCreateBulk is the builder for creating many User entities in bulk.
type UserCreateBulk struct {
	config
	err      error
	builders []*UserCreate
	conflict []sql.ConflictOption
}

// Save creates the User entities in the database.
func (ucb *UserCreateBulk) Save(ctx context.Context) ([]*User, error) {
	if ucb.err != nil {
		return nil, ucb.err
	}
	specs := make([]*sqlgraph.CreateSpec, len(ucb.builders))
	nodes := make([]*User, len(ucb.builders))
	mutators := make([]Mutator, len(ucb.builders))
	for i := range ucb.builders {
		func(i int, root context.Context) {
			builder := ucb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*UserMutation)
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
					_, err = mutators[i+1].Mutate(root, ucb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					spec.OnConflict = ucb.conflict
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, ucb.driver, spec); err != nil {
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
		if _, err := mutators[0].Mutate(ctx, ucb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (ucb *UserCreateBulk) SaveX(ctx context.Context) []*User {
	v, err := ucb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (ucb *UserCreateBulk) Exec(ctx context.Context) error {
	_, err := ucb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (ucb *UserCreateBulk) ExecX(ctx context.Context) {
	if err := ucb.Exec(ctx); err != nil {
		panic(err)
	}
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.User.CreateBulk(builders...).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.UserUpsert) {
//			SetCreateTime(v+v).
//		}).
//		Exec(ctx)
func (ucb *UserCreateBulk) OnConflict(opts ...sql.ConflictOption) *UserUpsertBulk {
	ucb.conflict = opts
	return &UserUpsertBulk{
		create: ucb,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.User.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (ucb *UserCreateBulk) OnConflictColumns(columns ...string) *UserUpsertBulk {
	ucb.conflict = append(ucb.conflict, sql.ConflictColumns(columns...))
	return &UserUpsertBulk{
		create: ucb,
	}
}

// UserUpsertBulk is the builder for "upsert"-ing
// a bulk of User nodes.
type UserUpsertBulk struct {
	create *UserCreateBulk
}

// UpdateNewValues updates the mutable fields using the new values that
// were set on create. Using this option is equivalent to using:
//
//	client.User.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//			sql.ResolveWith(func(u *sql.UpdateSet) {
//				u.SetIgnore(user.FieldID)
//			}),
//		).
//		Exec(ctx)
func (u *UserUpsertBulk) UpdateNewValues() *UserUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(s *sql.UpdateSet) {
		for _, b := range u.create.builders {
			if _, exists := b.mutation.ID(); exists {
				s.SetIgnore(user.FieldID)
			}
			if _, exists := b.mutation.CreateTime(); exists {
				s.SetIgnore(user.FieldCreateTime)
			}
		}
	}))
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.User.Create().
//		OnConflict(sql.ResolveWithIgnore()).
//		Exec(ctx)
func (u *UserUpsertBulk) Ignore() *UserUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *UserUpsertBulk) DoNothing() *UserUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the UserCreateBulk.OnConflict
// documentation for more info.
func (u *UserUpsertBulk) Update(set func(*UserUpsert)) *UserUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&UserUpsert{UpdateSet: update})
	}))
	return u
}

// SetUpdateTime sets the "update_time" field.
func (u *UserUpsertBulk) SetUpdateTime(v time.Time) *UserUpsertBulk {
	return u.Update(func(s *UserUpsert) {
		s.SetUpdateTime(v)
	})
}

// UpdateUpdateTime sets the "update_time" field to the value that was provided on create.
func (u *UserUpsertBulk) UpdateUpdateTime() *UserUpsertBulk {
	return u.Update(func(s *UserUpsert) {
		s.UpdateUpdateTime()
	})
}

// SetEmail sets the "email" field.
func (u *UserUpsertBulk) SetEmail(v string) *UserUpsertBulk {
	return u.Update(func(s *UserUpsert) {
		s.SetEmail(v)
	})
}

// UpdateEmail sets the "email" field to the value that was provided on create.
func (u *UserUpsertBulk) UpdateEmail() *UserUpsertBulk {
	return u.Update(func(s *UserUpsert) {
		s.UpdateEmail()
	})
}

// SetUserName sets the "user_name" field.
func (u *UserUpsertBulk) SetUserName(v string) *UserUpsertBulk {
	return u.Update(func(s *UserUpsert) {
		s.SetUserName(v)
	})
}

// UpdateUserName sets the "user_name" field to the value that was provided on create.
func (u *UserUpsertBulk) UpdateUserName() *UserUpsertBulk {
	return u.Update(func(s *UserUpsert) {
		s.UpdateUserName()
	})
}

// SetGivenName sets the "given_name" field.
func (u *UserUpsertBulk) SetGivenName(v string) *UserUpsertBulk {
	return u.Update(func(s *UserUpsert) {
		s.SetGivenName(v)
	})
}

// UpdateGivenName sets the "given_name" field to the value that was provided on create.
func (u *UserUpsertBulk) UpdateGivenName() *UserUpsertBulk {
	return u.Update(func(s *UserUpsert) {
		s.UpdateGivenName()
	})
}

// ClearGivenName clears the value of the "given_name" field.
func (u *UserUpsertBulk) ClearGivenName() *UserUpsertBulk {
	return u.Update(func(s *UserUpsert) {
		s.ClearGivenName()
	})
}

// SetFamilyName sets the "family_name" field.
func (u *UserUpsertBulk) SetFamilyName(v string) *UserUpsertBulk {
	return u.Update(func(s *UserUpsert) {
		s.SetFamilyName(v)
	})
}

// UpdateFamilyName sets the "family_name" field to the value that was provided on create.
func (u *UserUpsertBulk) UpdateFamilyName() *UserUpsertBulk {
	return u.Update(func(s *UserUpsert) {
		s.UpdateFamilyName()
	})
}

// ClearFamilyName clears the value of the "family_name" field.
func (u *UserUpsertBulk) ClearFamilyName() *UserUpsertBulk {
	return u.Update(func(s *UserUpsert) {
		s.ClearFamilyName()
	})
}

// SetPhotoURL sets the "photo_url" field.
func (u *UserUpsertBulk) SetPhotoURL(v string) *UserUpsertBulk {
	return u.Update(func(s *UserUpsert) {
		s.SetPhotoURL(v)
	})
}

// UpdatePhotoURL sets the "photo_url" field to the value that was provided on create.
func (u *UserUpsertBulk) UpdatePhotoURL() *UserUpsertBulk {
	return u.Update(func(s *UserUpsert) {
		s.UpdatePhotoURL()
	})
}

// ClearPhotoURL clears the value of the "photo_url" field.
func (u *UserUpsertBulk) ClearPhotoURL() *UserUpsertBulk {
	return u.Update(func(s *UserUpsert) {
		s.ClearPhotoURL()
	})
}

// SetType sets the "type" field.
func (u *UserUpsertBulk) SetType(v enum.UserType) *UserUpsertBulk {
	return u.Update(func(s *UserUpsert) {
		s.SetType(v)
	})
}

// UpdateType sets the "type" field to the value that was provided on create.
func (u *UserUpsertBulk) UpdateType() *UserUpsertBulk {
	return u.Update(func(s *UserUpsert) {
		s.UpdateType()
	})
}

// Exec executes the query.
func (u *UserUpsertBulk) Exec(ctx context.Context) error {
	if u.create.err != nil {
		return u.create.err
	}
	for i, b := range u.create.builders {
		if len(b.conflict) != 0 {
			return fmt.Errorf("ent: OnConflict was set for builder %d. Set it on the UserCreateBulk instead", i)
		}
	}
	if len(u.create.conflict) == 0 {
		return errors.New("ent: missing options for UserCreateBulk.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *UserUpsertBulk) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}
