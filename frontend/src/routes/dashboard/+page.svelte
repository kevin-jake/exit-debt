<script lang="ts">
	import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '$lib/components/ui/card/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { authStore } from '$lib/stores/auth.svelte.js';
	import { APP_CONFIG } from '$lib/config/api.js';
	import { PlusCircle, TrendingDown, TrendingUp, Users } from 'lucide-svelte';

	// Mock data for now
	const stats = {
		totalOwed: 2450.00,
		totalOwing: 1200.00,
		totalContacts: 12,
		recentDebts: 3
	};
</script>

<svelte:head>
	<title>Dashboard - {APP_CONFIG.NAME}</title>
	<meta name="description" content="Your debt tracking dashboard" />
</svelte:head>

<div class="space-y-8">
	<!-- Welcome Header -->
	<div>
		<h1 class="text-3xl font-bold text-foreground">
			Welcome back, {authStore.user?.firstName}!
		</h1>
		<p class="text-muted-foreground">
			Here's an overview of your financial situation.
		</p>
	</div>

	<!-- Stats Grid -->
	<div class="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
		<Card>
			<CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
				<CardTitle class="text-sm font-medium">Total Owed by Me</CardTitle>
				<TrendingDown class="h-4 w-4 text-muted-foreground" />
			</CardHeader>
			<CardContent>
				<div class="text-2xl font-bold text-destructive">
					${stats.totalOwed.toLocaleString()}
				</div>
				<p class="text-xs text-muted-foreground">
					Money you owe to others
				</p>
			</CardContent>
		</Card>

		<Card>
			<CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
				<CardTitle class="text-sm font-medium">Total Owed to Me</CardTitle>
				<TrendingUp class="h-4 w-4 text-muted-foreground" />
			</CardHeader>
			<CardContent>
				<div class="text-2xl font-bold text-green-600">
					${stats.totalOwing.toLocaleString()}
				</div>
				<p class="text-xs text-muted-foreground">
					Money others owe you
				</p>
			</CardContent>
		</Card>

		<Card>
			<CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
				<CardTitle class="text-sm font-medium">Net Position</CardTitle>
				<TrendingUp class="h-4 w-4 text-muted-foreground" />
			</CardHeader>
			<CardContent>
				<div class="text-2xl font-bold {stats.totalOwing - stats.totalOwed >= 0 ? 'text-green-600' : 'text-destructive'}">
					${Math.abs(stats.totalOwing - stats.totalOwed).toLocaleString()}
				</div>
				<p class="text-xs text-muted-foreground">
					{stats.totalOwing - stats.totalOwed >= 0 ? 'You are owed more' : 'You owe more'}
				</p>
			</CardContent>
		</Card>

		<Card>
			<CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
				<CardTitle class="text-sm font-medium">Contacts</CardTitle>
				<Users class="h-4 w-4 text-muted-foreground" />
			</CardHeader>
			<CardContent>
				<div class="text-2xl font-bold">{stats.totalContacts}</div>
				<p class="text-xs text-muted-foreground">
					People you track debts with
				</p>
			</CardContent>
		</Card>
	</div>

	<!-- Quick Actions -->
	<div class="grid gap-4 md:grid-cols-2">
		<Card>
			<CardHeader>
				<CardTitle>Quick Actions</CardTitle>
				<CardDescription>
					Common tasks to manage your debts
				</CardDescription>
			</CardHeader>
			<CardContent class="space-y-2">
				<Button href="/debts/new" class="w-full justify-start">
					<PlusCircle class="mr-2 h-4 w-4" />
					Add New Debt
				</Button>
				<Button href="/contacts/new" variant="outline" class="w-full justify-start">
					<Users class="mr-2 h-4 w-4" />
					Add Contact
				</Button>
			</CardContent>
		</Card>

		<Card>
			<CardHeader>
				<CardTitle>Recent Activity</CardTitle>
				<CardDescription>
					Your latest debt transactions
				</CardDescription>
			</CardHeader>
			<CardContent>
				<div class="text-sm text-muted-foreground">
					No recent activity to show.
				</div>
				<Button href="/debts" variant="outline" class="mt-4 w-full">
					View All Debts
				</Button>
			</CardContent>
		</Card>
	</div>
</div>
