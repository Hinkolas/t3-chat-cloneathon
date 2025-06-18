<script>
	import { Eye, EyeOff } from '@lucide/svelte';
	import { enhance } from '$app/forms';
	import { notificationState, showNotification, hideNotification } from '$lib/store';
	import Notification from '$lib/components/Notification.svelte';

	let username = $state('');
	let password = $state('');
	let showPassword = $state(false);
	let isSubmitting = $state(false);
	let isDisabled = $derived(!username || !password || isSubmitting);

	function togglePasswordVisibility() {
		showPassword = !showPassword;
	}

	function clear() {
		password = '';
		showPassword = false;
		isSubmitting = false;
	}

	$effect(() => {
		if ($notificationState.show) {
			const timer = setTimeout(() => {
			hideNotification();
			}, 5000);
			return () => clearTimeout(timer);
		}
	});
</script>

<div class="container">
	<div class="form-wrapper">
		<form
			method="POST"
			action="/auth/login"
			use:enhance={() => {
				isSubmitting = true;

				return async ({ result }) => {
					isSubmitting = false;

					if (result.type === 'failure') {
						// Handle different error cases
						if (result.status === 401) {
							showNotification({
								is_error: true,
								title: 'Login Failed',
								description: 'Invalid username or password. Please try again.'
							});
							clear();
							return;
						} else if (result.status === 429) {
							showNotification({
								is_error: true,
								title: 'Too Many Attempts',
								description: 'Too many login attempts. Please try again later.'
							});
							clear();
							return;
						} else {
							showNotification({
								is_error: true,
								title: 'Error',
								description: 'Login failed. Please try again.'
							});
							clear();
							return;
						}
					}
					// If success, the redirect will happen automatically
					window.location.href = '/';
				};
			}}
		>
			<h1>Kamino Chat</h1>
			<h2>Sign in to your account</h2>

			{#if $notificationState.show}
				<Notification />
			{/if}

			<div class="form-group">
				<label for="username">Username</label>
				<input
					type="text"
					name="username"
					bind:value={username}
					placeholder="Username"
					disabled={isSubmitting}
				/>
			</div>

			<div class="form-group">
				<label for="password">Password</label>
				<div class="password-input-container">
					<input
						type={showPassword ? 'text' : 'password'}
						name="password"
						bind:value={password}
						placeholder="Password"
						disabled={isSubmitting}
					/>
					<button
						type="button"
						class="password-toggle"
						onclick={togglePasswordVisibility}
						aria-label={showPassword ? 'Hide password' : 'Show password'}
						disabled={isSubmitting}
					>
						{#if showPassword}
							<EyeOff size={16} />
						{:else}
							<Eye size={16} />
						{/if}
					</button>
				</div>
			</div>

			<button type="submit" disabled={isDisabled}>
				{#if isSubmitting}
					Signing In...
				{:else}
					Sign In
				{/if}
			</button>
		</form>
	</div>
</div>

<style>
	.container {
		display: flex;
		justify-content: center;
		align-items: center;
		padding: 16px;
		width: 100%;
		height: 100dvh;
	}

	@media (max-width: 768px) {
		.container {
			padding: 0px;
		}
	}

	form {
		display: flex;
		flex-direction: column;
		align-items: start;
		justify-content: center;
		width: 100%;
		max-width: 450px;
	}

	.form-wrapper {
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		width: 100%;
		max-width: 450px;

		background-color: var(--sidebar-background);
		padding: 32px;
		border: 1px solid var(--primary-border);
		border-radius: 8px;

	}

	@media (max-width: 768px) {
		.form-wrapper {
			padding: 16px;
			max-width: 100%;
			height: 100%;
			border-radius: 0;
			border: none;
		}
		form {
			border: none;
		}
	}

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

	.form-group {
		width: 100%;
		display: flex;
		flex-direction: column;
		justify-content: center;
	}

	.password-input-container {
		position: relative;
		display: flex;
		align-items: center;
	}

	.password-toggle {
		width: 24px;
		height: 24px;
		position: absolute;
		bottom: 50%;
		transform: translateY(70%);
		right: 10px;
		background: none;
		border: none;
		cursor: pointer;
		color: var(--text);
		display: flex;
		align-items: center;
		justify-content: center;
		padding: 4px;
		border-radius: 4px;
		transition: background-color 0.2s ease;
	}

	.password-toggle:hover:not(:disabled) {
		background-color: var(--primary-border);
	}

	.password-toggle:disabled {
		opacity: 0.5;
		cursor: not-allowed;
	}

	label {
		margin-top: 20px;
		font-size: 14px;
		font-weight: 500;
		color: var(--text);
	}

	input {
		width: 100%;
		height: 40px;
		color: var(--text);
		background-color: var(--sidebar-background);
		border: 1px solid var(--primary-border);
		border-radius: 4px;
		padding: 0 10px;
		margin-top: 8px;
		font-size: 14px;
		outline: none;
	}

	.password-input-container input {
		padding-right: 45px;
	}

	::placeholder {
		color: var(--placeholder);
	}

	:focus {
		border-color: var(--primary-background);
	}

	:disabled {
		background-color: var(--primary-background-light);
		color: var(--placeholder);
		cursor: not-allowed;
	}

	button[type='submit'] {
		margin-top: 30px;
		width: 100%;
		background-color: var(--primary-background-light);
		color: var(--text);
		padding: 12px 0;
		font-size: 18px;
		font-weight: 500;
		border-radius: 4px;
		border: none;
		cursor: pointer;
	}

	button[type='submit']:disabled {
		opacity: 0.6;
		cursor: not-allowed;
	}

	h1 {
		font-size: 32px;
		font-weight: 600;
		margin: 0;
		color: var(--text);
	}

	h2 {
		font-size: 16px;
		font-weight: 400;
		margin-top: 4px;
		color: var(--text);
	}

	@media (max-width: 768px) {
		.form-wrapper {
			padding: 16px;
			max-width: 100%;
			height: 100%;
			border-radius: 0;
			border: none;
		}
		label {
			margin-top: 20px;
		}
		input {
			margin-top: 2px;
		}
	}
</style>
