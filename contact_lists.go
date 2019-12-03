package activecampaign

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

type ContactAddedToList struct {
	Contacts    []Contact   `json:"contacts"`
	ContactList ContactList `json:"contactList"`
}

type ContactLists struct {
	ContactLists []ContactList `json:"contactLists"`
}

type ContactList struct {
	Contact               string          `json:"contact"`
	List                  string          `json:"list"`
	Form                  *string         `json:"form"`
	Seriesid              string          `json:"seriesid"`
	Sdate                 string          `json:"sdate"`
	Status                interface{}     `json:"status"`
	Responder             string          `json:"responder"`
	Sync                  string          `json:"sync"`
	Unsubreason           string          `json:"unsubreason"`
	Campaign              *string         `json:"campaign"`
	Message               *string         `json:"message"`
	First_name            string          `json:"first_name"`
	Last_name             string          `json:"last_name"`
	Ip4Sub                string          `json:"ip4Sub"`
	Sourceid              string          `json:"sourceid"`
	AutosyncLog           *string         `json:"autosyncLog"`
	Ip4_last              string          `json:"ip4_last"`
	Ip4Unsub              string          `json:"ip4Unsub"`
	UnsubscribeAutomation *string         `json:"unsubscribeAutomation"`
	Links                 ContactListLink `json:"links"`
	ID                    string          `json:"id"`
	Automation            *string         `json:"automation"`
}

type ContactListLink struct {
	Automation            string `json:"automation"`
	List                  string `json:"list"`
	Contact               string `json:"contact"`
	Form                  string `json:"form"`
	AutosyncLog           string `json:"autosyncLog"`
	Campaign              string `json:"campaign"`
	UnsubscribeAutomation string `json:"unsubscribeAutomation"`
	Message               string `json:"message"`
}

type AddContactToListRequest struct {
	List    int `json:"list"`
	Contact int `json:"contact"`
	Status  int `json:"status"`
}

func (a *ActiveCampaign) AddContactToList(ctx context.Context, contactID string, listID string) (*ContactAddedToList, error) {
	b := new(bytes.Buffer)

	cID, err := strconv.Atoi(contactID)
	if err != nil {
		return nil, &Error{Op: "converting contactID to int", Err: err}
	}
	lID, err := strconv.Atoi(listID)
	if err != nil {
		return nil, &Error{Op: "converting listID to int", Err: err}
	}

	err = json.NewEncoder(b).Encode(struct {
		ContactList AddContactToListRequest `json:"contactList"`
	}{
		ContactList: AddContactToListRequest{
			List:    lID,
			Contact: cID,
			Status:  1,
		},
	})
	if err != nil {
		return nil, &Error{Op: "add contact to list", Err: err}
	}

	res, err := a.send(ctx, http.MethodPost, "contactLists", nil, b)
	if err != nil {
		return nil, &Error{Op: "add contact to list", Err: err}
	}
	if res.StatusCode != http.StatusCreated {
		return nil, errors.New("add contact to list: " + res.Status)
	}

	var contactAddedToList ContactAddedToList
	err = json.NewDecoder(res.Body).Decode(&contactAddedToList)
	if err != nil {
		return nil, err
	}

	return &contactAddedToList, nil
}

func (a *ActiveCampaign) ContactLists(ctx context.Context, contactID string) (*ContactLists, error) {
	res, err := a.send(ctx, http.MethodGet, fmt.Sprintf("contacts/%s/contactLists", contactID), nil, nil)
	if err != nil {
		return nil, &Error{Op: "contactLists", Err: err}
	}

	var contactLists ContactLists
	err = json.NewDecoder(res.Body).Decode(&contactLists)
	if err != nil {
		return nil, &Error{Op: "contactLists", Err: err}
	}

	return &contactLists, nil
}
