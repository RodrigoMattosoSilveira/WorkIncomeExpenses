package clients

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"
)

type PeopleClientConfig struct {
	BaseURL string
	Timeout time.Duration
}

var _ PeopleAPI = (*PeopleClient)(nil)

type PeopleClient struct {
	base string
	hc   *http.Client
}

func NewPeopleClient(cfg PeopleClientConfig) *PeopleClient {
	return &PeopleClient{
		base: cfg.BaseURL,
		hc:   &http.Client{Timeout: cfg.Timeout},
	}
}

type PersonDTO struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type CreatePersonRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UpdatePersonRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// for validation errors from people-svc
type ValidationErrorResponse struct {
	Errors map[string]string `json:"errors"`
}

func (pc *PeopleClient) ListPeople(ctx context.Context) ([]PersonDTO, error) {
	req, _ := http.NewRequestWithContext(ctx, "GET", pc.base+"/api/people", nil)
	resp, err := pc.hc.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, errors.New("non-200")
	}

	var out []PersonDTO
	return out, json.NewDecoder(resp.Body).Decode(&out)
}

func (pc *PeopleClient) GetPerson(ctx context.Context, id int64) (PersonDTO, error) {
	req, _ := http.NewRequestWithContext(ctx, "GET", pc.base+"/api/people/"+itoa(id), nil)
	resp, err := pc.hc.Do(req)
	if err != nil {
		return PersonDTO{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return PersonDTO{}, errors.New("non-200")
	}

	var out PersonDTO
	return out, json.NewDecoder(resp.Body).Decode(&out)
}

func (pc *PeopleClient) CreatePerson(ctx context.Context, in CreatePersonRequest) (PersonDTO, map[string]string, error) {
	b, _ := json.Marshal(in)
	req, _ := http.NewRequestWithContext(ctx, "POST", pc.base+"/api/people", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")

	resp, err := pc.hc.Do(req)
	if err != nil {
		return PersonDTO{}, nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 422 {
		var v ValidationErrorResponse
		_ = json.NewDecoder(resp.Body).Decode(&v)
		return PersonDTO{}, v.Errors, nil
	}
	if resp.StatusCode != 201 {
		return PersonDTO{}, nil, errors.New("non-201")
	}

	var out PersonDTO
	return out, nil, json.NewDecoder(resp.Body).Decode(&out)
}

func (pc *PeopleClient) UpdatePerson(ctx context.Context, id int64, in UpdatePersonRequest) (PersonDTO, map[string]string, error) {
	b, _ := json.Marshal(in)
	req, _ := http.NewRequestWithContext(ctx, "PATCH", pc.base+"/api/people/"+itoa(id), bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")

	resp, err := pc.hc.Do(req)
	if err != nil {
		return PersonDTO{}, nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 422 {
		var v ValidationErrorResponse
		_ = json.NewDecoder(resp.Body).Decode(&v)
		return PersonDTO{}, v.Errors, nil
	}
	if resp.StatusCode != 200 {
		return PersonDTO{}, nil, errors.New("non-200")
	}

	var out PersonDTO
	return out, nil, json.NewDecoder(resp.Body).Decode(&out)
}

func (pc *PeopleClient) DeletePerson(ctx context.Context, id int64) error {
	req, _ := http.NewRequestWithContext(ctx, "DELETE", pc.base+"/api/people/"+itoa(id), nil)
	resp, err := pc.hc.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 204 {
		return errors.New("non-204")
	}
	return nil
}

func itoa(v int64) string {
	// tiny helper to avoid fmt overhead in hot paths
	// ok to replace with strconv.FormatInt(v, 10)
	return strconv.FormatInt(v, 10)
}
