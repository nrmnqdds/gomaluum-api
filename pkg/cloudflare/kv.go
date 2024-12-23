package cloudflare

import (
	"log"
	"os"

	"github.com/cloudflare/cloudflare-go"
)

type AppCloudflare struct {
	*cloudflare.API
}

func New() *AppCloudflare {
	cloudflareClient, err := cloudflare.NewWithAPIToken(os.Getenv("CLOUDFLARE_API_TOKEN"))
	if err != nil {
		log.Fatalf("Error initiating cloudflare client: %s", err.Error())
	}

	return &AppCloudflare{
		API: cloudflareClient,
	}
}

func (s *AppCloudflare) GetClient() *cloudflare.API {
	return s.API
}
