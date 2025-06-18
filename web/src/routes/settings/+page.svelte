<script>
</script>

<div class="container">
	<h2>Welcome to Kamino Chat</h2>

	<p>
		This is our take on T3 Chat, built for the official Cloneathon. We're Nicholas and Ertu, two
		student developers from Germany who took on the challenge.
	</p>

	<p>The name is a nod to the Star Wars clone origin planet -- seemed fitting.</p>
	<p>All props to @theo and his team for the original application. Keep up the great work!</p>
	<br>
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
