package jwt_test

import (
	"time"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/isutare412/web-memo/api/internal/core/model"
	"github.com/isutare412/web-memo/api/internal/jwt"
)

var _ = Describe("Client", func() {
	Context("client methods", func() {
		var jwtClient *jwt.Client

		BeforeEach(func() {
			var (
				givenConfig = jwt.Config{
					ActiveKeyPair: "unit-test",
					KeyPairs: []jwt.RSAKeyBytesPair{{
						Name:    "unit-test",
						Private: "-----BEGIN RSA PRIVATE KEY-----\nMIIEowIBAAKCAQEA2fHrrI8PPNAADU1Dwyt7d1P3OLyErjHfNvsrJPJPUCpXk6qT\nnq+RljQKgVotgZJgD2Uq1qFkpE34frRo8sajmJLWYIlqAqGWyVm+083/K9hYUAKy\nKfgcpMbB4Bal+iQe6PmM+FrS9iNiYtfUE/VOEnOWJ5i9MzkpzutB6lYfH5UNqEHM\npFi3nkOillLAAFBx3KS6twF/KmomG6eKcZvUF6VDzPkT3rYYA+OGmPR+TQC91WqM\nh4PFMnmMFOWl/NYgqYpU1gu5zNfcR7PH4+Nrd9tlB0JZWvKc2CZRr/JWR9cNspfc\nh6tiuxcDhlevW3xkjqkPXSBIsZYUSEo+DMbS+QIDAQABAoIBAQCzmS4y1vkhje5d\nB8zNammmHeVeNcvImIRvFG+rVJqTXzEoGKrpW5jHhD8b9SoG2o6aYS36DKmY1D/w\nec76MULjGi3bA6H+ZRaS+ofMlraudOvhkzcqarF//+IIPeGszlNCfVLz3jR2bRZI\nib7Ua5NYlTOpka5oJbUUNL3u0+V16aIRbrAfDj1m8bJGjZEyPvfNDj5nwAXlWNHa\n+hM11v+kvoWC7QD3rVlA82s/sjkL5f3oZzh/8lUJ36CpnFl58Ti9X9oHNNYieNmt\n4PQj5XCXrxy5IgY9L8QPPeqZvIH6yJyHU6urmzZs+U9o/GrDDcs82aPK/p0YPspn\nBAQgM6J5AoGBAP8YWd+f9wBSGfoyo4xgzWRZ//l0InCavEBgQ4UffN3+u2Xs+Lf+\ntT6RxdFbMWguk4IAmgrF61xwnwl5NcJrwbOW38qz70egYCXna0hCR1CeqJXDMdDs\n8NWq5GhYoMzMxqkWTcNajbmCOEdXG1BDhoc4rOYdXTiUkpVuvZaRuvfnAoGBANq3\n1X0QX2tpGeyIlTxnIYvdiD/1Ku4wRjuRme9iK5nmwHbmtgrVRR0ReCjoG+Bm9N8B\nrhRC4EG3QsEim5OixsTGBgm8G2pudMquAp5Em9i6Al9IKoMktEng8KyolFjSL7XK\nrR3fWmvCeGHwXMlG4snhBVNeG1z8bhJH5k4i1QIfAoGAblBkhT1S/nOCwlzltw4h\npDT7ai/buOBhamF2sXn1cLb46VH6GO9wB5fYePm7uvbxWTXTZ7dBWd9mFx2wrtwU\njwo+yxTW9B2ZlqqmDUCEQIvsEZ+wyk28tFnLnog9OXOQsYxwontlcISsu8Uijao1\n4gITWwv1xUMSxMZ3/EYXGZ8CgYAEdkWjtIbN0SoXOj4ZKl1z1gQmkeDbVR7JrlG+\noXkUPbHlexVxqSIs0qNp6jpPXKpYNleP0EF09cEl4Yfc/jAh6YxL/Ituo2w8ikpB\nYlLvm/Pab8V2QXRwIWenjhTgrwEMK2NWvazBkAkWrmmmLY0I409RgRT706aHNvJK\n28kOrwKBgBYr8VZv8JmsLlZoABB8BuPhWJ6/9m7rwGc4ipRoeCjFZaUL8E9OWHvs\nA0hpq1tVb8tvecDFlW2vHKcNN/5p1tMwxmGrs07SFl+L9fLWRAsegwqqDZXCh1CQ\njbvvwZjPvY/DlTkU7RspYC7cvs1IafaWEVfIUkEuWDcEim6cQ2h8\n-----END RSA PRIVATE KEY-----\n",
						Public:  "-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA2fHrrI8PPNAADU1Dwyt7\nd1P3OLyErjHfNvsrJPJPUCpXk6qTnq+RljQKgVotgZJgD2Uq1qFkpE34frRo8saj\nmJLWYIlqAqGWyVm+083/K9hYUAKyKfgcpMbB4Bal+iQe6PmM+FrS9iNiYtfUE/VO\nEnOWJ5i9MzkpzutB6lYfH5UNqEHMpFi3nkOillLAAFBx3KS6twF/KmomG6eKcZvU\nF6VDzPkT3rYYA+OGmPR+TQC91WqMh4PFMnmMFOWl/NYgqYpU1gu5zNfcR7PH4+Nr\nd9tlB0JZWvKc2CZRr/JWR9cNspfch6tiuxcDhlevW3xkjqkPXSBIsZYUSEo+DMbS\n+QIDAQAB\n-----END PUBLIC KEY-----\n",
					}},
					Expiration: time.Hour,
				}
			)

			client, err := jwt.NewClient(givenConfig)
			Expect(err).ShouldNot(HaveOccurred())

			jwtClient = client
		})

		It("verifies signed app ID token", func() {
			var (
				givenAppIDToken = &model.AppIDToken{
					UserID:     uuid.New(),
					Email:      "fake-email",
					UserName:   "Alice Bob",
					FamilyName: "Bob",
					GivenName:  "Alice",
					PhotoURL:   "https://google.com",
				}
			)

			signedToken, tokenString, err := jwtClient.SignAppIDToken(givenAppIDToken)
			Expect(err).ShouldNot(HaveOccurred())

			gotAppIDToken, err := jwtClient.VerifyAppIDToken(tokenString)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(gotAppIDToken).Should(Equal(signedToken))
		})
	})
})
