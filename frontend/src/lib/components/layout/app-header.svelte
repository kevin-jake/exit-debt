<script lang="ts">
	import { Button } from '$lib/components/ui/button/index.js';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu/index.js';
	import ThemeToggle from '$lib/components/ui/theme-toggle.svelte';
	import { authStore } from '$lib/stores/auth.svelte.js';
	import { 
		Menu, 
		Search, 
		Bell, 
		User, 
		Settings, 
		LogOut,
		Plus
	} from 'lucide-svelte';

	function getInitials(firstName?: string, lastName?: string): string {
		if (!firstName && !lastName) return 'U';
		return `${firstName?.[0] || ''}${lastName?.[0] || ''}`.toUpperCase();
	}

	async function handleLogout() {
		await authStore.logout();
	}
</script>

<header class="border-b bg-background px-4 py-3 md:px-6">
	<div class="flex items-center justify-between">
		<!-- Left side - Mobile menu button and search -->
		<div class="flex items-center space-x-4">
			<!-- Mobile menu button -->
			<Button variant="ghost" size="sm" class="md:hidden">
				<Menu class="h-5 w-5" />
				<span class="sr-only">Toggle menu</span>
			</Button>
			
			<!-- Search -->
			<div class="hidden md:flex">
				<div class="relative">
					<Search class="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
					<input
						type="search"
						placeholder="Search debts..."
						class="h-9 w-64 rounded-md border border-input bg-background pl-10 pr-3 text-sm ring-offset-background placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2"
					/>
				</div>
			</div>
		</div>

		<!-- Right side - Actions and user menu -->
		<div class="flex items-center space-x-2">
			<!-- Quick add button -->
			<Button variant="default" size="sm" href="/debts/new" class="hidden md:flex">
				<Plus class="mr-2 h-4 w-4" />
				Add Debt
			</Button>
			
			<!-- Notifications -->
			<Button variant="ghost" size="sm" href="/notifications">
				<Bell class="h-4 w-4" />
				<span class="sr-only">Notifications</span>
			</Button>
			
			<!-- Theme toggle -->
			<ThemeToggle />
			
			<!-- User menu -->
			<DropdownMenu.Root>
				<DropdownMenu.Trigger asChild let:builder>
					<Button builders={[builder]} variant="ghost" size="sm" class="relative h-8 w-8 rounded-full">
						<div class="flex h-8 w-8 items-center justify-center rounded-full bg-primary text-primary-foreground text-xs font-medium">
							{getInitials(authStore.user?.firstName, authStore.user?.lastName)}
						</div>
						<span class="sr-only">User menu</span>
					</Button>
				</DropdownMenu.Trigger>
				<DropdownMenu.Content align="end" class="w-56">
					<DropdownMenu.Label class="font-normal">
						<div class="flex flex-col space-y-1">
							<p class="text-sm font-medium leading-none">
								{authStore.user?.firstName} {authStore.user?.lastName}
							</p>
							<p class="text-xs leading-none text-muted-foreground">
								{authStore.user?.email}
							</p>
						</div>
					</DropdownMenu.Label>
					<DropdownMenu.Separator />
					<DropdownMenu.Item href="/profile" class="cursor-pointer">
						<User class="mr-2 h-4 w-4" />
						<span>Profile</span>
					</DropdownMenu.Item>
					<DropdownMenu.Item href="/settings" class="cursor-pointer">
						<Settings class="mr-2 h-4 w-4" />
						<span>Settings</span>
					</DropdownMenu.Item>
					<DropdownMenu.Separator />
					<DropdownMenu.Item on:click={handleLogout} class="cursor-pointer text-destructive">
						<LogOut class="mr-2 h-4 w-4" />
						<span>Log out</span>
					</DropdownMenu.Item>
				</DropdownMenu.Content>
			</DropdownMenu.Root>
		</div>
	</div>
</header>
