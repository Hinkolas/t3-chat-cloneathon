<script lang="ts">
	import { env } from '$env/dynamic/public';
	import { hideNotification, notificationState, showNotification } from '$lib/store';
	import type { ProfileResponse } from '$lib/types.js';

	interface Props {
		data: {
			SESSION_TOKEN: string;
			profile: ProfileResponse;
		};
	}

	let { data }: Props = $props();

	let profile = $state(data.profile);
	let SESSION_TOKEN = $state(data.SESSION_TOKEN || '');

	async function savePreferences() {
		const res = await fetch(`${env.PUBLIC_API_URL}/v1/profile/`, {
			method: 'PATCH',
			headers: {
				'Content-Type': 'application/json',
				Authorization: `Bearer ${SESSION_TOKEN}`
			},
			body: JSON.stringify(profile)
		});

		if (!res.ok) {
			showNotification({
				is_error: true,
				title: 'Error Saving Preferences',
				description: 'There was an error saving your preferences. Please try again.'
			});
			setTimeout(() => {
				hideNotification();
			}, 2000);
			console.error(`Failed to save preferences: ${res.status} ${res.statusText}`);
			return;
		}

		const updatedProfile = await res.json();
		profile = { ...profile, ...updatedProfile };
		showNotification({
			is_error: false,
			title: 'Preferences Saved',
			description: ''
		});
		setTimeout(() => {
			hideNotification();
		}, 2000);
	}
</script>

<div class="container">
	<div class="title">Customize Your Chat</div>
	<form action="">
		<div class="form-group">
			<label for="name">What should Your Chat call you?</label>
			<div class="input-container">
				<input bind:value={profile.custom_user_name} type="text" placeholder="Enter your name" />
				<div class="chars">{profile.custom_user_name.length}/50</div>
			</div>
		</div>
		<div class="form-group">
			<label for="name">What do you do?</label>
			<div class="input-container">
				<input
					bind:value={profile.custom_user_profession}
					type="text"
					placeholder="Engineer, student, etc."
				/>
				<div class="chars">{profile.custom_user_profession.length}/100</div>
			</div>
		</div>
		<div class="form-group">
			<label for="name">What traits should Your Chat have?</label>
			<div class="input-container">
				<input
					bind:value={profile.custom_assistant_trait}
					type="text"
					placeholder="Type a trait and press Enter or Tab"
				/>
				<div class="chars">{profile.custom_assistant_trait.length}/50</div>
			</div>
		</div>
		<div class="form-group">
			<label for="name">Anything else Your Chat should know about you?</label>
			<div class="input-container">
				<textarea
					bind:value={profile.custom_context}
					rows="4"
					maxlength="3000"
					cols="50"
					placeholder="Interests, values, or preferences to keep in mind"
				></textarea>
				<div class="chars">{profile.custom_context.length}/3000</div>
			</div>
		</div>
		<button onclick={savePreferences}>Save Preferences</button>
	</form>
</div>

<style>
	.container {
		display: flex;
		flex-direction: column;
		gap: 32px;
	}

	.title {
		color: #ffffff;
		font-size: 22px;
		font-weight: 800;
	}

	button {
		all: unset;
		margin-top: 16px;
		width: max-content;
		font-size: 14px;
		font-weight: 500;
		color: var(--text);
		padding: 8px 16px;
		cursor: pointer;
		border-radius: 8px;
		align-self: flex-end;
		border: 1px solid hsl(var(--primary) / 0.2);
		background-color: hsl(var(--primary) / 0.2);
		transition: background-color 0.15s ease-out;
	}

	button:hover:not(:disabled) {
		background-color: hsl(var(--primary) / 0.8);
	}

	button:disabled {
		cursor: not-allowed;
		border: 1px solid hsl(var(--primary) / 0.1);
		background-color: hsl(var(--primary) / 0.1);
	}

	form {
		display: flex;
		flex-direction: column;
		gap: 16px;
	}

	.form-group {
		display: flex;
		flex-direction: column;
		gap: 8px;
	}

	.form-group label {
		font-size: 15px;
		font-weight: 600;
	}

	.form-group .input-container {
		position: relative;
		display: flex;
	}

	.form-group .input-container input {
		padding: 8px 12px;
		font-size: 14px;
		background-color: transparent;
		border: none;
		border-radius: 8px;
		box-shadow: 0 0 2px #483039;
		width: 100%;
		color: var(--text);
		outline: none;
	}

	.form-group .input-container input:focus {
		color: var(--text);
		box-shadow: 0 0 2px var(--primary-background);
	}

	.form-group .input-container input::placeholder {
		color: var(--placeholder);
	}

	.form-group .input-container .chars {
		position: absolute;
		bottom: 0;
		right: 6px;
		font-size: 11px;
		color: var(--placeholder);
	}

	.form-group .input-container textarea {
		all: unset;
		padding: 8px 12px;
		font-size: 14px;
		background-color: transparent;
		border: none;
		border-radius: 8px;
		box-shadow: 0 0 2px #483039;
		width: 100%;
		color: var(--primary-border);
		outline: none;
		min-height: 120px;
		max-height: 500px;
		resize: vertical;
	}

	.form-group .input-container textarea:focus {
		color: var(--text);
		box-shadow: 0 0 2px var(--primary-background);
	}

	.form-group .input-container textarea::placeholder {
		color: var(--placeholder);
	}
</style>
