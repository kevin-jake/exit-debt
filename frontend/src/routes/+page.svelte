<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { authStore } from '$lib/stores/auth.svelte.js';
	import { APP_CONFIG } from '$lib/config/api.js';

	onMount(() => {
		if (authStore.isAuthenticated) {
			goto('/dashboard');
		} else {
			goto('/auth/login');
		}
	});
</script>

<svelte:head>
	<title>{APP_CONFIG.NAME} - Take control of your debts</title>
	<meta name="description" content="Track your debts and money lent to others with Exit-Debt" />
</svelte:head>

<div class="flex min-h-screen items-center justify-center">
	<div class="text-center">
		<div class="mx-auto h-16 w-16 rounded-lg bg-primary flex items-center justify-center">
			<span class="text-primary-foreground font-bold text-xl">ED</span>
		</div>
		<h1 class="mt-4 text-4xl font-bold">{APP_CONFIG.NAME}</h1>
		<p class="mt-2 text-lg text-muted-foreground">Loading...</p>
	</div>
</div>
