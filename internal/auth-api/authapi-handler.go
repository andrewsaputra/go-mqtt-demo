package authapi

import (
	"andrewsaputra/go-mqtt-demo/internal/util"
	"encoding/base64"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewAuthApiHandler(cryptKey []byte, sessionManager SessionManager) (*AuthApiHandler, error) {
	crypt, err := util.NewAesCrypt(cryptKey)
	if err != nil {
		return nil, err
	}

	return &AuthApiHandler{
		Crypt:          crypt,
		SessionManager: sessionManager,
	}, nil
}

type AuthApiHandler struct {
	Crypt          util.Crypt
	SessionManager SessionManager
}

func (this AuthApiHandler) NewSession(c *gin.Context) {
	sessionData, err := this.SessionManager.GenerateNewSession()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	jsonBytes, err := json.Marshal(sessionData)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	encryptedToken, err := this.Crypt.Encrypt(jsonBytes)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	authToken := base64.StdEncoding.EncodeToString(encryptedToken)

	body := make(map[string]string)
	body["user_id"] = sessionData.UserId
	body["auth_token"] = authToken
	c.JSON(http.StatusOK, body)
}
