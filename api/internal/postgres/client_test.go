package postgres

import (
	"context"
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/isutare412/web-memo/api/internal/core/ent"
	"github.com/isutare412/web-memo/api/internal/core/enum"
	"github.com/isutare412/web-memo/api/internal/pkgerr"
)

var _ = Describe("Client", func() {
	var (
		client         *Client
		userRepository *UserRepository
	)

	BeforeEach(func(ctx SpecContext) {
		entClient, err := ent.Open("sqlite3", "file:ent?mode=memory")
		Expect(err).NotTo(HaveOccurred())

		client = &Client{entClient: entClient}
		err = client.MigrateSchemas(ctx)
		Expect(err).NotTo(HaveOccurred())

		userRepository = NewUserRepository(client)
	})

	Context("BeginTx", func() {
		It("commits as expected", func(ctx SpecContext) {
			var (
				givenEmail = "test-user@gmail.com"
				givenUser  = &ent.User{
					Email:    givenEmail,
					UserName: "tester-lee",
					Type:     enum.UserTypeClient,
				}
			)

			ctxWithTx, commit, _ := client.BeginTx(ctx)

			userCreated, err := userRepository.Upsert(ctxWithTx, givenUser)
			Expect(err).NotTo(HaveOccurred())

			err = commit()
			Expect(err).NotTo(HaveOccurred())

			userFound, err := userRepository.FindByEmail(ctx, givenEmail)
			Expect(err).NotTo(HaveOccurred())
			Expect(userFound.ID).To(Equal(userCreated.ID))
		})

		It("rollbacks as expected", func(ctx SpecContext) {
			var (
				givenEmail = "test-user@gmail.com"
				givenUser  = &ent.User{
					Email:    givenEmail,
					UserName: "tester-lee",
					Type:     enum.UserTypeClient,
				}
			)

			ctxWithTx, _, rollback := client.BeginTx(ctx)

			_, err := userRepository.Upsert(ctxWithTx, givenUser)
			Expect(err).NotTo(HaveOccurred())

			err = rollback()
			Expect(err).NotTo(HaveOccurred())

			_, err = userRepository.FindByEmail(ctx, givenEmail)
			Expect(pkgerr.IsErrNotFound(err)).To(BeTrue())
		})
	})

	Context("WithTx", func() {
		It("commits if no error", func(ctx SpecContext) {
			var (
				givenEmail = "test-user@gmail.com"
				givenUser  = &ent.User{
					Email:    givenEmail,
					UserName: "tester-lee",
					Type:     enum.UserTypeClient,
				}
			)

			var userCreated *ent.User
			err := client.WithTx(ctx, func(ctx context.Context) error {
				user, err := userRepository.Upsert(ctx, givenUser)
				Expect(err).NotTo(HaveOccurred())

				userCreated = user
				return nil
			})
			Expect(err).NotTo(HaveOccurred())

			userFound, err := userRepository.FindByEmail(ctx, givenEmail)
			Expect(err).NotTo(HaveOccurred())
			Expect(userFound.ID).To(Equal(userCreated.ID))
		})

		It("rollbacks if error occurred", func(ctx SpecContext) {
			var (
				givenEmail = "test-user@gmail.com"
				givenUser  = &ent.User{
					Email:    givenEmail,
					UserName: "tester-lee",
					Type:     enum.UserTypeClient,
				}
			)

			err := client.WithTx(ctx, func(ctx context.Context) error {
				_, err := userRepository.Upsert(ctx, givenUser)
				Expect(err).NotTo(HaveOccurred())

				return fmt.Errorf("some unexpected error")
			})
			Expect(err).To(HaveOccurred())

			_, err = userRepository.FindByEmail(ctx, givenEmail)
			Expect(pkgerr.IsErrNotFound(err)).To(BeTrue())
		})
	})
})
