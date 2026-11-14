<script lang="ts">
	import { enhance } from '$app/forms';
	import type { ActionData } from './$types';

	let { form }: { form: ActionData } = $props();
	let loading = $state(false);
</script>

<svelte:head>
	<title>Admin — Payments Gateway</title>
</svelte:head>

<main>
	<div class="card">
		<div class="logo">Payments Gateway</div>
		<h1>Acesso ao backoffice</h1>

		{#if form?.error}
			<p class="error">{form.error}</p>
		{/if}

		<form method="POST" use:enhance={() => {
			loading = true;
			return async ({ update }) => { await update(); loading = false; };
		}}>
			<label for="password">Password de administrador</label>
			<input
				id="password"
				name="password"
				type="password"
				placeholder="••••••••"
				required
				autocomplete="current-password"
			/>
			<button type="submit" disabled={loading}>
				{loading ? 'A entrar…' : 'Entrar'}
			</button>
		</form>
	</div>
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

	main { width: 100%; max-width: 380px; padding: 1.5rem; }

	.card {
		background: #fff;
		border-radius: 16px;
		padding: 2rem 1.75rem;
		box-shadow: 0 2px 12px rgba(0,0,0,.08);
	}

	.logo {
		font-size: 0.8rem;
		font-weight: 700;
		text-transform: uppercase;
		letter-spacing: 0.08em;
		color: #999;
		margin-bottom: 1.25rem;
	}

	h1 { font-size: 1.1rem; font-weight: 700; color: #171717; margin-bottom: 1.5rem; }

	.error {
		background: #fde8e8;
		color: #c0392b;
		font-size: 0.85rem;
		padding: 0.6rem 0.75rem;
		border-radius: 8px;
		margin-bottom: 1rem;
	}

	form { display: flex; flex-direction: column; gap: 0.6rem; }

	label { font-size: 0.82rem; font-weight: 600; color: #555; }

	input {
		padding: 0.75rem 0.9rem;
		border: 1px solid #d8d7d0;
		border-radius: 8px;
		font-size: 0.95rem;
		outline: none;
		transition: border-color 0.15s;
	}

	input:focus { border-color: #171717; }

	button {
		margin-top: 0.5rem;
		padding: 0.85rem;
		background: #171717;
		color: #ffe600;
		border: none;
		border-radius: 10px;
		font-size: 0.95rem;
		font-weight: 600;
		cursor: pointer;
		transition: opacity 0.15s;
	}

	button:disabled { opacity: 0.6; cursor: not-allowed; }
	button:hover:not(:disabled) { opacity: 0.85; }
</style>
