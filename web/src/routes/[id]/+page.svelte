<script lang="ts">
	import type { PageData, ActionData } from './$types';

	let { data, form }: { data: PageData; form: ActionData } = $props();

	const amount = $derived(new Intl.NumberFormat('pt-AO', {
		style: 'currency',
		currency: 'AOA',
		minimumFractionDigits: 2,
	}).format(data.charge.amount));

	const expiresAt = $derived(new Date(data.charge.expires_at).toLocaleString('pt-AO'));

	let loading = $state(false);
</script>

<svelte:head>
	<title>Aprovar Pagamento — Proxypay Sandbox</title>
</svelte:head>

<main>
	<header>
		<div class="logo">Multicaixa Express</div>
		<div class="badge">Sandbox</div>
	</header>

	<section class="card">
		{#if data.done || form?.success}
			{@const decision = data.charge.status !== 'pending' ? data.charge.status : form?.decision}
			<div class="result {decision}">
				{#if decision === 'accepted'}
					<span class="icon">✓</span>
					<h2>Pagamento aprovado</h2>
					<p>A transacção foi confirmada com sucesso.</p>
				{:else}
					<span class="icon">✕</span>
					<h2>Pagamento recusado</h2>
					<p>A transacção foi rejeitada.</p>
				{/if}
			</div>
		{:else}
			<h2>Confirmar pagamento</h2>

			<dl class="details">
				<dt>Valor</dt>
				<dd class="amount">{amount}</dd>
				<dt>Referência</dt>
				<dd class="ref">{data.charge.id}</dd>
				<dt>Válido até</dt>
				<dd>{expiresAt}</dd>
			</dl>

			{#if form?.error}
				<p class="error-msg">{form.error}</p>
			{/if}

			<form
				method="POST"
				onsubmit={() => { loading = true; }}
			>
				<button
					type="submit"
					name="decision"
					value="accepted"
					class="btn accept"
					disabled={loading}
				>
					{loading ? 'A processar…' : 'Aprovar'}
				</button>
				<button
					type="submit"
					name="decision"
					value="rejected"
					class="btn reject"
					disabled={loading}
				>
					Recusar
				</button>
			</form>
		{/if}
	</section>

	<footer>
		<p>Esta é uma sandbox de teste. Nenhum dinheiro real é movimentado.</p>
	</footer>
</main>

<style>
	:global(*, *::before, *::after) { box-sizing: border-box; margin: 0; padding: 0; }
	:global(body) {
		font-family: system-ui, -apple-system, sans-serif;
		background: #f4f3ee;
		min-height: 100vh;
		display: flex;
		align-items: center;
		justify-content: center;
	}

	main {
		width: 100%;
		max-width: 420px;
		padding: 1.5rem;
		display: flex;
		flex-direction: column;
		gap: 1.5rem;
	}

	header {
		display: flex;
		align-items: center;
		gap: 0.75rem;
	}

	.logo {
		font-size: 1.1rem;
		font-weight: 700;
		color: #171717;
	}

	.badge {
		font-size: 0.65rem;
		font-weight: 600;
		text-transform: uppercase;
		letter-spacing: 0.05em;
		background: #ffe600;
		color: #171717;
		padding: 0.2rem 0.5rem;
		border-radius: 4px;
	}

	.card {
		background: #fff;
		border-radius: 16px;
		padding: 2rem 1.5rem;
		box-shadow: 0 2px 12px rgba(0,0,0,.08);
	}

	h2 {
		font-size: 1.15rem;
		font-weight: 700;
		color: #171717;
		margin-bottom: 1.25rem;
	}

	dl.details {
		display: grid;
		grid-template-columns: auto 1fr;
		gap: 0.5rem 1rem;
		margin-bottom: 1.75rem;
		font-size: 0.9rem;
	}

	dt { color: #666; }
	dd { color: #171717; font-weight: 500; word-break: break-all; }
	dd.amount { font-size: 1.4rem; font-weight: 800; color: #171717; }
	dd.ref { font-family: monospace; font-size: 0.8rem; }

	.error-msg {
		color: #c0392b;
		font-size: 0.85rem;
		margin-bottom: 1rem;
	}

	form {
		display: flex;
		flex-direction: column;
		gap: 0.75rem;
	}

	.btn {
		width: 100%;
		padding: 0.9rem;
		border: none;
		border-radius: 10px;
		font-size: 1rem;
		font-weight: 600;
		cursor: pointer;
		transition: opacity 0.15s;
	}

	.btn:disabled { opacity: 0.6; cursor: not-allowed; }

	.btn.accept { background: #171717; color: #ffe600; }
	.btn.accept:hover:not(:disabled) { opacity: 0.85; }
	.btn.reject { background: #f0f0f0; color: #171717; }
	.btn.reject:hover:not(:disabled) { background: #e4e4e4; }

	.result {
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 0.75rem;
		padding: 1rem 0;
		text-align: center;
	}

	.result .icon {
		width: 56px;
		height: 56px;
		border-radius: 50%;
		display: flex;
		align-items: center;
		justify-content: center;
		font-size: 1.5rem;
		font-weight: 700;
	}

	.result.accepted .icon { background: #d4f5e2; color: #1a7a46; }
	.result.rejected .icon { background: #fde8e8; color: #c0392b; }

	.result h2 { margin: 0; }
	.result p { color: #666; font-size: 0.9rem; }

	footer { text-align: center; }
	footer p { font-size: 0.75rem; color: #999; }
</style>
