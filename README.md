# ClubHub Hotel Management

## Descripción
ClubHub Hotel Management es una API de gestión hotelera diseñada para facilitar el manejo de franquicias de hoteles. Permite a los usuarios crear, actualizar, y consultar información detallada sobre diversas franquicias, incluyendo datos de ubicación, información de dominio, y mas.

## Características
- Creación y actualización de franquicias.
- Consultas por ubicación, rango de fechas y nombre.
- Información detallada de cada franquicia, incluyendo datos WHOIS, SSL y más.

## Tecnologías Utilizadas
- Go
- Gin Web Framework
- MongoDB
- Docker y Docker Compose
- Swagger para documentación de la API

## Requisitos
- Go versión 1.18 o superior
- Docker y Docker Compose
- MongoDB

## Instalación y Ejecución

### Clonar el Repositorio
```bash
git clone https://github.com/[tu-usuario]/clubhub-hotel-management.git
cd clubhub-hotel-management
```
## Ejecución con Docker
```bash
docker-compose up --build
```
 - Una vez que los contenedores estén en funcionamiento, la API estará disponible en http://localhost:8080.

## Documentación de la API
Accede a la documentación de la API mediante Swagger en:
http://localhost:8080/swagger/index.html

#### POSTMAN COLLECTION - "HotelMagnament.postman_collection.json"
