# Étape 1 : Build de l'application Go
FROM golang:1.23-bookworm AS builder

WORKDIR /app

# Copier les fichiers go.mod et go.sum pour télécharger les dépendances
COPY go.* ./
RUN go mod download

# Copier tout le code source dans l'image
COPY . ./

# Build l'application
RUN go build -v -o server

# Étape 2 : Image de production
FROM debian:bookworm-slim

# Installer les certificats CA pour les connexions sécurisées
RUN set -x && apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y \
    ca-certificates && \
    rm -rf /var/lib/apt/lists/*

# Copier le binaire de l'application dans l'image finale
COPY --from=builder /app/server /app/server
# Définir le répertoire de travail
WORKDIR /app

# Exposer le port que Cloud Run attend (8080)
EXPOSE 8080

# Lancer le binaire au démarrage
CMD ["/app/server"]
