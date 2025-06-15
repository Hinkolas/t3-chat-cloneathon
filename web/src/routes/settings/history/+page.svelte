<script>
	import { Upload, Trash2, Download, Pin } from '@lucide/svelte';
</script>

<div class="container">
	<div class="head">
		<div class="title">Message History</div>
		<div class="description">
			Save your history as JSON, or import someone else's. Importing will NOT delete existing
			messages
		</div>
		<div class="history-container">
			<div class="header">
				<div class="head">
					<button>
						<div class="square"></div>
						Select All
					</button>
				</div>
				<div class="tail">
					<button disabled>
						<Upload size="14" />
						Export
					</button>
					<button class="delete-button" disabled>
						<Trash2 size="14" />
						Delete
					</button>
					<button>
						<Download size="14" />
						Import
					</button>
				</div>
			</div>
			<div class="history-box">
				<div class="histories">
					{#each Array.from({ length: 20 })}
						<div class="history">
							<div class="head">
								<div class="square"></div>
								<div class="title">Lucide Svelte ArrowUp Component Question</div>
							</div>
							<div class="tail">
								<Pin color="#9E174D" size="16" />
								<div class="date">6/14/2025, 9:26:07 PM</div>
							</div>
						</div>
					{/each}
				</div>
			</div>
		</div>
	</div>
	<div class="tail">
		<div class="title">Danger Zone</div>
		<div class="description">Permanently delete your account and all associated data.</div>
		<button class="main-delete-button">Delete Chat History</button>
	</div>
</div>

<style>
	.container {
		display: flex;
		flex-direction: column;
		gap: 48px;
	}

	.container > .head,
	.container > .tail {
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
		color: hsl(var(--secondary-foreground));
	}

	.history-container {
		display: flex;
		flex-direction: column;
		gap: 8px;
	}

	.header {
		display: flex;
		justify-content: space-between;
	}

	.header .head,
	.header .tail {
		display: flex;
		gap: 8px;
	}

	.header button {
		all: unset;
		display: flex;
		justify-content: center;
		align-items: center;
		gap: 8px;
		font-size: 12px;
		font-weight: 500;
		padding: 4px 16px;
		border-radius: 5px;
		color: hsl(var(--secondary-foreground));
		border: 1px solid #302029;
		background-color: #21141e;
		cursor: pointer;
		transition: background-color 0.15s ease-out;
	}

	.header button:hover:not(:disabled) {
		background-color: #261923;
	}

	.header button:disabled {
		color: hsl(var(--secondary-foreground) / 0.6);
		cursor: not-allowed;
	}

	.header button.delete-button {
		border: 1px solid #9e174d;
		background-color: #9e174d;
	}

	.header button.delete-button:disabled {
		border: 1px solid #621734;
		background-color: #621734;
	}

	.history-box {
		display: flex;
	}

	.histories {
		width: 100%;
		display: flex;
		flex-direction: column;
		border: 1px solid #26242d;
		border-radius: 5px;
		max-height: 240px;
		overflow-y: auto;
	}

	.history {
		display: flex;
		justify-content: space-between;
		padding: 8px 16px;
		cursor: pointer;
		transition: background-color 0.15s ease-out;
	}

	.history:hover {
		background-color: #2e272f;
	}

	.history:not(:last-child) {
		border-bottom: 1px solid #26242d;
	}

	.history .head {
		display: flex;
		align-items: center;
		gap: 10px;
	}

	.square {
		width: 15px;
		height: 15px;
		border-radius: 4px;
		border: 1px solid #614052;
	}

	.history .head .title {
		font-size: 14px;
		font-weight: 500;
		color: #f7f5f8;
	}

	.history .tail {
		display: flex;
		align-items: center;
		gap: 16px;
	}

	.history .tail .date {
		font-size: 12px;
		color: hsl(var(--secondary-foreground));
	}

	.main-delete-button {
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

	.main-delete-button:hover {
		background-color: var(--button-hover-danger);
	}
</style>
