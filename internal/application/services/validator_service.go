package services

import (
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"analizador-backend/internal/domain/entities"
)

type ValidatorService struct {
	validEmails   []string
	chiapasLadas  []string
}

// NewValidatorService crea una nueva instancia del servicio de validación
func NewValidatorService() *ValidatorService {
	return &ValidatorService{
		validEmails: []string{
			"gmail.com", "yahoo.com", "hotmail.com", "outlook.com",
			"live.com", "icloud.com", "protonmail.com",
		},
		chiapasLadas: []string{
			"961", "962", "963", "964", "965", "966", "967", "968", "994",
		},
	}
}

// ValidateContact valida todos los campos de un contacto
func (v *ValidatorService) ValidateContact(contact *entities.Contact) []entities.ValidationError {
	var errors []entities.ValidationError

	// Validar clave cliente
	if clientKeyErrors := v.validateClientKey(contact.ClientKey); len(clientKeyErrors) > 0 {
		errors = append(errors, clientKeyErrors...)
	}

	// Validar nombre
	if nameErrors := v.validateName(contact.Name); len(nameErrors) > 0 {
		errors = append(errors, nameErrors...)
	}

	// Validar email
	if emailErrors := v.validateEmail(contact.Email); len(emailErrors) > 0 {
		errors = append(errors, emailErrors...)
	}

	// Validar teléfono
	if phoneErrors := v.validatePhone(contact.Phone); len(phoneErrors) > 0 {
		errors = append(errors, phoneErrors...)
	}

	return errors
}

// validateClientKey valida que la clave cliente sea solo números
func (v *ValidatorService) validateClientKey(clientKey string) []entities.ValidationError {
	var errors []entities.ValidationError

	if clientKey == "" {
		errors = append(errors, entities.ValidationError{
			Field:   "client_key",
			Value:   clientKey,
			Message: "La clave cliente no puede estar vacía",
			Type:    "REQUIRED",
		})
		return errors
	}

	if _, err := strconv.Atoi(clientKey); err != nil {
		errors = append(errors, entities.ValidationError{
			Field:   "client_key",
			Value:   clientKey,
			Message: "La clave cliente debe contener solo números",
			Type:    "INVALID_FORMAT",
		})
	}

	return errors
}

// validateName valida que el nombre contenga solo letras, espacios y acentos
func (v *ValidatorService) validateName(name string) []entities.ValidationError {
	var errors []entities.ValidationError

	if name == "" {
		errors = append(errors, entities.ValidationError{
			Field:   "name",
			Value:   name,
			Message: "El nombre no puede estar vacío",
			Type:    "REQUIRED",
		})
		return errors
	}

	// Permitir letras (incluye acentos), espacios y apostrofes para nombres como "O'Connor"
	for _, char := range name {
		if !unicode.IsLetter(char) && !unicode.IsSpace(char) && char != '\'' && char != '.' {
			errors = append(errors, entities.ValidationError{
				Field:   "name",
				Value:   name,
				Message: "El nombre debe contener solo letras, espacios y apostrofes",
				Type:    "INVALID_CHARACTER",
			})
			break
		}
	}

	return errors
}

// validateEmail valida el formato y dominio del email
func (v *ValidatorService) validateEmail(email string) []entities.ValidationError {
	var errors []entities.ValidationError

	if email == "" {
		errors = append(errors, entities.ValidationError{
			Field:   "email",
			Value:   email,
			Message: "El email no puede estar vacío",
			Type:    "REQUIRED",
		})
		return errors
	}

	// Validar formato básico de email
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	if matched, _ := regexp.MatchString(emailRegex, email); !matched {
		errors = append(errors, entities.ValidationError{
			Field:   "email",
			Value:   email,
			Message: "El formato del email no es válido",
			Type:    "INVALID_FORMAT",
		})
		return errors
	}

	// Validar dominio conocido
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		errors = append(errors, entities.ValidationError{
			Field:   "email",
			Value:   email,
			Message: "El formato del email no es válido",
			Type:    "INVALID_FORMAT",
		})
		return errors
	}

	domain := strings.ToLower(parts[1])
	isValidDomain := false
	for _, validDomain := range v.validEmails {
		if domain == validDomain {
			isValidDomain = true
			break
		}
	}

	if !isValidDomain {
		errors = append(errors, entities.ValidationError{
			Field:   "email",
			Value:   email,
			Message: "El dominio del email no es reconocido (use gmail.com, yahoo.com, hotmail.com, etc.)",
			Type:    "INVALID_DOMAIN",
		})
	}

	return errors
}

// validatePhone valida el formato del teléfono con lada de Chiapas
func (v *ValidatorService) validatePhone(phone string) []entities.ValidationError {
	var errors []entities.ValidationError

	if phone == "" {
		errors = append(errors, entities.ValidationError{
			Field:   "phone",
			Value:   phone,
			Message: "El teléfono no puede estar vacío",
			Type:    "REQUIRED",
		})
		return errors
	}

	// Remover espacios y caracteres especiales para la validación
	cleanPhone := regexp.MustCompile(`[^\d]`).ReplaceAllString(phone, "")

	// Validar que solo contenga números después de limpieza
	if regexp.MustCompile(`[a-zA-Z]`).MatchString(phone) {
		errors = append(errors, entities.ValidationError{
			Field:   "phone",
			Value:   phone,
			Message: "El teléfono no debe contener letras",
			Type:    "INVALID_CHARACTER",
		})
		return errors
	}

	// Validar longitud (10 dígitos)
	if len(cleanPhone) != 10 {
		errors = append(errors, entities.ValidationError{
			Field:   "phone",
			Value:   phone,
			Message: "El teléfono debe tener exactamente 10 dígitos",
			Type:    "INVALID_LENGTH",
		})
		return errors
	}

	// Validar lada de Chiapas
	lada := cleanPhone[:3]
	isValidLada := false
	for _, validLada := range v.chiapasLadas {
		if lada == validLada {
			isValidLada = true
			break
		}
	}

	if !isValidLada {
		errors = append(errors, entities.ValidationError{
			Field:   "phone",
			Value:   phone,
			Message: "La lada debe ser de Chiapas (961, 962, 963, 964, 965, 966, 967, 968, 994)",
			Type:    "INVALID_AREA_CODE",
		})
	}

	return errors
}