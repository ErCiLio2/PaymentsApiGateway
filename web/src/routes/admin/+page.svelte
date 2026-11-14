<script lang="ts">
	import { enhance } from '$app/forms';
	import { invalidateAll } from '$app/navigation';
	import type { PageData, ActionData } from './$types';

	let { data, form }: { data: PageData; form: ActionData } = $props();

	let loadingId = $state<string | null>(null);

	$effect(() => {
		const t = setInterval(() => invalidateAll(), 10_000);
		return () => clearInterval(t);
	});

	const statusLabel: Record<string, string> = {
		pending: 'Pendente',
		accepted: 'Aprovado',
		rejected: 'Rejeitado',
		expired: 'Expirado',
	};

	function formatAmount(n: number) {
		return new Intl.NumberFormat('pt-AO', { style: 'currency', currency: 'AOA' }).format(n);
	}

	function formatDate(s: string) {
		return new Date(s).toLocaleString('pt-PT', {
			day: '2-digit', month: '2-digit', year: 'numeric',
			hour: '2-digit', minute: '2-digit',
		});
	}
</script>

<svelte:head>
	<title>Backoffice — Payments Gateway</title>
</svelte:head>

<header>
	<span class="brand">Payments Gateway</span>
	<form method="POST" action="?/logout" use:enhance>
		<button type="submit" class="logout">Sair</button>
	</form>
</header>

<main>
	<h1>Charges</h1>

	{#if form?.error}
		<p class="error">{form.error}</p>
	{/if}

	{#if data.charges.length === 0}
		<p class="empty">Nenhuma charge registada.</p>
	{:else}
		<div class="table-wrap">
			<table>
				<thead>
					<tr>
						<th>ID</th>
						<th>Montante</th>
						<th>Estado</th>
						<th>Criado em</th>
						<th>Expira em</th>
						<th>Ações</th>
					</tr>
				</thead>
				<tbody>
					{#each data.charges as c (c.id)}
						<tr>
							<td class="id" title={c.id}>{c.id.slice(0, 8)}…</td>
							<td class="amount">{formatAmount(c.amount)}</td>
							<td><span class="badge {c.status}">{statusLabel[c.status] ?? c.status}</span></td>
							<td>{formatDate(c.created_at)}</td>
							<td>{formatDate(c.expires_at)}</td>
							<td class="actions">
								{#if c.status === 'pending'}
									<form method="POST" action="?/pay" use:enhance={() => {
										loadingId = c.id;
										return async ({ update }) => { await update(); loadingId = null; };
									}}>
										<input type="hidden" name="id" value={c.id} />
										<input type="hidden" name="decision" value="accepted" />
										<button type="submit" class="btn accept" disabled={loadingId === c.id}>
											{loadingId === c.id ? '…' : 'Aprovar'}
										</button>
									</form>
									<form method="POST" action="?/pay" use:enhance={() => {
										loadingId = c.id;
										return async ({ update }) => { await update(); loadingId = null; };
									}}>
										<input type="hidden" name="id" value={c.id} />
										<input type="hidden" name="decision" value="rejected" />
										<button type="submit" class="btn reject" disabled={loadingId === c.id}>
											{loadingId === c.id ? '…' : 'Rejeitar'}
										</button>
									</form>
								{:else}
									<span class="no-action">—</span>
								{/if}
							</td>
						</tr>
					{/each}
				</tbody>
			</table>
		</div>
	{/if}
</main>

<style>
	:global(*, *::before, *::after) { box-sizing: border-box; margin: 0; padding: 0; }
	:global(body) {
		font-family: system-ui, -apple-system, sans-serif;
		background: #f4f3ee;
		min-height: 100vh;
		color: #171717;
	}

	header {
		background: #171717;
		color: #ffe600;
		padding: 0.9rem 1.5rem;
		display: flex;
		align-items: center;
		justify-content: space-between;
	}

	.brand { font-size: 0.8rem; font-weight: 700; text-transform: uppercase; letter-spacing: 0.08em; }

	.logout {
		background: transparent;
		border: 1px solid #ffe600;
		color: #ffe600;
		padding: 0.35rem 0.9rem;
		border-radius: 6px;
		font-size: 0.82rem;
		font-weight: 600;
		cursor: pointer;
	}
	.logout:hover { background: rgba(255,230,0,.12); }

	main { max-width: 1100px; margin: 0 auto; padding: 2rem 1.5rem; }

	h1 { font-size: 1.3rem; font-weight: 700; margin-bottom: 1.25rem; }

	.error {
		background: #fde8e8;
		color: #c0392b;
		font-size: 0.85rem;
		padding: 0.6rem 0.75rem;
		border-radius: 8px;
		margin-bottom: 1rem;
	}

	.empty { color: #888; font-size: 0.9rem; }

	.table-wrap { overflow-x: auto; border-radius: 12px; box-shadow: 0 2px 12px rgba(0,0,0,.08); }

	table { width: 100%; border-collapse: collapse; background: #fff; }

	thead th {
		text-align: left;
		padding: 0.75rem 1rem;
		font-size: 0.78rem;
		font-weight: 700;
		text-transform: uppercase;
		letter-spacing: 0.05em;
		color: #666;
		border-bottom: 1px solid #e8e7e0;
	}

	tbody tr { border-bottom: 1px solid #f0efe8; }
	tbody tr:last-child { border-bottom: none; }
	tbody tr:hover { background: #fafaf7; }

	td { padding: 0.75rem 1rem; font-size: 0.88rem; vertical-align: middle; }

	.id { font-family: monospace; color: #555; font-size: 0.82rem; }
	.amount { font-weight: 600; }

	.badge {
		display: inline-block;
		padding: 0.25rem 0.65rem;
		border-radius: 999px;
		font-size: 0.75rem;
		font-weight: 600;
	}
	.badge.pending  { background: #fff3cd; color: #856404; }
	.badge.accepted { background: #d1f7c4; color: #1a7a2e; }
	.badge.rejected { background: #fde8e8; color: #c0392b; }
	.badge.expired  { background: #e8e7e0; color: #666; }

	.actions { display: flex; gap: 0.4rem; }

	.btn {
		padding: 0.35rem 0.75rem;
		border: none;
		border-radius: 6px;
		font-size: 0.8rem;
		font-weight: 600;
		cursor: pointer;
		white-space: nowrap;
	}
	.btn:disabled { opacity: 0.5; cursor: not-allowed; }
	.btn.accept { background: #d1f7c4; color: #1a7a2e; }
	.btn.accept:hover:not(:disabled) { background: #b8f0a8; }
	.btn.reject { background: #fde8e8; color: #c0392b; }
	.btn.reject:hover:not(:disabled) { background: #f8d0d0; }

	.no-action { color: #bbb; }
</style>
