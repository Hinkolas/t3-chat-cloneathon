<script>
    let { form } = $props();

    let email = $state('');
    let password = $state('');

    let isDisabled = $derived(!email && !password);
</script>

<div class="container">
    <form method="POST" action="/login">
        <h1> Kamino Chat</h1>
        <h2>Sign in to your account</h2>

        <label for="email">Email</label>
        <input
            type="email"
            name="email"
            bind:value={email}
            placeholder="Email"
        />

        <label for="password">Password</label>
        <input
            type="password"
            name="password"
            bind:value={password}
            placeholder="Password"
        />

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

	img {
		height: 50px;
		margin: 16px;
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

	::placeholder {
		color: var(--placeholder);
	}

    :focus {
        border-color: var(--primary-background);
    }

	button {
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
