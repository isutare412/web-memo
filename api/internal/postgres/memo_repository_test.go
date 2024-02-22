package postgres

import (
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/samber/lo"

	"github.com/isutare412/web-memo/api/internal/core/ent"
	"github.com/isutare412/web-memo/api/internal/core/enum"
	"github.com/isutare412/web-memo/api/internal/core/model"
	"github.com/isutare412/web-memo/api/internal/pkgerr"
)

var _ = Describe("MemoRepository", func() {
	var (
		tagRepository  *TagRepository
		memoRepository *MemoRepository
		userRepository *UserRepository
	)

	BeforeEach(func(ctx SpecContext) {
		entClient, err := ent.Open("sqlite3", "file:ent?mode=memory&cache=shared")
		Expect(err).NotTo(HaveOccurred())

		client := &Client{entClient: entClient}
		err = client.MigrateSchemas(ctx)
		Expect(err).NotTo(HaveOccurred())

		tagRepository = NewTagRepository(client)
		memoRepository = NewMemoRepository(client)
		userRepository = NewUserRepository(client)
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
			}
			fakeMemos = [...]*ent.Memo{
				{
					Title:   "memo-one",
					Content: "content-one",
				},
				{
					Title:   "memo-two",
					Content: "content-two",
				},
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

			Expect(fakeMemos).To(HaveLen(len(fakeTags)))
			for i, memo := range fakeMemos {
				memoCreated, err := memoRepository.Create(ctx, memo, fakeUsers[0].ID, []int{fakeTags[i].ID})
				Expect(err).NotTo(HaveOccurred())
				fakeMemos[i] = memoCreated

				if err := memoRepository.RegisterSubscriber(ctx, memoCreated.ID, fakeUsers[1].ID); err != nil {
					Expect(err).NotTo(HaveOccurred())
				}
			}
		})

		// Delete fake data
		AfterEach(func(ctx SpecContext) {
			memos, err := memoRepository.FindAllByUserIDWithTags(
				ctx, fakeUsers[0].ID, model.MemoSortParams{}, model.PaginationParams{})
			Expect(err).NotTo(HaveOccurred())
			for _, memo := range memos {
				_ = memoRepository.Delete(ctx, memo.ID)
			}

			_, err = tagRepository.DeleteAllWithoutMemo(ctx, nil)
			Expect(err).NotTo(HaveOccurred())
		})

		Context("FindByID", func() {
			It("finds memo by ID", func(ctx SpecContext) {
				memo, err := memoRepository.FindByID(ctx, fakeMemos[0].ID)
				Expect(err).NotTo(HaveOccurred())
				Expect(memo).NotTo(BeNil())
				Expect(memo.ID).To(Equal(fakeMemos[0].ID))
			})

			It("returns not found error if unknown ID", func(ctx SpecContext) {
				_, err := memoRepository.FindByID(ctx, uuid.Must(uuid.NewRandom()))
				Expect(pkgerr.IsErrNotFound(err)).To(BeTrue())
			})
		})

		Context("FindByIDWithTags", func() {
			It("finds memo with tags eager loaded", func(ctx SpecContext) {
				memo, err := memoRepository.FindByIDWithTags(ctx, fakeMemos[0].ID)
				Expect(err).NotTo(HaveOccurred())
				Expect(memo).NotTo(BeNil())
				Expect(memo.ID).To(Equal(fakeMemos[0].ID))
				Expect(len(memo.Edges.Tags)).NotTo(BeZero())
			})

			It("returns not found error if unknown ID", func(ctx SpecContext) {
				_, err := memoRepository.FindByIDWithTags(ctx, uuid.Must(uuid.NewRandom()))
				Expect(pkgerr.IsErrNotFound(err)).To(BeTrue())
			})
		})

		Context("FindAllByUserIDWithTags", func() {
			It("finds memos of user", func(ctx SpecContext) {
				memos, err := memoRepository.FindAllByUserIDWithTags(
					ctx, fakeUsers[0].ID, model.MemoSortParams{}, model.PaginationParams{})
				Expect(err).NotTo(HaveOccurred())
				Expect(memos).To(HaveLen(len(fakeMemos)))
				Expect(memos[0].Edges.Tags).NotTo(HaveLen(0))
			})

			It("finds memos by subscriber", func(ctx SpecContext) {
				memos, err := memoRepository.FindAllByUserIDWithTags(
					ctx, fakeUsers[1].ID, model.MemoSortParams{}, model.PaginationParams{})
				Expect(err).NotTo(HaveOccurred())
				Expect(memos).To(HaveLen(len(fakeMemos)))
			})

			It("finds memos with pagination", func(ctx SpecContext) {
				var (
					givenPageParams = model.PaginationParams{
						PageOffset: 1,
						PageSize:   1,
					}
					givenSortParams = model.MemoSortParams{
						Order: enum.SortOrderAsc,
					}
				)

				memos, err := memoRepository.FindAllByUserIDWithTags(ctx, fakeUsers[0].ID, givenSortParams, givenPageParams)
				Expect(err).NotTo(HaveOccurred())
				Expect(memos).To(HaveLen(1))
				Expect(memos[0].ID).To(Equal(fakeMemos[0].ID))
				Expect(memos[0].OwnerID).To(Equal(fakeUsers[0].ID))
			})

			It("finds nothing if unknown ID", func(ctx SpecContext) {
				var (
					givenPageParams = model.PaginationParams{}
					givenSortParams = model.MemoSortParams{}
				)

				memos, err := memoRepository.FindAllByUserIDWithTags(
					ctx, uuid.Must(uuid.NewRandom()), givenSortParams, givenPageParams)
				Expect(err).NotTo(HaveOccurred())
				Expect(memos).To(HaveLen(0))
			})
		})

		Context("FindAllByUserIDAndTagNamesWithTags", func() {
			It("finds memos of user with tag name", func(ctx SpecContext) {
				var (
					givenTagNames   = []string{fakeTags[0].Name}
					givenPageParams = model.PaginationParams{}
					givenSortParams = model.MemoSortParams{}
				)

				memos, err := memoRepository.FindAllByUserIDAndTagNamesWithTags(
					ctx, fakeUsers[0].ID, givenTagNames, givenSortParams, givenPageParams)
				Expect(err).NotTo(HaveOccurred())
				Expect(memos).To(HaveLen(1))
				Expect(memos[0].ID).To(Equal(fakeMemos[0].ID))
				Expect(memos[0].OwnerID).To(Equal(fakeUsers[0].ID))
			})

			It("finds nothing if tag name does not match", func(ctx SpecContext) {
				var (
					givenTagNames   = []string{"some-complex-tag"}
					givenPageParams = model.PaginationParams{}
					givenSortParams = model.MemoSortParams{}
				)

				memos, err := memoRepository.FindAllByUserIDAndTagNamesWithTags(
					ctx, fakeUsers[0].ID, givenTagNames, givenSortParams, givenPageParams)
				Expect(err).NotTo(HaveOccurred())
				Expect(memos).To(HaveLen(0))
			})
		})

		Context("CountByUserIDAndTagNames", func() {
			It("counts memos of user", func(ctx SpecContext) {
				count, err := memoRepository.CountByUserIDAndTagNames(ctx, fakeUsers[0].ID, nil)
				Expect(err).NotTo(HaveOccurred())
				Expect(count).To(Equal(len(fakeMemos)))
			})

			It("counts memos of user with matching tags", func(ctx SpecContext) {
				var (
					givenTagNames = []string{fakeTags[0].Name}
				)

				count, err := memoRepository.CountByUserIDAndTagNames(ctx, fakeUsers[0].ID, givenTagNames)
				Expect(err).NotTo(HaveOccurred())
				Expect(count).To(Equal(1))
			})
		})

		Context("Create", func() {
			It("creates memo", func(ctx SpecContext) {
				var (
					givenMemo = &ent.Memo{
						Title:   "new-title",
						Content: "new-content",
					}
					givenTagIDs = lo.Map(fakeTags[:], func(t *ent.Tag, _ int) int { return t.ID })
				)

				fakeMemoIDs := lo.Map(fakeMemos[:], func(m *ent.Memo, _ int) uuid.UUID { return m.ID })

				memo, err := memoRepository.Create(ctx, givenMemo, fakeUsers[0].ID, givenTagIDs)
				Expect(err).NotTo(HaveOccurred())
				Expect(memo).NotTo(BeNil())
				Expect(memo.ID).NotTo(BeElementOf(fakeMemoIDs))
				Expect(memo.Content).To(Equal(givenMemo.Content))
			})

			It("returns error if title is empty", func(ctx SpecContext) {
				var (
					givenMemo = &ent.Memo{}
				)

				_, err := memoRepository.Create(ctx, givenMemo, fakeUsers[0].ID, nil)
				Expect(err).To(HaveOccurred())
			})
		})

		Context("Update", func() {
			It("updates memo", func(ctx SpecContext) {
				var (
					givenMemo = &ent.Memo{
						ID:          fakeMemos[0].ID,
						Title:       "new-title",
						Content:     "new-content",
						IsPublished: true,
					}
				)

				memo, err := memoRepository.Update(ctx, givenMemo)
				Expect(err).NotTo(HaveOccurred())
				Expect(memo).NotTo(BeNil())
				Expect(memo.ID).To(Equal(fakeMemos[0].ID))
				Expect(memo.Title).To(Equal(givenMemo.Title))
				Expect(memo.Content).To(Equal(givenMemo.Content))
				Expect(memo.IsPublished).To(Equal(givenMemo.IsPublished))
			})

			It("returns not found error if id does not exist", func(ctx SpecContext) {
				var (
					givenMemo = &ent.Memo{
						ID:      uuid.Must(uuid.NewRandom()),
						Title:   "new-title",
						Content: "new-content",
					}
				)

				_, err := memoRepository.Update(ctx, givenMemo)
				Expect(pkgerr.IsErrNotFound(err)).To(BeTrue())
			})
		})

		Context("Delete", func() {
			It("deletes memo", func(ctx SpecContext) {
				err := memoRepository.Delete(ctx, fakeMemos[0].ID)
				Expect(err).NotTo(HaveOccurred())
			})

			It("returns not found error if id does not exist", func(ctx SpecContext) {
				err := memoRepository.Delete(ctx, uuid.Must(uuid.NewRandom()))
				Expect(pkgerr.IsErrNotFound(err)).To(BeTrue())
			})
		})

		Context("ReplaceTags", func() {
			It("replaces tags of memo", func(ctx SpecContext) {
				err := memoRepository.ReplaceTags(ctx, fakeMemos[0].ID, []int{fakeTags[1].ID})
				Expect(err).NotTo(HaveOccurred())

				memo, err := memoRepository.FindByIDWithTags(ctx, fakeMemos[0].ID)
				Expect(err).NotTo(HaveOccurred())
				Expect(memo).NotTo(BeNil())
				Expect(memo.Edges.Tags).To(HaveLen(1))
				Expect(memo.Edges.Tags[0].Name).To(Equal(fakeTags[1].Name))
			})

			It("removes tags of memo", func(ctx SpecContext) {
				err := memoRepository.ReplaceTags(ctx, fakeMemos[0].ID, nil)
				Expect(err).NotTo(HaveOccurred())

				memo, err := memoRepository.FindByIDWithTags(ctx, fakeMemos[0].ID)
				Expect(err).NotTo(HaveOccurred())
				Expect(memo).NotTo(BeNil())
				Expect(memo.Edges.Tags).To(HaveLen(0))
			})

			It("returns not found error if id does not exist", func(ctx SpecContext) {
				err := memoRepository.ReplaceTags(ctx, uuid.Must(uuid.NewRandom()), nil)
				Expect(pkgerr.IsErrNotFound(err)).To(BeTrue())
			})
		})

		Context("RegisterSubscriber", func() {
			It("registers subscriber", func(ctx SpecContext) {
				err := memoRepository.RegisterSubscriber(ctx, fakeMemos[0].ID, fakeUsers[0].ID)
				Expect(err).NotTo(HaveOccurred())

				subs, err := userRepository.FindAllBySubscribingMemoID(ctx, fakeMemos[0].ID)
				Expect(err).NotTo(HaveOccurred())

				_, ok := lo.Find(subs, func(u *ent.User) bool { return u.ID == fakeUsers[0].ID })
				Expect(ok).To(BeTrue())
			})

			It("emits error if memo does not exist", func(ctx SpecContext) {
				err := memoRepository.RegisterSubscriber(ctx, uuid.Must(uuid.NewRandom()), fakeUsers[0].ID)
				Expect(err).To(HaveOccurred())
			})

			It("emits error if user does not exist", func(ctx SpecContext) {
				err := memoRepository.RegisterSubscriber(ctx, fakeMemos[0].ID, uuid.Must(uuid.NewRandom()))
				Expect(err).To(HaveOccurred())
			})

			It("emits error if user already subscribed", func(ctx SpecContext) {
				err := memoRepository.RegisterSubscriber(ctx, fakeMemos[0].ID, fakeUsers[1].ID)
				Expect(err).To(HaveOccurred())
			})
		})

		Context("UnregisterSubscriber", func() {
			It("unregisters subscriber", func(ctx SpecContext) {
				err := memoRepository.UnregisterSubscriber(ctx, fakeMemos[0].ID, fakeUsers[1].ID)
				Expect(err).NotTo(HaveOccurred())

				subs, err := userRepository.FindAllBySubscribingMemoID(ctx, fakeMemos[0].ID)
				Expect(err).NotTo(HaveOccurred())

				_, ok := lo.Find(subs, func(u *ent.User) bool { return u.ID == fakeUsers[1].ID })
				Expect(ok).To(BeFalse())
			})
		})

		Context("ClearSubscribers", func() {
			It("deletes all subscribers", func(ctx SpecContext) {
				err := memoRepository.ClearSubscribers(ctx, fakeMemos[0].ID)
				Expect(err).NotTo(HaveOccurred())

				subs, err := userRepository.FindAllBySubscribingMemoID(ctx, fakeMemos[0].ID)
				Expect(err).NotTo(HaveOccurred())

				_, ok := lo.Find(subs, func(u *ent.User) bool { return u.ID == fakeUsers[1].ID })
				Expect(ok).To(BeFalse())
			})
		})
	})
})
