import type { PageServerLoad } from './$types';
import type { ModelsResponse } from './types';

export const load = (async () => {
	const url = 'http://192.168.69.178:3141';

	try {
		const response = await fetch(`${url}/v1/models/`);

		if (!response.ok) {
			throw new Error('Something happened during Record');
		}

		const models: ModelsResponse = await response.json();

		return { models };
	} catch (error) {
		console.log(error);
	}

	return {};
}) satisfies PageServerLoad;
