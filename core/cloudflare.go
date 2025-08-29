package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

// Config almacena la configuración de Cloudflare
type Config struct {
	APIToken   string
	ZoneID     string
	BaseURL    string
	DomainName string
}

// NewConfig crea una nueva configuración desde variables de entorno
func NewConfig() *Config {
	return &Config{
		APIToken:   os.Getenv("CLOUDFLARE_API_TOKEN"),
		ZoneID:     os.Getenv("CLOUDFLARE_ZONE_ID"),
		DomainName: os.Getenv("CLOUDFLARE_DOMAIN_NAME"),
		BaseURL:    "https://api.cloudflare.com/client/v4",
	}
}

// Validate verifica que todas las configuraciones necesarias estén presentes
func (c *Config) Validate() error {
	if c.APIToken == "" {
		return fmt.Errorf("CLOUDFLARE_API_TOKEN no está configurado")
	}
	if c.ZoneID == "" {
		return fmt.Errorf("CLOUDFLARE_ZONE_ID no está configurado")
	}
	if c.DomainName == "" {
		return fmt.Errorf("CLOUDFLARE_DOMAIN_NAME no está configurado")
	}
	return nil
}

// DNSRecord representa un registro DNS
type DNSRecord struct {
	ID      string `json:"id,omitempty"`
	Name    string `json:"name"`
	Type    string `json:"type"`
	Content string `json:"content"`
	TTL     int    `json:"ttl"`
	Proxied bool   `json:"proxied"`
}

// CloudflareClient representa un cliente para interactuar con la API de Cloudflare
type CloudflareClient struct {
	config *Config
}

// NewCloudflareClient crea un nuevo cliente de Cloudflare
func NewCloudflareClient(config *Config) *CloudflareClient {
	return &CloudflareClient{
		config: config,
	}
}

// makeRequest realiza una solicitud HTTP a la API de Cloudflare
func (c *CloudflareClient) makeRequest(method, url string, body io.Reader) ([]byte, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.config.APIToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("error en la solicitud: %s - %s", resp.Status, string(respBody))
	}

	return respBody, nil
}

// CreateDNSRecord crea un nuevo registro DNS
func (c *CloudflareClient) CreateDNSRecord(record *DNSRecord) error {
	// Validar configuración
	if err := c.config.Validate(); err != nil {
		return err
	}
	
	url := fmt.Sprintf("%s/zones/%s/dns_records", c.config.BaseURL, c.config.ZoneID)
	
	jsonData, err := json.Marshal(record)
	if err != nil {
		return err
	}

	_, err = c.makeRequest("POST", url, bytes.NewBuffer(jsonData))
	return err
}

// UpdateDNSRecord actualiza un registro DNS existente
func (c *CloudflareClient) UpdateDNSRecord(recordID string, record *DNSRecord) error {
	// Validar configuración
	if err := c.config.Validate(); err != nil {
		return err
	}
	
	url := fmt.Sprintf("%s/zones/%s/dns_records/%s", c.config.BaseURL, c.config.ZoneID, recordID)
	
	jsonData, err := json.Marshal(record)
	if err != nil {
		return err
	}

	_, err = c.makeRequest("PATCH", url, bytes.NewBuffer(jsonData))
	return err
}

// DeleteDNSRecord elimina un registro DNS
func (c *CloudflareClient) DeleteDNSRecord(recordID string) error {
	// Validar configuración
	if err := c.config.Validate(); err != nil {
		return err
	}
	
	url := fmt.Sprintf("%s/zones/%s/dns_records/%s", c.config.BaseURL, c.config.ZoneID, recordID)
	
	_, err := c.makeRequest("DELETE", url, nil)
	return err
}

// GetDNSRecordByName obtiene un registro DNS por su nombre
func (c *CloudflareClient) GetDNSRecordByName(name string) (*DNSRecord, error) {
	// Validar configuración
	if err := c.config.Validate(); err != nil {
		return nil, err
	}
	
	// Construir el nombre completo si solo se proporciona el subdominio
	fullName := name
	if c.config.DomainName != "" {
		// Verificar si el nombre ya incluye el dominio
		if len(name) <= len(c.config.DomainName) || (len(name) > len(c.config.DomainName) && name[len(name)-len(c.config.DomainName)-1:] != "."+c.config.DomainName) {
			fullName = name + "." + c.config.DomainName
		}
	}
	
	url := fmt.Sprintf("%s/zones/%s/dns_records?name=%s", c.config.BaseURL, c.config.ZoneID, fullName)
	
	respBody, err := c.makeRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Parsear la respuesta para obtener el registro
	var result map[string]interface{}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, err
	}

	// Verificar si hay resultados
	resultArray, ok := result["result"].([]interface{})
	if !ok || len(resultArray) == 0 {
		return nil, fmt.Errorf("no se encontró el registro DNS para %s", name)
	}

	// Tomar el primer resultado
	recordData := resultArray[0].(map[string]interface{})
	
	record := &DNSRecord{
		ID:      recordData["id"].(string),
		Name:    recordData["name"].(string),
		Type:    recordData["type"].(string),
		Content: recordData["content"].(string),
		TTL:     int(recordData["ttl"].(float64)),
		Proxied: recordData["proxied"].(bool),
	}

	return record, nil
}

// ListDNSRecords lista todos los registros DNS de la zona
func (c *CloudflareClient) ListDNSRecords() ([]*DNSRecord, error) {
	// Validar configuración
	if err := c.config.Validate(); err != nil {
		return nil, err
	}
	
	url := fmt.Sprintf("%s/zones/%s/dns_records", c.config.BaseURL, c.config.ZoneID)
	
	respBody, err := c.makeRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Parsear la respuesta para obtener los registros
	var result map[string]interface{}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, err
	}

	// Verificar si hay resultados
	resultArray, ok := result["result"].([]interface{})
	if !ok || len(resultArray) == 0 {
		return []*DNSRecord{}, nil
	}

	// Convertir los resultados a registros DNS
	records := make([]*DNSRecord, len(resultArray))
	for i, item := range resultArray {
		recordData := item.(map[string]interface{})
		
		record := &DNSRecord{
			ID:      recordData["id"].(string),
			Name:    recordData["name"].(string),
			Type:    recordData["type"].(string),
			Content: recordData["content"].(string),
			TTL:     int(recordData["ttl"].(float64)),
			Proxied: recordData["proxied"].(bool),
		}
		
		records[i] = record
	}

	return records, nil
}