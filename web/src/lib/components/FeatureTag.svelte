<script lang="ts">
	let { feature, wText = false } = $props();

	import { Globe, Eye, Brain, CircleHelp, Icon, Settings2, FileText, Zap, ImagePlus } from '@lucide/svelte';

	interface Feature {
		title: string;
		backgroundColor: String;
		color: String;
		icon: typeof Icon;
	}

	const featuresList: Record<string, Feature> = {
		has_vision: {
			title: 'Vision',
			backgroundColor: '#262B2B',
			color: '#9AE1D2',
			icon: Eye
		},
		has_reasoning: {
			title: 'Reasoning',
			backgroundColor: '#24202F',
			color: '#979DF0',
			icon: Brain
		},
		has_pdf: {
			title: 'PDFs',
			backgroundColor: '#24202F',
			color: '#979DF0',
			icon: FileText
		},
		has_effort_control: {
			title: 'Effort Control',
			backgroundColor: '#2E202C',
			color: '#D798D3',
			icon: Settings2
		},
		has_image_generation: {
			title: 'Image Generation',
			backgroundColor: '#2E1A1A',
			color: '#BA6450',
			icon: ImagePlus
		},
		has_fast: {
			title: 'Fast',
			backgroundColor: '#312B24',
			color: '#F3DFA0',
			icon: Zap
		},
		has_web_search: {
			title: 'Web Search',
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
		style="background-color: {featuresList[feature].backgroundColor}; color: {featuresList[feature]
			.color};"
	>
		<Icon size={15} />
		{#if wText}
			{featuresList[feature].title}
		{/if}
	</div>
{:else}
	<div
		class="feature"
		style="background-color: {featuresList[feature].backgroundColor}; color: {featuresList[feature]
			.color};"
	>
		<CircleHelp size={15} />
	</div>
{/if}

<style>
	.feature {
		height: 24px;
		border-radius: 4px;
		display: flex;
		justify-content: center;
		align-items: center;
		padding: 2px 6px;
		gap: 4px;
		font-size: 12px;
	}
</style>
