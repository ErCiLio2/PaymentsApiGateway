<svelte:head>
	<title>Documentação — Payments API Gateway</title>
</svelte:head>

<main>
	<header>
		<div class="logo">Payments API Gateway</div>
		<div class="badge">Docs</div>
	</header>

	<article>
		<h1>Referência da API</h1>
		<p class="lead">Sandbox do Proxypay OPG — simula fluxos de aprovação Multicaixa Express. Não processa pagamentos reais.</p>

		<section>
			<h2>Autenticação</h2>
			<p>O endpoint de criação usa <strong>HTTP Basic Auth</strong>: username = API key, password vazio.</p>
			<pre><code># Gerar header Authorization:
echo -n 'SUA_API_KEY:' | base64 -w0
# Usar: -H "Authorization: Basic &lt;resultado&gt;"</code></pre>
			<p>Os restantes endpoints (<code>GET /api/charges/*</code>, <code>POST /api/charges/*/pay</code>) não requerem autenticação.</p>
		</section>

		<section>
			<h2>Endpoints</h2>

			<div class="endpoint">
				<div class="method post">POST</div>
				<code>/opg/v1/charges</code>
				<span class="auth-badge">Basic Auth</span>
			</div>
			<p>Criar uma nova charge.</p>
			<pre><code>curl -X POST https://API_URL/opg/v1/charges \
  -H "Content-Type: application/json" \
  -H "Authorization: Basic $(printf 'API_KEY:' | base64 -w0)" \
  -d '{
    "pos_id": 1,
    "amount": 5000.00,
    "callback_url": "https://teu-sistema.com/webhook"
  }'</code></pre>
			<p><strong>Resposta 201:</strong></p>
			<pre><code>{
  "id": "abc123...",
  "qrcode_url": "https://API_URL/api/charges/abc123/qrcode.png",
  "deeplink_redirect": "https://WEB_URL/abc123",
  "expires_at": "2026-11-14T18:10:00Z"
}</code></pre>

			<div class="endpoint">
				<div class="method get">GET</div>
				<code>/api/charges/{'{'}id{'}'}</code>
			</div>
			<p>Consultar estado de uma charge. Estados: <code>pending</code> · <code>accepted</code> · <code>rejected</code> · <code>expired</code></p>

			<div class="endpoint">
				<div class="method get">GET</div>
				<code>/api/charges/{'{'}id{'}'}/qrcode.png</code>
			</div>
			<p>QR code em PNG que codifica o link de aprovação. Pode ser embutido directamente em <code>&lt;img src="..."&gt;</code>.</p>

			<div class="endpoint">
				<div class="method post">POST</div>
				<code>/api/charges/{'{'}id{'}'}/pay</code>
			</div>
			<p>Aprovar ou rejeitar uma charge.</p>
			<pre><code>curl -X POST https://API_URL/api/charges/abc123/pay \
  -H "Content-Type: application/json" \
  -d '{"decision": "accepted"}'
# ou: "decision": "rejected"</code></pre>

			<div class="endpoint">
				<div class="method get">GET</div>
				<code>/health</code>
			</div>
			<p>Estado do servidor. Devolve <code>{"{''}status":"ok"{'}'}</code>.</p>
		</section>

		<section>
			<h2>Webhooks</h2>
			<p>Após aprovação/rejeição, o gateway envia um <code>POST</code> para o <code>callback_url</code>.</p>

			<table>
				<thead>
					<tr><th>Header</th><th>Valor</th></tr>
				</thead>
				<tbody>
					<tr><td><code>Content-Type</code></td><td><code>application/json</code></td></tr>
					<tr><td><code>x-proxypay-signature</code></td><td>HMAC-SHA256 hex do body</td></tr>
				</tbody>
			</table>

			<pre><code>{`{
  "id": "txn_xyz789",
  "charge_id": "abc123",
  "status": "accepted",
  "amount": 5000.00,
  "currency": "AOA"
}`}</code></pre>

			<h3>Validar a assinatura</h3>
			<pre><code>// Node.js
const crypto = require('crypto');
const expected = crypto.createHmac('sha256', WEBHOOK_SECRET)
  .update(rawBody).digest('hex');
const valid = crypto.timingSafeEqual(
  Buffer.from(expected), Buffer.from(signature)
);</code></pre>
			<p class="warn">O webhook é enviado <strong>uma única vez</strong>, sem retry. Implementa idempotência no teu endpoint.</p>
		</section>

		<section>
			<h2>Limites</h2>
			<ul>
				<li>Charges expiram ao fim de <strong>30 minutos</strong></li>
				<li>Sem persistência — o histórico perde-se se o Redis reiniciar</li>
				<li>O <code>callback_url</code> deve ser acessível publicamente</li>
			</ul>
		</section>
	</article>
</main>

<style>
	:global(*, *::before, *::after) { box-sizing: border-box; margin: 0; padding: 0; }
	:global(body) {
		font-family: system-ui, -apple-system, sans-serif;
		background: #f4f3ee;
		color: #171717;
		line-height: 1.6;
	}

	main {
		max-width: 760px;
		margin: 0 auto;
		padding: 2rem 1.5rem 4rem;
	}

	header {
		display: flex;
		align-items: center;
		gap: 0.75rem;
		margin-bottom: 2.5rem;
		padding-bottom: 1.25rem;
		border-bottom: 1px solid #e0dfd8;
	}

	.logo { font-size: 1rem; font-weight: 700; }

	.badge {
		font-size: 0.65rem;
		font-weight: 600;
		text-transform: uppercase;
		letter-spacing: 0.05em;
		background: #171717;
		color: #ffe600;
		padding: 0.2rem 0.5rem;
		border-radius: 4px;
	}

	h1 { font-size: 1.75rem; font-weight: 800; margin-bottom: 0.5rem; }

	.lead {
		color: #555;
		margin-bottom: 2.5rem;
		font-size: 1rem;
	}

	section { margin-bottom: 2.5rem; }

	h2 {
		font-size: 1.15rem;
		font-weight: 700;
		margin-bottom: 1rem;
		padding-bottom: 0.4rem;
		border-bottom: 1px solid #e0dfd8;
	}

	h3 { font-size: 0.95rem; font-weight: 700; margin: 1rem 0 0.5rem; }

	p { margin-bottom: 0.75rem; color: #333; }

	pre {
		background: #1e1e1e;
		color: #d4d4d4;
		border-radius: 8px;
		padding: 1rem 1.25rem;
		overflow-x: auto;
		margin: 0.75rem 0 1rem;
		font-size: 0.82rem;
		line-height: 1.6;
	}

	code {
		font-family: 'Menlo', 'Consolas', monospace;
		font-size: 0.85em;
		background: #eeecea;
		padding: 0.1em 0.35em;
		border-radius: 3px;
	}

	pre code { background: none; padding: 0; font-size: inherit; }

	.endpoint {
		display: flex;
		align-items: center;
		gap: 0.75rem;
		margin: 1.25rem 0 0.5rem;
	}

	.method {
		font-size: 0.7rem;
		font-weight: 700;
		padding: 0.2rem 0.5rem;
		border-radius: 4px;
		text-transform: uppercase;
		letter-spacing: 0.04em;
		min-width: 48px;
		text-align: center;
	}

	.method.get { background: #d4edda; color: #155724; }
	.method.post { background: #cce5ff; color: #004085; }

	.auth-badge {
		font-size: 0.7rem;
		background: #fff3cd;
		color: #856404;
		padding: 0.15rem 0.5rem;
		border-radius: 4px;
		font-weight: 600;
	}

	table {
		width: 100%;
		border-collapse: collapse;
		margin: 0.75rem 0 1rem;
		font-size: 0.88rem;
	}

	th, td {
		text-align: left;
		padding: 0.5rem 0.75rem;
		border: 1px solid #e0dfd8;
	}

	th { background: #eeecea; font-weight: 600; }

	ul {
		padding-left: 1.25rem;
		color: #333;
	}

	ul li { margin-bottom: 0.4rem; }

	.warn {
		background: #fff8e1;
		border-left: 3px solid #f59e0b;
		padding: 0.6rem 1rem;
		border-radius: 0 6px 6px 0;
		font-size: 0.88rem;
		color: #333;
	}
</style>
