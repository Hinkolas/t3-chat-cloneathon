<script lang="ts">
	import { ArrowLeft } from '@lucide/svelte';
	import SettingsSidebar from './SettingsSidebar.svelte';
	import TopBar from './TopBar.svelte';
	import type { ProfileResponse } from '$lib/types';

	interface Props {
		data: {
			profile: ProfileResponse;
		};
	}
	let { data } = $props();

	let profile = $derived(data.profile);
</script>

<div class="background-overlay">
	<div class="container">
		<div class="head">
			<a href="/">
				<ArrowLeft size="16" />
				Back to Chat
			</a>

			<div class="buttons">
				<a href="/auth/logout">Sign out</a>
			</div>
		</div>
		<div class="body">
			<SettingsSidebar {profile} />
			<div class="content">
				<TopBar />
				<slot />
			</div>
		</div>
	</div>
</div>

<style>
	.background-overlay {
		width: 100%;
		display: flex;
		justify-content: center;
		background-color: var(--settings-background);
	}
	.container {
		width: 100%;
		max-width: 1200px;
		height: 100dvh;
		padding: 24px 32px 96px 32px;
		overflow: hidden;
		color: white;
		display: flex;
		flex-direction: column;
		gap: 24px;
	}

	@media (max-width: 1200px) {
		.container {
			padding-inline: 16px;
		}
	}

	.head {
		width: 100%;

		display: flex;
		justify-content: space-between;
		align-items: center;
	}

	.head a {
		all: unset;
		display: flex;
		justify-content: center;
		align-items: center;
		gap: 12px;
		font-size: 14px;
		font-weight: 600;
		padding: 8px 16px;
		border-radius: 8px;
		cursor: pointer;
		transition: background-color 0.15s ease-out;
	}

	.head a:hover {
		background-color: #88888822;
		border-radius: 8px;
	}

	.body {
		display: flex;
		gap: 48px;
		flex: 1;
		min-height: 0;
	}

	.content {
		display: flex;
		flex-direction: column;
		gap: 32px;
		flex: 1;
		min-width: 0;
	}
</style>
