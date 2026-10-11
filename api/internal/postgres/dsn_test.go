package postgres

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("buildDSN", func() {
	givenConfig := Config{
		Host:     "db.example",
		Port:     5432,
		User:     "tester",
		Password: "secret",
		Database: "memo",
	}

	It("defaults sslmode to disable", func() {
		Expect(buildDSN(givenConfig)).To(Equal(
			"host=db.example port=5432 user=tester password=secret dbname=memo sslmode=disable"))
	})

	It("uses the configured sslmode", func() {
		cfg := givenConfig
		cfg.SSLMode = "require"

		Expect(buildDSN(cfg)).To(Equal(
			"host=db.example port=5432 user=tester password=secret dbname=memo sslmode=require"))
	})
})
