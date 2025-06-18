<script>
	import { Eye, EyeOff } from "@lucide/svelte";

	let username = $state('');
	let password = $state('');
	let showPassword = $state(false);

	let isDisabled = $derived(!username || !password);

	function togglePasswordVisibility() {
		showPassword = !showPassword;
	}
</script>

<div class="container">
	<form method="POST" action="/auth/login">
		<h1>Kamino Chat</h1>
		<h2>Sign in to your account</h2>

		<div class="form-group">
			<label for="username">Username</label>
			<input type="text" name="username" bind:value={username} placeholder="Username" />
		</div>

		<div class="form-group">
			<label for="password">Password</label>
			<div class="password-input-container">
				<input 
					type={showPassword ? "text" : "password"} 
					name="password" 
					bind:value={password} 
					placeholder="Password" 
				/>
				<button 
					type="button" 
					class="password-toggle" 
					onclick={togglePasswordVisibility}
					aria-label={showPassword ? "Hide password" : "Show password"}
				>
					{#if showPassword}
						<EyeOff size={16} />
					{:else}
						<Eye size={16} />
					{/if}
				</button>
			</div>
		</div>

		<button type="submit" disabled={isDisabled}>Sign In</button>
	</form>
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

	form {
		display: flex;
		flex-direction: column;
		align-items: start;
		justify-content: center;
		width: fit-content;

		background-color: var(--sidebar-background);
		padding: 32px;
		border: 1px solid var(--primary-border);
		border-radius: 8px;
	}

	.form-group {
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

	.password-toggle:hover {
		background-color: var(--primary-border);
	}

	label {
		margin-top: 20px;
		font-size: 14px;
		font-weight: 500;
		color: var(--text);
	}

	input {
		height: 40px;
		width: 384px;
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
		padding-right: 45px; /* Make space for the toggle button */
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

	button[type="submit"] {
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
</style>
