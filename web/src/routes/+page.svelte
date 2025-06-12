<script lang="ts">
	let { data }: PageProps = $props();

	console.log(data);

	import {
		ArrowUp,
		ChevronDown,
		Globe,
		Paperclip,
		PanelLeft,
		Search,
	} from '@lucide/svelte';
	import { onMount } from 'svelte';
	import ModelRow from './ModelRow.svelte';
	import Sidebar from './Sidebar.svelte';
	import SearchInput from './SearchInput.svelte';
	import type { PageProps } from './$types';

	interface ModelSelections {
		icon: string;
		modelName: string;
		features: string[];
	}

	const models: ModelSelections[] = [
		{
			icon: 'gemini',
			modelName: 'Gemini 2.5 Flash',
			features: ['vision', 'think']
		},
		{
			icon: 'openai',
			modelName: 'GPT 4o-mini',
			features: ['vision']
		},
		{
			icon: 'meta',
			modelName: 'Llama 4 Scout',
			features: ['search']
		},
		{
			icon: 'meta',
			modelName: 'Llama 3.3 70b',
			features: ['think']
		},
		{
			icon: 'anthropic',
			modelName: 'Claude 4 Sonnet (Resonning)',
			features: ['vision', 'search', 'think']
		},
		{
			icon: 'anthropic',
			modelName: 'Claude 4 Sonnet',
			features: ['search']
		}
	];

	const iconSize = 16;

	let textarea: HTMLElement;
	let message = $state('');
	let sidebarCollapsed = $state(false);
	let modelSelectionOpen = $state(false);
	let modelSearchTerm: string = $state('');
	let filteredModels: ModelSelections[] = $state(models);

	function sendMessage() {
		console.log('do something');
	}

	function toggleModelSelection() {
		if (modelSelectionOpen) {
			closeModelSelection();
		} else {
			modelSelectionOpen = true;
		}
	}

	function closeModelSelection() {
		modelSelectionOpen = false;

		setTimeout(() => {
			modelSearchTerm = '';
			filteredModels = models;
		}, 150);
	}

	function toggleSidebar() {
		sidebarCollapsed = !sidebarCollapsed;
	}

	function modelSearchFilter() {
		console.log(modelSearchTerm);
		filteredModels = models.filter((model) =>
			model.modelName.toLowerCase().includes(modelSearchTerm.toLowerCase())
		);
	}

	function changeModel() {
		closeModelSelection();
	}

	function newChat() {

	}

	onMount(() => {
		autoResize();
	});

	function autoResize() {
		if (textarea) {
			textarea.style.height = 'auto';
			textarea.style.height = textarea.scrollHeight + 'px';
		}
	}

	function clickOutside(node: Element) {
		const handleClick = (event: Event) => {
			if (!node.contains(<Node>event.target)) {
				node.dispatchEvent(new CustomEvent('outsideclick'));
			}
		};

		document.addEventListener('click', handleClick, true);

		return {
			destroy() {
				document.removeEventListener('click', handleClick, true);
			}
		};
	}
</script>

<div class="container">
	<button onclick={toggleSidebar} class="sidebar-button">
		<PanelLeft size={iconSize} />
	</button>
	<Sidebar {sidebarCollapsed} {newChat}/>
	<div class="content">
		<div class="chat"></div>
		<div class="input-wrapper">
			<div class="input-container">
				<textarea
					bind:this={textarea}
					bind:value={message}
					oninput={autoResize}
					placeholder="Type your message here..."
					name="message"
					id="Message"
				></textarea>
				<div class="buttons">
					<div class="button-group">
						<div class="selection-container" use:clickOutside onoutsideclick={closeModelSelection}>
							<div class="selection-box {modelSelectionOpen ? 'visible' : ''}">
								<SearchInput bind:value={modelSearchTerm} onInputFunction={modelSearchFilter} placeholder="Search Models..."/>
								<div class="model-container">
									{#each filteredModels as model}
										<ModelRow {model} {changeModel} />
									{/each}
								</div>
							</div>
							<button onclick={toggleModelSelection} class="selection-button non-selectable">
								<span>Gemini 2.5 Flash</span>
								<ChevronDown size={iconSize} />
							</button>
						</div>
						<button>
							<Globe size={iconSize} />
							<span>Search</span>
						</button>
						<button>
							<Paperclip size={iconSize} />
						</button>
					</div>
					<div class="button-group">
						<button
							class={message.length == 0 ? '' : 'active'}
							onclick={sendMessage}
							disabled={message.length == 0}
							id="SendButton"
						>
							<ArrowUp size="20" />
						</button>
					</div>
				</div>
			</div>
		</div>
	</div>
</div>

<style>
	.container {
		position: relative;
		width: 100%;
		height: 100dvh;
		display: flex;
	}

	.sidebar-button {
		all: unset;
		position: absolute;
		top: 16px;
		left: 16px;

		z-index: 1;
		display: flex;
		justify-self: center;
		align-items: center;
		border: 1px solid #88888822;
		border-radius: 8px;
		padding: 8px;
		cursor: pointer;
		color: hsl(var(--secondary-foreground));
		transition: background-color .15s ease-out;
	}

	.content {
		position: relative;
		flex: 1;
		background-color: var(--chat-background);
	}

	.input-wrapper {
		position: absolute;
		bottom: 0;
		left: 50%;
		transform: translateX(-50%);

		width: 100%;
		max-width: 768px;
		padding: 8px;
		padding-bottom: 0px;
		border-top-left-radius: 20px;
		border-top-right-radius: 20px;
		border: 1px solid #88888811;
		border-bottom: none;
		background-color: hsl(var(--chat-input-gradient));
	}

	.input-container {
		display: flex;
		flex-direction: column;
		gap: 14px;

		border-top-left-radius: 12px;
		border-top-right-radius: 12px;
		padding-block: 12px;
		border: 1px solid #88888811;

		background-color: var(--chat-input-background);
	}

	textarea {
		all: unset;
		font-size: 15px;
		min-height: 1.5rem;
		max-height: 200px;
		overflow-y: auto;
		resize: none;
		padding-inline: 12px;
		color: hsl(var(--secondary-foreground));
	}

	textarea::-webkit-scrollbar {
		display: none;
	}

	textarea::placeholder {
		color: #888888;
	}

	.buttons {
		width: 100%;
		display: flex;
		justify-content: space-between;
		padding-left: 6px;
		padding-right: 12px;
	}

	.button-group {
		display: flex;
		gap: 8px;
	}

	button {
		all: unset;

		display: flex;
		justify-content: center;
		align-items: center;
		gap: 4px;
		padding: 8px;
		border-radius: 9999px;
		white-space: nowrap;
		font-size: 12px;
		line-height: 1rem;
		cursor: pointer;
		border: 1px solid #88888833;
		color: hsl(var(--secondary-foreground));
		transition: background-color 0.15s ease;
	}

	button:hover {
		background-color: var(--button-hover);
	}

	.selection-container {
		position: relative;
		display: flex;
		justify-content: center;
		align-items: center;
	}

	.selection-box {
		position: absolute;
		bottom: 110%;
		left: 0;
		background-color: #0f0a0e;
		border-radius: 8px;
		display: flex;
		flex-direction: column;

		opacity: 0;
		transform: translateY(10px) scale(0.95);
		pointer-events: none;
		transition:
			opacity 0.2s ease,
			transform 0.2s ease;
	}

	.selection-box.visible {
		opacity: 1;
		transform: translateY(0) scale(1);
		pointer-events: auto;
	}

	.model-container {
		display: flex;
		flex-direction: column;
		padding-inline: 8px;
		padding-bottom: 16px;
	}

	

	.selection-button {
		border: none;
		border-radius: 8px;
		font-size: 14px;
		font-weight: 600;
	}

	#SendButton {
		all: unset;
		display: flex;
		justify-content: center;
		align-items: center;
		padding: 8px;
		background-color: hsl(var(--primary) / 0.2);
		box-shadow: 0px 0px 2px hsl(var(--primary));
		border-radius: 8px;
		cursor: pointer;
		color: hsl(var(--secondary-foreground));
		transition: background-color 0.15s ease;
	}

	#SendButton.active {
		background-color: hsl(var(--primary) / 0.4);
	}

	#SendButton:hover {
		background-color: hsl(var(--primary) / 0.8);
	}
</style>
