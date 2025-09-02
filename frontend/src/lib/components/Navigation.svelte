<script lang="ts">
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import { themeStore } from '$lib/stores/theme.svelte';
	import { authStore } from '$lib/stores/auth';
	import { cn } from '$lib/utils';
	import { createEventDispatcher } from 'svelte';

	export let mobile = false;

	const dispatch = createEventDispatcher();

	type NavItem = {
		href: string;
		label: string;
		icon: string;
	};

	const navItems: NavItem[] = [
		{ href: '/', label: 'Dashboard', icon: 'dashboard' },
		{ href: '/debts', label: 'Debts', icon: 'money' },
		{ href: '/contacts', label: 'Contacts', icon: 'people' },
		{ href: '/settings', label: 'Settings', icon: 'settings' }
	];

	$: isActive = (href: string): boolean => {
		// Special case for dashboard: both "/" and "/dashboard" should highlight the dashboard nav item
		if (href === '/') {
			return $page.url.pathname === '/' || $page.url.pathname === '/dashboard';
		}
		return $page.url.pathname === href;
	};

	function handleLogout() {
		authStore.logout();
		goto('/login');
	}

	$: handleThemeToggle = () => {
		themeStore.toggle();
	};

	function handleNavigation(href: string) {
		if (mobile) {
			dispatch('navigate');
		}
		goto(href);
	}

	// Get user initials for avatar - ensure reactivity
	$: userInitials = $authStore.user && $authStore.user.first_name && $authStore.user.last_name
		? `${$authStore.user.first_name[0]}${$authStore.user.last_name[0]}`
		: 'U';
</script>

<nav class={cn(
	"bg-card border-r border-border w-64 h-screen z-40",
	mobile ? "fixed left-0 top-0" : "fixed left-0 top-0"
)}>
	<div class="p-6">
		<!-- Logo (hide on mobile since it's in the header) -->
		<div class="flex items-center space-x-3 mb-8" class:hidden={mobile}>
			<div class="w-8 h-8 bg-primary rounded-lg flex items-center justify-center">
				<svg class="w-5 h-5 text-primary-foreground" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1"></path>
				</svg>
			</div>
			<h1 class="text-xl font-bold text-card-foreground">DebtTracker</h1>
		</div>

		<!-- Mobile Header Spacing -->
		{#if mobile}
			<div class="mb-8"></div>
		{/if}

		<!-- Theme Toggle -->
		<div class="mb-6">
			<button
				on:click={handleThemeToggle}
				class="w-full flex items-center justify-center space-x-2 px-3 py-2 bg-secondary text-secondary-foreground rounded-lg hover:bg-secondary/80 transition-colors duration-200"
			>
				{#if $themeStore === 'dark'}
					<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 3v1m0 16v1m9-9h-1M4 12H3m15.364 6.364l-.707-.707M6.343 6.343l-.707-.707m12.728 0l-.707.707M6.343 17.657l-.707.707M16 12a4 4 0 11-8 0 4 4 0 018 0z"></path>
					</svg>
					<span class="text-sm">Light Mode</span>
				{:else}
					<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20.354 15.354A9 9 0 018.646 3.646 9.003 9.003 0 0012 21a9.003 9.003 0 008.354-5.646z"></path>
					</svg>
					<span class="text-sm">Dark Mode</span>
				{/if}
			</button>
		</div>

		<!-- Navigation Items -->
		<ul class="space-y-2">
			{#each navItems as item (item.href)}
				<li>
					<button
						on:click={() => handleNavigation(item.href)}
						class={cn(
							"flex items-center space-x-3 px-3 py-2 rounded-lg text-sm font-medium transition-colors duration-200 w-full text-left",
							isActive(item.href)
								? 'bg-primary text-primary-foreground'
								: 'text-muted-foreground hover:bg-secondary hover:text-secondary-foreground'
						)}
					>
						{#if item.icon === 'dashboard'}
							<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2H5a2 2 0 00-2-2v0"></path>
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 7l9 6 9-6"></path>
							</svg>
						{:else if item.icon === 'money'}
							<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1"></path>
							</svg>
						{:else if item.icon === 'people'}
							<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197m13.5-9a2.5 2.5 0 11-5 0 2.5 2.5 0 015 0z"></path>
							</svg>
						{:else if item.icon === 'chart'}
							<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z"></path>
							</svg>
						{:else if item.icon === 'settings'}
							<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z"></path>
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"></path>
							</svg>
						{/if}
						<span>{item.label}</span>
					</button>
				</li>
			{/each}
		</ul>
	</div>

	<!-- User Profile -->
	<div class="absolute bottom-0 left-0 right-0 p-6 border-t border-border">
		<div class="flex items-center space-x-3 mb-3">
			<div class="w-8 h-8 bg-primary rounded-full flex items-center justify-center">
				<span class="text-primary-foreground text-sm font-medium">{userInitials}</span>
			</div>
			<div class="flex-1 min-w-0">
				{#if $authStore.user && $authStore.user.first_name && $authStore.user.last_name}
					<p class="text-sm font-medium text-card-foreground truncate">
						{$authStore.user.first_name} {$authStore.user.last_name}
					</p>
					<p class="text-xs text-muted-foreground truncate">{$authStore.user.email}</p>
				{:else if $authStore.isLoading}
					<p class="text-sm font-medium text-card-foreground truncate">User</p>
					<p class="text-xs text-muted-foreground truncate">Loading...</p>
				{:else}
					<p class="text-sm font-medium text-card-foreground truncate">User</p>
					<p class="text-xs text-muted-foreground truncate">Not signed in</p>
				{/if}
			</div>
		</div>
		<button
			on:click={handleLogout}
			class="w-full text-left px-3 py-2 text-sm text-muted-foreground hover:bg-secondary rounded-lg transition-colors duration-200"
		>
			Sign out
		</button>
	</div>
</nav>
