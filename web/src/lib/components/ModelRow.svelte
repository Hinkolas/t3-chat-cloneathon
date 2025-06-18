<script lang="ts">
	let { model, changeModel, disabled } = $props();

	import FeatureTag from '$lib/components/FeatureTag.svelte';
	import ModelIcon from '$lib/components/ModelIcon.svelte';
	import { Gem, Key, Sparkles } from '@lucide/svelte';
</script>

<button
	class:disabled
	{disabled}
	onclick={() => {
		changeModel(model);
	}}
	class="model"
>
	<div class="details">
		<div class="provider-icon">
			<ModelIcon {model} size="16" />
		</div>
		<div class="title">{model.title}</div>
		{#each Object.entries(model.flags) as [flag]}
			{#if flag === 'is_premium'}
				<Gem size="16" />
			{/if}
			{#if flag === 'is_key_required'}
				<Key size="16" />
			{/if}
			{#if flag === 'is_new'}
				<Sparkles size="16" color="#F3DF9F"/>
			{/if}
		{/each}
	</div>
	<div class="feature-container">
		{#each Object.entries(model.features) as [feature]}
			<FeatureTag {feature} />
		{/each}
	</div>
</button>

<style>
	.model {
		all: unset;
		display: flex;
		flex-direction: row;
		justify-content: space-between;
		align-items: center;
		gap: 32px;
		padding-inline: 8px;
		border-radius: 4px;
		cursor: pointer;
		color: var(--text);
		transition: background-color 0.15s ease-out;
	}
	.model:not(:disabled):hover {
		background-color: var(--model-hover);
	}

	.model.disabled {
		opacity: 0.5;
		cursor: default;
	}

	.details {
		display: flex;
		flex-direction: row;
		justify-content: flex-start;
		align-items: center;
		gap: 12px;
		padding-block: 8px;
	}

	.model .provider-icon {
		display: flex;
		justify-content: center;
		align-items: center;
		width: 18px;
		height: 18px;
		color: var(--text);
	}

	.model .title {
		font-size: 14px;
		white-space: nowrap;
		color: var(--text);
	}

	.feature-container {
		display: flex;
		align-items: center;
		gap: 8px;
	}
</style>
