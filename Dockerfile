# Etapa de compilacion: descarga dependencias y genera un binario estatico.
FROM golang:1.27 AS compilacion

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o /servidor .

# Etapa de ejecucion: imagen minima y usuario sin privilegios.
FROM gcr.io/distroless/static-debian13:nonroot

COPY --from=compilacion /servidor /servidor

EXPOSE 8080
ENTRYPOINT ["/servidor"]
