package server

import (
	"encoding/json"
	"net/http"

	pb "github.com/nrmnqdds/gomaluum/internal/proto"
)

// @Title LoginHandler
// @Description Logs in the user. Save the token and use it in the Authorization header for future requests.
// @Tags auth
// @Accept json
// @Produce json
// @Param body body dtos.LoginProps true "Login properties"
// @Success 200 {object} dtos.LoginResponseDTO
// @Router /auth/login [post]
func (s *Server) LoginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	logger := s.log.GetLogger()

	user := &pb.LoginRequest{}

	if err := json.NewDecoder(r.Body).Decode(user); err != nil {
		logger.Sugar().Errorf("Failed to decode request body: %v", err)
		_, _ = w.Write([]byte("Failed to decode request body"))
		return
	}

	// casCookie, respUsername, respPassword, err := s.Login(r.Context(), user)
	resp, err := s.Login(r.Context(), user)
	if err != nil {
		logger.Sugar().Errorf("Failed to login: %v", err)
		_, _ = w.Write([]byte("Failed to login"))
		return
	}

	newCookie, _, err := s.GeneratePasetoToken(resp.Token, resp.Username, resp.Password)
	if err != nil {
		logger.Sugar().Errorf("Failed to generate PASETO token: %v", err)
		_, _ = w.Write([]byte("Failed to generate PASETO token"))
		return
	}

	result := &pb.LoginResponse{
		Token:    newCookie,
		Username: resp.Username,
	}

	jsonResp, err := json.Marshal(result)
	if err != nil {
		logger.Sugar().Errorf("Failed to marshal JSON: %v", err)
		_, _ = w.Write([]byte("Failed to marshal JSON"))
		return
	}

	_, _ = w.Write(jsonResp)
}
