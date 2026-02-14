package people

import (
	"context"
	"net/mail"
	"strings"
)

type ValidationErrors map[string]string

type Service struct{ repo Repo }

func NewService(repo Repo) *Service { return &Service{repo: repo} }

func (s *Service) List(ctx context.Context) ([]Person, error) {
	return s.repo.List(ctx)
}

func (s *Service) Get(ctx context.Context, id uint) (Person, error) {
	return s.repo.Get(ctx, id)
}

func (s *Service) Create(ctx context.Context, name, email string) (Person, ValidationErrors, error) {
	name = strings.TrimSpace(name)
	email = strings.TrimSpace(email)

	verrs := validate(name, email)
	if len(verrs) > 0 {
		return Person{}, verrs, nil
	}

	exists, err := s.repo.EmailExists(ctx, email, 0)
	if err != nil {
		return Person{}, nil, err
	}
	if exists {
		return Person{}, ValidationErrors{"email": "email is already taken"}, nil
	}

	p := Person{Name: name, Email: email}
	if err := s.repo.Create(ctx, &p); err != nil {
		return Person{}, nil, err
	}
	return p, nil, nil
}

func (s *Service) Update(ctx context.Context, id uint, name, email string) (Person, ValidationErrors, error) {
	name = strings.TrimSpace(name)
	email = strings.TrimSpace(email)

	verrs := validate(name, email)
	if len(verrs) > 0 {
		return Person{}, verrs, nil
	}

	p, err := s.repo.Get(ctx, id)
	if err != nil {
		return Person{}, nil, err
	}

	exists, err := s.repo.EmailExists(ctx, email, id)
	if err != nil {
		return Person{}, nil, err
	}
	if exists {
		return Person{}, ValidationErrors{"email": "email is already taken"}, nil
	}

	p.Name = name
	p.Email = email
	if err := s.repo.Update(ctx, &p); err != nil {
		return Person{}, nil, err
	}
	return p, nil, nil
}

func (s *Service) Delete(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}

func validate(name, email string) ValidationErrors {
	verrs := ValidationErrors{}
	if name == "" {
		verrs["name"] = "name is required"
	}
	if email == "" {
		verrs["email"] = "email is required"
	} else if _, err := mail.ParseAddress(email); err != nil {
		verrs["email"] = "email must be a valid address"
	}
	return verrs
}
