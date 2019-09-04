package activecampaign

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
)

type ContactCreate struct {
	Email     string `json:"email"`
	FirstName string `json:"firstName,omitempty"`
	LastName  string `json:"lastName,omitempty"`
	Phone     int    `json:"phone,omitempty"`
}

func (a *ActiveCampaign) ContactCreate(ctx context.Context, contact ContactCreate) error {
	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(struct {
		Contact ContactCreate `json:"contact"`
	}{
		Contact: contact,
	})
	if err != nil {
		return err
	}

	res, err := a.send(ctx, http.MethodPost, "contacts", nil, b)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusCreated {
		return errors.New(res.Status)
	}

	return nil
}

//func (a *ActiveCampaign) ContactUpdate(ctx context.Context) error {
//
//}
