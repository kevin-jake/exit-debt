<script lang="ts">
	import { onMount } from 'svelte';
	import { themeStore } from '$lib/stores/theme.svelte.js';
	import CreateContactModal from '$lib/components/CreateContactModal.svelte';
	import DebtDetailsModal from '$lib/components/DebtDetailsModal.svelte';
	import EditDebtListModal from '$lib/components/EditDebtListModal.svelte';

	// Type definitions
	interface Contact {
		name: string;
		avatar: string | null;
	}

	interface Payment {
		amount: number;
		date: string;
	}

	interface DebtList {
		id: number;
		title: string;
		description: string;
		total_amount: number;
		is_owed_by_me: boolean;
		contact: Contact;
		created_at: string;
		updated_at: string;
		due_date: string;
		last_payment: Payment | null;
	}

	// State management
	let debtLists = $state<DebtList[]>([]);
	let searchQuery = $state('');
	let sortBy = $state('updated_at');
	let sortOrder = $state('desc');
	let isLoading = $state(true);
	let showCreateContactModal = $state(false);
	let showDebtDetailsModal = $state(false);
	let showEditDebtModal = $state(false);
	let selectedDebt = $state<DebtList | null>(null);

	// Quick overview calculations
	let totalIOwe = $derived(
		debtLists
			.filter((debt) => debt.total_amount > 0 && debt.is_owed_by_me)
			.reduce((sum, debt) => sum + debt.total_amount, 0)
	);

	let totalOwedToMe = $derived(
		debtLists
			.filter((debt) => debt.total_amount > 0 && !debt.is_owed_by_me)
			.reduce((sum, debt) => sum + debt.total_amount, 0)
	);

	// Filtered and sorted debt lists
	let filteredDebtLists = $derived(
		debtLists
			.filter((debt) => {
				if (!searchQuery) return true;
				const query = searchQuery.toLowerCase();
				return (
					debt.contact?.name?.toLowerCase().includes(query) ||
					debt.title?.toLowerCase().includes(query) ||
					debt.description?.toLowerCase().includes(query) ||
					debt.total_amount.toString().includes(query)
				);
			})
			.sort((a, b) => {
				let aValue, bValue;

				switch (sortBy) {
					case 'amount':
						aValue = a.total_amount;
						bValue = b.total_amount;
						break;
					case 'created_at':
						aValue = new Date(a.created_at);
						bValue = new Date(b.created_at);
						break;
					case 'contact_name':
						aValue = a.contact?.name || '';
						bValue = b.contact?.name || '';
						break;
					default: // updated_at
						aValue = new Date(a.updated_at);
						bValue = new Date(b.updated_at);
				}

				if (sortOrder === 'asc') {
					return aValue > bValue ? 1 : -1;
				} else {
					return aValue < bValue ? 1 : -1;
				}
			})
	);

	// Upcoming due dates calculation
	let upcomingDueDates = $derived(
		debtLists
			.filter((debt) => {
				if (!debt.due_date) return false;
				const dueDate = new Date(debt.due_date);
				const thirtyDaysFromNow = new Date();
				thirtyDaysFromNow.setDate(thirtyDaysFromNow.getDate() + 30);
				return dueDate <= thirtyDaysFromNow;
			})
			.sort((a, b) => new Date(a.due_date).getTime() - new Date(b.due_date).getTime())
	);

	// Methods
	function handleSearch() {
		// Debounced search implementation
		// The reactive statement will handle the filtering
	}

	function handleSort(newSortBy: string) {
		if (sortBy === newSortBy) {
			sortOrder = sortOrder === 'asc' ? 'desc' : 'asc';
		} else {
			sortBy = newSortBy;
			sortOrder = 'desc';
		}
	}

	function getDaysUntilDue(dueDate: string): number {
		const today = new Date();
		const due = new Date(dueDate);
		const diffTime = due.getTime() - today.getTime();
		return Math.ceil(diffTime / (1000 * 60 * 60 * 24));
	}

	function getDueDateColor(daysUntilDue: number): string {
		if (daysUntilDue < 0) return 'text-destructive';
		if (daysUntilDue <= 3) return 'text-orange-500';
		return 'text-green-600';
	}

	function formatCurrency(amount: number): string {
		return new Intl.NumberFormat('en-US', {
			style: 'currency',
			currency: 'USD',
			minimumFractionDigits: 2,
			maximumFractionDigits: 2
		}).format(amount);
	}

	function formatDate(dateString: string): string {
		return new Date(dateString).toLocaleDateString('en-US', {
			month: 'short',
			day: 'numeric',
			year: 'numeric'
		});
	}

	function handleContactCreated(event: CustomEvent) {
		// In a real app, this would refresh the dashboard data
		showCreateContactModal = false;
		console.log('Contact created:', event.detail);
	}

	function handleDebtDetails(debt: DebtList) {
		console.log(debt);
		selectedDebt = debt;
		showDebtDetailsModal = true;
	}

	function handleEditDebt(debt: DebtList) {
		selectedDebt = debt;
		showEditDebtModal = true;
	}

	function handleDebtUpdated(event: CustomEvent) {
		// In a real app, this would refresh the dashboard data
		showEditDebtModal = false;
		selectedDebt = null;
		console.log('Debt updated:', event.detail);
	}

	function handleDebtDeleted(event: CustomEvent) {
		// In a real app, this would refresh the dashboard data
		showDebtDetailsModal = false;
		selectedDebt = null;
		console.log('Debt deleted:', event.detail);
	}

	function mapDebtListToDebtModal(debtList: DebtList) {
		return {
			id: debtList.id,
			type: debtList.is_owed_by_me ? 'i_owe' : 'owed_to_me',
			contactName: debtList.contact?.name || 'Unknown Contact',
			totalAmount: debtList.total_amount,
			remainingBalance: debtList.total_amount, // Assuming no payments made yet
			status: 'active' as const, // Default status
			dueDate: debtList.due_date,
			installmentPlan: 'one_time' as const, // Default to one-time
			nextPayment: debtList.due_date,
			currency: 'USD', // Default currency
			description: debtList.description,
			createdAt: debtList.created_at,
			// Additional properties that might be needed
			notes: debtList.description || '',
			numberOfPayments: 1,
			updatedAt: debtList.updated_at
		};
	}

	onMount(async () => {
		try {
			// Mock data for demonstration
			debtLists = [
				{
					id: 1,
					title: 'Car Loan',
					description: 'Monthly car payment',
					total_amount: 15000,
					is_owed_by_me: true,
					contact: { name: 'Bank of America', avatar: null },
					created_at: '2024-01-15T10:00:00Z',
					updated_at: '2024-03-20T14:30:00Z',
					due_date: '2024-04-01T00:00:00Z',
					last_payment: { amount: 500, date: '2024-03-15T10:00:00Z' }
				},
				{
					id: 2,
					title: 'Personal Loan',
					description: 'Loan from friend',
					total_amount: 2500,
					is_owed_by_me: false,
					contact: { name: 'Alice Johnson', avatar: null },
					created_at: '2024-02-10T09:00:00Z',
					updated_at: '2024-03-18T16:45:00Z',
					due_date: '2024-03-25T00:00:00Z',
					last_payment: { amount: 200, date: '2024-03-10T10:00:00Z' }
				},
				{
					id: 3,
					title: 'Credit Card',
					description: 'Monthly credit card payment',
					total_amount: 3200,
					is_owed_by_me: true,
					contact: { name: 'Chase Bank', avatar: null },
					created_at: '2024-01-20T11:30:00Z',
					updated_at: '2024-03-19T12:15:00Z',
					due_date: '2024-03-28T00:00:00Z',
					last_payment: { amount: 300, date: '2024-03-05T14:20:00Z' }
				},
				{
					id: 4,
					title: 'Business Investment',
					description: 'Investment in startup',
					total_amount: 8000,
					is_owed_by_me: false,
					contact: { name: 'Bob Smith', avatar: null },
					created_at: '2024-02-05T08:15:00Z',
					updated_at: '2024-03-15T13:30:00Z',
					due_date: '2024-04-15T00:00:00Z',
					last_payment: { amount: 1000, date: '2024-03-01T09:45:00Z' }
				},
				{
					id: 5,
					title: 'Student Loan',
					description: 'Federal student loan',
					total_amount: 25000,
					is_owed_by_me: true,
					contact: { name: 'Department of Education', avatar: null },
					created_at: '2024-01-01T00:00:00Z',
					updated_at: '2024-03-10T10:00:00Z',
					due_date: '2024-04-05T00:00:00Z',
					last_payment: { amount: 400, date: '2024-03-01T11:00:00Z' }
				},
				{
					id: 6,
					title: 'Rent Money',
					description: 'Rent payment from roommate',
					total_amount: 1200,
					is_owed_by_me: false,
					contact: { name: 'Carol Davis', avatar: null },
					created_at: '2024-03-01T07:00:00Z',
					updated_at: '2024-03-20T15:20:00Z',
					due_date: '2024-03-30T00:00:00Z',
					last_payment: null
				}
			];
		} catch (error) {
			console.error('Failed to load dashboard data:', error);
		} finally {
			isLoading = false;
		}
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
	</div>

	<!-- Quick Overview Section -->
	<div class="grid grid-cols-1 md:grid-cols-2 gap-6 mb-8">
		<!-- Total Amount I Owe -->
		<div class="card p-6">
			<div class="flex items-center justify-between">
				<div>
					<h3 class="text-sm font-medium text-muted-foreground">Total Amount I Owe</h3>
					<p class="text-3xl font-bold text-destructive mt-2">
						{formatCurrency(totalIOwe)}
					</p>
				</div>
				<div class="w-12 h-12 bg-destructive/10 rounded-lg flex items-center justify-center">
					<svg class="w-6 h-6 text-destructive" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1" />
					</svg>
				</div>
			</div>
		</div>

		<!-- Total Amount Owed to Me -->
		<div class="card p-6">
			<div class="flex items-center justify-between">
				<div>
					<h3 class="text-sm font-medium text-muted-foreground">Total Amount Owed to Me</h3>
					<p class="text-3xl font-bold text-green-600 mt-2">
						{formatCurrency(totalOwedToMe)}
					</p>
				</div>
				<div class="w-12 h-12 bg-green-100 dark:bg-green-900/20 rounded-lg flex items-center justify-center">
					<svg class="w-6 h-6 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1" />
					</svg>
				</div>
			</div>
		</div>
	</div>

		<!-- Quick Actions Section -->
		<div class="card p-6 mt-8">
			<h2 class="text-xl font-semibold text-foreground mb-4">Quick Actions</h2>
	
			<div class="grid grid-cols-2 gap-4">
				<a href="/debts/new" class="btn-primary text-center">
					<svg class="w-5 h-5 mb-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
					</svg>
					Add New Debt
				</a>
	
				<button 
					on:click={() => showCreateContactModal = true}
					class="btn-secondary text-center"
				>
					<svg class="w-5 h-5 mb-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
					</svg>
					Add New Contact
				</button>
			</div>
		</div>

	<!-- Recent Debt Lists Section -->
	<div class="card p-6 mb-8">
		<div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4 mb-6">
			<h2 class="text-xl font-semibold text-foreground">Recent Debt Lists</h2>

			<!-- Search and Sort Controls -->
			<div class="flex flex-col sm:flex-row gap-3">
				<!-- Search Bar -->
				<div class="relative">
					<input
						type="text"
						bind:value={searchQuery}
						on:input={handleSearch}
						placeholder="Search debt lists..."
						class="input pl-10 pr-4 w-full sm:w-64"
					/>
					<svg class="absolute left-3 top-1/2 transform -translate-y-1/2 w-4 h-4 text-muted-foreground" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
					</svg>
					{#if searchQuery}
						<button
							on:click={() => (searchQuery = '')}
							class="absolute right-3 top-1/2 transform -translate-y-1/2 text-muted-foreground hover:text-foreground"
						>
							×
						</button>
					{/if}
				</div>

				<!-- Sort Dropdown -->
				<select bind:value={sortBy} on:change={() => handleSort(sortBy)} class="input">
					<option value="updated_at">Date Updated</option>
					<option value="created_at">Date Created</option>
					<option value="amount">Amount</option>
					<option value="contact_name">Contact Name</option>
				</select>

				<button
					on:click={() => (sortOrder = sortOrder === 'asc' ? 'desc' : 'asc')}
					class="btn-secondary px-3"
					title={sortOrder === 'asc' ? 'Sort Ascending' : 'Sort Descending'}
				>
					<svg class="w-4 h-4 {sortOrder === 'asc' ? 'rotate-180' : ''}" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 16V4m0 0L3 8m4-4l4 4m6 0v12m0 0l4-4m-4 4l-4-4" />
					</svg>
				</button>
			</div>
		</div>

		<!-- Debt Lists Grid -->
		{#if isLoading}
			<div class="flex justify-center py-8">
				<div class="animate-spin rounded-full h-8 w-8 border-b-2 border-primary"></div>
			</div>
		{:else if filteredDebtLists.length === 0}
			<div class="text-center py-8 text-muted-foreground">
				{searchQuery ? 'No debt lists found matching your search.' : 'No debt lists found.'}
			</div>
		{:else}
			<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
				{#each filteredDebtLists as debt}
					<div class="card p-4 hover:shadow-md transition-shadow cursor-pointer" on:click={() => handleDebtDetails(debt)}>
						<div class="flex items-start justify-between mb-3">
							<div class="flex items-center space-x-3">
								<div class="w-10 h-10 bg-primary/10 rounded-full flex items-center justify-center">
									{#if debt.contact?.avatar}
										<img src={debt.contact.avatar} alt={debt.contact.name} class="w-10 h-10 rounded-full" />
									{:else}
										<span class="text-primary font-medium">
											{debt.contact?.name?.charAt(0) || '?'}
										</span>
									{/if}
								</div>
								<div>
									<h3 class="font-medium text-foreground">{debt.contact?.name || 'Unknown Contact'}</h3>
									<p class="text-sm text-muted-foreground">{debt.title}</p>
								</div>
							</div>
							<div class="text-right">
								<p class="font-semibold {debt.is_owed_by_me ? 'text-destructive' : 'text-green-600'}">
									{formatCurrency(debt.total_amount)}
								</p>
								<p class="text-xs text-muted-foreground">
									{debt.is_owed_by_me ? 'You owe' : 'Owed to you'}
								</p>
							</div>
						</div>

						{#if debt.last_payment}
							<div class="text-sm text-muted-foreground mb-3">
								Last payment: {formatCurrency(debt.last_payment.amount)} on {formatDate(debt.last_payment.date)}
							</div>
						{/if}

						<div class="flex items-center justify-between text-xs text-muted-foreground">
							<span>Updated {formatDate(debt.updated_at)}</span>
							<div class="flex space-x-2" on:click|stopPropagation>
								<button class="text-primary hover:text-primary/80" on:click={() => handleDebtDetails(debt)}>View</button>
								<button class="text-muted-foreground hover:text-foreground" on:click={() => handleEditDebt(debt)}>Edit</button>
							</div>
						</div>
					</div>
				{/each}
			</div>
		{/if}
	</div>

	<!-- Upcoming Due Dates Section -->
	<div class="card p-6">
		<h2 class="text-xl font-semibold text-foreground mb-4">Upcoming Due Dates</h2>

		{#if upcomingDueDates.length === 0}
			<div class="text-center py-8 text-muted-foreground">
				No upcoming due dates in the next 30 days.
			</div>
		{:else}
			<div class="space-y-3">
				{#each upcomingDueDates as debt}
					{@const daysUntilDue = getDaysUntilDue(debt.due_date)}
					{@const dueDateColor = getDueDateColor(daysUntilDue)}

					<div class="flex items-center justify-between p-3 border border-border rounded-lg hover:bg-muted/50 transition-colors cursor-pointer" on:click={() => handleDebtDetails(debt)}>
						<div class="flex items-center space-x-3">
							<div class="w-8 h-8 bg-primary/10 rounded-full flex items-center justify-center">
								<span class="text-primary text-sm font-medium">
									{debt.contact?.name?.charAt(0) || '?'}
								</span>
							</div>
							<div>
								<h4 class="font-medium text-foreground">{debt.contact?.name || 'Unknown Contact'}</h4>
								<p class="text-sm text-muted-foreground">{debt.title}</p>
							</div>
						</div>

						<div class="text-right">
							<p class="font-semibold {debt.is_owed_by_me ? 'text-destructive' : 'text-green-600'}">
								{formatCurrency(debt.total_amount)}
							</p>
							<p class="text-sm {dueDateColor}">
								{#if daysUntilDue < 0}
									Overdue by {Math.abs(daysUntilDue)} days
								{:else if daysUntilDue === 0}
									Due today
								{:else}
									Due in {daysUntilDue} days
								{/if}
							</p>
						</div>
					</div>
				{/each}
			</div>
		{/if}
	</div>
</div>

<!-- Create Contact Modal -->
{#if showCreateContactModal}
	<CreateContactModal
		on:contact-created={handleContactCreated}
		on:close={() => showCreateContactModal = false}
	/>
{/if}

<!-- Debt Details Modal -->
{#if showDebtDetailsModal && selectedDebt}
	<DebtDetailsModal
		debt={mapDebtListToDebtModal(selectedDebt)}
		on:close={() => { showDebtDetailsModal = false; selectedDebt = null; }}
		on:edit={() => selectedDebt && handleEditDebt(selectedDebt)}
		on:delete={handleDebtDeleted}
	/>
{/if}

<!-- Edit Debt Modal -->
{#if showEditDebtModal && selectedDebt}
	<EditDebtListModal
		debt={mapDebtListToDebtModal(selectedDebt)}
		on:close={() => { showEditDebtModal = false; selectedDebt = null; }}
		on:debt-updated={handleDebtUpdated}
	/>
{/if}
