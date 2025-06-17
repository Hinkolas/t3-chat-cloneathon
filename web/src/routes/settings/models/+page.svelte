<script lang="ts">
	import type { ModelsResponse, ModelData, ModelFeatures } from '$lib/types.js';

	interface Props {
		data: {
			models: ModelsResponse;
		};
	}

	let { data }: Props = $props();

	import { Check, Link } from '@lucide/svelte';
	import FeatureTag from '$lib/components/FeatureTag.svelte';
	import ModelIcon from '$lib/components/ModelIcon.svelte';
	import { fade } from 'svelte/transition';

	// State for filtering
	let filteredModels: ModelsResponse = $state(data.models || {});
	let activeFilters: Set<string> = $state(new Set());
	let showFilterDropdown = $state(false);

	// Filter the models based on active filters
	$effect(() => {
		if (activeFilters.size === 0) {
			filteredModels = data.models;
		} else {
			const filtered: ModelsResponse = {};

			Object.entries(data.models).forEach(([modelId, model]: [string, ModelData]) => {
				const hasAllFeatures = Array.from(activeFilters).every((filter) => {
					return model.features && model.features[filter as keyof ModelFeatures] === true;
				});

				if (hasAllFeatures) {
					filtered[modelId] = model;
				}
			});

			filteredModels = filtered;
		}
	});

	function toggleFilter(feature: string) {
		if (activeFilters.has(feature)) {
			activeFilters.delete(feature);
		} else {
			activeFilters.add(feature);
		}
		activeFilters = new Set(activeFilters);
	}

	function clearFilters() {
		activeFilters.clear();
		activeFilters = new Set(activeFilters);
		showFilterDropdown = false;
	}

	function toggleFilterDropdown() {
		showFilterDropdown = !showFilterDropdown;
	}

	function handleClickOutside(event: MouseEvent) {
		const target = event.target as HTMLElement;
		if (!target.closest('.filter-container')) {
			showFilterDropdown = false;
		}
	}
</script>

<svelte:window onclick={handleClickOutside} />

<div class="container">
	<div class="title">Available Models</div>
	<div class="description">
		Choose which models appear in your model selector. This won't affect existing conversations.
	</div>
	<div class="header">
		<div class="head">
			<div class="filter-container">
				<button onclick={toggleFilterDropdown} class:active={activeFilters.size > 0}>
					Filter by features
					{#if activeFilters.size > 0}
						<span class="filter-count" transition:fade={{ duration: 100 }}
							>{activeFilters.size}</span
						>
					{/if}
				</button>
				{#if showFilterDropdown}
					<div class="filter-box">
						<button
							class="filter"
							class:active={activeFilters.has('has_vision')}
							onclick={() => toggleFilter('has_vision')}
						>
							<FeatureTag feature="has_vision" />
							<div class="name">Vision</div>
							{#if activeFilters.has('has_vision')}
								<div class="checkmark">
									<Check size="14" />
								</div>
							{/if}
						</button>
						<button
							class="filter"
							class:active={activeFilters.has('has_pdf')}
							onclick={() => toggleFilter('has_pdf')}
						>
							<FeatureTag feature="has_pdf" />
							<div class="name">PDFs</div>
							{#if activeFilters.has('has_pdf')}
								<div class="checkmark">
									<Check size="14" />
								</div>
							{/if}
						</button>
						<button
							class="filter"
							class:active={activeFilters.has('has_reasoning')}
							onclick={() => toggleFilter('has_reasoning')}
						>
							<FeatureTag feature="has_reasoning" />
							<div class="name">Reasoning</div>
							{#if activeFilters.has('has_reasoning')}
								<div class="checkmark">
									<Check size="14" />
								</div>
							{/if}
						</button>
						<button
							class="filter"
							class:active={activeFilters.has('has_fast')}
							onclick={() => toggleFilter('has_fast')}
						>
							<FeatureTag feature="has_fast" />
							<div class="name">Fast</div>
							{#if activeFilters.has('has_fast')}
								<div class="checkmark">
									<Check size="14" />
								</div>
							{/if}
						</button>
						<button
							class="filter"
							class:active={activeFilters.has('has_effort_control')}
							onclick={() => toggleFilter('has_effort_control')}
						>
							<FeatureTag feature="has_effort_control" />
							<div class="name">Effort Control</div>
							{#if activeFilters.has('has_effort_control')}
								<div class="checkmark">
									<Check size="14" />
								</div>
							{/if}
						</button>
						<button
							class="filter"
							class:active={activeFilters.has('has_web_search')}
							onclick={() => toggleFilter('has_web_search')}
						>
							<FeatureTag feature="has_web_search" />
							<div class="name">Search</div>
							{#if activeFilters.has('has_web_search')}
								<div class="checkmark">
									<Check size="14" />
								</div>
							{/if}
						</button>
						<button
							class="filter"
							class:active={activeFilters.has('has_image_generation')}
							onclick={() => toggleFilter('has_image_generation')}
						>
							<FeatureTag feature="has_image_generation" />
							<div class="name">Image Generation</div>
							{#if activeFilters.has('has_image_generation')}
								<div class="checkmark">
									<Check size="14" />
								</div>
							{/if}
						</button>
					</div>
				{/if}
			</div>
			{#if activeFilters.size > 0}
				<button
					class="clear-button"
					onclick={clearFilters}
					disabled={activeFilters.size === 0}
					transition:fade={{ duration: 100 }}
				>
					Clear
				</button>
			{/if}
		</div>
		<div class="tail"></div>
	</div>
	<div class="models">
		{#if Object.keys(filteredModels).length === 0}
			<div class="empty-state">
				<div class="empty-title">No models match your filters</div>
				<div class="empty-description">
					Try adjusting your filter criteria or clearing all filters.
				</div>
			</div>
		{:else}
			{#each Object.entries(filteredModels) as [modelId, model]}
				<div class="model">
					<ModelIcon {model} size="40" />
					<div class="details">
						<div class="name">{model.title}</div>
						<div class="description">
							{model.description}
						</div>
						<div class="tail">
							<div class="features">
								{#each Object.entries(model.features || {}) as [feature]}
									<FeatureTag {feature} wText={true} />
								{/each}
							</div>
							<!-- TODO: Copy Serach URL function and animation -->
							<button class="search-url-button">
								<Link size="14" />
								Search URL
							</button>
						</div>
					</div>
				</div>
			{/each}
		{/if}
	</div>
</div>

<style>
	.container {
		display: flex;
		flex-direction: column;
		gap: 20px;
		width: 100%;
		height: 100%;
		padding-bottom: 0;
	}

	.title {
		color: #ffffff;
		font-size: 22px;
		font-weight: 800;
	}

	.description {
		font-size: 14px;
		font-weight: 500;
		color: hsl(var(--secondary-foreground));
	}

	.header {
		display: flex;
		justify-content: space-between;
	}

	.header .head,
	.header .tail {
		display: flex;
		gap: 8px;
	}

	.header button {
		all: unset;
		display: flex;
		justify-content: center;
		align-items: center;
		gap: 8px;
		font-size: 12px;
		font-weight: 600;
		letter-spacing: 0.24px;
		padding: 4px 16px;
		border-radius: 5px;
		color: hsl(var(--secondary-foreground));
		border: 1px solid #302029;
		background-color: #21141e;
		cursor: pointer;
		transition:
			background-color 0.15s ease-out,
			border 0.15s ease-out;
	}

	.header button:hover:not(:disabled) {
		background-color: #261923;
	}

	.header button:disabled {
		color: hsl(var(--secondary-foreground) / 0.6);
		cursor: not-allowed;
	}

	.header button.clear-button {
		background-color: transparent;
		border: 1px solid transparent;
		color: #a3a3a3;
	}

	.header button.clear-button:hover:not(:disabled) {
		background-color: #2a232b;
		border: 1px solid #2a232b;
		color: hsl(var(--secondary-foreground));
	}

	.filter-count {
		color: #ffffff;
		font-size: 10px;
		margin-left: 4px;
	}

	.models {
		display: flex;
		flex-direction: column;
		gap: 16px;
		flex: 1;
		min-height: 0;
		overflow-y: auto;
	}

	.model {
		display: flex;
		gap: 16px;
		padding: 16px;
		border: 1px solid #302029;
		border-radius: 5px;
	}

	.model .details {
		width: 100%;
		display: flex;
		flex-direction: column;
		gap: 12px;
	}

	.model .name {
		font-size: 15px;
		font-weight: 600;
		letter-spacing: 0.24px;
	}

	.model .description {
		font-size: 14px;
		font-weight: 500;
	}

	.model .tail {
		display: flex;
		justify-content: space-between;
		align-items: center;
		gap: 16px;
	}

	.model .tail .features {
		display: flex;
		align-items: center;
		gap: 8px;
	}

	.model .search-url-button {
		all: unset;
		display: flex;
		align-items: center;
		padding: 4px 10px;
		gap: 12px;
		font-size: 12px;
		background-color: transparent;
		border-radius: 6px;
		cursor: pointer;
		transition: background-color 0.15s ease-out;
	}

	.model .search-url-button:hover {
		background-color: #29212a;
	}

	.filter-container {
		position: relative;
	}

	.filter-box {
		display: flex;
		flex-direction: column;
		position: absolute;
		top: 105%;
		left: 0;
		width: 100%;
		min-width: 240px;
		max-width: 240px;
		padding: 4px;
		border-radius: 8px;
		border: 1px solid #261f25;
		background-color: #0f0a0e;
		z-index: 10;
		box-shadow: 0 4px 12px rgba(0, 0, 0, 0.3);
	}

	button.filter {
		all: unset;
		box-sizing: border-box;
		display: flex;
		align-items: center;
		width: 100%;
		flex-direction: row;
		gap: 8px;
		font-size: 13px;
		padding: 6px;
		cursor: pointer;
		border-radius: 4px;
		transition: background-color 0.1s ease-out;
	}

	button.filter:hover {
		background-color: #201823;
	}

	.empty-state {
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 8px;
		padding: 40px 20px;
		text-align: center;
	}

	.empty-title {
		font-size: 16px;
		font-weight: 600;
		color: hsl(var(--foreground));
	}

	.empty-description {
		font-size: 14px;
		color: hsl(var(--secondary-foreground));
	}

	.checkmark {
		display: flex;
		align-items: center;
		margin-left: auto;
	}
</style>
