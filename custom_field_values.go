package activecampaign

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
)

type FieldValues struct {
	FieldValues []FieldValue    `json:"fieldValues"`
	Meta        FieldValuesMeta `json:"meta"`
}
type FieldValuesMeta struct {
	Total string `json:"total"`
}
type FieldValue struct {
	Contact    string `json:"contact"`
	Field      string `json:"field"`
	Value      string `json:"value"`
	CreateDate string `json:"cdate"`
	UpdateDate string `json:"udate"`
	ID         string `json:"id"`
}

func (a *ActiveCampaign) fieldValues(ctx context.Context, pof *POF, url string) (*FieldValues, error) {
	res, err := a.send(ctx, http.MethodGet, url, pof, nil)
	if err != nil {
		return nil, err
	}

	var values FieldValues
	err = json.NewDecoder(res.Body).Decode(&values)
	if err != nil {
		return nil, err
	}

	return &values, nil
}

func (a *ActiveCampaign) FieldValues(ctx context.Context, pof *POF) (*FieldValues, error) {
	return a.fieldValues(ctx, pof, "fieldValues")
}

type ChangeFieldValue struct {
	Contact string `json:"contact"`
	Field   string `json:"field"`
	Value   string `json:"value"`
}

func (a *ActiveCampaign) CreateFieldValue(ctx context.Context, create ChangeFieldValue) error {
	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(struct {
		FieldValue ChangeFieldValue `json:"fieldValue"`
	}{
		FieldValue: create,
	})
	if err != nil {
		return err
	}

	res, err := a.send(ctx, http.MethodPost, "fieldValues", nil, b)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusCreated {
		return errors.New(res.Status)
	}

	return nil
}

func (a *ActiveCampaign) UpdateFieldValue(ctx context.Context, id string, update ChangeFieldValue) error {
	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(struct {
		FieldValue ChangeFieldValue `json:"fieldValue"`
	}{
		FieldValue: update,
	})
	if err != nil {
		return err
	}

	res, err := a.send(ctx, http.MethodPut, "fieldValues/"+id, nil, b)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusCreated {
		return errors.New(res.Status)
	}

	return nil
}
