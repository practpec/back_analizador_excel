package handlers

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
	"analizador-backend/internal/application/services"
	"analizador-backend/internal/domain/entities"
)

type ContactHandler struct {
	contactService *services.ContactService
}

// NewContactHandler crea una nueva instancia del handler de contactos
func NewContactHandler(contactService *services.ContactService) *ContactHandler {
	return &ContactHandler{
		contactService: contactService,
	}
}

// UploadExcel maneja la carga de archivos Excel
func (h *ContactHandler) UploadExcel(c *gin.Context) {
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No se pudo leer el archivo"})
		return
	}
	defer file.Close()

	// Leer archivo Excel
	f, err := excelize.OpenReader(file)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No se pudo abrir el archivo Excel"})
		return
	}
	defer f.Close()

	// Obtener la primera hoja
	sheets := f.GetSheetList()
	if len(sheets) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "El archivo no contiene hojas"})
		return
	}

	rows, err := f.GetRows(sheets[0])
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No se pudieron leer las filas"})
		return
	}

	if len(rows) < 2 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "El archivo debe contener al menos un registro además del header"})
		return
	}

	// Procesar filas (saltar header)
	var contacts []*entities.Contact
	for i, row := range rows[1:] {
		if len(row) < 4 {
			continue
		}

		i=i
		// Limpiar y procesar cada campo
		clientKey := strings.TrimSpace(fmt.Sprintf("%v", row[0]))
		name := strings.TrimSpace(fmt.Sprintf("%v", row[1]))
		email := strings.TrimSpace(fmt.Sprintf("%v", row[2]))
		phone := strings.TrimSpace(fmt.Sprintf("%v", row[3]))

		// Remover caracteres no deseados del teléfono (espacios, guiones, etc.)
		phoneClean := strings.ReplaceAll(phone, " ", "")
		phoneClean = strings.ReplaceAll(phoneClean, "-", "")
		phoneClean = strings.ReplaceAll(phoneClean, "(", "")
		phoneClean = strings.ReplaceAll(phoneClean, ")", "")

		contact := &entities.Contact{
			ClientKey: clientKey,
			Name:      name,
			Email:     email,
			Phone:     phoneClean,
		}

		contacts = append(contacts, contact)
	}

	// Guardar contactos
	err = h.contactService.SaveContactsBatch(contacts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudieron guardar los contactos"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Archivo cargado exitosamente",
		"count":   len(contacts),
	})
}

// GetContacts obtiene contactos con paginación
func (h *ContactHandler) GetContacts(c *gin.Context) {
	// Parámetros de paginación
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "50")
	
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}
	
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 100 {
		limit = 50
	}

	allContacts, err := h.contactService.GetAllContacts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudieron obtener los contactos"})
		return
	}

	total := len(allContacts)
	start := (page - 1) * limit
	end := start + limit

	if start > total {
		start = total
	}
	if end > total {
		end = total
	}

	var paginatedContacts []*entities.Contact
	if start < total {
		paginatedContacts = allContacts[start:end]
	}

	totalPages := (total + limit - 1) / limit

	c.JSON(http.StatusOK, gin.H{
		"data":        paginatedContacts,
		"total":       total,
		"page":        page,
		"limit":       limit,
		"total_pages": totalPages,
		"has_next":    page < totalPages,
		"has_prev":    page > 1,
	})
}

// SearchContacts busca contactos por diferentes campos
func (h *ContactHandler) SearchContacts(c *gin.Context) {
	field := c.Query("field")
	value := c.Query("value")

	if field == "" || value == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Se requieren los parámetros 'field' y 'value'"})
		return
	}

	contacts, err := h.contactService.SearchContacts(field, value)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error en la búsqueda"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": contacts})
}

// UpdateContact actualiza un contacto específico
func (h *ContactHandler) UpdateContact(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var contact entities.Contact
	if err := c.ShouldBindJSON(&contact); err != nil {
		fmt.Printf("Error binding JSON: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	// Debug: mostrar los datos recibidos
	fmt.Printf("DEBUG: Actualizando contacto ID %d\n", id)
	fmt.Printf("DEBUG: Datos recibidos: %+v\n", contact)

	contact.ID = id
	
	// Limpiar y validar datos antes de actualizar
	contact.ClientKey = strings.TrimSpace(contact.ClientKey)
	contact.Name = strings.TrimSpace(contact.Name)
	contact.Email = strings.TrimSpace(contact.Email)
	contact.Phone = strings.TrimSpace(contact.Phone)
	
	fmt.Printf("DEBUG: Datos después de limpiar: %+v\n", contact)
	
	err = h.contactService.UpdateContact(&contact)
	if err != nil {
		fmt.Printf("Error updating contact: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo actualizar el contacto"})
		return
	}

	fmt.Printf("DEBUG: Contacto actualizado exitosamente\n")
	c.JSON(http.StatusOK, gin.H{
		"message": "Contacto actualizado exitosamente",
		"contact": contact,
	})
}

// ValidateContacts valida todos los contactos y retorna errores con paginación
func (h *ContactHandler) ValidateContacts(c *gin.Context) {
	// Parámetros de paginación
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "50")
	
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}
	
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 100 {
		limit = 50
	}

	// Obtener todas las validaciones
	allResults, err := h.contactService.ValidateAllContacts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo validar los contactos"})
		return
	}

	total := len(allResults)
	start := (page - 1) * limit
	end := start + limit

	if start > total {
		start = total
	}
	if end > total {
		end = total
	}

	var paginatedResults []*entities.ContactWithValidation
	if start < total {
		paginatedResults = allResults[start:end]
	}

	totalPages := (total + limit - 1) / limit

	// Calcular estadísticas de todos los resultados
	validCount := 0
	invalidCount := 0
	for _, result := range allResults {
		if result.IsValid {
			validCount++
		} else {
			invalidCount++
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"data":        paginatedResults,
		"total":       total,
		"page":        page,
		"limit":       limit,
		"total_pages": totalPages,
		"has_next":    page < totalPages,
		"has_prev":    page > 1,
		"stats": gin.H{
			"total":   total,
			"valid":   validCount,
			"invalid": invalidCount,
		},
	})
}

// DownloadExcel genera y descarga un archivo Excel con todos los contactos actuales
func (h *ContactHandler) DownloadExcel(c *gin.Context) {
	contacts, err := h.contactService.GetAllContacts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudieron obtener los contactos"})
		return
	}

	// Debug: verificar si hay contactos
	fmt.Printf("DEBUG: Número de contactos encontrados: %d\n", len(contacts))

	if len(contacts) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No hay contactos para descargar"})
		return
	}

	// Crear archivo Excel
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	
	sheetName := "Sheet1" // Cambiar a Sheet1 que es el default
	
	// Headers exactos como en tu estructura
	headers := []string{"Clave cliente", "   Nombre Contacto ", "Correo ", "Teléfono Contacto  "}
	for i, header := range headers {
		cell := fmt.Sprintf("%c1", 'A'+i)
		err := f.SetCellValue(sheetName, cell, header)
		if err != nil {
			fmt.Printf("Error setting header %s: %v\n", header, err)
		}
	}

	// Datos de contactos
	for i, contact := range contacts {
		row := i + 2
		fmt.Printf("DEBUG: Procesando contacto %d: %+v\n", i+1, contact)
		
		// Asegurar que todos los valores se escriban correctamente
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), contact.ClientKey)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), contact.Name)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), contact.Email)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), contact.Phone)
	}

	// Ajustar anchos de columna
	f.SetColWidth(sheetName, "A", "A", 15)
	f.SetColWidth(sheetName, "B", "B", 35)
	f.SetColWidth(sheetName, "C", "C", 40)
	f.SetColWidth(sheetName, "D", "D", 18)

	// Crear buffer temporal para escribir el archivo
	buf := new(bytes.Buffer)
	if err := f.Write(buf); err != nil {
		fmt.Printf("Error writing to buffer: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo generar el archivo Excel"})
		return
	}

	// Verificar que el buffer tiene contenido
	fmt.Printf("DEBUG: Tamaño del archivo Excel: %d bytes\n", buf.Len())

	// Configurar headers para descarga
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename=\"contactos_corregidos.xlsx\"")
	c.Header("Content-Length", fmt.Sprintf("%d", buf.Len()))

	// Escribir directamente el buffer al response
	if _, err := c.Writer.Write(buf.Bytes()); err != nil {
		fmt.Printf("Error writing response: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo enviar el archivo"})
		return
	}
}