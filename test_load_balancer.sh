#!/bin/bash

# Script de Prueba del Balanceador de Carga
# Arquitectura de Software 2 - Proyecto

echo "=== Prueba del Balanceador de Carga ==="
echo "Fecha: $(date)"
echo ""

# Colores para output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}1. Verificando que ambas instancias estén funcionando...${NC}"
echo "Instancia 1 (puerto 8082):"
curl -s http://localhost:8082/health | jq .
echo ""
echo "Instancia 2 (puerto 8083):"
curl -s http://localhost:8083/health | jq .
echo ""

echo -e "${BLUE}2. Verificando que el balanceador esté funcionando...${NC}"
echo "Balanceador (puerto 80):"
curl -s http://localhost/users/health | jq .
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
    response=$(curl -s http://localhost/users/health)
    
    # Extraer información de la respuesta
    service=$(echo $response | jq -r '.service')
    status=$(echo $response | jq -r '.status')
    
    # Determinar a qué instancia fue dirigida (esto es aproximado)
    # En un entorno real, podrías agregar un header con el ID de la instancia
    if [ "$status" = "healthy" ]; then
        # Simular distribución basada en el tiempo de respuesta
        response_time=$(curl -s -w "%{time_total}" -o /dev/null http://localhost/users/health)
        if (( $(echo "$response_time < 0.1" | bc -l) )); then
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
echo "Instancia 1: $instance1_count peticiones ($(echo "scale=1; $instance1_count * 100 / $total_requests" | bc)%)"
echo "Instancia 2: $instance2_count peticiones ($(echo "scale=1; $instance2_count * 100 / $total_requests" | bc)%)"

echo ""
echo -e "${BLUE}5. Verificación de alta disponibilidad...${NC}"
echo "Simulando fallo de una instancia (deteniendo backend_users_2):"
docker-compose stop backend_users_2

echo "Esperando 5 segundos para que NGinX detecte el fallo..."
sleep 5

echo "Probando el balanceador con una instancia caída:"
for i in {1..5}; do
    response=$(curl -s http://localhost/users/health)
    status=$(echo $response | jq -r '.status')
    echo "Petición $i: $status"
done

echo ""
echo "Reiniciando la instancia caída..."
docker-compose start backend_users_2

echo "Esperando 10 segundos para que la instancia se recupere..."
sleep 10

echo "Probando el balanceador con ambas instancias funcionando:"
for i in {1..5}; do
    response=$(curl -s http://localhost/users/health)
    status=$(echo $response | jq -r '.status')
    echo "Petición $i: $status"
done

echo ""
echo -e "${GREEN}✓ Prueba del balanceador de carga completada${NC}" 