package postgres

import (
	"slices"
	"strings"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/samber/lo"

	"github.com/isutare412/web-memo/api/internal/core/ent"
	"github.com/isutare412/web-memo/api/internal/core/enum"
)

var _ = Describe("TagRepository", func() {
	var (
		tagRepository           *TagRepository
		memoRepository          *MemoRepository
		userRepository          *UserRepository
		collaborationRepository *CollaborationRepository
	)

	BeforeEach(func(ctx SpecContext) {
		entClient, err := ent.Open("sqlite3", "file:ent?mode=memory")
		Expect(err).NotTo(HaveOccurred())

		client := &Client{entClient: entClient}
		err = client.MigrateSchemas(ctx)
		Expect(err).NotTo(HaveOccurred())

		tagRepository = NewTagRepository(client)
		memoRepository = NewMemoRepository(client)
		userRepository = NewUserRepository(client)
		collaborationRepository = NewCollaborationRepository(client)
	})

	Context("with fake data", func() {
		var (
			fakeUsers = [...]*ent.User{
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
				{
					Email:      "third-user@gmail.com",
					UserName:   "Ergo God",
					GivenName:  "Ergo",
					FamilyName: "God",
					PhotoURL:   "google.com/picture3",
					Type:       enum.UserTypeClient,
				},
			}
			fakeMemo = &ent.Memo{
				Title:        "memo-one",
				Content:      "content-one",
				PublishState: enum.PublishStatePrivate,
			}
			fakeTags = [...]*ent.Tag{
				{
					Name: "tag-one",
				},
				{
					Name: "tag-two",
				},
			}
		)

		// Insert fake data
		BeforeEach(func(ctx SpecContext) {
			for i, u := range fakeUsers {
				userCreated, err := userRepository.Upsert(ctx, u)
				Expect(err).NotTo(HaveOccurred())
				fakeUsers[i] = userCreated
			}

			for i, tag := range fakeTags {
				tagCreated, err := tagRepository.CreateIfNotExist(ctx, tag.Name)
				Expect(err).NotTo(HaveOccurred())
				fakeTags[i] = tagCreated
			}

			tagIDs := lo.Map(fakeTags[:], func(t *ent.Tag, _ int) int { return t.ID })
			memoCreated, err := memoRepository.Create(ctx, fakeMemo, fakeUsers[0].ID, tagIDs)
			Expect(err).NotTo(HaveOccurred())
			fakeMemo = memoCreated

			err = memoRepository.RegisterSubscriber(ctx, memoCreated.ID, fakeUsers[1].ID, true)
			Expect(err).NotTo(HaveOccurred())

			_, err = collaborationRepository.Create(ctx, memoCreated.ID, fakeUsers[2].ID)
			Expect(err).NotTo(HaveOccurred())
		})

		// Delete fake data
		AfterEach(func(ctx SpecContext) {
			_ = memoRepository.Delete(ctx, fakeMemo.ID)
			_, err := tagRepository.DeleteAllWithoutMemo(ctx, nil)
			Expect(err).NotTo(HaveOccurred())
		})

		Context("FindAllByMemoID", func() {
			It("returns all tags", func(ctx SpecContext) {
				tags, err := tagRepository.FindAllByMemoID(ctx, fakeMemo.ID)
				Expect(err).NotTo(HaveOccurred())
				Expect(len(tags)).To(Equal(2))
				Expect(slices.IsSortedFunc(tags, func(a, b *ent.Tag) int { return strings.Compare(a.Name, b.Name) })).To(BeTrue())
			})

			It("does not return if memo ID is not known", func(ctx SpecContext) {
				tags, err := tagRepository.FindAllByMemoID(ctx, uuid.Must(uuid.NewRandom()))
				Expect(err).NotTo(HaveOccurred())
				Expect(len(tags)).To(Equal(0))
			})
		})

		Context("FindAllByUserIDAndNameContains", func() {
			It("returns all if name is zero-value", func(ctx SpecContext) {
				tags, err := tagRepository.FindAllByUserIDAndNameContains(ctx, fakeUsers[0].ID, "")
				Expect(err).NotTo(HaveOccurred())
				Expect(len(tags)).To(Equal(2))
				Expect(slices.IsSortedFunc(tags, func(a, b *ent.Tag) int { return strings.Compare(a.Name, b.Name) })).To(BeTrue())
			})

			It("returns tags by name", func(ctx SpecContext) {
				tags, err := tagRepository.FindAllByUserIDAndNameContains(ctx, fakeUsers[0].ID, "one")
				Expect(err).NotTo(HaveOccurred())
				Expect(len(tags)).To(Equal(1))
				Expect(tags[0].Name).To(Equal(fakeTags[0].Name))
			})

			It("returns by subscriber", func(ctx SpecContext) {
				tags, err := tagRepository.FindAllByUserIDAndNameContains(ctx, fakeUsers[1].ID, "")
				Expect(err).NotTo(HaveOccurred())
				Expect(len(tags) > 0).To(BeTrue())
			})
		})

		Context("CreateIfNotExist", func() {
			It("creates tag", func(ctx SpecContext) {
				fakeTagIDs := lo.Map(fakeTags[:], func(t *ent.Tag, _ int) int { return t.ID })

				tagCreated, err := tagRepository.CreateIfNotExist(ctx, "very-complex-tag-name")
				Expect(err).NotTo(HaveOccurred())
				Expect(tagCreated.ID).NotTo(BeElementOf(fakeTagIDs))
			})

			It("does not create if tag already exists", func(ctx SpecContext) {
				tagCreated, err := tagRepository.CreateIfNotExist(ctx, fakeTags[0].Name)
				Expect(err).NotTo(HaveOccurred())
				Expect(tagCreated.ID).To(Equal(fakeTags[0].ID))
			})
		})

		Context("DeleteAllWithoutMemo", func() {
			It("deletes tags", func(ctx SpecContext) {
				err := memoRepository.Delete(ctx, fakeMemo.ID)
				Expect(err).NotTo(HaveOccurred())

				deleteCount, err := tagRepository.DeleteAllWithoutMemo(ctx, []string{fakeTags[0].Name})
				Expect(err).NotTo(HaveOccurred())
				Expect(deleteCount).To(Equal(len(fakeTags) - 1))
			})
		})
	})
})
