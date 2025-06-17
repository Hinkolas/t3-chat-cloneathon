<script lang="ts">
	import { Key, Trash, Check } from '@lucide/svelte';
	import { writable } from 'svelte/store';

	// Simulate loading/saving keys from storage
	const initialKeys = {
		anthropic: '', // e.g. 'sk-...'
		openai: '',
		google: ''
	};

	const apiKeys = writable({ ...initialKeys });

	function saveKey(provider: string, value: string) {
		apiKeys.update((keys) => ({ ...keys, [provider]: value }));
	}

	function deleteKey(provider: string) {
		apiKeys.update((keys) => ({ ...keys, [provider]: '' }));
	}
</script>

<div class="container">
	<div class="title">API Keys</div>
	<div class="description">
		Bring your own API keys for select models. Messages sent using your API keys will not count
		towards your monthly limits.
	</div>
	<div class="api-keys">
		<!-- Anthropic -->
		<div class="api-key-container">
			<div class="head">
				<div class="title">
					<Key size="16" />
					Anthropic API Key
				</div>
				{#if $apiKeys.anthropic}
					<button class="delete-button" on:click={() => deleteKey('anthropic')}>
						<Trash size="16" />
					</button>
				{/if}
			</div>
			<div class="body">
				<div class="description">Used for the following models:</div>
				<div class="models">
					<div class="model">Claude 3.5 Sonnet</div>
					<div class="model">Claude 3.7 Sonnet</div>
					<div class="model">Claude 3.7 Sonnet (Reasoning)</div>
					<div class="model">Claude 4 Opus</div>
					<div class="model">Claude 4 Sonnet</div>
					<div class="model">Claude 4 Sonnet (Reasoning)</div>
				</div>
			</div>
			<div class="tail">
				<div class="input-group">
					<input
						type="text"
						bind:value={$apiKeys.anthropic}
						placeholder="Enter your API key here"
						on:input={(e) => {
							if (e.target) saveKey('anthropic', (e.target as HTMLInputElement).value);
						}}
					/>
					<div class="label">Get your API key from Anthropic</div>
				</div>
				{#if !$apiKeys.anthropic}
					<button
						class="save-button"
						disabled={!$apiKeys.anthropic}
						on:click={() => saveKey('anthropic', $apiKeys.anthropic)}
					>
						Save
					</button>
				{/if}
			</div>
		</div>

		<!-- OpenAI -->
		<div class="api-key-container">
			<div class="head">
				<div class="title">
					<Key size="16" />
					OpenAI API Key
				</div>
				{#if $apiKeys.openai}
					<button class="delete-button" on:click={() => deleteKey('openai')}>
						<Trash size="16" />
					</button>
				{/if}
			</div>
			<div class="body">
				<div class="description">Used for the following models:</div>
				<div class="models">
					<div class="model">GPT-4.5</div>
					<div class="model">o3</div>
				</div>
			</div>
			<div class="tail">
				<div class="input-group">
					<input
						type="text"
						bind:value={$apiKeys.openai}
						placeholder="Enter your API key here"
						on:input={(e) => {
							if (e.target) saveKey('openai', (e.target as HTMLInputElement).value);
						}}
					/>
					<div class="label">Get your API key from OpenAI's Dashboard</div>
				</div>
				{#if !$apiKeys.openai}
					<button
						class="save-button"
						disabled={!$apiKeys.openai}
						on:click={() => saveKey('openai', $apiKeys.openai)}
					>
						Save
					</button>
				{/if}
			</div>
		</div>

		<!-- Google -->
		<div class="api-key-container">
			<div class="head">
				<div class="title">
					<Key size="16" />
					Google API Key
				</div>
				{#if $apiKeys.google}
					<button class="delete-button" on:click={() => deleteKey('google')}>
						<Trash size="16" />
					</button>
				{/if}
			</div>
			<div class="body">
				<div class="description">Used for the following models:</div>
				<div class="models">
					<div class="model">Gemini 2.0 Flash</div>
					<div class="model">Gemini 2.0 Flash Lite</div>
					<div class="model">Gemini 2.5 Flash</div>
					<div class="model">Gemini 2.5 Flash (Thinking)</div>
					<div class="model">Gemini 2.5 Pro</div>
				</div>
			</div>
			<div class="tail">
				<div class="input-group">
					<input
						type="text"
						bind:value={$apiKeys.google}
						placeholder="Enter your API key here"
						on:input={(e) => {
							if (e.target) saveKey('google', (e.target as HTMLInputElement).value);
						}}
					/>
					<div class="label">Get your API key from Google's Dashboard</div>
				</div>
				{#if !$apiKeys.google}
					<button
						class="save-button"
						disabled={!$apiKeys.google}
						on:click={() => saveKey('google', $apiKeys.google)}
					>
						Save
					</button>
				{/if}
			</div>
		</div>
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

	.api-keys {
		display: flex;
		flex-direction: column;
		gap: 16px;
		flex: 1;
		min-height: 0;
		overflow-y: auto;
	}

	.api-key-container {
		display: flex;
		flex-direction: column;
		gap: 12px;
		padding: 16px;
		border: 1px solid hsl(var(--primary) / 0.2);
		border-radius: 8px;
	}

	.head {
		display: flex;
		justify-content: space-between;
		align-items: center;
	}

	.head .title {
		display: flex;
		align-items: center;
		gap: 8px;
		font-size: 16px;
		font-weight: 600;
	}

	.delete-button {
		all: unset;
		cursor: pointer;
		color: hsl(var(--secondary-foreground));
		display: flex;
		align-items: center;
		gap: 4px;
		transition: color 0.1s ease-out;
	}
	.delete-button:hover {
		color: hsl(var(--primary) / 0.8);
	}

	.body {
		display: flex;
		flex-direction: column;
		gap: 8px;
	}

	.body .description {
		font-size: 14px;
		color: hsl(var(--secondary-foreground));
	}

	.models {
		display: flex;
		flex-wrap: wrap;
		gap: 8px;
	}

	.model {
		padding: 4px 12px;
		background-color: hsl(var(--primary) / 0.1);
		border-radius: 999px;
		border: 1px solid hsl(var(--primary) / 0.2);
		font-size: 12px;
		font-weight: 500;
		color: hsl(var(--secondary-foreground));
	}

	.tail {
		display: flex;
		flex-direction: column;
		justify-content: space-between;
		align-items: flex-end;
		gap: 16px;
	}

	.input-group {
		width: 100%;
		display: flex;
		flex-direction: column;
		gap: 4px;
	}

	.input-group input {
		width: 100%;
		padding: 8px;
		border: 1px solid hsl(var(--primary) / 0.3);
		border-radius: 4px;
		background-color: hsl(var(--background) / 0.8);
		color: hsl(var(--foreground));
		font-size: 14px;
		outline: none;
	}

	.input-group input:focus {
		border-color: hsl(var(--primary) / 0.8);
	}

	.input-group .label {
		font-size: 12px;
		color: hsl(var(--secondary-foreground));
	}

	.save-button {
		all: unset;
		display: flex;
		justify-content: center;
		align-items: center;
		padding: 6px 16px;
		background-color: hsl(var(--primary));
		color: #ffffff;
		font-size: 14px;
		font-weight: 600;
		border-radius: 4px;
		cursor: pointer;
		transition: background-color 0.1s ease-out;
	}
	.save-button:hover {
		background-color: hsl(var(--primary) / 0.8);
	}
</style>
