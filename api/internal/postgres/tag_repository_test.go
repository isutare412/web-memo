package postgres

import (
	"slices"
	"strings"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/samber/lo"

	"github.com/isutare412/web-memo/api/internal/core/ent"
	"github.com/isutare412/web-memo/api/internal/core/ent/enttest"
	"github.com/isutare412/web-memo/api/internal/core/enum"
)

var _ = Describe("TagRepository", func() {
	var (
		tagRepository  *TagRepository
		memoRepository *MemoRepository
		userRepository *UserRepository
	)

	BeforeEach(func() {
		client := enttest.Open(GinkgoT(), "sqlite3", "file:ent?mode=memory&_fk=1")
		tagRepository = NewTagRepository(&Client{entClient: client})
		memoRepository = NewMemoRepository(&Client{entClient: client})
		userRepository = NewUserRepository(&Client{entClient: client})
	})

	Context("with fake data", func() {
		var (
			fakeUser = &ent.User{
				Email:      "faker-user@gmail.com",
				UserName:   "Alice Bob",
				GivenName:  "Alice",
				FamilyName: "Bob",
				PhotoURL:   "google.com/picture",
				Type:       enum.UserTypeClient,
			}
			fakeMemo = &ent.Memo{
				Title:   "memo-one",
				Content: "content-one",
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
			userCreated, err := userRepository.Upsert(ctx, fakeUser)
			Expect(err).NotTo(HaveOccurred())
			fakeUser = userCreated

			for i, tag := range fakeTags {
				tagCreated, err := tagRepository.CreateIfNotExist(ctx, tag.Name)
				Expect(err).NotTo(HaveOccurred())
				fakeTags[i] = tagCreated
			}

			tagIDs := lo.Map(fakeTags[:], func(t *ent.Tag, _ int) int { return t.ID })
			memoCreated, err := memoRepository.Create(ctx, fakeMemo, fakeUser.ID, tagIDs)
			Expect(err).NotTo(HaveOccurred())
			fakeMemo = memoCreated
		})

		// Delete fake data
		AfterEach(func(ctx SpecContext) {
			_ = memoRepository.Delete(ctx, fakeMemo.ID)
			_, err := tagRepository.DeleteAllWithoutMemo(ctx)
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
				tags, err := tagRepository.FindAllByUserIDAndNameContains(ctx, fakeUser.ID, "")
				Expect(err).NotTo(HaveOccurred())
				Expect(len(tags)).To(Equal(2))
				Expect(slices.IsSortedFunc(tags, func(a, b *ent.Tag) int { return strings.Compare(a.Name, b.Name) })).To(BeTrue())
			})

			It("returns tags by name", func(ctx SpecContext) {
				tags, err := tagRepository.FindAllByUserIDAndNameContains(ctx, fakeUser.ID, "one")
				Expect(err).NotTo(HaveOccurred())
				Expect(len(tags)).To(Equal(1))
				Expect(tags[0].Name).To(Equal(fakeTags[0].Name))
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

				deleteCount, err := tagRepository.DeleteAllWithoutMemo(ctx)
				Expect(err).NotTo(HaveOccurred())
				Expect(deleteCount).To(Equal(len(fakeTags)))
			})
		})
	})
})
