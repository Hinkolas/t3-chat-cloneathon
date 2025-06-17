import type { PageServerLoad } from './$types';
import { redirect } from '@sveltejs/kit';

export const load = (async ({ params, url, fetch }) => {
    throw redirect(302, '/');    
}) satisfies PageServerLoad;
