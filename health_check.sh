#!/bin/bash

# Script de Health Check para Microservicios - Windows Compatible
# Arquitectura de Software 2 - Proyecto

echo "=== Health Check de Microservicios ==="
echo "Fecha: $(date)"
echo ""

# Colores para output (funciona en Git Bash)
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Función para verificar un servicio
check_service() {
    local service_name=$1
    local service_url=$2
    local expected_status=${3:-200}
    
    echo -n "Verificando $service_name... "
    
    # Intentar hacer una petición HTTP
    response=$(curl -s -o /dev/null -w "%{http_code}" --connect-timeout 5 --max-time 10 "$service_url" 2>/dev/null)
    
    if [ "$response" = "$expected_status" ]; then
        echo -e "${GREEN}✓ OK (Status: $response)${NC}"
        return 0
    else
        echo -e "${RED}✗ ERROR (Status: $response)${NC}"
        return 1
    fi
}

# Función para verificar un puerto (compatible con Windows)
check_port() {
    local service_name=$1
    local host=$2
    local port=$3
    
    echo -n "Verificando $service_name ($host:$port)... "
    
    # Verificar si el puerto está abierto usando netstat (Windows) o curl
    if command -v netstat >/dev/null 2>&1; then
        # Windows con netstat
        if netstat -an | grep -q ":$port.*LISTENING"; then
            echo -e "${GREEN}✓ OK${NC}"
            return 0
        else
            echo -e "${RED}✗ ERROR${NC}"
            return 1
        fi
    else
        # Fallback: usar curl para verificar
        if curl -s --connect-timeout 3 "http://$host:$port" >/dev/null 2>&1; then
            echo -e "${GREEN}✓ OK${NC}"
            return 0
        else
            echo -e "${RED}✗ ERROR${NC}"
            return 1
        fi
    fi
}

# Función para extraer JSON sin jq
extract_json_value() {
    local json="$1"
    local key="$2"
    
    # Extraer valor usando grep y sed (compatible con Windows)
    echo "$json" | grep -o "\"$key\":[^,}]" | sed 's/."'"$key"'":"\([^,"]\)"*/\1/'
}

# Contadores
total_checks=0
passed_checks=0

# Verificar servicios de base de datos
echo -e "${BLUE}--- Servicios de Base de Datos ---${NC}"
((total_checks++))
if check_port "MySQL" "localhost" "3307"; then
    ((passed_checks++))
fi

((total_checks++))
if check_port "MongoDB" "localhost" "27017"; then
    ((passed_checks++))
fi

((total_checks++))
if check_port "RabbitMQ" "localhost" "5672"; then
    ((passed_checks++))
fi

((total_checks++))
if check_port "Memcached" "localhost" "11211"; then
    ((passed_checks++))
fi

echo ""

# Verificar microservicios
echo -e "${BLUE}--- Microservicios ---${NC}"
((total_checks++))
if check_service "Users Service" "http://localhost/users/health" "200"; then
    ((passed_checks++))
fi

((total_checks++))
if check_service "Subscriptions Service" "http://localhost/subscriptions/health" "200"; then
    ((passed_checks++))
fi

((total_checks++))
if check_service "Admin Service" "http://localhost/admin/health" "200"; then
    ((passed_checks++))
fi

echo ""

# Verificar servicios de infraestructura
echo -e "${BLUE}--- Servicios de Infraestructura ---${NC}"
((total_checks++))
if check_port "Solr" "localhost" "8983"; then
    ((passed_checks++))
fi

((total_checks++))
if check_service "NGinX Load Balancer" "http://localhost/health" "200"; then
    ((passed_checks++))
fi

((total_checks++))
if check_service "Frontend" "http://localhost/" "200"; then
    ((passed_checks++))
fi

echo ""

# Verificar balanceo de carga para usuarios
echo -e "${BLUE}--- Verificación de Balanceo de Carga ---${NC}"
echo "Verificando balanceo de carga para el servicio de usuarios..."

# Hacer múltiples peticiones para ver el balanceo
for i in {1..10}; do
    response=$(curl -s http://localhost/users/health 2>/dev/null)
    if [ $? -eq 0 ]; then
        status=$(extract_json_value "$response" "status")
        echo "Petición $i: Status $status"
    else
        echo "Petición $i: Error de conexión"
    fi
done

echo ""

# Resumen final
echo -e "${BLUE}=== Resumen ===${NC}"
echo "Total de verificaciones: $total_checks"
echo "Verificaciones exitosas: $passed_checks"
echo "Verificaciones fallidas: $((total_checks - passed_checks))"

if [ $passed_checks -eq $total_checks ]; then
    echo -e "${GREEN}✓ Todos los servicios están funcionando correctamente${NC}"
    exit 0
else
    echo -e "${RED}✗ Algunos servicios tienen problemas${NC}"
    exit 1
fi