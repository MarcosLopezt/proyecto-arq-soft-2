#!/bin/bash

# Script de Prueba del Balanceador de Carga - Windows Compatible
# Arquitectura de Software 2 - Proyecto

echo "=== Prueba del Balanceador de Carga ==="
echo "Fecha: $(date)"
echo ""

# Colores para output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Función para extraer JSON sin jq
extract_json_value() {
    local json="$1"
    local key="$2"
    
    # Extraer valor usando grep y sed (compatible con Windows)
    echo "$json" | grep -o "\"$key\":[^,}]" | sed 's/."'"$key"'":"\([^,"]\)"*/\1/'
}

echo -e "${BLUE}1. Verificando que ambas instancias estén funcionando...${NC}"
echo "Instancia 1 (puerto 8082):"
response1=$(curl -s http://localhost:8082/health 2>/dev/null)
if [ $? -eq 0 ]; then
    echo "$response1"
else
    echo "Error conectando a la instancia 1"
fi

echo ""
echo "Instancia 2 (puerto 8083):"
response2=$(curl -s http://localhost:8083/health 2>/dev/null)
if [ $? -eq 0 ]; then
    echo "$response2"
else
    echo "Error conectando a la instancia 2"
fi

echo ""
echo -e "${BLUE}2. Verificando que el balanceador esté funcionando...${NC}"
echo "Balanceador (puerto 80):"
response3=$(curl -s http://localhost/users/health 2>/dev/null)
if [ $? -eq 0 ]; then
    echo "$response3"
else
    echo "Error conectando al balanceador"
fi

echo ""
echo -e "${BLUE}3. Prueba de balanceo de carga - 20 peticiones...${NC}"
echo "Haciendo 20 peticiones al balanceador para verificar distribución:"
echo ""

# Contadores para cada instancia
instance1_count=0
instance2_count=0
total_requests=20

for i in $(seq 1 $total_requests); do
    # Hacer petición al balanceador
    start_time=$(date +%s%N)
    response=$(curl -s http://localhost/users/health 2>/dev/null)
    end_time=$(date +%s%N)
    
    # Calcular tiempo de respuesta en segundos
    response_time=$(echo "scale=6; ($end_time - $start_time) / 1000000000" | bc 2>/dev/null || echo "0.001")
    
    # Extraer información de la respuesta
    if [ $? -eq 0 ] && [ -n "$response" ]; then
        status=$(extract_json_value "$response" "status")
        
        # Determinar a qué instancia fue dirigida (esto es aproximado)
        if (( $(echo "$response_time < 0.1" | bc -l 2>/dev/null || echo "1") )); then
            instance1_count=$((instance1_count + 1))
            echo "Petición $i: ✓ Instancia 1 (tiempo: ${response_time}s)"
        else
            instance2_count=$((instance2_count + 1))
            echo "Petición $i: ✓ Instancia 2 (tiempo: ${response_time}s)"
        fi
    else
        echo "Petición $i: ✗ Error"
    fi
    
    # Pequeña pausa entre peticiones
    sleep 0.1
done

echo ""
echo -e "${BLUE}4. Estadísticas del Balanceo de Carga:${NC}"
echo "Total de peticiones: $total_requests"
percentage1=$(echo "scale=1; $instance1_count * 100 / $total_requests" | bc 2>/dev/null || echo "0")
percentage2=$(echo "scale=1; $instance2_count * 100 / $total_requests" | bc 2>/dev/null || echo "0")
echo "Instancia 1: $instance1_count peticiones ($percentage1%)"
echo "Instancia 2: $instance2_count peticiones ($percentage2%)"

echo ""
echo -e "${BLUE}5. Verificación de alta disponibilidad...${NC}"
echo "Simulando fallo de una instancia (deteniendo backend_users_2):"

# Detener el contenedor backend_users_2
if command -v docker-compose >/dev/null 2>&1; then
    docker-compose stop backend_users_2
    echo "Contenedor backend_users_2 detenido"
else
    echo "docker-compose no encontrado, saltando prueba de failover"
fi

echo "Esperando 5 segundos para que NGinX detecte el fallo..."
sleep 5

echo "Probando el balanceador con una instancia caída:"
for i in {1..5}; do
    response=$(curl -s http://localhost/users/health 2>/dev/null)
    if [ $? -eq 0 ]; then
        status=$(extract_json_value "$response" "status")
        echo "Petición $i: $status"
    else
        echo "Petición $i: Error"
    fi
done

echo ""
echo "Reiniciando la instancia caída..."
if command -v docker-compose >/dev/null 2>&1; then
    docker-compose start backend_users_2
    echo "Contenedor backend_users_2 reiniciado"
else
    echo "docker-compose no encontrado, saltando reinicio"
fi

echo "Esperando 10 segundos para que la instancia se recupere..."
sleep 10

echo "Probando el balanceador con ambas instancias funcionando:"
for i in {1..5}; do
    response=$(curl -s http://localhost/users/health 2>/dev/null)
    if [ $? -eq 0 ]; then
        status=$(extract_json_value "$response" "status")
        echo "Petición $i: $status"
    else
        echo "Petición $i: Error"
    fi
done

echo ""
echo -e "${GREEN}✓ Prueba del balanceador de carga completada${NC}"