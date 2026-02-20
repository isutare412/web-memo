package postgres

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/isutare412/web-memo/api/internal/core/ent"
	"github.com/isutare412/web-memo/api/internal/core/enum"
	"github.com/isutare412/web-memo/api/internal/pkgerr"
)

var _ = Describe("CollaborationRepository", func() {
	var (
		collaborationRepository *CollaborationRepository
		memoRepository          *MemoRepository
		userRepository          *UserRepository
	)

	BeforeEach(func(ctx SpecContext) {
		entClient, err := ent.Open("sqlite3", "file:ent?mode=memory")
		Expect(err).NotTo(HaveOccurred())

		client := &Client{entClient: entClient}
		err = client.MigrateSchemas(ctx)
		Expect(err).NotTo(HaveOccurred())

		collaborationRepository = NewCollaborationRepository(client)
		memoRepository = NewMemoRepository(client)
		userRepository = NewUserRepository(client)
	})

	Context("with fake data", func() {
		var (
			fakeUsers = []*ent.User{
				{
					Email:      "faker-user@gmail.com",
					UserName:   "Alice Bob",
					GivenName:  "Alice",
					FamilyName: "Bob",
					PhotoURL:   "google.com/picture",
					Type:       enum.UserTypeClient,
				},
				{
					Email:      "another-user@gmail.com",
					UserName:   "Charlie Decart",
					GivenName:  "Charlie",
					FamilyName: "Decart",
					PhotoURL:   "google.com/picture2",
					Type:       enum.UserTypeClient,
				},
			}
			fakeMemo = &ent.Memo{
				Title:        "memo-one",
				Content:      "content-one",
				PublishState: enum.PublishStatePrivate,
			}
		)

		BeforeEach(func(ctx SpecContext) {
			for i, u := range fakeUsers {
				userCreated, err := userRepository.Upsert(ctx, u)
				Expect(err).NotTo(HaveOccurred())
				fakeUsers[i] = userCreated
			}

			memoCreated, err := memoRepository.Create(ctx, fakeMemo, fakeUsers[0].ID, nil)
			Expect(err).NotTo(HaveOccurred())
			fakeMemo = memoCreated

			_, err = collaborationRepository.Create(ctx, fakeMemo.ID, fakeUsers[0].ID)
			Expect(err).NotTo(HaveOccurred())
		})

		Context("Find", func() {
			It("finds collaboration", func(ctx SpecContext) {
				collabo, err := collaborationRepository.Find(ctx, fakeMemo.ID, fakeUsers[0].ID)
				Expect(err).NotTo(HaveOccurred())
				Expect(collabo).NotTo(BeNil())
				Expect(collabo.UserID).To(Equal(fakeUsers[0].ID))
			})

			It("emits not found", func(ctx SpecContext) {
				_, err := collaborationRepository.Find(ctx, fakeMemo.ID, fakeUsers[1].ID)
				Expect(pkgerr.IsErrNotFound(err)).To(BeTrue())
			})
		})

		Context("Create", func() {
			It("creates collaboration", func(ctx SpecContext) {
				collabo, err := collaborationRepository.Create(ctx, fakeMemo.ID, fakeUsers[1].ID)
				Expect(err).NotTo(HaveOccurred())
				Expect(collabo).NotTo(BeNil())
				Expect(collabo.UserID).To(Equal(fakeUsers[1].ID))
				Expect(collabo.MemoID).To(Equal(fakeMemo.ID))
			})

			It("fails if duplicate", func(ctx SpecContext) {
				_, err := collaborationRepository.Create(ctx, fakeMemo.ID, fakeUsers[0].ID)
				Expect(err).To(HaveOccurred())
			})
		})

		Context("UpdateApprovedStatus", func() {
			It("updates approved field", func(ctx SpecContext) {
				collabo, err := collaborationRepository.UpdateApprovedStatus(ctx, fakeMemo.ID, fakeUsers[0].ID, true)
				Expect(err).NotTo(HaveOccurred())
				Expect(collabo.Approved).To(BeTrue())
			})

			It("emits not found error if collaboration does not exist", func(ctx SpecContext) {
				_, err := collaborationRepository.UpdateApprovedStatus(ctx, fakeMemo.ID, fakeUsers[1].ID, true)
				Expect(err).To(HaveOccurred())
				Expect(pkgerr.IsErrNotFound(err)).To(BeTrue())
			})
		})

		Context("Delete", func() {
			It("deletes collaboration", func(ctx SpecContext) {
				err := collaborationRepository.Delete(ctx, fakeMemo.ID, fakeUsers[0].ID)
				Expect(err).NotTo(HaveOccurred())
			})

			It("fails if collaboration does not exist", func(ctx SpecContext) {
				err := collaborationRepository.Delete(ctx, fakeMemo.ID, fakeUsers[1].ID)
				Expect(pkgerr.IsErrNotFound(err)).To(BeTrue())
			})
		})

		Context("DeleteAllByMemoID", func() {
			It("deletes by memo ID", func(ctx SpecContext) {
				count, err := collaborationRepository.DeleteAllByMemoID(ctx, fakeMemo.ID)
				Expect(err).NotTo(HaveOccurred())
				Expect(count).To(Equal(1))
			})
		})
	})
})
