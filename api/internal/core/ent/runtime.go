// Code generated by ent, DO NOT EDIT.

package ent

import (
	"time"

	"github.com/google/uuid"
	"github.com/isutare412/web-memo/api/internal/core/ent/memo"
	"github.com/isutare412/web-memo/api/internal/core/ent/schema"
	"github.com/isutare412/web-memo/api/internal/core/ent/tag"
	"github.com/isutare412/web-memo/api/internal/core/ent/user"
)

// The init function reads all schema descriptors with runtime code
// (default values, validators, hooks and policies) and stitches it
// to their package variables.
func init() {
	memoMixin := schema.Memo{}.Mixin()
	memoMixinFields0 := memoMixin[0].Fields()
	_ = memoMixinFields0
	memoFields := schema.Memo{}.Fields()
	_ = memoFields
	// memoDescCreateTime is the schema descriptor for create_time field.
	memoDescCreateTime := memoMixinFields0[0].Descriptor()
	// memo.DefaultCreateTime holds the default value on creation for the create_time field.
	memo.DefaultCreateTime = memoDescCreateTime.Default.(func() time.Time)
	// memoDescUpdateTime is the schema descriptor for update_time field.
	memoDescUpdateTime := memoMixinFields0[1].Descriptor()
	// memo.DefaultUpdateTime holds the default value on creation for the update_time field.
	memo.DefaultUpdateTime = memoDescUpdateTime.Default.(func() time.Time)
	// memo.UpdateDefaultUpdateTime holds the default value on update for the update_time field.
	memo.UpdateDefaultUpdateTime = memoDescUpdateTime.UpdateDefault.(func() time.Time)
	// memoDescTitle is the schema descriptor for title field.
	memoDescTitle := memoFields[1].Descriptor()
	// memo.TitleValidator is a validator for the "title" field. It is called by the builders before save.
	memo.TitleValidator = memoDescTitle.Validators[0].(func(string) error)
	// memoDescContent is the schema descriptor for content field.
	memoDescContent := memoFields[2].Descriptor()
	// memo.ContentValidator is a validator for the "content" field. It is called by the builders before save.
	memo.ContentValidator = memoDescContent.Validators[0].(func(string) error)
	// memoDescID is the schema descriptor for id field.
	memoDescID := memoFields[0].Descriptor()
	// memo.DefaultID holds the default value on creation for the id field.
	memo.DefaultID = memoDescID.Default.(func() uuid.UUID)
	tagMixin := schema.Tag{}.Mixin()
	tagMixinFields0 := tagMixin[0].Fields()
	_ = tagMixinFields0
	tagFields := schema.Tag{}.Fields()
	_ = tagFields
	// tagDescCreateTime is the schema descriptor for create_time field.
	tagDescCreateTime := tagMixinFields0[0].Descriptor()
	// tag.DefaultCreateTime holds the default value on creation for the create_time field.
	tag.DefaultCreateTime = tagDescCreateTime.Default.(func() time.Time)
	// tagDescUpdateTime is the schema descriptor for update_time field.
	tagDescUpdateTime := tagMixinFields0[1].Descriptor()
	// tag.DefaultUpdateTime holds the default value on creation for the update_time field.
	tag.DefaultUpdateTime = tagDescUpdateTime.Default.(func() time.Time)
	// tag.UpdateDefaultUpdateTime holds the default value on update for the update_time field.
	tag.UpdateDefaultUpdateTime = tagDescUpdateTime.UpdateDefault.(func() time.Time)
	// tagDescName is the schema descriptor for name field.
	tagDescName := tagFields[0].Descriptor()
	// tag.NameValidator is a validator for the "name" field. It is called by the builders before save.
	tag.NameValidator = func() func(string) error {
		validators := tagDescName.Validators
		fns := [...]func(string) error{
			validators[0].(func(string) error),
			validators[1].(func(string) error),
		}
		return func(name string) error {
			for _, fn := range fns {
				if err := fn(name); err != nil {
					return err
				}
			}
			return nil
		}
	}()
	userMixin := schema.User{}.Mixin()
	userMixinFields0 := userMixin[0].Fields()
	_ = userMixinFields0
	userFields := schema.User{}.Fields()
	_ = userFields
	// userDescCreateTime is the schema descriptor for create_time field.
	userDescCreateTime := userMixinFields0[0].Descriptor()
	// user.DefaultCreateTime holds the default value on creation for the create_time field.
	user.DefaultCreateTime = userDescCreateTime.Default.(func() time.Time)
	// userDescUpdateTime is the schema descriptor for update_time field.
	userDescUpdateTime := userMixinFields0[1].Descriptor()
	// user.DefaultUpdateTime holds the default value on creation for the update_time field.
	user.DefaultUpdateTime = userDescUpdateTime.Default.(func() time.Time)
	// user.UpdateDefaultUpdateTime holds the default value on update for the update_time field.
	user.UpdateDefaultUpdateTime = userDescUpdateTime.UpdateDefault.(func() time.Time)
	// userDescEmail is the schema descriptor for email field.
	userDescEmail := userFields[1].Descriptor()
	// user.EmailValidator is a validator for the "email" field. It is called by the builders before save.
	user.EmailValidator = func() func(string) error {
		validators := userDescEmail.Validators
		fns := [...]func(string) error{
			validators[0].(func(string) error),
			validators[1].(func(string) error),
		}
		return func(email string) error {
			for _, fn := range fns {
				if err := fn(email); err != nil {
					return err
				}
			}
			return nil
		}
	}()
	// userDescUserName is the schema descriptor for user_name field.
	userDescUserName := userFields[2].Descriptor()
	// user.UserNameValidator is a validator for the "user_name" field. It is called by the builders before save.
	user.UserNameValidator = func() func(string) error {
		validators := userDescUserName.Validators
		fns := [...]func(string) error{
			validators[0].(func(string) error),
			validators[1].(func(string) error),
		}
		return func(user_name string) error {
			for _, fn := range fns {
				if err := fn(user_name); err != nil {
					return err
				}
			}
			return nil
		}
	}()
	// userDescGivenName is the schema descriptor for given_name field.
	userDescGivenName := userFields[3].Descriptor()
	// user.GivenNameValidator is a validator for the "given_name" field. It is called by the builders before save.
	user.GivenNameValidator = func() func(string) error {
		validators := userDescGivenName.Validators
		fns := [...]func(string) error{
			validators[0].(func(string) error),
			validators[1].(func(string) error),
		}
		return func(given_name string) error {
			for _, fn := range fns {
				if err := fn(given_name); err != nil {
					return err
				}
			}
			return nil
		}
	}()
	// userDescFamilyName is the schema descriptor for family_name field.
	userDescFamilyName := userFields[4].Descriptor()
	// user.FamilyNameValidator is a validator for the "family_name" field. It is called by the builders before save.
	user.FamilyNameValidator = func() func(string) error {
		validators := userDescFamilyName.Validators
		fns := [...]func(string) error{
			validators[0].(func(string) error),
			validators[1].(func(string) error),
		}
		return func(family_name string) error {
			for _, fn := range fns {
				if err := fn(family_name); err != nil {
					return err
				}
			}
			return nil
		}
	}()
	// userDescPhotoURL is the schema descriptor for photo_url field.
	userDescPhotoURL := userFields[5].Descriptor()
	// user.PhotoURLValidator is a validator for the "photo_url" field. It is called by the builders before save.
	user.PhotoURLValidator = userDescPhotoURL.Validators[0].(func(string) error)
	// userDescID is the schema descriptor for id field.
	userDescID := userFields[0].Descriptor()
	// user.DefaultID holds the default value on creation for the id field.
	user.DefaultID = userDescID.Default.(func() uuid.UUID)
}
