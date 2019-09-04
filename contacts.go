package activecampaign

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
)

type Contacts struct {
	Contacts []Contact       `json:"contacts"`
	Meta     FieldValuesMeta `json:"meta"`
}
type Contact struct {
	CreateDate string `json:"cdate"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`
	FirstName  string `json:"firstName,omitempty"`
	LastName   string `json:"lastName,omitempty"`
	ID         string `json:"id"`
	UpdateDate string `json:"udate"`
}

func (a *ActiveCampaign) Contacts(ctx context.Context, pof *POF) (*Contacts, error) {
	res, err := a.send(ctx, http.MethodGet, "contacts", pof, nil)
	if err != nil {
		return nil, err
	}

	var contacts Contacts
	err = json.NewDecoder(res.Body).Decode(&contacts)
	if err != nil {
		return nil, err
	}

	return &contacts, err
}

type ContactCreate struct {
	Email     string `json:"email"`
	FirstName string `json:"firstName,omitempty"`
	LastName  string `json:"lastName,omitempty"`
	Phone     string `json:"phone,omitempty"`
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
