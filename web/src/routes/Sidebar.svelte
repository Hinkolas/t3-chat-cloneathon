<script lang="ts">
	let { sidebarCollapsed, newChat } = $props();

	import SearchInput from './SearchInput.svelte';
	import { X } from '@lucide/svelte';

	let chatSearchTerm = $state('');

	function chatSearchFilter() {}
</script>

<div class="sidebar {sidebarCollapsed ? 'collapsed' : ''}">
	<div class="head">
		<div class="title">Chat</div>
		<div class="newChatButton">
			<button onclick={newChat}>New Chat</button>
		</div>
		<div class="search-container">
			<SearchInput bind:value={chatSearchTerm} onInputFunction={chatSearchFilter} placeholder="Search your threads..."/>
		</div>
		<div class="chats-container">
			<div class="day-title">Today</div>
			<div class="chats">
				<div class="chat">
					<span>Greeting Title</span>
					<div class="buttons">
						<button>
							<X size="14" />
						</button>
					</div>
				</div>
                <div class="chat">
					<span>Greeting Title</span>
					<div class="buttons">
						<button>
							<X size="14" />
						</button>
					</div>
				</div>
			</div>
		</div>
	</div>
	<div class="foot">Login</div>
</div>

<style>
	.sidebar {
		flex: 0 0 256px;
		padding-block: 20px;
		background-color: var(--sidebar-background);
		transition: margin 0.15s ease-in-out;

		display: flex;
		flex-direction: column;
		justify-content: space-between;
		align-items: center;
	}

	.sidebar.collapsed {
		margin-left: -256px;
	}

	.head {
		width: 100%;
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 16px;
	}

	.title {
		color: hsl(var(--secondary-foreground));
		font-size: 18px;
		font-weight: 700;
		text-shadow: 0 0 4px hsl(var(--primary) / 0.3);
	}

	.newChatButton {
		width: 100%;
		padding-inline: 16px;
	}

	.newChatButton button {
		all: unset;
		box-sizing: border-box;
		padding: 8px 16px;
		width: 100%;
		display: flex;
		justify-content: center;
		align-items: center;
		background-color: hsl(var(--primary) / 0.2);
		box-shadow: 0px 0px 2px hsl(var(--primary));
		border-radius: 8px;
		cursor: pointer;
		font-size: 14px;
		font-weight: 600;
		line-height: 20px;
		color: hsl(var(--secondary-foreground));
		text-shadow: 0 0 4px hsl(var(--primary) / 0.5);
		transition: background-color 0.15s ease;
	}

	.newChatButton:hover button {
		background-color: hsl(var(--primary) / 0.8);
	}

	.search-container {
		padding-inline: 2px;
	}

	.chats-container {
		width: 100%;
		display: flex;
		flex-direction: column;

		padding-inline: 16px;
	}

	.day-title {
		font-size: 12px;
		font-weight: 600;
		color: var(--secondary);
		text-shadow: 0 0 4px hsl(var(--primary) / 0.2);
	}

	.chats {
		display: flex;
		flex-direction: column;
		gap: 4px;
	}

	.chat {
		position: relative;
		padding: 8px 8px;
		border-radius: 8px;
		cursor: pointer;
		transition: background-color 0.15s ease-out;
		overflow: hidden;
	}

	.chat span {
		color: hsl(var(--secondary-foreground));
		font-size: 14px;
		font-weight: 500;
	}

	.chat:hover {
		background-color: var(--sidebar-chat-hover);
	}

	.buttons {
		position: absolute;
		top: 50%;
		left: 100%;
		transform: translateY(-50%);
		width: 100%;

		padding-right: 8px;
		display: flex;
		justify-content: flex-end;
		transition: left 0.15s ease;
	}

	.buttons button {
		display: flex;
		justify-content: center;
		align-items: center;

		background: none;
		border: none;
		border-radius: 4px;
		padding: 4px;
		cursor: pointer;
		color: hsl(var(--secondary-foreground));
		transition: background-color 0.15s ease-out;
	}

	.buttons button:hover {
		background-color: hsl(var(--primary) / 0.5);
	}

	.chat:hover .buttons {
		left: 0;
	}
</style>
