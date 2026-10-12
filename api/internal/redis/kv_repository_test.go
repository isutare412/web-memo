package redis

import (
	"time"

	"github.com/go-redis/redismock/v9"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/isutare412/web-memo/api/internal/pkgerr"
)

var _ = Describe("KvRepository", func() {
	var (
		kvRepository *KVRepository
		mockRedis    redismock.ClientMock
	)

	BeforeEach(func() {
		redisClient, mock := redismock.NewClientMock()
		client := &Client{innerClient: redisClient}

		kvRepository = NewKVRepository(client)
		mockRedis = mock
	})

	Context("Get", func() {
		It("returns value if found", func(ctx SpecContext) {
			var (
				givenKey   = "foo"
				givenValue = "bar"
			)

			mockRedis.
				ExpectGet(givenKey).
				SetVal(givenValue)

			value, err := kvRepository.Get(ctx, givenKey)
			Expect(err).NotTo(HaveOccurred())
			Expect(value).To(Equal(givenValue))
		})

		It("emits not found error if key does not exist", func(ctx SpecContext) {
			var (
				givenKey = "foo"
			)

			mockRedis.
				ExpectGet(givenKey).
				RedisNil()

			_, err := kvRepository.Get(ctx, givenKey)
			Expect(pkgerr.IsErrNotFound(err)).To(BeTrue())
		})
	})

	Context("GetThenDelete", func() {
		It("returns value if found", func(ctx SpecContext) {
			var (
				givenKey   = "foo"
				givenValue = "bar"
			)

			mockRedis.
				ExpectGetDel(givenKey).
				SetVal(givenValue)

			value, err := kvRepository.GetThenDelete(ctx, givenKey)
			Expect(err).NotTo(HaveOccurred())
			Expect(value).To(Equal(givenValue))
		})

		It("emits not found error if key does not exist", func(ctx SpecContext) {
			var (
				givenKey = "foo"
			)

			mockRedis.
				ExpectGetDel(givenKey).
				RedisNil()

			_, err := kvRepository.GetThenDelete(ctx, givenKey)
			Expect(pkgerr.IsErrNotFound(err)).To(BeTrue())
		})
	})

	Context("Set", func() {
		It("sets value by key", func(ctx SpecContext) {
			var (
				givenKey   = "foo"
				givenValue = "bar"
				givenTTL   = time.Hour
			)

			mockRedis.
				ExpectSet(givenKey, givenValue, givenTTL).
				SetVal("OK")

			err := kvRepository.Set(ctx, givenKey, givenValue, givenTTL)
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Context("Delete", func() {
		It("works as expected", func(ctx SpecContext) {
			var (
				givenKeys = []string{
					"foo",
					"bar",
				}
			)

			mockRedis.
				ExpectDel(givenKeys...).
				SetVal(int64(len(givenKeys)))

			count, err := kvRepository.Delete(ctx, givenKeys...)
			Expect(err).NotTo(HaveOccurred())
			Expect(count).To(BeEquivalentTo(len(givenKeys)))
		})
	})

	Context("with key prefix", func() {
		var givenPrefix = "webmemo:test:"

		BeforeEach(func() {
			redisClient, mock := redismock.NewClientMock()
			client := &Client{innerClient: redisClient, keyPrefix: givenPrefix}

			kvRepository = NewKVRepository(client)
			mockRedis = mock
		})

		It("prefixes key on Get", func(ctx SpecContext) {
			mockRedis.
				ExpectGet(givenPrefix + "foo").
				SetVal("bar")

			value, err := kvRepository.Get(ctx, "foo")
			Expect(err).NotTo(HaveOccurred())
			Expect(value).To(Equal("bar"))
		})

		It("prefixes key on GetThenDelete", func(ctx SpecContext) {
			mockRedis.
				ExpectGetDel(givenPrefix + "foo").
				SetVal("bar")

			value, err := kvRepository.GetThenDelete(ctx, "foo")
			Expect(err).NotTo(HaveOccurred())
			Expect(value).To(Equal("bar"))
		})

		It("prefixes key on Set", func(ctx SpecContext) {
			mockRedis.
				ExpectSet(givenPrefix+"foo", "bar", time.Hour).
				SetVal("OK")

			err := kvRepository.Set(ctx, "foo", "bar", time.Hour)
			Expect(err).NotTo(HaveOccurred())
		})

		It("prefixes keys on Delete", func(ctx SpecContext) {
			mockRedis.
				ExpectDel(givenPrefix+"foo", givenPrefix+"bar").
				SetVal(2)

			count, err := kvRepository.Delete(ctx, "foo", "bar")
			Expect(err).NotTo(HaveOccurred())
			Expect(count).To(BeEquivalentTo(2))
		})
	})
})
