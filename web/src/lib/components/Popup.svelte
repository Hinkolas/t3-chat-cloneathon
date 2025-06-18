<script lang="ts">
	import { fade, scale } from 'svelte/transition';
	import { popup, hidePopup, notificationState, showNotification, hideNotification } from '$lib/store';

	function cancelPopup() {
		if ($popup.secondaryButtonFunction) {
			$popup.secondaryButtonFunction();
		}
		hidePopup();
	}

	function handlePrimaryAction() {
		$popup.primaryButtonFunction();
		showNotification({
			is_error: false,
			title: $popup.primaryButtonName,
			description: $popup.description || 'Action completed successfully.'
		});
		setTimeout(() => {
			hideNotification();
		}, 2000);
		hidePopup();
	}

	function updateInputValue(event: Event) {
		const target = event.target as HTMLInputElement;
		popup.update((state) => ({
			...state,
			inputValue: target.value
		}));
	}

	function handleKeydown(event: KeyboardEvent) {
		if (event.key === 'Escape') {
			cancelPopup();
		} else if (event.key === 'Enter' && !event.shiftKey) {
			event.preventDefault();
			handlePrimaryAction();
		}
	}
</script>

{#if $popup.show}
	<div
		class="popup-overlay"
		transition:fade={{ duration: 100 }}
		on:click={cancelPopup}
		on:keydown={handleKeydown}
		role="dialog"
		aria-modal="true"
		tabindex="0"
	>
		<!-- TODO: handle svelte ignores -->
		<!-- svelte-ignore a11y_click_events_have_key_events -->
		<!-- svelte-ignore a11y_no_static_element_interactions -->
		<!-- svelte-ignore a11y_autofocus -->
		<div class="popup" transition:scale={{ duration: 100, start: 0.9 }} on:click|stopPropagation>
			<div class="title">{$popup.title}</div>

			{#if $popup.description}
				<div class="description">{$popup.description}</div>
			{/if}

			<!-- Input field for rename and input popups -->
			{#if $popup.type === 'rename' || $popup.type === 'input'}
				<div class="input-container">
					{#if $popup.inputLabel}
						<label for="popup-input" class="input-label">{$popup.inputLabel}</label>
					{/if}
					<input
						id="popup-input"
						type={$popup.type === 'input' ? $popup.inputType || 'text' : 'text'}
						value={$popup.inputValue}
						placeholder={$popup.inputPlaceholder || ''}
						on:input={updateInputValue}
						on:keydown={handleKeydown}
						class="input-field"
						autofocus
					/>
				</div>
			{/if}

			<div class="buttons">
				<button on:click={cancelPopup} class="secondary">
					{$popup.secondaryButtonName || 'Cancel'}
				</button>
				<button
					on:click={handlePrimaryAction}
					class="primary"
					disabled={($popup.type === 'rename' || $popup.type === 'input') &&
						!$popup.inputValue.trim()}
				>
					{$popup.primaryButtonName}
				</button>
			</div>
		</div>
	</div>
{/if}

<style>
	.popup-overlay {
		z-index: 9999;
		position: fixed;
		top: 0;
		left: 0;
		width: 100%;
		height: 100%;
		background-color: var(--background-overlay);
		display: flex;
		align-items: center;
		justify-content: center;
	}

	.popup {
		display: flex;
		flex-direction: column;
		gap: 16px;
		max-width: 512px;
		min-width: 320px;
		padding: 24px;
		border: 1px solid var(--border);
		border-radius: 12px;
		background-color: var(--popup-background);
		margin: 20px;
	}

	.title {
		font-size: 18px;
		color: var(--white);
		font-weight: 600;
		margin: 0;
	}

	.description {
		font-size: 14px;
		color: var(--text);
		margin: 0;
	}

	.input-container {
		display: flex;
		flex-direction: column;
		gap: 8px;
	}

	.input-label {
		font-size: 14px;
		color: var(--white);
		font-weight: 500;
	}

	.input-field {
		all: unset;
		box-sizing: border-box;
		padding: 12px 16px;
		border: 1px solid var(--border);
		border-radius: 8px;
		background-color: var(--popup-input-background);
		color: var(--white);
		font-size: 14px;
		transition: border-color 0.15s ease-out;
	}

	.input-field:focus {
		border-color: var(--primary-background);
		outline: none;
	}

	.input-field::placeholder {
		color: var(--placeholder);
	}

	.buttons {
		display: flex;
		justify-content: flex-end;
		gap: 8px;
		margin-top: 8px;
	}

	button {
		all: unset;
		box-sizing: border-box;
		color: var(--white);
		letter-spacing: 0.24px;
		font-size: 14px;
		font-weight: 600;
		padding: 10px 20px;
		border-radius: 8px;
		cursor: pointer;
		transition: all 0.15s ease-out;
	}

	button:disabled {
		opacity: 0.5;
		cursor: not-allowed;
	}

	button.secondary:hover:not(:disabled) {
		background-color: var(--border);
	}

	button.primary {
		background-color: var(--primary-background);
	}

	button.primary:hover:not(:disabled) {
		background-color: var(--primary-background);
	}
</style>
