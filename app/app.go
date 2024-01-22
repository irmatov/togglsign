package app

import (
	"context"
	"time"
)

type Storage interface {
	SaveResponseSet(ctx context.Context, rs ResponseSet) error
	LoadResponseSet(ctx context.Context, email, sig string) (*ResponseSet, error)
}

type App struct {
	storage Storage
	jwtKey  []byte
	signKey string
}

type Response struct {
	Question string
	Answer   string
}

type ResponseSet struct {
	Email     string
	SignedAt  time.Time
	Sig       string
	Responses []Response
}

func New(storage Storage, jwtKey []byte, signKey string) *App {
	return &App{storage, jwtKey, signKey}
}

func (a *App) Sign(ctx context.Context, token string, responses []Response) (string, error) {
	email, err := verifyToken(token, a.jwtKey)
	if err != nil {
		return "", err
	}
	rs := ResponseSet{
		Email:     email,
		SignedAt:  time.Now(),
		Sig:       "",
		Responses: responses,
	}
	sig, err := createSignature(rs, a.signKey)
	if err != nil {
		return "", err
	}
	rs.Sig = sig
	return sig, a.storage.SaveResponseSet(ctx, rs)
}

type VerifyResponse struct {
	Ok        bool
	SignedAt  time.Time
	Responses []Response
}

func (a *App) Verify(ctx context.Context, email, sig string) (VerifyResponse, error) {
	rs, err := a.storage.LoadResponseSet(ctx, email, sig)
	if err != nil || rs == nil {
		return VerifyResponse{}, err
	}
	return VerifyResponse{
		Ok:        true,
		SignedAt:  rs.SignedAt,
		Responses: rs.Responses,
	}, nil
}
