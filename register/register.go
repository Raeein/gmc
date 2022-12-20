package register

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Raeein/gmc"
	"github.com/Raeein/gmc/firestore"
	"github.com/Raeein/gmc/webadvisor"
	"io"
	"log"
)

type RegisterRequest struct {
	Section gmc.Section `json:"section"`
	User    gmc.User    `json:"user"`
}

func (r RegisterRequest) Valid() error {
	err := r.Section.Valid()
	if err != nil {
		return err
	}
	return r.User.Valid()
}

type RegisterService struct {
	ctx               context.Context
	webadvisorService webadvisor.WebAdvisor
	firestoreService  firestore.Service
}

func New(ctx context.Context, w webadvisor.WebAdvisor, f firestore.Service) RegisterService {
	return RegisterService{
		ctx:               ctx,
		webadvisorService: w,
		firestoreService:  f,
	}
}

func (r RegisterService) ValidateRegister(body io.ReadCloser) error {
	data, err := decode(body)
	if err != nil {
		return err
	}
	fmt.Println(data)
	if err := data.Valid(); err != nil {
		log.Printf("register request invalid: %s", err)
		return fmt.Errorf("invalid request")
	}
	if err := r.Register(r.ctx, data.Section, data.User); err != nil {
		log.Printf("registration failed: %s", err)
		return fmt.Errorf("registration failed, please ensure the course you are registering for exists")
	}
	return nil
}

func decode(body io.ReadCloser) (RegisterRequest, error) {
	data := RegisterRequest{}
	err := json.NewDecoder(body).Decode(&data)
	if err != nil {
		log.Println(err)
		return RegisterRequest{}, err
	}
	fmt.Println(data)
	return data, nil
}
func (r RegisterService) Register(ctx context.Context, section gmc.Section, user gmc.User) error {

	err := r.webadvisorService.Exists(ctx, section)
	if err != nil {
		return fmt.Errorf("failed to check if section exists: %w", err)
	}

	if err := r.firestoreService.AddWatcher(ctx, section, user); err != nil {
		return fmt.Errorf("failed to persist %s to %s: %w", user, section, err)
	}

	return nil
}
