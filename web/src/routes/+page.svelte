<script lang="ts">
	import {
		ArrowUp,
		ChevronDown,
		Globe,
		Paperclip,
		PanelLeft,
		Search,
		Eye,
		Brain
	} from '@lucide/svelte';
	import { onMount } from 'svelte';
	import FeatureTag from './FeatureTag.svelte';
	import googleIcon from '$lib/assets/modelIcons/googleai.svg';
	import chatgptIcon from '$lib/assets/modelIcons/chatgpt.svg';
	import anthropicIcon from '$lib/assets/modelIcons/anthropic.svg';
	import metaIcon from '$lib/assets/modelIcons/meta.svg';

	interface ModelSelections {
		icon: string;
		modelName: string;
		features: string[];
	}

	const models: ModelSelections[] = [
		{
			icon: googleIcon,
			modelName: 'Gemini 2.5 Flash',
			features: ['vision', 'think']
		},
		{
			icon: chatgptIcon,
			modelName: 'GPT 4o-mini',
			features: ['vision']
		},
		{
			icon: metaIcon,
			modelName: 'Llama 4 Scout',
			features: ['search']
		},
		{
			icon: metaIcon,
			modelName: 'Llama 3.3 70b',
			features: ['think']
		},
		{
			icon: anthropicIcon,
			modelName: 'Claude 4 Sonnet (Resonning)',
			features: ['vision', 'search', 'think']
		},
		{
			icon: anthropicIcon,
			modelName: 'Claude 4 Sonnet',
			features: ['search']
		}
	];

	const iconSize = 16;

	let textarea: HTMLElement;
	let message = $state('');
	let sidebarCollapsed = $state(true);
	let modelSelectionOpen = $state(false);

	function autoResize() {
		if (textarea) {
			textarea.style.height = 'auto';
			textarea.style.height = textarea.scrollHeight + 'px';
		}
	}

	function sendMessage() {
		console.log('do something');
	}

	function toggleModelSelection() {
		modelSelectionOpen = !modelSelectionOpen;
	}

	function closeModelSelection() {
		modelSelectionOpen = false;
	}

	function toggleSidebar() {
		sidebarCollapsed = !sidebarCollapsed;
	}

	onMount(() => {
		autoResize();
	});

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
	<div class="sidebar {sidebarCollapsed ? 'collapsed' : ''}"></div>
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
								<div class="search-container">
									<Search size={iconSize} />
									<input placeholder="Search Models..." type="text" />
								</div>
								<div class="model-container">
									{#each models as model}
										<div class="model">
											<div class="details">
												<img src={model.icon} alt={model.modelName} />
												<div class="title">{model.modelName}</div>
											</div>
											<div class="feature-container">
												{#each model.features as featureName}
													<FeatureTag tag={featureName} />
												{/each}
											</div>
										</div>
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
		top: 20px;
		left: 20px;

		z-index: 1;
		display: flex;
		justify-self: center;
		align-items: center;
		border: 1px solid #88888833;
		border-radius: 8px;
		padding: 8px;
		cursor: pointer;
		color: hsl(var(--secondary-foreground));
	}

	.sidebar {
		flex: 0 0 256px;
		padding: 16px;
		background-color: var(--sidebar-background);
		transition: margin 0.15s ease-in-out;
	}

	.sidebar.collapsed {
		margin-left: -256px;
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
	.model {
		display: flex;
		flex-direction: row;
		justify-content: space-between;
		align-items: center;
		gap: 32px;
		padding-inline: 8px;
		border-radius: 4px;
		cursor: pointer;
		transition: background-color 0.15s ease;
	}
	.model:hover {
		background-color: var(--button-hover);
	}
	.details {
		display: flex;
		flex-direction: row;
		justify-content: flex-start;
		align-items: center;
		gap: 8px;
		padding-block: 8px;
	}

	.model img {
		width: 20px;
		height: 20px;
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
		flex: 1;
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
