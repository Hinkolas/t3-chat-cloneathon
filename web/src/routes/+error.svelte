<script>
	import { Skull, Zap } from '@lucide/svelte';
	import { page } from '$app/state';

	const errorMessages = [
		'Oops! The hamsters powering our servers took a coffee break ☕',
		'404: Brain not found. Please try turning it off and on again 🧠',
		"Error 418: I'm a teapot... wait, that's not right 🫖",
		'Something went wrong! Even our error messages are confused 🤔',
		'The models escaped! Send help (and snacks) 🍪',
		'Critical failure: Too much awesome detected ⚡',
		'Whoopsie! Our AI is having an existential crisis 🤖'
	];

	const randomError = errorMessages[Math.floor(Math.random() * errorMessages.length)];
</script>

<div class="error-container">
	<div class="error-detail">
		<div class="status">{page.status}</div>
		<div class="message">{page.error?.message}</div>
	</div>
	<div class="error-icon">
		<Skull size="120" />
	</div>
	<div class="error-content">
		<h2 class="error-title">Well, this is embarrassing...</h2>
		<p class="error-message">{randomError}</p>
		<div class="error-details">
			<p>What happened?</p>
			<ul>
				<li>The models decided to play hide and seek</li>
				<li>Our servers are having a moment</li>
				<li>The internet gremlins struck again</li>
			</ul>
		</div>
		<button class="retry-button" onclick={() => window.location.reload()}>
			<Zap size="16" />
			Try Again (Pretty Please?)
		</button>
	</div>
</div>

<style>
	.error-container {
		position: relative;
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		width: 100%;
		height: 100dvh;
		color: white;
		text-align: center;
		padding: 2rem;
		animation: fadeIn 0.5s ease-out;
	}

	.error-detail {
		position: absolute;
		bottom: 20px;
		left: 50%;
		transform: translateX(-50%);
		color: #88888866;
	}

	.error-icon {
		margin-bottom: 2rem;
		animation: bounce 2s infinite;
		filter: drop-shadow(0 4px 8px rgba(0, 0, 0, 0.3));
	}

	.error-content {
		max-width: 600px;
		display: flex;
		flex-direction: column;
		align-items: center;
	}

	.error-title {
		font-size: 2.5rem;
		font-weight: bold;
		margin-bottom: 1rem;
		text-shadow: 2px 2px 4px rgba(0, 0, 0, 0.3);
	}

	.error-message {
		font-size: 1.25rem;
		margin-bottom: 2rem;
		opacity: 0.9;
		line-height: 1.3;
	}

	.error-details {
		width: max-content;
		background: #88888811;
		border-radius: 12px;
		padding: 1.5rem;
		margin-bottom: 2rem;
		backdrop-filter: blur(10px);
		border: 1px solid #88888822;
	}

	.error-details p {
		font-weight: bold;
		margin-bottom: 1rem;
		font-size: 1.1rem;
	}

	.error-details ul {
		list-style: none;
		padding: 0;
		text-align: left;
	}

	.error-details li {
		padding: 0.5rem 0;
		position: relative;
		padding-left: 1.5rem;
	}

	.error-details li::before {
		content: '💥';
		position: absolute;
		left: 0;
	}

	.retry-button {
		background: linear-gradient(45deg, hsl(var(--primary) / 0.5), hsl(var(--primary) / 0.8));
		color: white;
		border: none;
		padding: 1rem 2rem;
		border-radius: 50px;
		font-size: 1rem;
		font-weight: bold;
		cursor: pointer;
		display: flex;
		align-items: center;
		gap: 0.5rem;
		margin: 0 auto;
		box-shadow: 0 4px 15px hsl(var(--primary) / 0.5);
	}

	.retry-button:hover {
		box-shadow: 0 6px 20px hsl(var(--primary) / 0.5);
		background: linear-gradient(45deg, hsl(var(--primary) / 0.8), hsl(var(--primary)));
	}

	.retry-button:active {
		transform: translateY(0);
	}

	@keyframes fadeIn {
		from {
			opacity: 0;
			transform: translateY(20px);
		}
		to {
			opacity: 1;
			transform: translateY(0);
		}
	}

	@keyframes bounce {
		0%,
		20%,
		50%,
		80%,
		100% {
			transform: translateY(0);
		}
		40% {
			transform: translateY(-10px);
		}
		60% {
			transform: translateY(-5px);
		}
	}
</style>
