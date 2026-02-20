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

			It("passes if approved collaborator", func() {
				var (
					givenUserID = uuid.Must(uuid.NewRandom())
					givenToken  = &model.AppIDToken{
						UserID:   givenUserID,
						UserType: enum.UserTypeClient,
					}
					givenMemo = &ent.Memo{
						Edges: ent.MemoEdges{
							Collaborations: []*ent.Collaboration{{
								UserID:   givenUserID,
								Approved: true,
							}},
						},
					}
				)

				Expect(givenToken.CanWriteMemo(givenMemo)).To(BeTrue())
			})

			It("fails if non-approved collaborator", func() {
				var (
					givenUserID = uuid.Must(uuid.NewRandom())
					givenToken  = &model.AppIDToken{
						UserID:   givenUserID,
						UserType: enum.UserTypeClient,
					}
					givenMemo = &ent.Memo{
						Edges: ent.MemoEdges{
							Collaborations: []*ent.Collaboration{{
								UserID: givenUserID,
							}},
						},
					}
				)

				Expect(givenToken.CanWriteMemo(givenMemo)).To(BeFalse())
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
						PublishState: enum.PublishStatePublished,
					}
				)

				Expect(givenToken.CanReadMemo(givenMemo)).To(BeTrue())
			})

			It("fails if shared memo without token", func() {
				var (
					givenToken *model.AppIDToken = nil
					givenMemo                    = &ent.Memo{
						PublishState: enum.PublishStateShared,
					}
				)

				Expect(givenToken.CanReadMemo(givenMemo)).To(BeFalse())
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

		Context("IsOwner", func() {
			It("returns true if owner", func() {
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

				Expect(givenToken.IsOwner(givenMemo)).To(BeTrue())
			})

			It("returns false if nil token", func() {
				var (
					givenUserID                   = uuid.Must(uuid.NewRandom())
					givenToken  *model.AppIDToken = nil
					givenMemo                     = &ent.Memo{
						OwnerID: givenUserID,
					}
				)

				Expect(givenToken.IsOwner(givenMemo)).To(BeFalse())
			})
		})
	})
})
