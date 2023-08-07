package services

import (
	"errors"
	"time"

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
	content.ContentId = uuid.New().String()
	content.PublicationDate = time.Now()
	return svc.repo.CreateContent(content)
}

func (svc *ContentManagementService) ReadContent(id string) (*domain.Content, error) {
	return svc.repo.ReadContent(id)
}

func (svc *ContentManagementService) ReadCreatorContents(creator_id string) ([]domain.Content, error) {
	contents, err := svc.repo.ReadContents()
	if err != nil {
		return nil, err
	}

	results := []domain.Content{}
	for _, content := range contents {
		if content.User.Id == creator_id {
			results = append(results, content)
		}
	}

	return results, nil

}
func (svc *ContentManagementService) ReadContents() ([]domain.Content, error) {
	return svc.repo.ReadContents()
}

func (svc *ContentManagementService) UpdateContent(content *domain.Content) (*domain.Content, error) {
	return svc.repo.UpdateContent(content)
}

func (svc *ContentManagementService) DeleteContent(id string) (string, error) {
	_, err := svc.ReadContent(id)
	if err != nil {
		return " ", errors.New("Error, item not found!")
	}
	return svc.repo.DeleteContent(id)
}

func (svc *ContentManagementService) DeleteAllContent() (string, error) {
	return svc.repo.DeleteAllContent()
}
