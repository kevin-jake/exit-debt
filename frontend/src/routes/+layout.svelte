<script lang="ts">
	import '../app.css';
	import favicon from '$lib/assets/favicon.svg';
	import { themeStore } from '$lib/stores/theme.svelte.js';
	import AppLayout from '$lib/components/layout/app-layout.svelte';
	
	let { children } = $props();
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
	<meta name="theme-color" content={themeStore.isDark ? '#0f172a' : '#ffffff'} />
</svelte:head>

<AppLayout>
	{@render children?.()}
</AppLayout>
