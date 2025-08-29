package cmd

import (
	"fmt"
	"os"

	"cloudflare-domain-controller/core"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Lista todos los registros DNS",
	Long:  `Lista todos los registros DNS configurados en la zona de Cloudflare.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Crear cliente de Cloudflare
		config := core.NewConfig()
		client := core.NewCloudflareClient(config)
		
		// Obtener todos los registros
		records, err := client.ListDNSRecords()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error al obtener los registros DNS: %v\n", err)
			os.Exit(1)
		}
		
		// Mostrar los registros
		if len(records) == 0 {
			fmt.Println("No se encontraron registros DNS.")
			return
		}
		
		fmt.Printf("Registros DNS encontrados (%d):\n", len(records))
		fmt.Println("----------------------------------------")
		for _, record := range records {
			// Mostrar solo el subdominio si pertenece al dominio principal
			displayName := record.Name
			if config.DomainName != "" && len(record.Name) > len(config.DomainName) {
				if record.Name[len(record.Name)-len(config.DomainName):] == config.DomainName {
					displayName = record.Name[:len(record.Name)-len(config.DomainName)-1]
				}
			}
			
			fmt.Printf("%-20s %-6s %-15s\n", displayName, record.Type, record.Content)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}