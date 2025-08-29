# Cloudflare Domain Controller

Una herramienta CLI para gestionar registros DNS en Cloudflare de forma rápida y sencilla.

## Descripción

Cloudflare Domain Controller es una utilidad de línea de comandos desarrollada en Go que permite administrar fácilmente los registros DNS de tu dominio en Cloudflare. Con esta herramienta, puedes agregar, actualizar, eliminar y listar registros DNS sin necesidad de acceder al panel web de Cloudflare.

## Requisitos

- Go 1.24.6 o superior (para compilación)
- Una cuenta en Cloudflare con un dominio configurado
- Un token de API de Cloudflare con permisos para gestionar registros DNS

## Instalación

### Desde el código fuente

1. Clona el repositorio:
   ```bash
   git clone <URL_DEL_REPOSITORIO>
   cd cloudflare-domain-controller
   ```

2. Construye el binario:
   ```bash
   make build
   ```

3. (Opcional) Instala el binario en tu sistema:
   ```bash
   make install
   ```

### Desde binarios precompilados

Descarga el binario precompilado para tu plataforma desde la sección de [releases](https://github.com/tu-usuario/cloudflare-domain-controller/releases) y colócalo en tu PATH.

## Configuración

La herramienta requiere las siguientes variables de entorno:

- `CLOUDFLARE_API_TOKEN`: Tu token de API de Cloudflare
- `CLOUDFLARE_ZONE_ID`: El ID de la zona de tu dominio en Cloudflare
- `CLOUDFLARE_DOMAIN_NAME`: El nombre de tu dominio principal (ejemplo.com)

### Configuración permanente de variables de entorno

Agrega las siguientes líneas a tu archivo de perfil de shell (`~/.zshrc` para zsh o `~/.bash_profile` para bash):

```bash
export CLOUDFLARE_API_TOKEN="tu_token_de_api"
export CLOUDFLARE_ZONE_ID="id_de_tu_zona"
export CLOUDFLARE_DOMAIN_NAME="tu_dominio.com"
```

Luego recarga la configuración:
```bash
source ~/.zshrc  # o source ~/.bash_profile
```

### Obtención de credenciales de Cloudflare

1. **CLOUDFLARE_ZONE_ID**: 
   - Ve a tu panel de Cloudflare
   - Selecciona tu dominio
   - En la página de "Overview", encontrarás el "Zone ID" en el panel derecho

2. **CLOUDFLARE_API_TOKEN**:
   - Ve a "User Profile" > "API Tokens" en Cloudflare
   - Crea un nuevo token con permisos para:
     - Zone: DNS: Edit (para gestionar registros DNS)
     - Zone Resources: Include: Specific zone (selecciona tu dominio)

## Uso

### Agregar un registro DNS

```bash
cloudflare-domain-controller add mipagina --type A --content 192.168.1.1
```

Parámetros:
- `mipagina`: Nombre del subdominio (la herramienta automáticamente lo combina con tu dominio principal)
- `--type`: Tipo de registro DNS (A, CNAME, etc.)
- `--content`: Valor del registro (IP para registros A, nombre de dominio para CNAME, etc.)

### Actualizar un registro DNS

```bash
cloudflare-domain-controller update mipagina --type A --content 192.168.1.2
```

### Eliminar un registro DNS

```bash
cloudflare-domain-controller delete mipagina
```

### Listar todos los registros DNS

```bash
cloudflare-domain-controller list
```

### Ayuda

Para ver todas las opciones disponibles:

```bash
cloudflare-domain-controller --help
```

Para obtener ayuda específica de un comando:

```bash
cloudflare-domain-controller add --help
cloudflare-domain-controller update --help
cloudflare-domain-controller delete --help
cloudflare-domain-controller list --help
```

## Compilación cruzada

Para compilar el binario para diferentes plataformas:

```bash
make cross-compile
```

Los binarios se guardarán en el directorio `dist/`.

## Desarrollo

### Estructura del proyecto

```
cloudflare-domain-controller/
├── cmd/              # Comandos de la CLI
├── core/             # Lógica principal y cliente de Cloudflare
├── utils/            # Funciones auxiliares
├── main.go           # Punto de entrada
├── go.mod            # Dependencias del módulo Go
├── Makefile          # Scripts de compilación
└── README.md         # Este archivo
```

### Dependencias

- `github.com/spf13/cobra`: Para la creación de comandos CLI

### Compilación local

```bash
go build -o cloudflare-domain-controller .
```

### Ejecución de pruebas

```bash
go test ./...
```

## Seguridad

⚠️ **Importante**: 
- Nunca commitees tus credenciales en el repositorio git
- Usa tokens de API con el mínimo de permisos necesarios
- Considera rotar regularmente tus tokens de API

## Licencia

MIT

## Autor

[Freddy Solorzano](https://github.com/freddysolorzano)