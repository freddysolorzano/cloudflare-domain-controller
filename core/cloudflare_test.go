package core

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

// mockCloudflareServer crea un servidor HTTP de prueba que simula la API de Cloudflare
func mockCloudflareServer() *httptest.Server {
	// Almacenamiento en memoria para los registros DNS simulados
	records := make(map[string]*DNSRecord)
	
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verificar que el token de autorización esté presente
		authHeader := r.Header.Get("Authorization")
		if authHeader != "Bearer test-token" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		
		// Manejar diferentes rutas de la API
		switch {
		// Crear un nuevo registro DNS
		case r.Method == http.MethodPost && r.URL.Path == "/client/v4/zones/test-zone-id/dns_records":
			var record DNSRecord
			body, _ := io.ReadAll(r.Body)
			json.Unmarshal(body, &record)
			
			// Asignar un ID simulado si no está presente
			if record.ID == "" {
				record.ID = "test-record-id"
			}
			// Guardar una copia del registro
			storedRecord := DNSRecord{
				ID:      record.ID,
				Name:    record.Name,
				Type:    record.Type,
				Content: record.Content,
				TTL:     record.TTL,
				Proxied: record.Proxied,
			}
			records[record.ID] = &storedRecord
			
			// Responder con el registro creado
			response := map[string]interface{}{
				"success": true,
				"errors":  []string{},
				"result":  storedRecord,
			}
			json.NewEncoder(w).Encode(response)
			
			// Actualizar el registro original con el ID asignado
			record.ID = storedRecord.ID
			
		// Actualizar un registro DNS existente
		case r.Method == http.MethodPatch && r.URL.Path == "/client/v4/zones/test-zone-id/dns_records/test-record-id":
			body, _ := io.ReadAll(r.Body)
			var requestData map[string]interface{}
			json.Unmarshal(body, &requestData)
			
			// Verificar si el registro existe
			record, exists := records["test-record-id"]
			if !exists {
				http.Error(w, "Record not found", http.StatusNotFound)
				return
			}
			
			// Actualizar solo los campos proporcionados
			if name, ok := requestData["name"].(string); ok {
				record.Name = name
			}
			if recordType, ok := requestData["type"].(string); ok {
				record.Type = recordType
			}
			if content, ok := requestData["content"].(string); ok {
				record.Content = content
			}
			if ttl, ok := requestData["ttl"].(float64); ok {
				record.TTL = int(ttl)
			}
			if proxied, ok := requestData["proxied"].(bool); ok {
				record.Proxied = proxied
			}
			
			// Responder con el registro actualizado
			response := map[string]interface{}{
				"success": true,
				"errors":  []string{},
				"result":  record,
			}
			json.NewEncoder(w).Encode(response)
			
		// Eliminar un registro DNS
		case r.Method == http.MethodDelete && r.URL.Path == "/client/v4/zones/test-zone-id/dns_records/test-record-id":
			delete(records, "test-record-id")
			
			// Responder con éxito
			response := map[string]interface{}{
				"success": true,
				"errors":  []string{},
				"result":  nil,
			}
			json.NewEncoder(w).Encode(response)
			
		// Obtener un registro DNS por nombre o listar todos los registros
		case r.Method == http.MethodGet && r.URL.Path == "/client/v4/zones/test-zone-id/dns_records":
			// Obtener el parámetro de consulta "name"
			name := r.URL.Query().Get("name")
			
			// Si se proporciona un nombre, buscar el registro con ese nombre
			if name != "" {
				var foundRecord *DNSRecord
				for _, record := range records {
					if record.Name == name {
						foundRecord = record
						break
					}
				}
				
				// Responder con el registro encontrado o una lista vacía
				var result []interface{}
				if foundRecord != nil {
					result = append(result, foundRecord)
				}
				
				response := map[string]interface{}{
					"success": true,
					"errors":  []string{},
					"result":  result,
				}
				json.NewEncoder(w).Encode(response)
			} else {
				// Si no se proporciona un nombre, listar todos los registros
				// Convertir el mapa de registros a una lista
				var result []interface{}
				for _, record := range records {
					result = append(result, record)
				}
				
				response := map[string]interface{}{
					"success": true,
					"errors":  []string{},
					"result":  result,
				}
				json.NewEncoder(w).Encode(response)
			}
			
		// Ruta no encontrada
		default:
			http.Error(w, "Not Found", http.StatusNotFound)
		}
	}))
}

func TestCloudflareClient(t *testing.T) {
	// Crear un servidor de prueba
	server := mockCloudflareServer()
	defer server.Close()
	
	// Establecer variables de entorno para la configuración
	os.Setenv("CLOUDFLARE_API_TOKEN", "test-token")
	os.Setenv("CLOUDFLARE_ZONE_ID", "test-zone-id")
	os.Setenv("CLOUDFLARE_DOMAIN_NAME", "test-domain.com")
	
	// Crear una configuración y cliente de prueba
	config := NewConfig()
	config.BaseURL = server.URL + "/client/v4" // Usar la URL del servidor de prueba
	
	// Validar la configuración
	if err := config.Validate(); err != nil {
		t.Fatalf("Error en la configuración: %v", err)
	}
	
	client := NewCloudflareClient(config)
	
	// Variable para almacenar el ID del registro de prueba
	_ = ""
	
	// Prueba: Crear un registro DNS
	t.Run("CreateDNSRecord", func(t *testing.T) {
		record := &DNSRecord{
			Name:    "test.test-domain.com",
			Type:    "A",
			Content: "192.168.1.1",
			TTL:     1,
			Proxied: false,
		}
		
		if err := client.CreateDNSRecord(record); err != nil {
			t.Errorf("Error al crear el registro DNS: %v", err)
		}
		
		// Verificar que el ID se haya asignado
		if record.ID == "" {
			t.Errorf("El ID del registro no se asignó correctamente: %v", record)
		}
		
		// Guardar el ID para usarlo en otras pruebas
			_ = record.ID
	})
	
	// Prueba: Obtener un registro DNS por nombre
	t.Run("GetDNSRecordByName", func(t *testing.T) {
		record, err := client.GetDNSRecordByName("test.test-domain.com")
		if err != nil {
			t.Errorf("Error al obtener el registro DNS: %v", err)
		}
		
		// Verificar los datos del registro
		if record.Name != "test.test-domain.com" {
			t.Errorf("Nombre del registro incorrecto: esperado 'test.test-domain.com', obtenido '%s'", record.Name)
		}
		if record.Type != "A" {
			t.Errorf("Tipo del registro incorrecto: esperado 'A', obtenido '%s'", record.Type)
		}
		if record.Content != "192.168.1.1" {
			t.Errorf("Contenido del registro incorrecto: esperado '192.168.1.1', obtenido '%s'", record.Content)
		}
	})
	
	// Prueba: Actualizar un registro DNS
	t.Run("UpdateDNSRecord", func(t *testing.T) {
		// Obtener el registro existente
		record, err := client.GetDNSRecordByName("test.test-domain.com")
		if err != nil {
			t.Fatalf("Error al obtener el registro DNS: %v", err)
		}
		
		// Actualizar el contenido del registro
		record.Content = "192.168.1.2"
		if err := client.UpdateDNSRecord(record.ID, record); err != nil {
			t.Errorf("Error al actualizar el registro DNS: %v", err)
		}
		
		// Verificar que el contenido se haya actualizado
		updatedRecord, err := client.GetDNSRecordByName("test.test-domain.com")
		if err != nil {
			t.Fatalf("Error al obtener el registro DNS actualizado: %v", err)
		}
		
		if updatedRecord.Content != "192.168.1.2" {
			t.Errorf("El contenido del registro no se actualizó correctamente: esperado '192.168.1.2', obtenido '%s'", updatedRecord.Content)
		}
	})
	
	// Prueba: Listar todos los registros DNS
	t.Run("ListDNSRecords", func(t *testing.T) {
		records, err := client.ListDNSRecords()
		if err != nil {
			t.Errorf("Error al listar los registros DNS: %v", err)
		}
		
		// Verificar que haya al menos un registro
		if len(records) == 0 {
			t.Error("No se encontraron registros DNS")
		}
		
		// Verificar los datos del primer registro
		record := records[0]
		if record.Name != "test.test-domain.com" {
			t.Errorf("Nombre del registro incorrecto: esperado 'test.test-domain.com', obtenido '%s'", record.Name)
		}
		if record.Type != "A" {
			t.Errorf("Tipo del registro incorrecto: esperado 'A', obtenido '%s'", record.Type)
		}
		if record.Content != "192.168.1.2" {
			t.Errorf("Contenido del registro incorrecto: esperado '192.168.1.2', obtenido '%s'", record.Content)
		}
	})
	
	// Prueba: Eliminar un registro DNS
	t.Run("DeleteDNSRecord", func(t *testing.T) {
		// Obtener el registro existente
		record, err := client.GetDNSRecordByName("test.test-domain.com")
		if err != nil {
			t.Fatalf("Error al obtener el registro DNS: %v", err)
		}
		
		// Eliminar el registro
		if err := client.DeleteDNSRecord(record.ID); err != nil {
			t.Errorf("Error al eliminar el registro DNS: %v", err)
		}
		
		// Verificar que el registro se haya eliminado
		_, err = client.GetDNSRecordByName("test.test-domain.com")
		if err == nil {
			t.Error("El registro DNS no se eliminó correctamente")
		}
	})
}