package repositories

import (
	"errors"
	"strings"
	"sync"
	"time"

	"analizador-backend/internal/domain/entities"
	"analizador-backend/internal/domain/repositories"
)

type InMemoryContactRepository struct {
	contacts map[int]*entities.Contact
	nextID   int
	mutex    sync.RWMutex
}

// NewInMemoryContactRepository crea una nueva instancia del repositorio en memoria
func NewInMemoryContactRepository() repositories.ContactRepository {
	return &InMemoryContactRepository{
		contacts: make(map[int]*entities.Contact),
		nextID:   1,
	}
}

// Save guarda un contacto en memoria
func (r *InMemoryContactRepository) Save(contact *entities.Contact) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if contact.ID == 0 {
		contact.ID = r.nextID
		r.nextID++
		contact.CreatedAt = time.Now()
	}
	contact.UpdatedAt = time.Now()
	
	r.contacts[contact.ID] = contact
	return nil
}

// FindAll obtiene todos los contactos
func (r *InMemoryContactRepository) FindAll() ([]*entities.Contact, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	contacts := make([]*entities.Contact, 0, len(r.contacts))
	for _, contact := range r.contacts {
		contacts = append(contacts, contact)
	}
	return contacts, nil
}

// FindByID busca un contacto por ID
func (r *InMemoryContactRepository) FindByID(id int) (*entities.Contact, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	contact, exists := r.contacts[id]
	if !exists {
		return nil, errors.New("contacto no encontrado")
	}
	return contact, nil
}

// Update actualiza un contacto existente
func (r *InMemoryContactRepository) Update(contact *entities.Contact) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.contacts[contact.ID]; !exists {
		return errors.New("contacto no encontrado")
	}

	contact.UpdatedAt = time.Now()
	r.contacts[contact.ID] = contact
	return nil
}

// Delete elimina un contacto por ID
func (r *InMemoryContactRepository) Delete(id int) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.contacts[id]; !exists {
		return errors.New("contacto no encontrado")
	}

	delete(r.contacts, id)
	return nil
}

// Search busca contactos por campo y valor
func (r *InMemoryContactRepository) Search(field, value string) ([]*entities.Contact, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var results []*entities.Contact
	searchValue := strings.ToLower(value)

	for _, contact := range r.contacts {
		var fieldValue string
		switch field {
		case "client_key":
			fieldValue = strings.ToLower(contact.ClientKey)
		case "name":
			fieldValue = strings.ToLower(contact.Name)
		case "email":
			fieldValue = strings.ToLower(contact.Email)
		case "phone":
			fieldValue = strings.ToLower(contact.Phone)
		default:
			continue
		}

		if strings.Contains(fieldValue, searchValue) {
			results = append(results, contact)
		}
	}

	return results, nil
}

// SaveBatch guarda múltiples contactos
func (r *InMemoryContactRepository) SaveBatch(contacts []*entities.Contact) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	for _, contact := range contacts {
		if contact.ID == 0 {
			contact.ID = r.nextID
			r.nextID++
			contact.CreatedAt = time.Now()
		}
		contact.UpdatedAt = time.Now()
		r.contacts[contact.ID] = contact
	}

	return nil
}