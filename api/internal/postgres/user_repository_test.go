package postgres

import (
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/samber/lo"

	"github.com/isutare412/web-memo/api/internal/core/ent"
	"github.com/isutare412/web-memo/api/internal/core/enum"
	"github.com/isutare412/web-memo/api/internal/pkgerr"
)

var _ = Describe("UserRepository", func() {
	var (
		userRepository          *UserRepository
		memoRepository          *MemoRepository
		collaborationRepository *CollaborationRepository
	)

	BeforeEach(func(ctx SpecContext) {
		entClient, err := ent.Open("sqlite3", "file:ent?mode=memory")
		Expect(err).NotTo(HaveOccurred())

		client := &Client{entClient: entClient}
		err = client.MigrateSchemas(ctx)
		Expect(err).NotTo(HaveOccurred())

		userRepository = NewUserRepository(client)
		memoRepository = NewMemoRepository(client)
		collaborationRepository = NewCollaborationRepository(client)
	})

	Context("with fake data", func() {
		var (
			fakeUsers = [...]*ent.User{
				{
					Email:      "faker-user-one@gmail.com",
					UserName:   "Alice Bob",
					GivenName:  "Alice",
					FamilyName: "Bob",
					PhotoURL:   "google.com/picture-one",
					Type:       enum.UserTypeOperator,
				},
				{
					Email:      "faker-user-two@gmail.com",
					UserName:   "Charlie Dave",
					GivenName:  "Cahrlie",
					FamilyName: "Dave",
					PhotoURL:   "google.com/picture-two",
					Type:       enum.UserTypeClient,
				},
			}
			fakeMemos = [...]*ent.Memo{
				{
					Title: "memo-one",
				},
			}
		)

		// Insert fake data
		BeforeEach(func(ctx SpecContext) {
			for i, user := range fakeUsers {
				userCreated, err := userRepository.Upsert(ctx, user)
				Expect(err).NotTo(HaveOccurred())
				fakeUsers[i] = userCreated
			}

			for i, memo := range fakeMemos {
				memoCreated, err := memoRepository.Create(ctx, memo, fakeUsers[0].ID, nil)
				Expect(err).NotTo(HaveOccurred())
				fakeMemos[i] = memoCreated
			}

			err := memoRepository.RegisterSubscriber(ctx, fakeMemos[0].ID, fakeUsers[0].ID)
			Expect(err).NotTo(HaveOccurred())

			_, err = collaborationRepository.Create(ctx, fakeMemos[0].ID, fakeUsers[0].ID)
			Expect(err).NotTo(HaveOccurred())
		})

		Context("FindByID", func() {
			It("finds user by id", func(ctx SpecContext) {
				user, err := userRepository.FindByID(ctx, fakeUsers[0].ID)
				Expect(err).NotTo(HaveOccurred())
				Expect(user).NotTo(BeNil())
				Expect(user.ID).To(Equal(fakeUsers[0].ID))
			})

			It("returns not found error if unknown id", func(ctx SpecContext) {
				_, err := userRepository.FindByID(ctx, uuid.Must(uuid.NewRandom()))
				Expect(pkgerr.IsErrNotFound(err)).To(BeTrue())
			})
		})

		Context("FindByEmail", func() {
			It("finds user by email", func(ctx SpecContext) {
				user, err := userRepository.FindByEmail(ctx, fakeUsers[0].Email)
				Expect(err).NotTo(HaveOccurred())
				Expect(user).NotTo(BeNil())
				Expect(user.ID).To(Equal(fakeUsers[0].ID))
			})

			It("returns not found error if unknown email", func(ctx SpecContext) {
				_, err := userRepository.FindByEmail(ctx, "complex-email@abc.com")
				Expect(pkgerr.IsErrNotFound(err)).To(BeTrue())
			})
		})

		Context("FindAllBySubscribingMemoID", func() {
			It("finds only subscribers", func(ctx SpecContext) {
				subscribers, err := userRepository.FindAllBySubscribingMemoID(ctx, fakeMemos[0].ID)
				Expect(err).NotTo(HaveOccurred())
				Expect(subscribers).To(HaveLen(1))
				Expect(subscribers[0].ID).To(Equal(fakeUsers[0].ID))
			})
		})

		Context("FindAllByCollaboratingMemoID", func() {
			It("finds only collaborators", func(ctx SpecContext) {
				collaborators, err := userRepository.FindAllByCollaboratingMemoIDWithEdges(ctx, fakeMemos[0].ID)
				Expect(err).NotTo(HaveOccurred())
				Expect(collaborators).To(HaveLen(1))
				Expect(collaborators[0].ID).To(Equal(fakeUsers[0].ID))
			})
		})

		Context("Upsert", func() {
			It("creates user if not exists", func(ctx SpecContext) {
				var (
					givenUser = &ent.User{
						Email:    "unknown-user@gmail.com",
						UserName: "Foo Bar",
						Type:     enum.UserTypeClient,
					}
				)

				fakeUserIDs := lo.Map(fakeUsers[:], func(u *ent.User, _ int) uuid.UUID { return u.ID })

				userCreated, err := userRepository.Upsert(ctx, givenUser)
				Expect(err).NotTo(HaveOccurred())
				Expect(userCreated.ID).NotTo(BeZero())
				Expect(userCreated.ID).NotTo(BeElementOf(fakeUserIDs))
			})

			It("updates user if email exists", func(ctx SpecContext) {
				var (
					givenUser = &ent.User{
						Email:      fakeUsers[0].Email,
						UserName:   "new-user-name",
						GivenName:  "new-given-name",
						FamilyName: "new-family-name",
						PhotoURL:   "new-photo-url",
						Type:       enum.UserTypeClient,
					}
				)

				fakeUserIDs := lo.Map(fakeUsers[:], func(u *ent.User, _ int) uuid.UUID { return u.ID })

				userUpdated, err := userRepository.Upsert(ctx, givenUser)
				Expect(err).NotTo(HaveOccurred())
				Expect(userUpdated.ID).NotTo(BeZero())
				Expect(userUpdated.ID).To(BeElementOf(fakeUserIDs))
				Expect(userUpdated.UserName).To(Equal(givenUser.UserName))
				Expect(userUpdated.GivenName).To(Equal(givenUser.GivenName))
				Expect(userUpdated.FamilyName).To(Equal(givenUser.FamilyName))
				Expect(userUpdated.PhotoURL).To(Equal(givenUser.PhotoURL))
				Expect(userUpdated.Type).To(Equal(givenUser.Type))
			})
		})
	})
})
