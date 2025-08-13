<script lang="ts">
	import AppSidebar from './app-sidebar.svelte';
	import AppHeader from './app-header.svelte';
	import { authStore } from '$lib/stores/auth.svelte.js';
	
	let { children } = $props();
</script>

{#if authStore.isAuthenticated}
	<div class="flex min-h-screen bg-background">
		<!-- Sidebar -->
		<AppSidebar />
		
		<!-- Main content area -->
		<div class="flex flex-1 flex-col min-h-screen">
			<!-- Header -->
			<AppHeader />
			
			<!-- Page content -->
			<main class="flex-1 bg-background p-4 md:p-6">
				{@render children?.()}
			</main>
		</div>
	</div>
{:else}
	<!-- Authentication layout -->
	<div class="flex min-h-screen w-full bg-background">
		{@render children?.()}
	</div>
{/if}
