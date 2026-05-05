.PHONY: help install build run dev up down logs clean docker-build docker-up docker-down docker-logs test

BINARY_NAME=server
GO=go
DOCKER_COMPOSE=docker-compose
PORT=8080
DB_HOST=localhost
DB_PORT=5432

help:
	@echo "Comandos disponibles:"
	@echo ""
	@echo "📦 Desarrollo LOCAL (sin Docker):"
	@echo "  make install        - Instalar dependencias Go"
	@echo "  make build          - Compilar la aplicación"
	@echo "  make run            - Ejecutar la aplicación compilada"
	@echo "  make dev            - Ejecutar en modo desarrollo (con hot reload)"
	@echo ""
	@echo "🐳 Docker:"
	@echo "  make docker-build   - Construir imagen Docker"
	@echo "  make docker-up      - Levantar contenedores (app + postgres)"
	@echo "  make docker-down    - Detener y eliminar contenedores"
	@echo "  make docker-logs    - Ver logs de los contenedores"
	@echo "  make docker-clean   - Eliminar volúmenes de Docker"
	@echo ""
	@echo "🧪 Tests:"
	@echo "  make test           - Ejecutar tests"
	@echo ""
	@echo "🔧 Utilidades:"
	@echo "  make clean          - Limpiar archivos compilados"
	@echo "  make env-setup      - Crear archivo .env desde .env.example"
	@echo ""

env-setup:
	@if [ ! -f .env ]; then \
		cp .env.example .env; \
		echo "✅ Archivo .env creado desde .env.example"; \
		echo "⚠️  Recuerda actualizar las credenciales en .env si es necesario"; \
	else \
		echo "✅ El archivo .env ya existe"; \
	fi

install:
	@echo "📥 Instalando dependencias Go..."
	$(GO) mod download
	$(GO) mod verify
	@echo "✅ Dependencias instaladas"

build: install
	@echo "🔨 Compilando aplicación..."
	$(GO) build -o $(BINARY_NAME) ./cmd/server
	@echo "✅ Compilación completada: ./$(BINARY_NAME)"

run: build
	@echo "▶️  Ejecutando aplicación..."
	@echo "La aplicación estará disponible en http://localhost:$(PORT)"
	./$(BINARY_NAME)

dev: env-setup
	@echo "🚀 Ejecutando en modo desarrollo..."
	@echo "Asegúrate de tener PostgreSQL ejecutándose en $(DB_HOST):$(DB_PORT)"
	@command -v air > /dev/null 2>&1 || { \
		echo "📦 Instalando air (hot reload)..."; \
		go install github.com/cosmtrek/air@latest; \
	}
	air

# Docker targets
docker-build:
	@echo "🐳 Construyendo imagen Docker..."
	$(DOCKER_COMPOSE) build
	@echo "✅ Imagen construida"

docker-up: env-setup
	@echo "🚀 Levantando contenedores..."
	@echo "App disponible en http://localhost:$(PORT)"
	@echo "Base de datos en localhost:$(DB_PORT)"
	$(DOCKER_COMPOSE) up -d
	@echo "✅ Contenedores levantados"
	@echo "💡 Usa 'make docker-logs' para ver los logs"

docker-down:
	@echo "⏹️  Deteniendo contenedores..."
	$(DOCKER_COMPOSE) down
	@echo "✅ Contenedores detenidos"

docker-logs:
	$(DOCKER_COMPOSE) logs -f

docker-clean: docker-down
	@echo "🗑️  Eliminando volúmenes..."
	$(DOCKER_COMPOSE) down -v
	@echo "✅ Volúmenes eliminados"

# Test targets
test:
	@echo "🧪 Ejecutando tests..."
	$(GO) test -v -cover ./...

# Cleanup
clean:
	@echo "🧹 Limpiando archivos compilados..."
	@rm -f $(BINARY_NAME)
	$(GO) clean
	@echo "✅ Limpieza completada"

# Información del proyecto
info:
	@echo "📋 Información del proyecto:"
	@echo "Nombre: reservar-club-back"
	@echo "Lenguaje: Go 1.25"
	@echo "Framework: Gin"
	@echo "Base de datos: PostgreSQL 16"
	@echo "Puerto: $(PORT)"
	@echo ""
	@echo "Variables de entorno disponibles en: .env"
