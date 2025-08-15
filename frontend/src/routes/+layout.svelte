<script lang="ts">
	import '../app.css';
	import { page } from '$app/stores';
	import { onMount } from 'svelte';
	import Navigation from '$lib/components/Navigation.svelte';
	import { getTheme, setTheme } from '$lib/utils';

	let { children } = $props();

	// Check if current page should show navigation
	const showNavigation = $derived(!$page.url.pathname.startsWith('/login') && !$page.url.pathname.startsWith('/register'));

	onMount(() => {
		// Initialize theme
		const theme = getTheme();
		setTheme(theme);
	});
</script>

{#if showNavigation}
	<div class="min-h-screen bg-background">
		<Navigation />
		<main class="ml-64">
			<div class="p-6">
				{@render children()}
			</div>
		</main>
	</div>
{:else}
	<div class="min-h-screen bg-background">
		{@render children()}
	</div>
{/if}
