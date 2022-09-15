package service

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/sergera/star-notary-listener/internal/conf"
	"github.com/sergera/star-notary-listener/internal/domain"
	"github.com/sergera/star-notary-listener/internal/logger"
)

type StarNotaryAPIService struct {
	host        string
	port        string
	contentType string
	client      *http.Client
}

func NewStarNotaryAPIService() *StarNotaryAPIService {
	conf := conf.GetConf()
	return &StarNotaryAPIService{
		conf.StarNotaryAPIHost(),
		conf.StarNotaryAPIPort(),
		"application/json; charset=UTF-8",
		&http.Client{},
	}
}

func (b StarNotaryAPIService) Post(route string, jsonData []byte) error {
	request, err := http.NewRequest("POST", b.host+":"+b.port+"/"+route, bytes.NewBuffer(jsonData))
	if err != nil {
		logger.Error(
			"failed to create post request",
			logger.String("message", err.Error()),
		)
		return err
	}

	request.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		logger.Error(
			"failed post request",
			logger.String("message", err.Error()),
			logger.String("status", response.Status),
		)
		return err
	}

	defer response.Body.Close()
	return nil
}

func (b StarNotaryAPIService) Put(route string, jsonData []byte) error {
	request, err := http.NewRequest("PUT", b.host+":"+b.port+"/"+route, bytes.NewBuffer(jsonData))
	if err != nil {
		logger.Error(
			"failed to create put request",
			logger.String("message", err.Error()),
		)
		return err
	}

	request.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		logger.Error(
			"failed put request",
			logger.String("message", err.Error()),
			logger.String("status", response.Status),
		)
		return err
	}

	defer response.Body.Close()
	return nil
}

func (b StarNotaryAPIService) CreateStar(e domain.CreateEvent) error {
	m, err := json.Marshal(e)
	if err != nil {
		logger.Error(
			"failed to marshal event model into json",
			logger.String("message", err.Error()),
			logger.Object("event", &e),
		)
		return err
	}

	err = b.Post("create", m)
	if err != nil {
		return err
	}

	return nil
}

func (b StarNotaryAPIService) ChangeName(e domain.ChangeNameEvent) error {
	m, err := json.Marshal(e)
	if err != nil {
		logger.Error(
			"failed to marshal event model into json",
			logger.String("message", err.Error()),
			logger.Object("event", &e),
		)
		return err
	}

	err = b.Put("set-name", m)
	if err != nil {
		return err
	}

	return nil
}

func (b StarNotaryAPIService) PutForSale(e domain.PutForSaleEvent) error {
	m, err := json.Marshal(e)
	if err != nil {
		logger.Error(
			"failed to marshal event model into json",
			logger.String("message", err.Error()),
			logger.Object("event", &e),
		)
		return err
	}

	err = b.Put("set-price", m)
	if err != nil {
		return err
	}

	return nil
}

func (b StarNotaryAPIService) RemoveFromSale(e domain.RemoveFromSaleEvent) error {
	m, err := json.Marshal(e)
	if err != nil {
		logger.Error(
			"failed to marshal event model into json",
			logger.String("message", err.Error()),
			logger.Object("event", &e),
		)
		return err
	}

	err = b.Put("remove-from-sale", m)
	if err != nil {
		return err
	}

	return nil
}

func (b StarNotaryAPIService) Purchase(e domain.PurchaseEvent) error {
	m, err := json.Marshal(e)
	if err != nil {
		logger.Error(
			"failed to marshal event model into json",
			logger.String("message", err.Error()),
			logger.Object("event", &e),
		)
		return err
	}

	err = b.Put("purchase", m)
	if err != nil {
		return err
	}

	return nil
}
