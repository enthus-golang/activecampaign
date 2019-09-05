package activecampaign

import (
	"context"
	"encoding/json"
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
