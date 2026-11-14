---
name: verify
description: Verifica o fluxo completo do payments gateway — health check, criar charge, aprovar/rejeitar e confirmar estado final. Usar após uma mudança para confirmar que tudo funciona end-to-end.
---

Verifica que o payments gateway está a funcionar correctamente:

1. Health check:
```bash
curl -sf http://localhost:8080/health && echo "✓ Backend OK" || echo "✗ Backend down"
```

2. Criar charge de teste:
```bash
CHARGE=$(curl -sf -X POST http://localhost:8080/opg/v1/charges \
  -H "Content-Type: application/json" \
  -H "Authorization: Basic $(printf 'dev-proxypay-key:' | base64 -w0)" \
  -d '{"pos_id":1,"amount":1000.00,"callback_url":"http://localhost:9999/webhook"}')
echo "Charge criada: $CHARGE"
CHARGE_ID=$(echo "$CHARGE" | grep -o '"id":"[^"]*"' | cut -d'"' -f4)
echo "ID: $CHARGE_ID"
```

3. Verificar detalhes da charge:
```bash
curl -sf "http://localhost:8080/api/charges/$CHARGE_ID"
```

4. Aprovar a charge:
```bash
curl -sf -X POST "http://localhost:8080/api/charges/$CHARGE_ID/pay" \
  -H "Content-Type: application/json" \
  -d '{"decision":"accepted"}'
```

5. Confirmar estado `accepted`:
```bash
curl -sf "http://localhost:8080/api/charges/$CHARGE_ID"
```

Reportar qualquer passo que falhe com o output de erro completo.
