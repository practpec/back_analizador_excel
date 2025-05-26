package entities

import "time"

// Contact representa una entidad de contacto del dominio
type Contact struct {
	ID           int       `json:"id"`
	ClientKey    string    `json:"client_key"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	Phone        string    `json:"phone"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// ValidationError representa un error de validación
type ValidationError struct {
	Field   string `json:"field"`
	Value   string `json:"value"`
	Message string `json:"message"`
	Type    string `json:"type"`
}

// ContactWithValidation representa un contacto con sus errores de validación
type ContactWithValidation struct {
	Contact Contact           `json:"contact"`
	Errors  []ValidationError `json:"errors"`
	IsValid bool              `json:"is_valid"`
}