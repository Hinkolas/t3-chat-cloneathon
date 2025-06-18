<script>
</script>

<div class="container">
	<div class="title">Danger Zone</div>
	<div class="description">Permanently delete your account and all associated data.</div>
	<button> Delete Account </button>
</div>

<style>
	.container {
		display: flex;
		flex-direction: column;
		gap: 12px;
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

	button {
		all: unset;
		width: max-content;
		font-size: 14px;
		font-weight: 600;
		color: #ffffff;
		padding: 8px 16px;
		cursor: pointer;
		border-radius: 8px;
		border: 1px solid var(--button-border-danger);
		background-color: var(--button-background-danger);
		transition: background-color 0.15s ease-out;
	}

	button:hover {
		background-color: var(--button-hover-danger);
	}
</style>
