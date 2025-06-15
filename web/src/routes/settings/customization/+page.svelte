<script>
	let form = $state({
		name: '',
		job: '',
		traits: '',
		information: ''
	});

	let initialForm = {
		name: '',
		job: '',
		traits: '',
		information: ''
	};

	let hasChanged = $derived(
		form.name !== initialForm.name ||
			form.job !== initialForm.job ||
			form.traits !== initialForm.traits ||
			form.information !== initialForm.information
	);
</script>

<div class="container">
	<div class="title">Customize Your Chat</div>
	<form action="">
		<div class="form-group">
			<label for="name">What should Your Chat call you?</label>
			<div class="input-container">
				<input bind:value={form.name} type="text" placeholder="Enter your name" />
				<div class="chars">{form.name.length}/50</div>
			</div>
		</div>
		<div class="form-group">
			<label for="name">What do you do?</label>
			<div class="input-container">
				<input bind:value={form.job} type="text" placeholder="Engineer, student, etc." />
				<div class="chars">{form.job.length}/100</div>
			</div>
		</div>
		<div class="form-group">
			<label for="name">What traits should Your Chat have?</label>
			<div class="input-container">
				<input
					bind:value={form.traits}
					type="text"
					placeholder="Type a trait and press Enter or Tab"
				/>
				<div class="chars">{form.traits.length}/50</div>
			</div>
		</div>
		<div class="form-group">
			<label for="name">Anything else Your Chat should know about you?</label>
			<div class="input-container">
				<input
					bind:value={form.information}
					type="text"
					placeholder="Interests, values, or preferences to keep in mind"
				/>
				<div class="chars">{form.information.length}/3000</div>
			</div>
		</div>
		<button disabled={!hasChanged}>Save Preferences</button>
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
		color: hsl(var(--secondary-foreground));
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
		color: hsl(var(--secondary-foreground));
		outline: none;
	}

	.form-group .input-container input:focus {
		box-shadow: 0 0 2px hsl(var(--primary));
	}

	.form-group .input-container input::placeholder {
		color: hsl(var(--secondary-foreground) / 0.6);
	}

	.form-group .input-container .chars {
		position: absolute;
		bottom: 0;
		right: 6px;
		font-size: 11px;
		color: hsl(var(--secondary-foreground) / 0.6);
	}
</style>
