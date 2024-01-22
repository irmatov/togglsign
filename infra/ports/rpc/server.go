package rpc

import (
	context "context"

	"github.com/irmatov/togglsign/app"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Application interface {
	Sign(ctx context.Context, token string, responses []app.Response) (string, error)
	Verify(ctx context.Context, email, sig string) (app.VerifyResponse, error)
}

type Server struct {
	app Application
	UnimplementedSignerServer
}

func New(app Application) *Server {
	return &Server{app: app}
}

func (s *Server) Sign(ctx context.Context, sr *SignRequest) (*SignResponse, error) {
	rs := make([]app.Response, 0, len(sr.Responses))
	for _, r := range sr.Responses {
		rs = append(rs, app.Response{Question: r.GetQuestion(), Answer: r.GetAnswer()})
	}
	sig, err := s.app.Sign(ctx, sr.GetJwtToken(), rs)
	if err != nil {
		return nil, err
	}
	return &SignResponse{Signature: sig}, nil
}

func (s *Server) Verify(ctx context.Context, vr *VerifyRequest) (*VerifyResponse, error) {
	ar, err := s.app.Verify(ctx, vr.Email, vr.Signature)
	if err != nil {
		return nil, err
	}
	rs := make([]*Response, 0, len(ar.Responses))
	for _, v := range ar.Responses {
		rs = append(rs, &Response{Question: v.Question, Answer: v.Answer})
	}
	return &VerifyResponse{
		Ok:        ar.Ok,
		SignedAt:  timestamppb.New(ar.SignedAt),
		Responses: rs,
	}, nil
}
