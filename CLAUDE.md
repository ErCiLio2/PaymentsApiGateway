# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Projecto

Sandbox simulator do Proxypay OPG — simula fluxos de aprovação de pagamentos Multicaixa Express. É um ambiente de testes; não é produção.

## Stack

- Backend: Go 1.24.3, stdlib `net/http` (sem framework), Redis 7
- Frontend: SvelteKit 3 + Svelte 5 (runes mode obrigatório), TypeScript strict, pnpm
- DB: Redis — charges expiram ao fim de 30 min (TTL). Sem persistência longa.
- Runtime: Docker + Docker Compose

## Arranque local

```bash
cp .env.example .env
# Gerar e preencher PROXYPAY_WEBHOOK_SECRET:
openssl rand -hex 32
docker-compose up
```

Backend: `http://localhost:8080` · Frontend: `http://localhost:3001`

## Variáveis de ambiente obrigatórias

| Var | Notas |
|-----|-------|
| `PROXYPAY_WEBHOOK_SECRET` | ≥ 32 chars; o servidor termina na inicialização se faltar ou for curta |
| `WEB_URL` | URL pública do frontend — usada no deeplink do QR code e no CORS |

## Backend (`services/proxypay/`)

```bash
# Run sem Docker:
cd services/proxypay && go run ./cmd/server

# Tests (ainda a criar):
go test ./...
```

Rotas registadas em `cmd/server/main.go` via stdlib `http.ServeMux`:

| Método | Path | Auth |
|--------|------|------|
| POST | `/opg/v1/charges` | Basic Auth — username = `PROXYPAY_API_KEY`, password vazio |
| GET | `/api/charges/{id}` | — |
| GET | `/api/charges/{id}/qrcode.png` | — |
| POST | `/api/charges/{id}/pay` | — |
| GET | `/health` | — |

**Webhooks:** disparados em goroutine (assíncrono), sem retry. Assinados com HMAC-SHA256 no header `x-proxypay-signature`.

## Frontend (`web/`)

```bash
pnpm check        # type check (svelte-check + svelte-kit sync)
```

**Svelte 5 runes** estão forçados via `vite.config.ts`. Usar sempre `$props()`, `$state()`, `$derived()` — nunca a API de componentes Svelte 4 (stores, `export let`, eventos `on:`).

O frontend usa `@sveltejs/adapter-node` — produz um servidor Node.js, não HTML estático. Env var obrigatória: `PROXYPAY_API_URL`.

## Git e PRs

Feature branches + PR para `main`. Nunca incluir assinaturas de AI em commits ou PRs.
