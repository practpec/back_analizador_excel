package repositories

import "analizador-backend/internal/domain/entities"

// ContactRepository define la interfaz para el repositorio de contactos
type ContactRepository interface {
	Save(contact *entities.Contact) error
	FindAll() ([]*entities.Contact, error)
	FindByID(id int) (*entities.Contact, error)
	Update(contact *entities.Contact) error
	Delete(id int) error
	Search(field, value string) ([]*entities.Contact, error)
	SaveBatch(contacts []*entities.Contact) error
}