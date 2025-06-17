<script>
</script>

<div class="container">
	<div class="account">
		<div class="img-wrapper">
			<img src="https://placehold.co/120" alt="Profile" />
		</div>
		<div class="infos">
			<div class="username">Ertu K.</div>
			<div class="email">ertugul.kilic@icloud.com</div>
			<div class="subscription">Free Plan</div>
		</div>
	</div>
	<div class="card">
		<div class="head">
			<div class="title">Message usage</div>
			<div class="content">Resets tomorrow at 2:00 AM</div>
		</div>
		<div class="standart">
			<div class="head">
				<div class="title">Standart</div>
				<div class="value">1/20</div>
			</div>
			<div class="body">
				<!-- loader width percentage -->
				<div class="slider"></div>
			</div>
			<div class="tail">
				<div class="info">19 messages remaining</div>
			</div>
		</div>
	</div>
</div>

<style>
	.container {
		width: 100%;
		max-width: 270px;
		display: flex;
		flex-direction: column;
		gap: 24px;
	}

	@media (max-width: 1024px) {
		.container {
			display: none;
		}
	}

	.account {
		display: flex;
		flex-direction: column;
		justify-content: center;
		align-items: center;
		gap: 20px;
	}

	.account .img-wrapper {
		width: 120px;
		height: 120px;
	}

	.account .img-wrapper img {
		width: 100%;
		height: 100%;
		border-radius: 50%;
		object-fit: cover;
	}

	.account .infos {
		display: flex;
		flex-direction: column;
		justify-content: center;
		align-items: center;
		gap: 4px;
	}

	.account .infos .username {
		font-size: 24px;
		font-weight: 700;
	}

	.account .infos .email {
		font-size: 14px;
		font-weight: 500;
		color: hsl(var(--secondary-foreground));
	}

	.account .infos .subscription {
		font-size: 12px;
		font-weight: 500;
		color: hsl(var(--secondary-foreground));
		padding: 2px 14px;
		border-radius: 9999px;
		background-color: #88888833;
	}

	.card {
		display: flex;
		flex-direction: column;
		gap: 16px;
		padding: 16px;
		border-radius: 8px;
		background-color: #0b080b;
	}

	.card .head {
		display: flex;
		justify-content: space-between;
		/* flex-wrap: wrap; */
	}

	.card .head .title {
		font-size: 12px;
		font-weight: 700;
	}

	.card .head .content {
		font-size: 12px;
		font-weight: 500;
		color: hsl(var(--secondary-foreground));
	}

	.card .standart {
		display: flex;
		flex-direction: column;
		gap: 4px;
	}

	.card .standart .head {
		display: flex;
		justify-content: space-between;
	}

	.card .standart .head .title {
		font-size: 12px;
		font-weight: 600;
	}

	.card .standart .head .value {
		font-size: 12px;
		font-weight: 500;
		color: hsl(var(--secondary-foreground));
	}

	.card .standart .body {
		position: relative;
		width: 100%;
		height: 8px;
		border-radius: 9999px;
		background-color: #362e3c;
	}

	.card .standart .body .slider {
		position: relative;
		width: 5%;
		height: 8px;
		border-radius: 9999px;
		background-color: hsl(var(--primary));
	}

	.card .standart .tail .info {
		font-size: 12px;
		font-weight: 500;
		color: hsl(var(--secondary-foreground));
	}
</style>
