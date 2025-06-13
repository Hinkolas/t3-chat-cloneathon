<script lang="ts">
	let { feature } = $props();

	import { Globe, Eye, Brain, CircleHelp, Icon } from '@lucide/svelte';

	interface Feature {
		backgroundColor: String;
		color: String;
		icon: typeof Icon;
	}

	const featuresList: Record<string, Feature> = {
		has_vision: {
			backgroundColor: '#262B2B',
			color: '#9AE1D2',
			icon: Eye
		},
		has_reasoning: {
			backgroundColor: '#24202F',
			color: '#979DF0',
			icon: Brain
		},
		has_pdf: {
			backgroundColor: '#24262E',
			color: '#98BFE2',
			icon: Globe
		},
		has_effort_control: {
			backgroundColor: '#24262E',
			color: '#98BFE2',
			icon: Globe
		},
		has_image_generation: {
			backgroundColor: '#24262E',
			color: '#98BFE2',
			icon: Globe
		},
		has_fast: {
			backgroundColor: '#24262E',
			color: '#98BFE2',
			icon: Globe
		},
		has_web_search: {
			backgroundColor: '#24262E',
			color: '#98BFE2',
			icon: Globe
		}
	};
</script>

{#if feature}
	{@const Icon = featuresList[feature].icon}
	<div
		class="feature"
		style="background-color: {featuresList[feature].backgroundColor}; color: {featuresList[feature].color};"
	>
		<Icon size={16} />
	</div>
{:else}
	<div
		class="feature"
		style="background-color: {featuresList[feature].backgroundColor}; color: {featuresList[feature].color};"
	>
		<CircleHelp size={16} />
	</div>
{/if}

<!-- {#if features == 'vision' || features == 'think' || features == 'search'}
	{@const Icon = featuresList[features].icon}
	<div
		class="feature"
		style="background-color: {featuresList[features].backgroundColor}; color: {featuresList[features].color};"
	>
		<Icon size={16} />
	</div>
{:else}
	<div
		class="feature"
		style="background-color: {featuresList[features].backgroundColor}; color: {featuresList[features].color};"
	>
		<CircleHelp size={16} />
	</div>
{/if} -->

<style>
	.feature {
		width: 24px;
		height: 24px;
		border-radius: 4px;
		display: flex;
		justify-content: center;
		align-items: center;
	}
</style>
