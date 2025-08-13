<script lang="ts">
	import { page } from '$app/stores';
	import { cn } from '$lib/utils.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { 
		LayoutDashboard, 
		CreditCard, 
		Users, 
		Settings, 
		PlusCircle,
		TrendingUp,
		Bell
	} from 'lucide-svelte';
	import { APP_CONFIG } from '$lib/config/api.js';

	const navigation = [
		{
			name: 'Dashboard',
			href: '/dashboard',
			icon: LayoutDashboard
		},
		{
			name: 'My Debts',
			href: '/debts',
			icon: CreditCard
		},
		{
			name: 'Owed to Me',
			href: '/credits',
			icon: TrendingUp
		},
		{
			name: 'Contacts',
			href: '/contacts',
			icon: Users
		},
		{
			name: 'Notifications',
			href: '/notifications',
			icon: Bell
		},
		{
			name: 'Settings',
			href: '/settings',
			icon: Settings
		}
	];

	function isActiveRoute(href: string): boolean {
		return $page.url.pathname === href || $page.url.pathname.startsWith(href + '/');
	}
</script>

<aside class="hidden w-64 flex-col border-r bg-sidebar md:flex">
	<!-- Logo/Brand -->
	<div class="border-b p-6">
		<div class="flex items-center space-x-2">
			<div class="h-8 w-8 rounded-lg bg-primary flex items-center justify-center">
				<span class="text-primary-foreground font-bold text-sm">ED</span>
			</div>
			<span class="text-xl font-bold text-sidebar-foreground">{APP_CONFIG.NAME}</span>
		</div>
	</div>

	<!-- Navigation -->
	<nav class="flex-1 space-y-1 p-4">
		{#each navigation as item}
			<a
				href={item.href}
				class={cn(
					"flex items-center space-x-3 rounded-lg px-3 py-2 text-sm font-medium transition-colors",
					isActiveRoute(item.href)
						? "bg-sidebar-accent text-sidebar-accent-foreground"
						: "text-sidebar-foreground hover:bg-sidebar-accent hover:text-sidebar-accent-foreground"
				)}
			>
				<svelte:component this={item.icon} class="h-4 w-4" />
				<span>{item.name}</span>
			</a>
		{/each}
	</nav>

	<!-- Quick Actions -->
	<div class="border-t p-4">
		<Button href="/debts/new" class="w-full justify-start" variant="default">
			<PlusCircle class="mr-2 h-4 w-4" />
			Add New Debt
		</Button>
	</div>
</aside>

<!-- Mobile sidebar overlay (will be implemented later) -->
<div class="md:hidden">
	<!-- Mobile navigation will be added here -->
</div>
