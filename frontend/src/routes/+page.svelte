<script lang="ts">
	import { onMount } from 'svelte';

	// Mock data for demonstration
	let totalOwedToMe = $state(15750.00);
	let totalIOwe = $state(8900.00);
	let netPosition = $derived(totalOwedToMe - totalIOwe);
	let overdueCount = $state(3);
	let dueSoonCount = $state(5);

	let recentDebts = $state([
		{
			id: 1,
			contact: 'Alice Johnson',
			type: 'owed_to_me',
			amount: 2500.00,
			dueDate: '2024-01-15',
			status: 'overdue'
		},
		{
			id: 2,
			contact: 'Bob Smith',
			type: 'i_owe',
			amount: 1200.00,
			dueDate: '2024-01-20',
			status: 'due_soon'
		},
		{
			id: 3,
			contact: 'Carol Davis',
			type: 'owed_to_me',
			amount: 800.00,
			dueDate: '2024-01-25',
			status: 'active'
		},
		{
			id: 4,
			contact: 'David Wilson',
			type: 'i_owe',
			amount: 450.00,
			dueDate: '2024-01-18',
			status: 'overdue'
		}
	]);

	let upcomingPayments = $state([
		{
			id: 1,
			contact: 'Emma Brown',
			amount: 500.00,
			dueDate: '2024-01-16',
			type: 'payment_due'
		},
		{
			id: 2,
			contact: 'Frank Miller',
			amount: 750.00,
			dueDate: '2024-01-18',
			type: 'payment_receive'
		},
		{
			id: 3,
			contact: 'Grace Taylor',
			amount: 300.00,
			dueDate: '2024-01-20',
			type: 'payment_due'
		}
	]);

	function formatCurrency(amount: number): string {
		return new Intl.NumberFormat('en-PH', {
			style: 'currency',
			currency: 'PHP'
		}).format(amount);
	}

	function formatDate(dateString: string): string {
		return new Date(dateString).toLocaleDateString('en-US', {
			month: 'short',
			day: 'numeric',
			year: 'numeric'
		});
	}

	function getStatusColor(status: string): string {
		switch (status) {
			case 'overdue':
				return 'text-destructive bg-destructive/10';
			case 'due_soon':
				return 'text-warning bg-warning/10';
			case 'active':
				return 'text-success bg-success/10';
			default:
				return 'text-muted-foreground bg-muted';
		}
	}

	onMount(() => {
		// Here you would fetch real data from your API
		console.log('Dashboard mounted');
	});
</script>

<svelte:head>
	<title>Dashboard - DebtTracker</title>
</svelte:head>

<div class="max-w-7xl mx-auto space-y-6">
	<!-- Header -->
	<div class="flex items-center justify-between">
		<div>
			<h1 class="text-3xl font-bold text-foreground">Dashboard</h1>
			<p class="text-muted-foreground mt-1">Overview of your debt tracking</p>
		</div>
		<div class="flex space-x-3">
			<button class="btn-secondary">
				<svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 10v6m0 0l-3-3m3 3l3-3m2 8H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"></path>
				</svg>
				Export Report
			</button>
			<button class="btn-primary">
				<svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6"></path>
				</svg>
				Add Debt
			</button>
		</div>
	</div>

	<!-- Summary Cards -->
	<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
		<!-- Total Owed to Me -->
		<div class="card p-6">
			<div class="flex items-center justify-between">
				<div>
					<p class="text-sm font-medium text-muted-foreground">Owed to Me</p>
					<p class="text-2xl font-bold text-success">{formatCurrency(totalOwedToMe)}</p>
				</div>
				<div class="w-12 h-12 bg-success/10 rounded-lg flex items-center justify-center">
					<svg class="w-6 h-6 text-success" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 7h8m0 0v8m0-8l-8 8-4-4-6 6"></path>
					</svg>
				</div>
			</div>
		</div>

		<!-- Total I Owe -->
		<div class="card p-6">
			<div class="flex items-center justify-between">
				<div>
					<p class="text-sm font-medium text-muted-foreground">I Owe</p>
					<p class="text-2xl font-bold text-destructive">{formatCurrency(totalIOwe)}</p>
				</div>
				<div class="w-12 h-12 bg-destructive/10 rounded-lg flex items-center justify-center">
					<svg class="w-6 h-6 text-destructive" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 17h8m0 0V9m0 8l-8-8-4 4-6-6"></path>
					</svg>
				</div>
			</div>
		</div>

		<!-- Net Position -->
		<div class="card p-6">
			<div class="flex items-center justify-between">
				<div>
					<p class="text-sm font-medium text-muted-foreground">Net Position</p>
					<p class="text-2xl font-bold {netPosition >= 0 ? 'text-success' : 'text-destructive'}">
						{formatCurrency(Math.abs(netPosition))}
					</p>
					<p class="text-xs text-muted-foreground mt-1">
						{netPosition >= 0 ? 'In your favor' : 'You owe more'}
					</p>
				</div>
				<div class="w-12 h-12 bg-primary/10 rounded-lg flex items-center justify-center">
					<svg class="w-6 h-6 text-primary" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 7h6m0 10v-3m-3 3h.01M9 17h.01M9 14h.01M12 14h.01M15 11h.01M12 11h.01M9 11h.01M7 21h10a2 2 0 002-2V5a2 2 0 00-2-2H7a2 2 0 00-2 2v14a2 2 0 002 2z"></path>
					</svg>
				</div>
			</div>
		</div>

		<!-- Alerts -->
		<div class="card p-6">
			<div class="flex items-center justify-between">
				<div>
					<p class="text-sm font-medium text-muted-foreground">Alerts</p>
					<div class="flex items-center space-x-4 mt-2">
						<div class="text-center">
							<p class="text-lg font-bold text-destructive">{overdueCount}</p>
							<p class="text-xs text-muted-foreground">Overdue</p>
						</div>
						<div class="text-center">
							<p class="text-lg font-bold text-warning">{dueSoonCount}</p>
							<p class="text-xs text-muted-foreground">Due Soon</p>
						</div>
					</div>
				</div>
				<div class="w-12 h-12 bg-warning/10 rounded-lg flex items-center justify-center">
					<svg class="w-6 h-6 text-warning" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.865-.833-2.635 0L4.178 16.5c-.77.833.192 2.5 1.732 2.5z"></path>
					</svg>
				</div>
			</div>
		</div>
	</div>

	<!-- Recent Debts & Upcoming Payments -->
	<div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
		<!-- Recent Debts -->
		<div class="card p-6">
			<div class="flex items-center justify-between mb-6">
				<h2 class="text-lg font-semibold text-foreground">Recent Debts</h2>
				<a href="/debts" class="text-sm text-primary hover:text-primary/80 font-medium">
					View all
				</a>
			</div>
			<div class="space-y-4">
				{#each recentDebts as debt (debt.id)}
					<div class="flex items-center justify-between p-4 bg-muted/50 rounded-lg">
						<div class="flex-1">
							<div class="flex items-center space-x-3">
								<div class="w-8 h-8 bg-primary rounded-full flex items-center justify-center">
									<span class="text-primary-foreground text-xs font-medium">
										{debt.contact.split(' ').map(n => n[0]).join('')}
									</span>
								</div>
								<div>
									<p class="font-medium text-foreground">{debt.contact}</p>
									<p class="text-sm text-muted-foreground">Due {formatDate(debt.dueDate)}</p>
								</div>
							</div>
						</div>
						<div class="text-right">
							<p class="font-semibold text-foreground {debt.type === 'owed_to_me' ? 'text-success' : 'text-destructive'}">
								{debt.type === 'owed_to_me' ? '+' : '-'}{formatCurrency(debt.amount)}
							</p>
							<span class="inline-flex px-2 py-1 text-xs font-medium rounded-full {getStatusColor(debt.status)}">
								{debt.status.replace('_', ' ')}
							</span>
						</div>
					</div>
				{/each}
			</div>
		</div>

		<!-- Upcoming Payments -->
		<div class="card p-6">
			<div class="flex items-center justify-between mb-6">
				<h2 class="text-lg font-semibold text-foreground">Upcoming Payments</h2>
				<a href="/debts" class="text-sm text-primary hover:text-primary/80 font-medium">
					View all
				</a>
			</div>
			<div class="space-y-4">
				{#each upcomingPayments as payment (payment.id)}
					<div class="flex items-center justify-between p-4 bg-muted/50 rounded-lg">
						<div class="flex items-center space-x-3">
							<div class="w-2 h-2 rounded-full {payment.type === 'payment_due' ? 'bg-destructive' : 'bg-success'}"></div>
							<div>
								<p class="font-medium text-foreground">{payment.contact}</p>
								<p class="text-sm text-muted-foreground">{formatDate(payment.dueDate)}</p>
							</div>
						</div>
						<div class="text-right">
							<p class="font-semibold {payment.type === 'payment_due' ? 'text-destructive' : 'text-success'}">
								{payment.type === 'payment_due' ? 'Pay' : 'Receive'} {formatCurrency(payment.amount)}
							</p>
						</div>
					</div>
				{/each}
			</div>
		</div>
	</div>

	<!-- Quick Actions -->
	<div class="card p-6">
		<h2 class="text-lg font-semibold text-foreground mb-6">Quick Actions</h2>
		<div class="grid grid-cols-1 md:grid-cols-3 gap-4">
			<a href="/debts/new" class="flex items-center space-x-3 p-4 bg-primary/10 hover:bg-primary/20 rounded-lg transition-colors duration-200">
				<div class="w-10 h-10 bg-primary rounded-lg flex items-center justify-center">
					<svg class="w-5 h-5 text-primary-foreground" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6"></path>
					</svg>
				</div>
				<div>
					<p class="font-medium text-primary">Add New Debt</p>
					<p class="text-sm text-primary/80">Record a new debt or loan</p>
				</div>
			</a>

			<a href="/contacts/new" class="flex items-center space-x-3 p-4 bg-success/10 hover:bg-success/20 rounded-lg transition-colors duration-200">
				<div class="w-10 h-10 bg-success rounded-lg flex items-center justify-center">
					<svg class="w-5 h-5 text-success-foreground" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"></path>
					</svg>
				</div>
				<div>
					<p class="font-medium text-success">Add Contact</p>
					<p class="text-sm text-success/80">Create a new contact</p>
				</div>
			</a>

			<a href="/debts/payment" class="flex items-center space-x-3 p-4 bg-warning/10 hover:bg-warning/20 rounded-lg transition-colors duration-200">
				<div class="w-10 h-10 bg-warning rounded-lg flex items-center justify-center">
					<svg class="w-5 h-5 text-warning-foreground" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 9V7a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2m2 4h10a2 2 0 002-2v-6a2 2 0 00-2-2H9a2 2 0 00-2 2v6a2 2 0 002 2zm7-5a2 2 0 11-4 0 2 2 0 014 0z"></path>
					</svg>
				</div>
				<div>
					<p class="font-medium text-warning">Record Payment</p>
					<p class="text-sm text-warning/80">Log a payment or receipt</p>
				</div>
			</a>
		</div>
	</div>
</div>
