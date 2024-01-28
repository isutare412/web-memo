package redis

import (
	"errors"

	"github.com/go-redis/redismock/v9"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Client", func() {
	var (
		client    *Client
		mockRedis redismock.ClientMock
	)

	BeforeEach(func() {
		redisClient, mock := redismock.NewClientMock()

		client = &Client{innerClient: redisClient}
		mockRedis = mock
	})

	AfterEach(func() {
		err := mockRedis.ExpectationsWereMet()
		Expect(err).NotTo(HaveOccurred())
	})

	Context("Shutdown", func() {
		It("return nil if redis close successful", func(ctx SpecContext) {
			err := client.Shutdown(ctx)
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Context("Ping", func() {
		It("success as expected", func(ctx SpecContext) {
			mockRedis.
				ExpectPing().
				SetVal("OK")

			err := client.Ping(ctx)
			Expect(err).NotTo(HaveOccurred())
		})

		It("emits error if ping fails", func(ctx SpecContext) {
			mockRedis.
				ExpectPing().
				SetErr(errors.New("ping failure"))

			err := client.Ping(ctx)
			Expect(err).To(HaveOccurred())
		})
	})
})
