<script lang="ts">
	import { Key, Trash, Check } from '@lucide/svelte';
	import type { ProfileResponse } from '$lib/types';
	import { env } from '$env/dynamic/public';

	interface Props {
		data: {
			SESSION_TOKEN: string;
			profile: ProfileResponse;
		};
	}

	let { data } = $props();

	let profile = $state(data.profile || {});
	let SESSION_TOKEN = $state(data.SESSION_TOKEN || '');

	// svelte-ignore state_referenced_locally
	let originalValues = $state({
		anthropic_api_key: profile.anthropic_api_key || '',
		openai_api_key: profile.openai_api_key || '',
		gemini_api_key: profile.gemini_api_key || '',
		ollama_base_url: profile.ollama_base_url || ''
	});

	// Track save states for animations
	let saveStates = $state({
		anthropic: false,
		openai: false,
		google: false,
		ollama: false
	});

	// Helper function to check if value has changed
	function hasChanged(provider: 'anthropic' | 'openai' | 'google' | 'ollama'): boolean {
		const providerKeyMap = {
			anthropic: 'anthropic_api_key',
			openai: 'openai_api_key',
			google: 'gemini_api_key',
			ollama: 'ollama_base_url'
		} as const;
		const profileKey = providerKeyMap[provider];
		return profile[profileKey] !== originalValues[profileKey];
	}

	// Show success animation
	function showSaveSuccess(provider: 'anthropic' | 'openai' | 'google' | 'ollama') {
		saveStates[provider] = true;
		setTimeout(() => {
			saveStates[provider] = false;
		}, 2000);
	}

	async function saveKey(provider: 'anthropic' | 'openai' | 'google' | 'ollama', keyValue: string) {
		if (!keyValue || keyValue.length === 0) return;

		// Map provider to profile key
		const providerKeyMap: Record<string, string> = {
			anthropic: 'anthropic_api_key',
			openai: 'openai_api_key',
			google: 'gemini_api_key',
			ollama: 'ollama_base_url'
		};
		const profileKey = providerKeyMap[provider];
		if (!profileKey) {
			console.error(`Unknown provider: ${provider}`);
			return;
		}

		const res = await fetch(`${env.PUBLIC_API_URL}/v1/profile/`, {
			method: 'PATCH',
			headers: {
				'Content-Type': 'application/json',
				Authorization: `Bearer ${SESSION_TOKEN}`
			},
			body: JSON.stringify({ [profileKey]: keyValue })
		});

		if (!res.ok) {
			console.error(`Failed to save ${provider} key`);
			return;
		}

		const updatedProfile = await res.json();
		profile = { ...profile, ...updatedProfile };

		// Update original values after successful save
		originalValues = { ...originalValues, [profileKey]: keyValue };

		// Show success animation
		showSaveSuccess(provider);
	}

	async function deleteKey(provider: 'anthropic' | 'openai' | 'google' | 'ollama') {
		// Map provider to profile key
		const providerKeyMap: Record<string, string> = {
			anthropic: 'anthropic_api_key',
			openai: 'openai_api_key',
			google: 'gemini_api_key',
			ollama: 'ollama_base_url'
		};
		const profileKey = providerKeyMap[provider];
		if (!profileKey) {
			console.error(`Unknown provider: ${provider}`);
			return;
		}

		const res = await fetch(`${env.PUBLIC_API_URL}/v1/profile/`, {
			method: 'PATCH',
			headers: {
				'Content-Type': 'application/json',
				Authorization: `Bearer ${SESSION_TOKEN}`
			},
			body: JSON.stringify({ [profileKey]: '' })
		});

		if (!res.ok) {
			console.error(`Failed to delete ${provider} key`);
			return;
		}

		const updatedProfile = await res.json();
		profile = { ...profile, ...updatedProfile };

		profile = {
			...profile,
			...updatedProfile,
			[profileKey]: ''
		};

		// Update original values after successful delete
		originalValues = { ...originalValues, [profileKey]: '' };
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
				{#if profile.anthropic_api_key}
					<button class="delete-button" onclick={() => deleteKey('anthropic')}>
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
						bind:value={profile.anthropic_api_key}
						placeholder="Enter your API key here"
					/>
					<div class="label">Get your API key from Anthropic</div>
				</div>
				{#if hasChanged('anthropic')}
					<button
						class="save-button"
						class:saved={saveStates.anthropic}
						onclick={() => saveKey('anthropic', profile.anthropic_api_key)}
					>
						{#if saveStates.anthropic}
							<Check size="16" />
							Saved!
						{:else}
							Save
						{/if}
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
				{#if profile.openai_api_key}
					<button class="delete-button" onclick={() => deleteKey('openai')}>
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
						bind:value={profile.openai_api_key}
						placeholder="Enter your API key here"
					/>
					<div class="label">Get your API key from OpenAI's Dashboard</div>
				</div>
				{#if hasChanged('openai')}
					<button
						class="save-button"
						class:saved={saveStates.openai}
						onclick={() => saveKey('openai', profile.openai_api_key)}
					>
						{#if saveStates.openai}
							<Check size="16" />
							Saved!
						{:else}
							Save
						{/if}
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
				{#if profile.gemini_api_key}
					<button class="delete-button" onclick={() => deleteKey('google')}>
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
						bind:value={profile.gemini_api_key}
						placeholder="Enter your API key here"
					/>
					<div class="label">Get your API key from Google's Dashboard</div>
				</div>
				{#if hasChanged('google')}
					<button
						class="save-button"
						class:saved={saveStates.google}
						onclick={() => saveKey('google', profile.gemini_api_key)}
					>
						{#if saveStates.google}
							<Check size="16" />
							Saved!
						{:else}
							Save
						{/if}
					</button>
				{/if}
			</div>
		</div>

		<!-- Ollama Base URL -->
		<div class="api-key-container">
			<div class="head">
				<div class="title">
					<Key size="16" />
					Ollama Base URL
				</div>
				{#if profile.ollama_base_url}
					<button class="delete-button" onclick={() => deleteKey('ollama')}>
						<Trash size="16" />
					</button>
				{/if}
			</div>
			<div class="body">
				<div class="description">Used for the following models:</div>
				<div class="models">
					<div class="model">Llama 3.2</div>
				</div>
			</div>
			<div class="tail">
				<div class="input-group">
					<input
						type="text"
						bind:value={profile.ollama_base_url}
						placeholder="Enter your Ollama base URL here"
					/>
					<div class="label"></div>
				</div>
				{#if hasChanged('ollama')}
					<button
						class="save-button"
						class:saved={saveStates.ollama}
						onclick={() => saveKey('ollama', profile.ollama_base_url)}
					>
						{#if saveStates.ollama}
							<Check size="16" />
							Saved!
						{:else}
							Save
						{/if}
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
		color: var(--text);
	}

	.api-keys {
		display: flex;
		flex-direction: column;
		gap: 16px;
		flex: 1;
		min-height: 0;
		overflow-y: auto;
	}

	.api-keys::-webkit-scrollbar {
		background-color: transparent;
		width: 6px ;
	}

	.api-keys::-webkit-scrollbar-thumb {
		background-color: var(--text-disabled);
		border-radius: 10px;
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
		color: var(--text);
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
		color: var(--text);
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
		color: var(--text);
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
		color: hsl(var(--secondary-foreground) / 0.8);
		font-size: 14px;
		outline: none;
	}

	.input-group input:focus {
		border-color: hsl(var(--primary) / 0.8);
		color: var(--text);
	}

	.input-group .label {
		font-size: 12px;
		color: var(--text);
	}

	.save-button {
		all: unset;
		display: flex;
		justify-content: center;
		align-items: center;
		gap: 6px;
		padding: 6px 16px;
		background-color: hsl(var(--primary));
		color: #ffffff;
		font-size: 14px;
		font-weight: 600;
		border-radius: 4px;
		cursor: pointer;
		transition: all 0.2s ease-out;
	}

	.save-button:hover {
		background-color: hsl(var(--primary) / 0.8);
	}

	.save-button.saved {
		background-color: #22c55e;
		animation: saveSuccess 0.3s ease-out;
	}

	.save-button.saved:hover {
		background-color: #16a34a;
	}

	@keyframes saveSuccess {
		0% {
			transform: scale(1);
		}
		50% {
			transform: scale(1.05);
		}
		100% {
			transform: scale(1);
		}
	}
</style>
