package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "cloudflare-domain-controller",
	Short: "Una herramienta CLI para gestionar registros DNS en Cloudflare",
	Long: `Una herramienta CLI que permite agregar, modificar y eliminar registros DNS
en Cloudflare mediante comandos simples.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	// Aquí se pueden agregar flags globales si son necesarios
}