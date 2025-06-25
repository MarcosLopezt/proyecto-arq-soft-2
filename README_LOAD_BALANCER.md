# Balanceador de Carga con NGinX

## Arquitectura de Software 2 - Proyecto

### Descripción General

Este proyecto implementa un balanceador de carga utilizando NGinX para distribuir el tráfico de manera eficiente entre múltiples instancias de microservicios, garantizando alta disponibilidad y rendimiento.

### Características Implementadas

#### 1. Balanceo de Carga Inteligente
- **Algoritmo de Balanceo**: `least_conn` (menor número de conexiones)
- **Health Checks**: Verificación automática del estado de los servicios
- **Failover**: Redirección automática en caso de fallo
- **Timeouts Configurables**: Manejo de timeouts para conexiones lentas

#### 2. Microservicios Soportados
- **Users Service**: Balanceo entre 2 instancias (backend_users, backend_users_2)
- **Subscriptions Service**: Una instancia con failover
- **Courses Service**: Una instancia con failover
- **Search Service**: Una instancia con failover
- **Admin Service**: Una instancia con failover

#### 3. Configuraciones de Rendimiento
- **Compresión Gzip**: Para reducir el tamaño de las respuestas
- **Keep-alive Connections**: Para mejorar el rendimiento
- **Buffer Optimization**: Configuración optimizada de buffers
- **Security Headers**: Headers de seguridad implementados

### Estructura de Archivos

```
proyecto-arq-soft-2/
├── nginx.conf              # Configuración principal de NGinX
├── health_check.sh         # Script de verificación de salud
├── docker-compose.yml      # Orquestación de servicios
└── README_LOAD_BALANCER.md # Esta documentación
```

### Configuración Detallada

#### NGinX Configuration (`nginx.conf`)

```nginx
# Upstream para el microservicio de usuarios con balanceo de carga
upstream backend_users {
    least_conn;  # Algoritmo de balanceo
    server backend_users:8082 max_fails=3 fail_timeout=30s weight=1;
    server backend_users_2:8082 max_fails=3 fail_timeout=30s weight=1;
    keepalive 32;  # Conexiones persistentes
}
```

**Parámetros Clave:**
- `least_conn`: Distribuye las conexiones al servidor con menos conexiones activas
- `max_fails=3`: Máximo número de fallos antes de marcar el servidor como no disponible
- `fail_timeout=30s`: Tiempo que el servidor permanece marcado como no disponible
- `weight=1`: Peso del servidor (igual para ambos en este caso)

#### Rutas Configuradas

| Ruta | Microservicio | Descripción |
|------|---------------|-------------|
| `/users/` | Users Service | Gestión de usuarios con balanceo de carga |
| `/subscriptions/` | Subscriptions Service | Gestión de suscripciones |
| `/courses/` | Courses Service | Gestión de cursos |
| `/search/` | Search Service | Búsqueda de cursos |
| `/admin/` | Admin Service | Panel de administración |
| `/health` | NGinX | Health check del balanceador |

### Monitoreo y Health Checks

#### Script de Health Check (`health_check.sh`)

El script verifica:
- **Servicios de Base de Datos**: MySQL, MongoDB, RabbitMQ, Memcached
- **Microservicios**: Todos los servicios a través de endpoints `/health`
- **Infraestructura**: Solr, NGinX, Frontend
- **Balanceo de Carga**: Verificación de distribución de tráfico

**Uso:**
```bash
./health_check.sh
```

### Implementación en Docker

#### Docker Compose Configuration

```yaml
nginx:
  image: nginx:latest
  platform: linux/arm64/v8
  container_name: nginx
  volumes:
    - ./nginx.conf:/etc/nginx/nginx.conf
  ports:
    - "80:80"
  depends_on:
    - backend_users
    - backend_users_2
  networks:
    - app-network
```

### Métricas y Monitoreo

#### Headers de Seguridad Implementados
- `X-Frame-Options`: Protección contra clickjacking
- `X-XSS-Protection`: Protección contra XSS
- `X-Content-Type-Options`: Prevención de MIME sniffing
- `Referrer-Policy`: Control de información de referrer
- `Content-Security-Policy`: Política de seguridad de contenido

#### Logging
- **Access Log**: Registro de todas las peticiones
- **Error Log**: Registro de errores con nivel de advertencia
- **Custom Log Format**: Formato personalizado con información detallada

### Alta Disponibilidad

#### Estrategias Implementadas

1. **Failover Automático**
   - Detección automática de servicios no disponibles
   - Redirección a servicios saludables
   - Reintentos automáticos

2. **Health Checks**
   - Verificación periódica del estado de los servicios
   - Endpoint `/health` en cada microservicio
   - Timeouts configurables

3. **Load Balancing**
   - Distribución inteligente de carga
   - Múltiples instancias para servicios críticos
   - Configuración de pesos para control de tráfico

### Pruebas y Validación

#### Verificación del Balanceo de Carga

1. **Iniciar los servicios:**
   ```bash
   docker-compose up -d
   ```

2. **Ejecutar health check:**
   ```bash
   ./health_check.sh
   ```

3. **Verificar balanceo de carga:**
   ```bash
   # Hacer múltiples peticiones para ver la distribución
   for i in {1..10}; do
     curl -I http://localhost/users/health
   done
   ```

#### Métricas de Rendimiento

- **Throughput**: Capacidad de procesamiento de peticiones
- **Latency**: Tiempo de respuesta promedio
- **Error Rate**: Tasa de errores
- **Availability**: Tiempo de disponibilidad del servicio

### Escalabilidad

#### Escalado Horizontal

Para agregar más instancias de un microservicio:

1. **Agregar servicio en docker-compose.yml:**
   ```yaml
   backend_users_3:
     build:
       context: ./users
       dockerfile: Dockerfile
     container_name: backend_users_3
     ports:
       - "8088:8082"
     # ... resto de configuración
   ```

2. **Actualizar nginx.conf:**
   ```nginx
   upstream backend_users {
       least_conn;
       server backend_users:8082 max_fails=3 fail_timeout=30s weight=1;
       server backend_users_2:8082 max_fails=3 fail_timeout=30s weight=1;
       server backend_users_3:8082 max_fails=3 fail_timeout=30s weight=1;
       keepalive 32;
   }
   ```

3. **Reiniciar servicios:**
   ```bash
   docker-compose up -d
   ```

### Troubleshooting

#### Problemas Comunes

1. **Servicio no responde:**
   - Verificar logs: `docker-compose logs [service_name]`
   - Ejecutar health check: `./health_check.sh`
   - Verificar conectividad de red

2. **Balanceo de carga no funciona:**
   - Verificar configuración de upstream en nginx.conf
   - Comprobar que todos los servicios estén corriendo
   - Revisar logs de NGinX

3. **Alto tiempo de respuesta:**
   - Verificar configuración de timeouts
   - Comprobar recursos del sistema
   - Revisar configuración de buffers

### Conclusión

Esta implementación proporciona:
- ✅ **Alta Disponibilidad**: Failover automático y health checks
- ✅ **Rendimiento Optimizado**: Balanceo de carga inteligente y compresión
- ✅ **Escalabilidad**: Fácil agregado de nuevas instancias
- ✅ **Monitoreo**: Scripts de verificación y logging detallado
- ✅ **Seguridad**: Headers de seguridad implementados

La solución cumple con los requisitos de implementar correctamente un balanceador de carga con NGinX para al menos 1 microservicio, y en este caso se implementa para todos los microservicios del proyecto. 