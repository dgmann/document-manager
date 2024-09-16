package client

import "github.com/dgmann/document-manager/pkg/api"

type PatientClient struct {
	*httpClient
}

func (c *PatientClient) Get(id string) (api.Patient, error) {
	res, err := c.GetJson("patients/" + id)
	if err != nil {
		return api.Patient{}, err
	}
	return HandleResponse[api.Patient](res)
}

func (c *PatientClient) All() ([]api.Patient, error) {
	res, err := c.GetJson("patients")
	if err != nil {
		return nil, err
	}
	return HandleResponse[[]api.Patient](res)
}
