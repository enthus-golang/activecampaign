package activecampaign

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type Contacts struct {
	Contacts []Contact       `json:"contacts"`
	Meta     FieldValuesMeta `json:"meta"`
}

type Contact struct {
	CreateDate string       `json:"cdate"`
	Email      string       `json:"email"`
	Phone      string       `json:"phone"`
	FirstName  string       `json:"firstName,omitempty"`
	LastName   string       `json:"lastName,omitempty"`
	ID         string       `json:"id"`
	UpdateDate string       `json:"udate"`
	Links      ContactLinks `json:"links"`
}

type ContactLinks struct {
	BounceLogs            string `json:"bounceLogs"`
	ContactAutomations    string `json:"contactAutomations"`
	ContactData           string `json:"contactData"`
	ContactGoals          string `json:"contactGoals"`
	ContactLists          string `json:"contactLists"`
	ContactLogs           string `json:"contactLogs"`
	ContactTags           string `json:"contactTags"`
	ContactDeals          string `json:"contactDeals"`
	Deals                 string `json:"deals"`
	FieldValues           string `json:"fieldValues"`
	GeoIps                string `json:"geoIps"`
	Notes                 string `json:"notes"`
	Organization          string `json:"organization"`
	PlusAppend            string `json:"plusAppend"`
	TrackingLogs          string `json:"trackingLogs"`
	ScoreValues           string `json:"scoreValues"`
	AccountContacts       string `json:"accountContacts"`
	AutomationEntryCounts string `json:"automationEntryCounts"`
}

func (a *ActiveCampaign) Contacts(ctx context.Context, pof *POF) (*Contacts, error) {
	res, err := a.send(ctx, http.MethodGet, "contacts", pof, nil)
	if err != nil {
		return nil, &Error{Op: "contacts", Err: err}
	}
	defer func() {
		if closeErr := res.Body.Close(); closeErr != nil && err == nil {
			err = &Error{Op: "contacts", Err: closeErr}
		}
	}()

	var contacts Contacts
	err = json.NewDecoder(res.Body).Decode(&contacts)
	if err != nil {
		return nil, &Error{Op: "contacts", Err: err}
	}

	return &contacts, nil
}

func (a *ActiveCampaign) ListContacts(ctx context.Context, listID string) (*Contacts, error) {
	res, err := a.send(ctx, http.MethodGet, "contacts?listid="+listID+"&status=1", nil, nil)
	if err != nil {
		return nil, &Error{Op: "list contacts", Err: err}
	}
	defer func() {
		if closeErr := res.Body.Close(); closeErr != nil && err == nil {
			err = &Error{Op: "list contacts", Err: closeErr}
		}
	}()

	var contacts Contacts
	err = json.NewDecoder(res.Body).Decode(&contacts)
	if err != nil {
		return nil, &Error{Op: "list contacts", Err: err}
	}

	return &contacts, nil
}

func (a *ActiveCampaign) ContactFieldValues(ctx context.Context, pof *POF, id string) (*FieldValues, error) {
	return a.fieldValues(ctx, pof, "contacts/"+id+"/fieldValues")
}

type ContactCreate struct {
	Email     string `json:"email"`
	FirstName string `json:"firstName,omitempty"`
	LastName  string `json:"lastName,omitempty"`
	Phone     string `json:"phone,omitempty"`
}

type ContactCreated struct {
	Email      string       `json:"email"`
	CreateDate string       `json:"cdate"`
	UpdateDate string       `json:"udate"`
	Links      ContactLinks `json:"links"`
	ID         string       `json:"id"`
}

func (a *ActiveCampaign) ContactCreate(ctx context.Context, contact ContactCreate) (*ContactCreated, error) {
	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(struct {
		Contact ContactCreate `json:"contact"`
	}{
		Contact: contact,
	})
	if err != nil {
		return nil, &Error{Op: "contact create", Err: err}
	}

	res, err := a.send(ctx, http.MethodPost, "contacts", nil, b)
	if err != nil {
		return nil, &Error{Op: "contact create", Err: err}
	}
	defer func() {
		if closeErr := res.Body.Close(); closeErr != nil && err == nil {
			err = &Error{Op: "contact create", Err: closeErr}
		}
	}()

	if res.StatusCode != http.StatusCreated {
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		return nil, &Error{Op: "contact create", Err: errors.New("wrong status code"), Body: body}
	}

	var contactCreated struct {
		Contact ContactCreated `json:"contact"`
	}
	err = json.NewDecoder(res.Body).Decode(&contactCreated)
	if err != nil {
		return nil, err
	}

	return &contactCreated.Contact, nil
}

func (a *ActiveCampaign) ContactDelete(ctx context.Context, id string) error {
	res, err := a.send(ctx, http.MethodDelete, "contacts/"+id, nil, nil)
	if err != nil {
		return &Error{Op: "contact delete", Err: err}
	}
	defer func() {
		if closeErr := res.Body.Close(); closeErr != nil && err == nil {
			err = &Error{Op: "contact delete", Err: closeErr}
		}
	}()

	if res.StatusCode == http.StatusOK {
		return nil
	}

	var message struct {
		Message string `json:"message"`
	}
	err = json.NewDecoder(res.Body).Decode(&message)
	if err != nil {
		return &Error{Op: "contact delete", Err: err}
	}

	return errors.New("contact delete: " + message.Message)
}

type ContactUpdate struct {
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

func (a *ActiveCampaign) ContactUpdate(ctx context.Context, id string, contact ContactUpdate) error {
	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(struct {
		Contact ContactUpdate `json:"contact"`
	}{
		Contact: contact,
	})
	if err != nil {
		return &Error{Op: "contact update", Err: err}
	}

	res, err := a.send(ctx, http.MethodPut, "contacts/"+id, nil, b)
	if err != nil {
		return &Error{Op: "contact update", Err: err}
	}
	defer func() {
		if closeErr := res.Body.Close(); closeErr != nil && err == nil {
			err = &Error{Op: "contact update", Err: closeErr}
		}
	}()

	if res.StatusCode == http.StatusOK {
		return nil
	}

	var message struct {
		Message string `json:"message"`
	}
	err = json.NewDecoder(res.Body).Decode(&message)
	if err != nil {
		return &Error{Op: "contact update", Err: err}
	}

	return errors.New("contact update: " + message.Message)
}

type ContactSync struct {
	Email       string                  `json:"email"`
	FirstName   string                  `json:"firstName,omitempty"`
	LastName    string                  `json:"lastName,omitempty"`
	Phone       string                  `json:"phone,omitempty"`
	FieldValues []ContactSyncFieldValue `json:"fieldValues,omitempty"`
}

type ContactSyncFieldValue struct {
	Field string `json:"field"`
	Value string `json:"value"`
}

// ContactSync creates the contact or updates the one with the same email. All
// field values are written in this one request, so an automation triggered by
// one of them already sees the others.
func (a *ActiveCampaign) ContactSync(ctx context.Context, contact ContactSync) (*ContactCreated, error) {
	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(struct {
		Contact ContactSync `json:"contact"`
	}{
		Contact: contact,
	})
	if err != nil {
		return nil, &Error{Op: "contact sync", Err: err}
	}

	res, err := a.send(ctx, http.MethodPost, "contact/sync", nil, b)
	if err != nil {
		return nil, &Error{Op: "contact sync", Err: err}
	}
	defer func() {
		if closeErr := res.Body.Close(); closeErr != nil && err == nil {
			err = &Error{Op: "contact sync", Err: closeErr}
		}
	}()

	if res.StatusCode != http.StatusCreated && res.StatusCode != http.StatusOK {
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, &Error{Op: "contact sync", Err: err}
		}

		return nil, &Error{Op: "contact sync", Err: errors.New("wrong status code"), Body: body}
	}

	var contactSynced struct {
		Contact ContactCreated `json:"contact"`
	}
	err = json.NewDecoder(res.Body).Decode(&contactSynced)
	if err != nil {
		return nil, &Error{Op: "contact sync", Err: err}
	}

	return &contactSynced.Contact, nil
}
