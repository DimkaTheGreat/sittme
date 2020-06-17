package routing

import (
	"net/http"
	"sittme/models"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo"
)

const (
	messageBadRequest = "Bad Request"
)

type handler struct {
	Translations map[string]*models.Translation
	Timeout      int
	mu           *sync.RWMutex
}

func Run(tr models.Translations, timeout int, port string) {
	h := handler{
		Translations: tr,
		Timeout:      timeout,
	}

	e := echo.New()
	e.HideBanner = true
	e.GET("/create", h.createID)
	e.GET("/list", h.listTranslations)
	e.DELETE("/delete", h.deleteTranslation)
	e.GET("/activate", h.activateTranslation)
	e.GET("/interrupt", h.interruptTranslation)

	e.Start(":" + port)

}

func (h *handler) createID(c echo.Context) error {
	id := uuid.New()

	tr := &models.Translation{
		ID:         id.String(),
		Type:       "stream",
		ActivateCh: make(chan struct{}, 1),
		Attributes: &models.Attributes{
			State:     "created",
			CreatedAt: time.Now(),
		},
	}

	h.mu.Lock()
	h.Translations[id.String()] = tr
	h.mu.Unlock()
	c.String(http.StatusOK, "OK")
	return nil
}

func (h *handler) listTranslations(c echo.Context) error {

	var translations []*models.Translation

	for _, translation := range h.Translations {
		translations = append(translations, translation)
	}
	c.JSONPretty(http.StatusOK, translations, "  ")

	return nil

}

func (h *handler) deleteTranslation(c echo.Context) error {
	id := c.QueryParam("id")

	if id == "" {
		c.JSON(http.StatusBadRequest, messageBadRequest)
		return nil
	}
	err := models.DeleteByID(h.Translations, id)

	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return nil
	}
	c.JSON(http.StatusOK, "OK")
	return nil

}

func (h *handler) activateTranslation(c echo.Context) error {

	id := c.QueryParam("id")

	if id == "" {
		c.JSON(http.StatusBadRequest, messageBadRequest)
		return nil
	}

	err := models.ActivateByID(h.Translations, id)

	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return nil
	}
	c.JSON(http.StatusOK, "OK")
	return nil

}

func (h *handler) interruptTranslation(c echo.Context) error {
	id := c.QueryParam("id")

	if id == "" {
		c.JSON(http.StatusBadRequest, messageBadRequest)
		return nil
	}

	err := models.InterruptTranslation(h.Translations, id, h.Timeout)

	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return nil
	}
	c.JSON(http.StatusOK, "OK")
	return nil
}
