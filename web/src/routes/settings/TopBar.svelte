<script lang="ts">
	import { page } from '$app/stores';

	let tabs = [
		{
			title: 'Account',
			destination: '/settings'
		},
		{
			title: 'Customization',
			destination: '/settings/customization'
		},
		{
			title: 'History & Sync',
			destination: '/settings/history'
		},
		{
			title: 'Models',
			destination: '/settings/models'
		},
		{
			title: 'API Keys',
			destination: '/settings/api-keys'
		},
		{
			title: 'Attachments',
			destination: '/settings/attachments'
		},
		{
			title: 'About Us',
			destination: '/settings/about-us'
		}
	];

	let activeIndex = $state(tabs.findIndex((tab) => $page.url.pathname === tab.destination));

	$effect(() => {
		activeIndex = tabs.findIndex((tab) => $page.url.pathname === tab.destination);
	});
</script>

<div class="tabs">
	{#each tabs as tab, index}
		<a href={tab.destination} class="tab" class:active={activeIndex === index}>
			{tab.title}
		</a>
	{/each}
</div>

<style>
	.tabs {
		position: relative;
		display: flex;
		justify-content: flex-start;
		align-items: center;
		gap: 4px;
		background-color: #312737;
		padding: 4px;
		border-radius: 8px;
		width: max-content;
	}

	.tab {
		all: unset;
		font-size: 12px;
		letter-spacing: 0.4px;
		font-weight: 500;
		padding: 4px 10px;
		border-radius: 6px;
		cursor: pointer;
		white-space: nowrap;
		transition: background-color 0.15s ease-out;
	}

	.tab:hover {
		background-color: #2c232e;
	}

	.tab.active {
		background-color: #21141e;
	}
</style>
