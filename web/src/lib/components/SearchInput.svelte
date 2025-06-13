<script lang="ts">
	let { value = $bindable(), onInputFunction, placeholder } = $props();

	import { Search } from '@lucide/svelte';
</script>

<div class="search-container">
	<Search size="16" />
	<input oninput={onInputFunction} bind:value {placeholder} type="text" />
</div>

<style>
	.search-container {
		display: flex;
		justify-self: flex-start;
		align-items: center;
		gap: 8px;
		color: #888888;
		padding-inline: 14px;
	}

	.search-container input {
		all: unset;
		background: none;
		padding-block: 8px;
		font-size: 14px;
		color: hsl(var(--secondary-foreground));
		width: 100%;
	}
</style>
