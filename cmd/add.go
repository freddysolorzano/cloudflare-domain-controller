package cmd

import (
	"fmt"
	"os"

	"cloudflare-domain-controller/core"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add [subdominio]",
	Short: "Agrega un nuevo registro DNS",
	Long: `Agrega un nuevo registro DNS para el subdominio especificado.
Ejemplo: cloudflare-domain-controller add mipagina --type A --content 192.168.1.1`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		subdomain := args[0]
		recordType, _ := cmd.Flags().GetString("type")
		content, _ := cmd.Flags().GetString("content")
		
		// Crear cliente de Cloudflare
		config := core.NewConfig()
		// Validar configuración
		if err := config.Validate(); err != nil {
			fmt.Fprintf(os.Stderr, "Error de configuración: %v\n", err)
			os.Exit(1)
		}
		
		client := core.NewCloudflareClient(config)
		
		// Construir el nombre completo del registro
		fullName := subdomain
		if config.DomainName != "" {
			// Verificar si el subdominio ya incluye el dominio
			if len(subdomain) <= len(config.DomainName) || (len(subdomain) > len(config.DomainName) && subdomain[len(subdomain)-len(config.DomainName)-1:] != "."+config.DomainName) {
				fullName = subdomain + "." + config.DomainName
			}
		}
		
		// Crear el registro DNS
		record := &core.DNSRecord{
			Name:    fullName,
			Type:    recordType,
			Content: content,
			TTL:     1, // Auto
			Proxied: false,
		}
		
		// Intentar crear el registro
		if err := client.CreateDNSRecord(record); err != nil {
			fmt.Fprintf(os.Stderr, "Error al agregar el registro DNS: %v\n", err)
			os.Exit(1)
		}
		
		fmt.Printf("Registro DNS para %s agregado exitosamente\n", subdomain)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.Flags().StringP("type", "t", "A", "Tipo de registro DNS (A, CNAME, etc.)")
	addCmd.Flags().StringP("content", "c", "", "Contenido del registro DNS (IP o CNAME)")
	// Requerir el flag 'content'
	addCmd.MarkFlagRequired("content")
}