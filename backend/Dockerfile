# syntax=docker/dockerfile:1
# ============================================================================
# Dockerfile do backend (build multi-stage)
# ----------------------------------------------------------------------------
# Didatico: um build "multi-stage" usa varias imagens. Na 1a (builder) temos o
# compilador Go; na final copiamos APENAS o binario para uma imagem minima.
# Resultado: imagem pequena e com muito menos superficie de ataque.
# ============================================================================

# --- Stage 1: builder -------------------------------------------------------
FROM golang:1.25 AS builder

WORKDIR /src

# Copiamos primeiro go.mod/go.sum e fazemos download das dependencias. Como
# estes ficheiros mudam pouco, o Docker reaproveita esta layer em cache e os
# rebuilds ficam muito mais rapidos (so re-baixa deps se go.mod/go.sum mudarem).
COPY go.mod go.sum ./
RUN go mod download

# Agora sim, o codigo-fonte.
COPY . .

# CGO_ENABLED=0 produz um binario estatico (sem dependencias de libc do SO),
# que corre numa imagem "scratch"/distroless sem bibliotecas instaladas.
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /out/server ./cmd/server

# --- Stage 2: imagem final --------------------------------------------------
# distroless/static: so o necessario para correr um binario estatico. Nao tem
# shell nem gestor de pacotes -> dificil de explorar se for comprometida.
FROM gcr.io/distroless/static-debian12:nonroot

# Corre como utilizador sem privilegios (nonroot), boa pratica de seguranca.
USER nonroot:nonroot

COPY --from=builder /out/server /server

EXPOSE 8080
ENTRYPOINT ["/server"]
