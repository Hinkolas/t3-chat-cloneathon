<script lang="ts">
	import { fade, scale } from 'svelte/transition';
	import { popupModule } from '$lib/store';

	function cancelPopup() {
		popupModule.update(() => {
			return {
				show: false,
				title: '',
				description: '',
				primaryButtonName: '',
				primaryButtonFunction: () => {}
			};
		});
	}
</script>

{#if $popupModule.show}
	<div class="popup-overlay" transition:fade={{ duration: 100 }}>
		<div class="popup" transition:scale={{ duration: 100, start: 0.9 }}>
			<div class="title">{$popupModule.title}</div>
			<div class="description">{$popupModule.description}</div>
			<div class="buttons">
				<button onclick={cancelPopup} class="secondary">Cancel</button>
				<button onclick={$popupModule.primaryButtonFunction} class="primary"
					>{$popupModule.primaryButtonName}</button
				>
			</div>
		</div>
	</div>
{/if}

<style>
	.popup-overlay {
		z-index: 9999;
		position: absolute;
		top: 0;
		left: 0;
		width: 100%;
		height: 100%;
		background-color: #00000099;
	}

	.popup {
		position: absolute;
		top: 50%;
		left: 50%;
		transform: translate(-50%, -50%);

		display: flex;
		flex-direction: column;
		gap: 12px;
		max-width: 512px;
		padding: 24px;
		border: 1px solid #88888833;
		border-radius: 12px;
		background-color: var(--popup-background);
	}

	.title {
		font-size: 18px;
		color: #ffffff;
		font-weight: 600;
	}

	.description {
		font-size: 14px;
		color: hsl(var(--secondary-foreground));
	}

	.buttons {
		display: flex;
		justify-content: flex-end;
		gap: 8px;
	}

	button {
		all: unset;
		box-sizing: border-box;
		color: #ffffff;
		letter-spacing: 0.24px;
		font-size: 14px;
		font-weight: 600;
		padding: 8px 16px;
		border-radius: 8px;
		cursor: pointer;
		transition: background-color 0.15s ease-out;
	}

	button.secondary:hover {
		background-color: #88888822;
	}

	button.primary {
		background-color: hsl(var(--primary) / 0.9);
	}

	button.primary:hover {
		background-color: hsl(var(--primary));
	}
</style>
