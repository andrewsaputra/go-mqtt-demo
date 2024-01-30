package main

import (
	"andrewsaputra/go-mqtt-demo/internal/util"
	"crypto/rand"
	"encoding/base64"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

const ADDRESS = ":5000"

var startTime time.Time = time.Now()

func main() {
	router := gin.Default()
	router.GET("/status", statusCheck)
	router.GET("/generatekey", generateEncryptionKey)
	router.POST("/encrypt", encrypt)
	router.POST("/decrypt", decrypt)

	router.Run(ADDRESS)
}

func statusCheck(c *gin.Context) {
	body := make(map[string]string)
	body["status"] = "Healthy"
	body["started_at"] = startTime.Format(time.RFC822Z)

	c.JSON(http.StatusOK, body)
}

func generateEncryptionKey(c *gin.Context) {
	newKey := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, newKey); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	encoded := base64.StdEncoding.EncodeToString(newKey)
	c.String(http.StatusOK, "%s", encoded)
}

type EncDecRequest struct {
	B64key string `json:"b64key"`
	Text   string `json:"text"`
}

func encrypt(c *gin.Context) {
	var request EncDecRequest
	if err := c.BindJSON(&request); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	key, err := base64.StdEncoding.DecodeString(request.B64key)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	crypt, err := util.NewAesCrypt(key)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	encrypted, err := crypt.EncryptString(request.Text)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	c.String(http.StatusOK, encrypted)
}

func decrypt(c *gin.Context) {
	var request EncDecRequest
	if err := c.Bind(&request); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	key, err := base64.StdEncoding.DecodeString(request.B64key)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	crypt, err := util.NewAesCrypt(key)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	decrypted, err := crypt.DecryptString(request.Text)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	c.String(http.StatusOK, decrypted)
}
