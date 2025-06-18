<script>
	import { notificationState, hideNotification } from '$lib/store';
    import { X } from '@lucide/svelte';
</script>

<div class="notification" class:error={$notificationState.is_error}>
	<div class="notification-content">
		<div class="notification-details">
			<h3>{$notificationState.title}</h3>
			<p>{$notificationState.description}</p>
		</div>
		<button type="button" class="close-btn" onclick={hideNotification}><X size=14/></button>
	</div>
</div>

<style>
	.notification {
		position: relative;
		width: 100%;
		margin-top: 20px;
		padding: 16px 20px;
		border-radius: 8px;
		border: 1px solid var(--notification-success-border);
		background: linear-gradient(
			135deg,
			var(--notification-success-bg) 0%,
			var(--notification-success-bg-secondary) 100%
		);
		color: var(--notification-success-text);
		box-shadow: 0 4px 12px var(--notification-shadow);
		animation: slideIn 0.3s ease-out;
		transition: all 0.3s ease;
	}

	.notification.error {
		border-color: var(--notification-error-border);
		background: linear-gradient(
			135deg,
			var(--notification-error-bg) 0%,
			var(--notification-error-bg-secondary) 100%
		);
		color: var(--notification-error-text);
		box-shadow: 0 4px 12px var(--notification-error-shadow);
	}

	.notification::before {
		content: '';
		position: absolute;
		left: 0;
		top: 0;
		bottom: 0;
		width: 4px;
		background: var(--notification-success-accent);
		border-radius: 4px 0 0 4px;
	}

	.notification.error::before {
		background: var(--notification-error-accent);
	}

	.notification-content {
		position: relative;
		display: flex;
		justify-content: space-between;
		align-items: center;
		gap: 16px;
	}

	.notification-details {
		display: flex;
		flex-direction: column;
	}

	.notification h3 {
		margin: 0 0 6px 0;
		font-size: 15px;
		font-weight: 600;
		line-height: 1.3;
		color: inherit;
	}

	.notification p {
		margin: 0;
		font-size: 14px;
		line-height: 1.4;
		color: inherit;
		opacity: 0.9;
	}

	.close-btn {
		background: none;
		border: none;
		font-size: 20px;
		font-weight: 600;
		cursor: pointer;
		color: inherit;
		padding: 4px;
		width: 28px;
		height: 28px;
		display: flex;
		align-items: center;
		justify-content: center;
		border-radius: 50%;
		opacity: 0.7;
		transition: all 0.2s ease;
		line-height: 1;
	}

	.close-btn:hover {
		opacity: 1;
		background: rgba(0, 0, 0, 0.1);
		transform: scale(1.1);
	}

	.close-btn:active {
		transform: scale(0.95);
	}

	.notification.error .close-btn:hover {
		background: var(--notification-error-close-hover);
	}

	/* Animations */
	@keyframes slideIn {
		from {
			opacity: 0;
			transform: translateY(-10px);
		}
		to {
			opacity: 1;
			transform: translateY(0);
		}
	}

	@keyframes slideOut {
		from {
			opacity: 1;
			transform: translateY(0);
		}
		to {
			opacity: 0;
			transform: translateY(-10px);
		}
	}

	/* Responsive Design */
	@media (max-width: 480px) {
		.notification {
			margin-bottom: 16px;
			padding: 14px 16px;
		}

		.notification h3 {
			font-size: 14px;
		}

		.notification p {
			font-size: 13px;
		}

		.close-btn {
			width: 24px;
			height: 24px;
			font-size: 18px;
		}
	}
</style>
