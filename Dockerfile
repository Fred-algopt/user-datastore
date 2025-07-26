# Étape 1 : build de l'application Go
FROM golang:1.21 as builder

WORKDIR /app

# Copier les fichiers go
COPY go.mod ./
COPY go.sum ./
RUN go mod download

# Copier le code source
COPY . .

# Compiler l'application
RUN go build -v -o server

# Étape 2 : image finale minimaliste
FROM gcr.io/distroless/base-debian11

WORKDIR /app

# Copier le binaire depuis l'étape précédente
COPY --from=builder /app/server .

# Port utilisé par l'application (Cloud Run détecte automatiquement si besoin)
EXPOSE 8080

# Commande de démarrage
CMD ["/app/server"]
