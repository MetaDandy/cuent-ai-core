#!/bin/bash

# Script auxiliar para gestionar screenshots de Selenium
# Uso: ./manage-screenshots.sh [comando] [opciones]

SCREENSHOTS_DIR="tests/selenium/screenshots"
COLORS_RESET='\033[0m'
COLORS_GREEN='\033[0;32m'
COLORS_YELLOW='\033[1;33m'
COLORS_BLUE='\033[0;34m'

# Crear directorio si no existe
mkdir -p "$SCREENSHOTS_DIR"

# Función para mostrar ayuda
show_help() {
    cat << EOF
${COLORS_BLUE}Selenium Screenshots Manager${COLORS_RESET}

Uso: $(basename "$0") [comando] [opciones]

${COLORS_GREEN}Comandos:${COLORS_RESET}

  list                     - Lista todos los directorios de screenshots
  list-latest             - Muestra el directorio de screenshots más reciente
  view <test_name>        - Abre el directorio de screenshots en el explorador
  view-latest             - Abre el último directorio de screenshots
  compare <dir1> <dir2>   - Compara dos directorios de screenshots
  clean                   - Elimina todos los screenshots
  clean-old [days]        - Elimina screenshots más antiguos que X días (default: 7)
  size                    - Muestra el tamaño total de screenshots
  zip                     - Comprime todos los screenshots en un archivo
  open-last               - Abre el último directorio de screenshots

${COLORS_GREEN}Ejemplos:${COLORS_RESET}

  ./$(basename "$0") list
  ./$(basename "$0") view-latest
  ./$(basename "$0") clean-old 7
  ./$(basename "$0") zip

EOF
}

# Función para listar screenshots
list_screenshots() {
    if [ ! -d "$SCREENSHOTS_DIR" ]; then
        echo "⚠ Directorio de screenshots no existe: $SCREENSHOTS_DIR"
        return 1
    fi

    echo -e "${COLORS_GREEN}Directorios de screenshots encontrados:${COLORS_RESET}"
    ls -lhd "$SCREENSHOTS_DIR"/*/ 2>/dev/null | awk '{print $NF, "(" $5 ")"}'

    if [ $? -ne 0 ]; then
        echo "⚠ No hay directorios de screenshots"
        return 1
    fi
}

# Función para obtener el último directorio
get_latest_dir() {
    latest=$(ls -td "$SCREENSHOTS_DIR"/*/ 2>/dev/null | head -1)
    if [ -z "$latest" ]; then
        echo ""
        return 1
    fi
    echo "$latest"
}

# Función para abrir directorio
open_dir() {
    local dir="$1"
    if [ ! -d "$dir" ]; then
        echo "⚠ Directorio no existe: $dir"
        return 1
    fi

    echo -e "${COLORS_GREEN}✓ Abriendo: $dir${COLORS_RESET}"
    
    # Detectar el sistema operativo
    case "$(uname)" in
        Linux*)
            if command -v xdg-open &> /dev/null; then
                xdg-open "$dir"
            elif command -v nautilus &> /dev/null; then
                nautilus "$dir" &
            else
                echo "⚠ No se puede abrir el directorio automáticamente"
                echo "   Abre manualmente: $dir"
            fi
            ;;
        Darwin*)
            open "$dir"
            ;;
        *)
            echo "⚠ Sistema operativo no soportado"
            ;;
    esac
}

# Función para limpiar screenshots
clean_screenshots() {
    echo "⚠ Eliminando todos los screenshots..."
    read -p "¿Estás seguro? (s/n): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Ss]$ ]]; then
        rm -rf "$SCREENSHOTS_DIR"/*
        echo -e "${COLORS_GREEN}✓ Screenshots eliminados${COLORS_RESET}"
    else
        echo "Cancelado"
    fi
}

# Función para limpiar screenshots antiguos
clean_old_screenshots() {
    local days=${1:-7}
    echo -e "${COLORS_YELLOW}Eliminando screenshots más antiguos que $days días...${COLORS_RESET}"
    
    find "$SCREENSHOTS_DIR" -maxdepth 1 -type d -mtime "+$days" -exec rm -rf {} \;
    
    echo -e "${COLORS_GREEN}✓ Limpieza completada${COLORS_RESET}"
}

# Función para mostrar tamaño
show_size() {
    if [ ! -d "$SCREENSHOTS_DIR" ]; then
        echo "⚠ Directorio de screenshots no existe"
        return 1
    fi

    local total_size=$(du -sh "$SCREENSHOTS_DIR" | cut -f1)
    echo -e "${COLORS_GREEN}Tamaño total de screenshots: $total_size${COLORS_RESET}"
    
    echo -e "\n${COLORS_BLUE}Desglose por directorio:${COLORS_RESET}"
    du -sh "$SCREENSHOTS_DIR"/*/ 2>/dev/null | sort -hr
}

# Función para comprimir screenshots
zip_screenshots() {
    if [ ! -d "$SCREENSHOTS_DIR" ]; then
        echo "⚠ Directorio de screenshots no existe"
        return 1
    fi

    local zip_file="screenshots_$(date +%Y%m%d_%H%M%S).zip"
    echo -e "${COLORS_YELLOW}Comprimiendo screenshots a $zip_file...${COLORS_RESET}"
    
    zip -r "$zip_file" "$SCREENSHOTS_DIR" > /dev/null
    
    echo -e "${COLORS_GREEN}✓ Archivo creado: $zip_file${COLORS_RESET}"
    ls -lh "$zip_file"
}

# Función para comparar directorios
compare_dirs() {
    local dir1="$1"
    local dir2="$2"

    if [ ! -d "$dir1" ]; then
        echo "⚠ Directorio no existe: $dir1"
        return 1
    fi

    if [ ! -d "$dir2" ]; then
        echo "⚠ Directorio no existe: $dir2"
        return 1
    fi

    echo -e "${COLORS_BLUE}Comparando:${COLORS_RESET}"
    echo "  Dir 1: $dir1"
    echo "  Dir 2: $dir2"
    echo

    diff -r "$dir1" "$dir2"
}

# Main
case "${1:-help}" in
    list)
        list_screenshots
        ;;
    list-latest)
        latest=$(get_latest_dir)
        if [ -n "$latest" ]; then
            echo -e "${COLORS_GREEN}Última carpeta de screenshots:${COLORS_RESET}"
            echo "$latest"
        else
            echo "⚠ No hay screenshots"
        fi
        ;;
    view)
        if [ -z "$2" ]; then
            echo "⚠ Debes proporcionar el nombre del directorio"
            echo "   Uso: $(basename "$0") view <test_name>"
            exit 1
        fi
        open_dir "$SCREENSHOTS_DIR/$2"
        ;;
    view-latest|open-last)
        latest=$(get_latest_dir)
        if [ -n "$latest" ]; then
            open_dir "$latest"
        else
            echo "⚠ No hay screenshots"
        fi
        ;;
    compare)
        if [ -z "$2" ] || [ -z "$3" ]; then
            echo "⚠ Debes proporcionar dos directorios"
            echo "   Uso: $(basename "$0") compare <dir1> <dir2>"
            exit 1
        fi
        compare_dirs "$2" "$3"
        ;;
    clean)
        clean_screenshots
        ;;
    clean-old)
        clean_old_screenshots "${2:-7}"
        ;;
    size)
        show_size
        ;;
    zip)
        zip_screenshots
        ;;
    help|--help|-h)
        show_help
        ;;
    *)
        echo "⚠ Comando no reconocido: $1"
        echo "Usa: $(basename "$0") help"
        exit 1
        ;;
esac
