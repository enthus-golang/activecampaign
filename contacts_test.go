package activecampaign

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestContactSync(t *testing.T) {
	var gotMethod, gotPath, gotToken string
	var gotBody map[string]any
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod, gotPath, gotToken = r.Method, r.URL.Path, r.Header.Get("Api-Token")
		if err := json.NewDecoder(r.Body).Decode(&gotBody); err != nil {
			t.Errorf("decode request body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"contact":{"email":"jane@example.com","id":"42"}}`))
	}))
	t.Cleanup(srv.Close)

	ac, err := New(srv.URL, "token")
	if err != nil {
		t.Fatal(err)
	}

	contact, err := ac.ContactSync(t.Context(), ContactSync{
		Email: "jane@example.com",
		FieldValues: []ContactSyncFieldValue{
			{Field: "218", Value: "123456"},
			{Field: "220", Value: "self"},
		},
	})
	if err != nil {
		t.Fatalf("ContactSync: %v", err)
	}

	if gotMethod != http.MethodPost || gotPath != "/api/3/contact/sync" || gotToken != "token" {
		t.Errorf("request = %s %s token %q, want POST /api/3/contact/sync token \"token\"", gotMethod, gotPath, gotToken)
	}
	if contact.ID != "42" {
		t.Errorf("contact ID = %q, want 42", contact.ID)
	}

	body, _ := json.Marshal(gotBody)
	want := `{"contact":{"email":"jane@example.com","fieldValues":[{"field":"218","value":"123456"},{"field":"220","value":"self"}]}}`
	if string(body) != want {
		t.Errorf("request body = %s, want %s", body, want)
	}
}

func TestContactSyncRejected(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		_, _ = w.Write([]byte(`{"errors":[{"title":"Email address is invalid"}]}`))
	}))
	t.Cleanup(srv.Close)

	ac, err := New(srv.URL, "token")
	if err != nil {
		t.Fatal(err)
	}

	_, err = ac.ContactSync(t.Context(), ContactSync{Email: "invalid"})

	var acErr *Error
	if !errors.As(err, &acErr) {
		t.Fatalf("error = %v, want *Error", err)
	}
	if string(acErr.Body) != `{"errors":[{"title":"Email address is invalid"}]}` {
		t.Errorf("error body = %s", acErr.Body)
	}
}
