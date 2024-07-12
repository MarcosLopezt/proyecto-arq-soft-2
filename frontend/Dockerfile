# Usa una imagen base oficial de Node.js versión 16
FROM node:16 AS build

# Establece el directorio de trabajo en el contenedor
WORKDIR /app

# Copia los archivos package.json y package-lock.json al contenedor
COPY package*.json ./

# Instala las dependencias
RUN npm install --silent

# Copia el resto de los archivos de la aplicación
COPY . .

# Construye la aplicación para producción
RUN npm run build

FROM nginx:alpine

# Copia los archivos construidos desde la imagen de construcción
COPY --from=build /app/build /usr/share/nginx/html

# Expone el puerto en el que Nginx está sirviendo la aplicación
EXPOSE 80

# Comando para correr Nginx en modo foreground
CMD ["nginx", "-g", "daemon off;"]