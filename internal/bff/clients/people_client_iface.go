package clients

import "context"

type PeopleAPI interface {
	ListPeople(ctx context.Context) ([]PersonDTO, error)
	GetPerson(ctx context.Context, id int64) (PersonDTO, error)
	CreatePerson(ctx context.Context, in CreatePersonRequest) (PersonDTO, map[string]string, error)
	UpdatePerson(ctx context.Context, id int64, in UpdatePersonRequest) (PersonDTO, map[string]string, error)
	DeletePerson(ctx context.Context, id int64) error
}
