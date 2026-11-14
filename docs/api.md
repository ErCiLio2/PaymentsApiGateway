# Payments API Gateway — Referência da API

Sandbox do Proxypay OPG. Simula fluxos de aprovação de pagamentos Multicaixa Express. Não processa pagamentos reais.

## Autenticação

O endpoint de criação usa **HTTP Basic Auth**: username = `PROXYPAY_API_KEY`, password vazio.

```bash
echo -n 'SUA_API_KEY:' | base64 -w0
# Usar: -H "Authorization: Basic <resultado>"
```

Os restantes endpoints não requerem autenticação.

---

## Endpoints

### POST /opg/v1/charges — Criar charge

```bash
curl -X POST https://API_URL/opg/v1/charges \
  -H "Content-Type: application/json" \
  -H "Authorization: Basic $(printf 'API_KEY:' | base64 -w0)" \
  -d '{
    "pos_id": 1,
    "amount": 5000.00,
    "callback_url": "https://teu-sistema.com/webhook"
  }'
```

**Resposta 201:**
```json
{
  "id": "abc123...",
  "qrcode_url": "https://API_URL/api/charges/abc123/qrcode.png",
  "deeplink_redirect": "https://WEB_URL/abc123",
  "expires_at": "2026-11-14T18:10:00Z"
}
```

---

### GET /api/charges/{id} — Consultar charge

```bash
curl https://API_URL/api/charges/abc123
```

Estados possíveis: `pending` · `accepted` · `rejected` · `expired`

---

### GET /api/charges/{id}/qrcode.png — QR code

Devolve PNG. Pode ser embutido directamente em `<img src="...">`.

---

### POST /api/charges/{id}/pay — Aprovar ou rejeitar

```bash
curl -X POST https://API_URL/api/charges/abc123/pay \
  -H "Content-Type: application/json" \
  -d '{"decision": "accepted"}'
# ou: "decision": "rejected"
```

---

### GET /health

```bash
curl https://API_URL/health
# {"status":"ok"}
```

---

## Webhooks

Após aprovação/rejeição, o gateway envia um `POST` para o `callback_url`.

| Header | Valor |
|--------|-------|
| `Content-Type` | `application/json` |
| `x-proxypay-signature` | HMAC-SHA256 hex do body |

**Body:**
```json
{
  "id": "txn_xyz789",
  "charge_id": "abc123",
  "status": "accepted",
  "amount": 5000.00,
  "currency": "AOA"
}
```

### Validar assinatura

```js
// Node.js
const crypto = require('crypto');
const expected = crypto.createHmac('sha256', WEBHOOK_SECRET)
  .update(rawBody).digest('hex');
const valid = crypto.timingSafeEqual(
  Buffer.from(expected), Buffer.from(signature)
);
```

```python
# Python
import hmac, hashlib
expected = hmac.new(secret.encode(), body, hashlib.sha256).hexdigest()
valid = hmac.compare_digest(expected, signature)
```

> O webhook é enviado **uma única vez**, sem retry.

---

## Limites

- Charges expiram ao fim de **30 minutos**
- Sem persistência longa — Redis reiniciado apaga o histórico
- `callback_url` deve ser acessível publicamente
