import { redirect, fail } from '@sveltejs/kit';
import type { Actions, PageServerLoad } from './$types';

const API_URL = process.env.PROXYPAY_API_URL ?? 'http://localhost:8080';
const ADMIN_KEY = process.env.ADMIN_KEY ?? 'admin';

export const load: PageServerLoad = async ({ cookies, fetch }) => {
	if (!cookies.get('admin_session')) redirect(302, '/admin/login');

	const res = await fetch(`${API_URL}/api/admin/charges`, {
		headers: { Authorization: `Bearer ${ADMIN_KEY}` },
	});

	if (!res.ok) {
		return { charges: [] };
	}

	const charges = await res.json();
	return { charges };
};

export const actions: Actions = {
	pay: async ({ request, cookies, fetch }) => {
		if (!cookies.get('admin_session')) redirect(302, '/admin/login');

		const data = await request.formData();
		const id = data.get('id')?.toString() ?? '';
		const decision = data.get('decision')?.toString() ?? '';

		if (!id || (decision !== 'accepted' && decision !== 'rejected')) {
			return fail(400, { error: 'Pedido inválido.' });
		}

		const res = await fetch(`${API_URL}/api/charges/${id}/pay`, {
			method: 'POST',
			headers: { 'Content-Type': 'application/json' },
			body: JSON.stringify({ decision }),
		});

		if (!res.ok) {
			const body = await res.json().catch(() => ({}));
			return fail(res.status, { error: body.error ?? 'Erro ao processar.' });
		}

		return { ok: true };
	},

	logout: async ({ cookies }) => {
		cookies.delete('admin_session', { path: '/' });
		redirect(302, '/admin/login');
	},
};
