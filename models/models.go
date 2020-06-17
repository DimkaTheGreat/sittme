package models

import (
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
)

const (
	streamType        = "stream"
	statusCreated     = "created"
	statusActivated   = "activate"
	statusFinished    = "finished"
	statusInterrupted = "interrupted"
)

var (
	wrongIDError = errors.New("Cant find translation with such ID")
)

type Translations map[string]*Translation

type Translation struct {
	Type       string        `json:"type"`
	ID         string        `json:"id"`
	Attributes *Attributes   `json:"attributes"`
	ActivateCh chan struct{} `json:"-"`
}

type Attributes struct {
	State     string    `json:"state"`
	CreatedAt time.Time `json:"created"`
}

func (translations Translations) LoadTestData() {
	for i := 0; i < 6; i++ {
		id := uuid.New()
		tr := &Translation{
			Type:       streamType,
			ID:         id.String(),
			ActivateCh: make(chan struct{}, 1),
			Attributes: &Attributes{
				State:     statusCreated,
				CreatedAt: time.Now(),
			},
		}
		translations[id.String()] = tr

	}

}

func DeleteByID(translations Translations, id string) error {

	rwMutex := &sync.RWMutex{}
	rwMutex.Lock()
	defer rwMutex.Unlock()

	_, ok := translations[id]
	if !ok {
		return wrongIDError
	}

	delete(translations, id)

	return nil

}

func ActivateByID(translations Translations, id string) error {

	tr, ok := translations[id]

	if !ok {
		return wrongIDError
	}

	if tr.Attributes.State == statusFinished {
		return errors.New("Cant activate already finished translation")

	}

	if tr.Attributes.State == statusInterrupted {
		tr.ActivateCh <- struct{}{}
		return nil

	}

	rwMutex := sync.RWMutex{}
	rwMutex.Lock()
	defer rwMutex.Unlock()
	tr.Attributes.State = statusActivated
	return nil

}

func InterruptTranslation(translations Translations, id string, timeout int) error {

	tr, ok := translations[id]

	if !ok {
		return wrongIDError
	}

	if tr.Attributes.State == statusFinished {
		return errors.New("Cant interrupt already finished translation")

	}

	if tr.Attributes.State == statusCreated {
		return errors.New("Cant interrupt not activated translation")

	}
	rwMutex := &sync.RWMutex{}
	rwMutex.Lock()
	tr.Attributes.State = statusInterrupted
	rwMutex.Unlock()

	tr.runTimer(timeout)

	return nil

}

func (tr *Translation) runTimer(timeout int) {
	go func() {
		timer := time.NewTimer(time.Second * time.Duration(timeout))
		select {
		case <-timer.C:
			rwMutex := &sync.RWMutex{}
			rwMutex.Lock()
			tr.Attributes.State = statusFinished
			rwMutex.Unlock()
		case <-tr.ActivateCh:
			timer.Stop()
			rwMutex := &sync.RWMutex{}
			rwMutex.Lock()
			tr.Attributes.State = statusActivated
			rwMutex.Unlock()
		}
	}()

}
