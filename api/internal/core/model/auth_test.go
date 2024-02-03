package model_test

import (
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/isutare412/web-memo/api/internal/core/ent"
	"github.com/isutare412/web-memo/api/internal/core/enum"
	"github.com/isutare412/web-memo/api/internal/core/model"
)

var _ = Describe("Auth", func() {
	Context("AppIDToken", func() {
		Context("CanWriteMemo", func() {
			It("passes if operator", func() {
				var (
					givenToken = &model.AppIDToken{
						UserType: enum.UserTypeOperator,
					}
					givenMemo = &ent.Memo{
						OwnerID: uuid.Must(uuid.NewRandom()),
					}
				)

				Expect(givenToken.CanWriteMemo(givenMemo)).To(BeTrue())
			})

			It("passes if owner", func() {
				var (
					givenUserID = uuid.Must(uuid.NewRandom())
					givenToken  = &model.AppIDToken{
						UserID:   givenUserID,
						UserType: enum.UserTypeClient,
					}
					givenMemo = &ent.Memo{
						OwnerID: givenUserID,
					}
				)

				Expect(givenToken.CanWriteMemo(givenMemo)).To(BeTrue())
			})

			It("fails if not owner", func() {
				var (
					givenUserID = uuid.Must(uuid.NewRandom())
					givenToken  = &model.AppIDToken{
						UserID:   givenUserID,
						UserType: enum.UserTypeClient,
					}
					givenMemo = &ent.Memo{}
				)

				Expect(givenToken.CanWriteMemo(givenMemo)).To(BeFalse())
			})

			It("fails if nil token", func() {
				var (
					givenToken *model.AppIDToken = nil
					givenMemo                    = &ent.Memo{}
				)

				Expect(givenToken.CanWriteMemo(givenMemo)).To(BeFalse())
			})
		})

		Context("CanReadMemo", func() {
			It("passes if operator", func() {
				var (
					givenToken = &model.AppIDToken{
						UserType: enum.UserTypeOperator,
					}
					givenMemo = &ent.Memo{
						OwnerID: uuid.Must(uuid.NewRandom()),
					}
				)

				Expect(givenToken.CanReadMemo(givenMemo)).To(BeTrue())
			})

			It("passes if owner", func() {
				var (
					givenUserID = uuid.Must(uuid.NewRandom())
					givenToken  = &model.AppIDToken{
						UserID:   givenUserID,
						UserType: enum.UserTypeClient,
					}
					givenMemo = &ent.Memo{
						OwnerID: givenUserID,
					}
				)

				Expect(givenToken.CanReadMemo(givenMemo)).To(BeTrue())
			})

			It("passes if published memo", func() {
				var (
					givenToken *model.AppIDToken = nil
					givenMemo                    = &ent.Memo{
						IsPublished: true,
					}
				)

				Expect(givenToken.CanReadMemo(givenMemo)).To(BeTrue())
			})

			It("fails if not owner", func() {
				var (
					givenUserID = uuid.Must(uuid.NewRandom())
					givenToken  = &model.AppIDToken{
						UserID:   givenUserID,
						UserType: enum.UserTypeClient,
					}
					givenMemo = &ent.Memo{}
				)

				Expect(givenToken.CanReadMemo(givenMemo)).To(BeFalse())
			})
		})

		Context("CanReadMemo", func() {
		})
	})
})
