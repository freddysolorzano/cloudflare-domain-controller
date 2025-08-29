package cmd

import (
	"fmt"
	"os"

	"cloudflare-domain-controller/core"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update [subdominio]",
	Short: "Actualiza un registro DNS existente",
	Long: `Actualiza un registro DNS existente para el subdominio especificado.
Ejemplo: cloudflare-domain-controller update mipagina --type A --content 192.168.1.2`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		subdomain := args[0]
		recordType, _ := cmd.Flags().GetString("type")
		content, _ := cmd.Flags().GetString("content")
		
		// Crear cliente de Cloudflare
		config := core.NewConfig()
		client := core.NewCloudflareClient(config)
		
		// Construir el nombre completo del registro
		fullName := subdomain
		if config.DomainName != "" {
			// Verificar si el subdominio ya incluye el dominio
			if len(subdomain) <= len(config.DomainName) || subdomain[len(subdomain)-len(config.DomainName):] != config.DomainName {
				fullName = subdomain + "." + config.DomainName
			}
		}
		
		// Obtener el registro existente
		record, err := client.GetDNSRecordByName(fullName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error al obtener el registro DNS: %v\n", err)
			os.Exit(1)
		}
		
		// Actualizar los campos
		record.Type = recordType
		record.Content = content
		
		// Actualizar el registro
		if err := client.UpdateDNSRecord(record.ID, record); err != nil {
			fmt.Fprintf(os.Stderr, "Error al actualizar el registro DNS: %v\n", err)
			os.Exit(1)
		}
		
		fmt.Printf("Registro DNS para %s actualizado exitosamente\n", subdomain)
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
	updateCmd.Flags().StringP("type", "t", "A", "Tipo de registro DNS (A, CNAME, etc.)")
	updateCmd.Flags().StringP("content", "c", "", "Nuevo contenido del registro DNS (IP o CNAME)")
	// Requerir el flag 'content'
	updateCmd.MarkFlagRequired("content")
}