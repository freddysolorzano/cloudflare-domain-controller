package cmd

import (
	"fmt"
	"os"

	"cloudflare-domain-controller/core"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete [subdominio]",
	Short: "Elimina un registro DNS",
	Long: `Elimina el registro DNS asociado al subdominio especificado.
Ejemplo: cloudflare-domain-controller delete mipagina`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		subdomain := args[0]
		
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
		
		// Eliminar el registro
		if err := client.DeleteDNSRecord(record.ID); err != nil {
			fmt.Fprintf(os.Stderr, "Error al eliminar el registro DNS: %v\n", err)
			os.Exit(1)
		}
		
		fmt.Printf("Registro DNS para %s eliminado exitosamente\n", subdomain)
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}