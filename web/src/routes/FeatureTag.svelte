<script lang="ts">
	export let tag;

	import { Globe, Eye, Brain, CircleHelp, Icon } from '@lucide/svelte';

	interface Feature {
		backgroundColor: String;
		color: String;
		icon: typeof Icon;
	}

	const features: Record<string, Feature> = {
		vision: {
			backgroundColor: '#262B2B',
			color: '#9AE1D2',
			icon: Eye
		},
		think: {
			backgroundColor: '#24202F',
			color: '#979DF0',
			icon: Brain
		},
		search: {
			backgroundColor: '#24262E',
			color: '#98BFE2',
			icon: Globe
		}
	};
</script>

{#if tag == 'vision' || tag == 'think' || tag == 'search'}
	{@const Icon = features[tag].icon}
	<div
		class="feature"
		style="background-color: {features[tag].backgroundColor}; color: {features[tag].color};"
	>
		<Icon size={16} />
	</div>
{:else}
	<div
		class="feature"
		style="background-color: {features[tag].backgroundColor}; color: {features[tag].color};"
	>
		<CircleHelp size={16} />
	</div>
{/if}

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
