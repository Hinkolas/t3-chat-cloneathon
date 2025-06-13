<script lang="ts">
	let { model, changeModel } = $props();

	import FeatureTag from '$lib/components/FeatureTag.svelte';
	import ModelIcon from '$lib/components/ModelIcon.svelte';
</script>

<button onclick={changeModel} class="model">
	<div class="details">
		<div class="provider-icon">
			<ModelIcon {model} size="16" />
		</div>
		<div class="title">{model.title}</div>
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
		transition: background-color 0.15s ease-out;
	}
	.model:hover {
		background-color: var(--button-hover);
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
		color: hsl(var(--secondary-foreground));
	}

	.model .title {
		font-size: 14px;
		white-space: nowrap;
		color: hsl(var(--secondary-foreground));
	}

	.feature-container {
		display: flex;
		align-items: center;
		gap: 8px;
	}
</style>
