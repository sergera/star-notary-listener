package backend

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/sergera/star-notary-listener/internal/conf"
	"github.com/sergera/star-notary-listener/internal/logger"
	"github.com/sergera/star-notary-listener/internal/models"
)

type Backend struct {
	host        string
	port        string
	contentType string
	client      *http.Client
}

func NewBackend() *Backend {
	conf := conf.GetConf()
	return &Backend{
		conf.BackendHost,
		conf.BackendPort,
		"application/json; charset=UTF-8",
		&http.Client{},
	}
}

func (b Backend) Post(route string, jsonData []byte) error {
	log.Println("host: ", b.host)
	log.Println("port: ", b.port)
	log.Println("route: ", route)
	log.Println("full url: ", b.host+":"+b.port+"/"+route)

	request, err := http.NewRequest("POST", b.host+":"+b.port+"/"+route, bytes.NewBuffer(jsonData))
	if err != nil {
		logger.Error(
			"Failed to create post request",
			logger.String("message", err.Error()),
			logger.String("route", route),
			logger.String("request", string(jsonData)),
		)
		return err
	}

	request.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		responseBody, _ := ioutil.ReadAll(response.Body)
		responseBodyString := string(responseBody)
		logger.Error(
			"Failed backend post request",
			logger.String("message", err.Error()),
			logger.String("route", route),
			logger.String("request", string(jsonData)),
			logger.String("status", response.Status),
			logger.String("response", responseBodyString),
		)
		return err
	}

	defer response.Body.Close()
	return nil
}

func (b Backend) Put(route string, jsonData []byte) error {
	request, err := http.NewRequest("PUT", b.host+":"+b.port+"/"+route, bytes.NewBuffer(jsonData))
	if err != nil {
		logger.Error(
			"Failed to create put request",
			logger.String("message", err.Error()),
			logger.String("route", route),
			logger.String("request", string(jsonData)),
		)
		return err
	}

	request.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		responseBody, _ := ioutil.ReadAll(response.Body)
		responseBodyString := string(responseBody)
		logger.Error(
			"Failed backend put request",
			logger.String("message", err.Error()),
			logger.String("route", route),
			logger.String("request", string(jsonData)),
			logger.String("status", response.Status),
			logger.String("response", responseBodyString),
		)
		return err
	}

	defer response.Body.Close()
	return nil
}

func (b Backend) CreateStar(e models.CreatedEvent) error {
	m, err := json.Marshal(e)
	if err != nil {
		logger.Error(
			"Failed to marshal event model into json",
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

func (b Backend) ChangeName(e models.ChangedNameEvent) error {
	m, err := json.Marshal(e)
	if err != nil {
		logger.Error(
			"Failed to marshal event model into json",
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

func (b Backend) PutForSale(e models.PutForSaleEvent) error {
	m, err := json.Marshal(e)
	if err != nil {
		logger.Error(
			"Failed to marshal event model into json",
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

func (b Backend) RemoveFromSale(e models.RemovedFromSaleEvent) error {
	m, err := json.Marshal(e)
	if err != nil {
		logger.Error(
			"Failed to marshal event model into json",
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

func (b Backend) Sell(e models.SoldEvent) error {
	m, err := json.Marshal(e)
	if err != nil {
		logger.Error(
			"Failed to marshal event model into json",
			logger.String("message", err.Error()),
			logger.Object("event", &e),
		)
		return err
	}

	err = b.Put("sell", m)
	if err != nil {
		return err
	}

	return nil
}
