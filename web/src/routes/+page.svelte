<svelte:head>
	<title>Proxypay Sandbox</title>
</svelte:head>

<main>
	<h1>Proxypay Sandbox</h1>
	<p>Acede a <code>/&#123;charge-id&#125;</code> para aprovar ou recusar um pagamento.</p>
</main>

<style>
	:global(body) {
		font-family: system-ui, sans-serif;
		display: flex;
		align-items: center;
		justify-content: center;
		min-height: 100vh;
		background: #f4f3ee;
	}
	main { text-align: center; }
	h1 { font-size: 1.5rem; margin-bottom: 0.5rem; }
	p { color: #666; }
	code { background: #fff; padding: 0.2em 0.4em; border-radius: 4px; }
</style>
