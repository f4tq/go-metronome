package metronome_test

import (
	"net/http"

	. "github.com/adobe-platform/go-metronome/metronome"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	ghttp "github.com/onsi/gomega/ghttp"
)

var _ = Describe("Client", func() {
	var (
		config_stub Config
		server      *ghttp.Server
	)

	BeforeEach(func() {
		server = ghttp.NewServer()

		config_stub = Config{
			URL:            server.URL(),
			Debug:          false,
			RequestTimeout: 5,
		}
	})

	AfterEach(func() {
		server.Close()
	})

	Describe("NewClient", func() {
		It("Returns a new client", func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/v1/jobs"),
				),
			)

			client, err := NewClient(config_stub)

			Expect(client).To(BeAssignableToTypeOf(new(Client)))
			Expect(err).To(BeNil())
		})

		It("Defaults to unverifiedtls being false", func() {
			test_config := Config{
				URL:            server.URL(),
				Debug:          false,
				RequestTimeout: 5,
			}

			Expect(test_config.AllowUnverifiedTLS).To(BeFalse())
		})

		It("Errors if it cannot hit metronome", func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/v1/jobs"),
					ghttp.RespondWith(http.StatusInternalServerError, nil),
				),
			)

			_, err := NewClient(config_stub)
			Expect(err).To(MatchError("Could not reach metronome cluster: 500 Internal Server Error"))
		})
	})
})
