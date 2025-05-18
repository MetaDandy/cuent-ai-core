package user

import (
	"errors"
	"regexp"
	"strings"
	"time"

	"github.com/MetaDandy/cuent-ai-core/helper"
	"github.com/MetaDandy/cuent-ai-core/src/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Service struct {
	repo *Repository
}

func NewService(r *Repository) *Service {
	return &Service{repo: r}
}

var (
	ErrInvalidEmail = errors.New("email no tiene un formato válido")
	ErrWeakPassword = errors.New("la contraseña debe tener al menos 8 caracteres")
	ErrEmailTaken   = errors.New("ya existe un usuario con ese email")
)

/* -------- regex simple RFC 5322 -------- */
var emailRx = regexp.MustCompile(`(?i)^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)

func (s *Service) FindAll(opts *helper.FindAllOptions) (*helper.PaginatedResponse[UserResponse], error) {
	users, total, err := s.repo.FindAll(opts)
	if err != nil {
		return nil, err
	}
	dtos := UsersToListDTO(users)
	pages := uint((total + int64(opts.Limit) - 1) / int64(opts.Limit))

	return &helper.PaginatedResponse[UserResponse]{
		Data:   dtos,
		Total:  total,
		Limit:  opts.Limit,
		Offset: opts.Offset,
		Pages:  pages,
	}, nil
}

func (s *Service) FindById(id string) (*UserResponse, error) {
	user, err := s.repo.FindById(id)
	if err != nil {
		return nil, err
	}

	dto := UserToDTO(user)
	return &dto, nil
}

func (s *Service) SignUp(in *Singup) (*UserResponse, string, error) {
	name := strings.TrimSpace(in.Name)
	email := strings.TrimSpace(strings.ToLower(in.Email))

	if !emailRx.MatchString(email) {
		return nil, "", ErrInvalidEmail
	}
	if len(in.Password) < 8 {
		return nil, "", ErrWeakPassword
	}

	if _, err := s.repo.FindByEmail(email); err == nil {
		return nil, "", ErrEmailTaken
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, "", err
	}

	hash, err := helper.HashPassword(in.Password)
	if err != nil {
		return nil, "", err
	}

	user := &model.User{
		ID:       uuid.New(),
		Name:     name,
		Email:    email,
		Password: string(hash),
	}

	free, err := s.repo.FindSubscriptionByName("Free")
	if err != nil {
		return nil, "", err
	}

	sub := &model.UserSubscribed{
		ID:              uuid.New(),
		SubscriptionID:  free.ID,
		UserID:          user.ID,
		StartDate:       time.Now(),
		EndDate:         time.Now().AddDate(0, 1, 0),
		TokensRemaining: free.Cuentokens,
	}

	if err := s.repo.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(user).Error; err != nil {
			return err
		}

		return tx.Create(sub).Error
	}); err != nil {
		return nil, "", err
	}

	dto := UserToDTO(user)
	token, err := helper.GenerateJwt(user.ID.String(), user.Email)
	if err != nil {
		return nil, "", err
	}

	return &dto, token, nil
}

func (s *Service) Signin(in *Signin) (*UserResponse, string, error) {
	user, err := s.repo.FindByEmail(in.Email)
	if err != nil {
		return nil, "", err
	}

	if !helper.CheckPasswordHash(in.Password, user.Password) {
		return nil, "", errors.New("la contraseña no coincide")
	}

	token, err := helper.GenerateJwt(user.ID.String(), user.Email)
	if err != nil {
		return nil, "", err
	}
	dto := UserToDTO(user)

	return &dto, token, nil
}

func (s *Service) ChangePassword(id string, in *ChangePassoword) (*UserResponse, error) {
	user, err := s.repo.FindById(id)
	if err != nil {
		return nil, err
	}

	if !helper.CheckPasswordHash(in.Old_Password, user.Password) {
		return nil, errors.New("la contraseña no coincide")
	}

	if in.New_Password != in.Confirm_Password {
		return nil, errors.New("la nueva contraseña no coincide con la confirmación")
	}

	hash, err := helper.HashPassword(in.New_Password)
	if err != nil {
		return nil, err
	}

	user.Password = string(hash)

	if err := s.repo.Update(user); err != nil {
		return nil, err
	}

	dto := UserToDTO(user)

	return &dto, nil
}

func (s *Service) AddSubscription(userID, subsID string) (*UserSubscriptionResponse, error) {
	user, err := s.repo.FindById(userID)
	if err != nil {
		return nil, err
	}

	sub, err := s.repo.FindSubscriptionById(subsID)
	if err != nil {
		return nil, err
	}

	subUser := &model.UserSubscribed{
		ID:              uuid.New(),
		UserID:          user.ID,
		SubscriptionID:  sub.ID,
		TokensRemaining: sub.Cuentokens,
		StartDate:       time.Now(),
		EndDate:         time.Now().AddDate(0, 1, 0),
	}

	if err := s.repo.AddSubscription(subUser); err != nil {
		return nil, err
	}

	dto := UserSubscriptionToDto(subUser)

	return &dto, nil
}

func (s *Service) GetActiveSubscription(id string) (*UserSubscriptionResponse, error) {
	sub, err := s.repo.GetActiveSubscription(id)
	if err != nil {
		return nil, err
	}

	dto := UserSubscriptionToDto(sub)
	return &dto, nil
}
