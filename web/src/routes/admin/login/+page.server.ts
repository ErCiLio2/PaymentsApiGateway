import { redirect, fail } from '@sveltejs/kit';
import type { Actions, PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ cookies }) => {
	if (cookies.get('admin_session')) redirect(302, '/admin');
	return {};
};

export const actions: Actions = {
	default: async ({ request, cookies }) => {
		const data = await request.formData();
		const password = data.get('password')?.toString() ?? '';

		const adminPassword = process.env.ADMIN_PASSWORD ?? 'admin';
		if (password !== adminPassword) {
			return fail(401, { error: 'Password incorrecta.' });
		}

		cookies.set('admin_session', '1', {
			path: '/',
			httpOnly: true,
			sameSite: 'lax',
			maxAge: 60 * 60 * 8, // 8 horas
		});

		redirect(302, '/admin');
	},
};
