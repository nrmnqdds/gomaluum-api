package server

import (
	"net/http"

	"github.com/MarceloPetrucio/go-scalar-api-reference"
)

func (s *Server) ScalarReference(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	logger := s.log.GetLogger()
	swaggerContent, err := DocsPath.ReadFile("docs/swagger/swagger.json")
	if err != nil {
		logger.Sugar().Fatalf("could not read swagger.json: %v", err)
	}

	htmlContent, err := scalar.ApiReferenceHTML(&scalar.Options{
		SpecContent: string(swaggerContent),
		CustomOptions: scalar.CustomOptions{
			PageTitle: "GoMaluum API",
		},
		DarkMode: true,
	})
	if err != nil {
		logger.Sugar().Errorf("%v", err)
	}

	_, _ = w.Write([]byte(htmlContent))
}
