# Použití oficiálního obrazu Go
FROM golang:latest

# Nastavení pracovního adresáře
WORKDIR /app

# Kopírování zdrojových souborů do kontejneru
COPY . .

# Kompilace Go kódu
RUN go build -o app

# Příkaz spuštění aplikace po vytvoření kontejneru
CMD ["./app"]
