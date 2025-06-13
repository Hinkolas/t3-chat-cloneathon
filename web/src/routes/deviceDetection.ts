import { writable } from 'svelte/store';
import { browser } from '$app/environment';

function createDeviceStore() {
	const checkIfMobile = () => {
		if (!browser) return false;

		const userAgent = navigator.userAgent;

		return /Android|webOS|iPhone|iPad|iPod|BlackBerry|IEMobile|Opera Mini/i.test(userAgent);
	};

	const { subscribe, set } = writable(checkIfMobile()); // Initialize immediately

	return {
		subscribe,
		init: () => set(checkIfMobile()) // Keep this for manual re-checking if needed
	};
}

export const isMobile = createDeviceStore();
