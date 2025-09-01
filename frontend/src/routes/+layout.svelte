<script lang="ts">
	import '../app.css';
	import { page } from '$app/stores';
	import { onMount } from 'svelte';
	import Navigation from '$lib/components/Navigation.svelte';
	import { getTheme, setTheme } from '$lib/utils';

	let { children } = $props();

	// Mobile navigation state
	let isMobileMenuOpen = $state(false);

	// Check if current page should show navigation
	const showNavigation = $derived(!$page.url.pathname.startsWith('/login') && !$page.url.pathname.startsWith('/register'));

	function toggleMobileMenu() {
		isMobileMenuOpen = !isMobileMenuOpen;
	}

	function closeMobileMenu() {
		isMobileMenuOpen = false;
	}

	onMount(() => {
		// Initialize theme
		const theme = getTheme();
		setTheme(theme);

		// Close mobile menu when clicking outside
		function handleClickOutside(event: MouseEvent) {
			const target = event.target as HTMLElement;
			if (isMobileMenuOpen && !target.closest('.mobile-nav') && !target.closest('.hamburger-btn')) {
				closeMobileMenu();
			}
		}

		document.addEventListener('click', handleClickOutside);
		return () => document.removeEventListener('click', handleClickOutside);
	});

	// Close mobile menu when route changes
	$effect(() => {
		$page.url.pathname;
		closeMobileMenu();
	});
</script>

{#if showNavigation}
	<div class="min-h-screen bg-background">
		<!-- Mobile Header (visible on small screens) -->
		<header class="lg:hidden bg-card border-b border-border px-4 py-3 flex items-center justify-between sticky top-0 z-40">
			<div class="flex items-center space-x-3">
				<div class="w-8 h-8 bg-primary rounded-lg flex items-center justify-center">
					<svg class="w-5 h-5 text-primary-foreground" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1"></path>
					</svg>
				</div>
				<h1 class="text-xl font-bold text-card-foreground">DebtTracker</h1>
			</div>
			
			<!-- Hamburger Menu Button -->
			<button 
				class="hamburger-btn p-2 rounded-lg hover:bg-secondary transition-colors"
				on:click={toggleMobileMenu}
				aria-label="Toggle navigation menu"
			>
				<svg class="w-6 h-6 text-foreground" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					{#if isMobileMenuOpen}
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
					{:else}
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16"></path>
					{/if}
				</svg>
			</button>
		</header>

		<!-- Desktop Navigation (hidden on small screens) -->
		<div class="hidden lg:block">
			<Navigation />
		</div>

		<!-- Mobile Navigation Overlay -->
		{#if isMobileMenuOpen}
			<div class="lg:hidden fixed inset-0 z-50">
				<!-- Backdrop -->
				<div class="fixed inset-0 bg-black/50" on:click={closeMobileMenu}></div>
				
				<!-- Mobile Navigation Panel -->
				<div class="mobile-nav fixed left-0 top-0 h-full max-w-[85vw] bg-card transform transition-transform duration-300 ease-in-out">
					<Navigation mobile={true} on:navigate={closeMobileMenu} />
				</div>
			</div>
		{/if}

		<!-- Main Content -->
		<main class="lg:ml-64 min-h-screen">
			<div class="p-4 lg:p-6">
				{@render children()}
			</div>
		</main>
	</div>
{:else}
	<div class="min-h-screen bg-background">
		{@render children()}
	</div>
{/if}
