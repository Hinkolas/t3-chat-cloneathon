<script>
	import { Gem } from "@lucide/svelte";

	let { profile } = $props();
</script>

<div class="container">
	<div class="account">
		<div class="img-wrapper">
			<img src="https://placehold.co/120" alt="Profile" />
		</div>
		<div class="infos">
			<div class="username">{profile.custom_user_name}</div>
			<div class="email">ertugul.kilic@icloud.com</div>
			<div class="subscription">Premium Plan<Gem size=14/></div>
		</div>
	</div>
	<div class="card">
		<div class="head">
			<div class="title">Message usage</div>
		</div>
		<div class="standart">
			<div class="head">
				<div class="title">Standart</div>
				<div class="value">{profile.usage_standard}/{profile.limit_standard}</div>
			</div>
			<div class="body">
				<!-- loader width percentage -->
				<div class="slider" style="width: {(profile.usage_standard/profile.limit_standard) * 100}%"></div>
			</div>
			<div class="tail">
				<div class="info">{profile.limit_standard - profile.usage_standard} messages remaining</div>
			</div>
		</div>
		<div class="standart">
			<div class="head">
				<div class="title">Premium <Gem size=14 /></div>
				<div class="value">{profile.usage_premium}/{profile.limit_premium}</div>
			</div>
			<div class="body">
				<!-- loader width percentage -->
				<div class="slider" style="width: {(profile.usage_premium/profile.limit_premium) * 100}%"></div>
			</div>
			<div class="tail">
				<div class="info">{profile.limit_premium - profile.usage_premium} messages remaining</div>
			</div>
		</div>
	</div>
</div>

<style>
	.container {
		flex: 0 0 270px;
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
		color: var(--text);
	}

	.account .infos .subscription {
		font-size: 12px;
		font-weight: 500;
		color: var(--text);
		padding: 2px 14px;
		border-radius: 9999px;
		background-color: #88888833;
		display: flex;
		align-items: center;
		gap: 6px;
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
		display: flex;
		align-items: center;
		gap: 4px;
		font-size: 12px;
		font-weight: 700;
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
		color: var(--text);
	}

	.card .standart .body {
		overflow: hidden;
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
		background-color: hsl(var(--primary));
	}

	.card .standart .tail .info {
		font-size: 12px;
		font-weight: 500;
		color: var(--text);
	}
</style>
