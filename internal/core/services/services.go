package services

import (
	"errors"

	"github.com/AntonyIS/notlify-content-svc/internal/core/domain"
	"github.com/AntonyIS/notlify-content-svc/internal/core/ports"
	"github.com/google/uuid"
)

type ContentManagementService struct {
	repo ports.ContentRepository
}

func NewContentManagementService(repo ports.ContentRepository) *ContentManagementService {
	svc := ContentManagementService{
		repo: repo,
	}
	return &svc
}

func (svc *ContentManagementService) CreateContent(content *domain.Content) (*domain.Content, error) {
	// Assign new content with a unique id
	content.ContentId = uuid.New().String()
	// hash content password
	return svc.repo.CreateContent(content)
}

func (svc *ContentManagementService) ReadContent(id string) (*domain.Content, error) {
	return svc.repo.ReadContent(id)
}

func (svc *ContentManagementService) ReadContentWithEmail(email string) (*domain.Content, error) {
	return svc.repo.ReadContent(email)
}

func (svc *ContentManagementService) ReadContents() ([]domain.Content, error) {
	return svc.repo.ReadContents()
}

func (svc *ContentManagementService) UpdateContent(content *domain.Content) (*domain.Content, error) {
	return svc.repo.UpdateContent(content)
}

func (svc *ContentManagementService) DeleteUser(id string) (string, error) {
	// Check if user exists
	_, err := svc.ReadContent(id)
	if err != nil {
		return " ", errors.New("Error, item not found!")
	}

	return svc.repo.DeleteContent(id)
}
