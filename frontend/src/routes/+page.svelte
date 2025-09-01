<script lang="ts">
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';

	onMount(() => {
		// Redirect to the new dashboard
		goto('/dashboard');
	});
</script>

<svelte:head>
	<title>DebtTracker</title>
</svelte:head>

<div class="flex items-center justify-center min-h-screen">
	<div class="text-center">
		<div class="animate-spin rounded-full h-8 w-8 border-b-2 border-primary mx-auto mb-4"></div>
		<p class="text-muted-foreground">Redirecting to dashboard...</p>
	</div>
</div>
