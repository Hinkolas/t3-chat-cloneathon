<script lang="ts">
	let { value = $bindable(), onInputFunction = () => {}, placeholder } = $props();

	import { Search } from '@lucide/svelte';

	function handleInput(event: Event) {
		const target = event.target as HTMLInputElement;
		value = target.value;
		onInputFunction?.(event);
	}
</script>

<div class="search-container">
	<Search size="16" />
	<input oninput={handleInput} bind:value {placeholder} type="text" />
</div>

<style>
	.search-container {
		display: flex;
		justify-self: flex-start;
		align-items: center;
		gap: 8px;
		color: var(--sidebar-search-icon);
		padding-inline: 14px;
	}

	.search-container input {
		all: unset;
		background: none;
		padding-block: 8px;
		font-size: 14px;
		color: var(--text);
		width: 100%;
	}
</style>
