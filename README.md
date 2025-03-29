# Lab-6-Backend-only---Parte-2

# ⚽ LaLigaTracker Backend API

## Descripción
API backend para gestionar partidos de fútbol de La Liga Española con:
- Registro de estadísticas 
- Documentación Swagger integrada
- Configuración optimizada para Docker

## Instalación Rápida

### Requisitos
- Docker 28.01
- Docker Compose 2.33.1

### Pasos
```bash
git clone https://github.com/SebastianUVG/Lab-6-Backend-only---Parte-1.git
cd laliga-tracker
docker-compose up --build
```

### colección de postman
https://sebastian-1545329.postman.co/workspace/My-Workspace~86169d7e-33f9-4de6-89f3-3b4d7c4ad9e3/collection/43574994-19b30329-dafd-487b-ad4d-49fd19fc138f

## 🔌 Endpoints de la API
```http
GET /api/matches
GET /api/matches/{id}
POST /api/matches
PUT /api/matches/{id}
DELETE /api/matches/{id}
PATCH /api/matches/{id}/goals
PATCH /api/matches/{id}/yellowcards
PATCH /api/matches/{id}/redcards
PATCH /api/matches/{id}/extratime
```

### Imagenes de la primera parte
![image](https://github.com/user-attachments/assets/2eb1935d-0d17-4d0d-8ea4-c214f3ef6eb5)

![image](https://github.com/user-attachments/assets/69bbd271-f9e0-4358-857d-aded34dd58d2)
