<script>
	import Notification from '$lib/components/Notification.svelte';
    import { notificationState } from '$lib/store';
	import { scale } from 'svelte/transition';
</script>

<slot />
{#if $notificationState.show}
    <div class="notification-container" transition:scale={{ duration: 100, start: 0.95 }}>
        <Notification />
    </div>
{/if}

<style>
    .notification-container {
        position: fixed;
        top: 20px;
        right: 20px;
        z-index: 1000;
        width: 100%;
        max-width: 400px;
    }
</style>
