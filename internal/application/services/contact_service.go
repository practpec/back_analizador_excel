package services

import (
	"analizador-backend/internal/domain/entities"
	"analizador-backend/internal/domain/repositories"
)

type ContactService struct {
	contactRepo      repositories.ContactRepository
	validatorService *ValidatorService
}

// NewContactService crea una nueva instancia del servicio de contactos
func NewContactService(contactRepo repositories.ContactRepository, validatorService *ValidatorService) *ContactService {
	return &ContactService{
		contactRepo:      contactRepo,
		validatorService: validatorService,
	}
}

// GetAllContacts obtiene todos los contactos
func (s *ContactService) GetAllContacts() ([]*entities.Contact, error) {
	return s.contactRepo.FindAll()
}

// SearchContacts busca contactos por campo y valor
func (s *ContactService) SearchContacts(field, value string) ([]*entities.Contact, error) {
	return s.contactRepo.Search(field, value)
}

// UpdateContact actualiza un contacto existente
func (s *ContactService) UpdateContact(contact *entities.Contact) error {
	return s.contactRepo.Update(contact)
}

// ValidateAllContacts valida todos los contactos y retorna los resultados
func (s *ContactService) ValidateAllContacts() ([]*entities.ContactWithValidation, error) {
	contacts, err := s.contactRepo.FindAll()
	if err != nil {
		return nil, err
	}

	var results []*entities.ContactWithValidation
	for _, contact := range contacts {
		errors := s.validatorService.ValidateContact(contact)
		result := &entities.ContactWithValidation{
			Contact: *contact,
			Errors:  errors,
			IsValid: len(errors) == 0,
		}
		results = append(results, result)
	}

	return results, nil
}

// SaveContactsBatch guarda múltiples contactos
func (s *ContactService) SaveContactsBatch(contacts []*entities.Contact) error {
	return s.contactRepo.SaveBatch(contacts)
}