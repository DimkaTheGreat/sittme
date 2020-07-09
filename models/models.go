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

	mu := &sync.Mutex{}

	_, ok := translations[id]
	if !ok {
		return wrongIDError
	}
	mu.Lock()
	delete(translations, id)
	mu.Unlock()

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

	mu := &sync.Mutex{}
	mu.Lock()
	tr.Attributes.State = statusActivated
	mu.Unlock()
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
	mu := &sync.Mutex{}
	mu.Lock()
	tr.Attributes.State = statusInterrupted
	mu.Unlock()

	tr.runTimer(timeout)

	return nil

}

func (tr *Translation) runTimer(timeout int) {
	go func() {
		timer := time.NewTimer(time.Second * time.Duration(timeout))
		select {
		case <-timer.C:
			mu := &sync.Mutex{}
			mu.Lock()
			tr.Attributes.State = statusFinished
			mu.Unlock()
			close(tr.ActivateCh)
		case <-tr.ActivateCh:
			timer.Stop()
			mu := &sync.Mutex{}
			mu.Lock()
			tr.Attributes.State = statusActivated
			mu.Unlock()
		}
	}()

}
