import { error, fail } from '@sveltejs/kit';
import type { Actions, PageServerLoad } from './$types';

interface Charge {
	id: string;
	amount: number;
	status: 'pending' | 'accepted' | 'rejected' | 'expired';
	expires_at: string;
}

export const load: PageServerLoad = async ({ params }) => {
	const apiUrl = process.env.PROXYPAY_API_URL ?? 'http://localhost:8080';
	const res = await fetch(`${apiUrl}/api/charges/${params.id}`);

	if (res.status === 404) error(404, 'Cobrança não encontrada ou expirada.');
	if (!res.ok) error(502, 'Serviço de pagamento indisponível.');

	const charge: Charge = await res.json();

	if (charge.status === 'expired') error(410, 'Esta cobrança já expirou.');
	if (charge.status !== 'pending') {
		return { charge, done: true };
	}

	return { charge, done: false };
};

export const actions: Actions = {
	default: async ({ params, request }) => {
		const apiUrl = process.env.PROXYPAY_API_URL ?? 'http://localhost:8080';
		const data = await request.formData();
		const decision = data.get('decision');

		if (decision !== 'accepted' && decision !== 'rejected') {
			return fail(400, { error: 'Decisão inválida.' });
		}

		const res = await fetch(`${apiUrl}/api/charges/${params.id}/pay`, {
			method: 'POST',
			headers: { 'Content-Type': 'application/json' },
			body: JSON.stringify({ decision }),
		});

		if (!res.ok) {
			const body = await res.json().catch(() => ({}));
			return fail(res.status, { error: (body as { error?: string }).error ?? 'Erro desconhecido.' });
		}

		return { success: true, decision };
	},
};
