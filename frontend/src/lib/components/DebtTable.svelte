<script lang="ts">
	import { onMount } from 'svelte';
	import { apiClient, type DebtList } from '../api';
	import DebtDetailsModal from './DebtDetailsModal.svelte';
	import EditDebtListModal from './EditDebtListModal.svelte';
	import DeleteDebtListModal from './DeleteDebtListModal.svelte';

	// Extended debt interface for display purposes
	type Debt = DebtList & {
		contactName: string;
		remainingBalance: number;
		status: 'active' | 'settled' | 'archived' | 'overdue';
		dueDate: string;
		nextPayment: string;
	};

	let debts: Debt[] = [];
	let filteredDebts: Debt[] = [];
	let selectedDebt: Debt | null = null;
	let showDetailsModal = false;
	let showEditModal = false;
	let showDeleteDialog = false;
	let debtToDelete: Debt | null = null;
	let isLoading = false;
	let error: string | null = null;

	// Filter and search state
	let searchQuery = '';
	let statusFilter = 'all';
	let typeFilter = 'all';
	let sortBy = 'dueDate';
	let sortOrder: 'asc' | 'desc' = 'asc';

	// Pagination
	let currentPage = 1;
	let itemsPerPage = 10;
	let totalPages = 1;

	onMount(() => {
		loadDebts();
	});

	async function loadDebts() {
		isLoading = true;
		error = null;
		
		try {
			// Load debts from backend (contacts are already included)
			const debtLists = await apiClient.getDebtLists();
			
			// Transform debt lists to include contact names and calculate derived fields
			debts = debtLists.map(debtList => {
				// Use contact information from the debt response
				const contactName = debtList.contact?.name || 'Unknown Contact';
				
				// Use actual remaining balance from backend
				const remainingBalance = parseFloat(debtList.total_remaining_debt || debtList.total_amount);
				
				// Use actual status from backend
				const status = debtList.status as 'active' | 'settled' | 'archived' | 'overdue' || 'active';
				
				// Use actual due date from backend
				const dueDate = debtList.due_date || debtList.created_at;
				const nextPayment = debtList.next_payment_date || debtList.created_at;

				return {
					...debtList,
					contactName,
					remainingBalance,
					status,
					dueDate,
					nextPayment,
				};
			});
			
			filterAndSortDebts();
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to load debts';
			console.error('Error loading debts:', err);
		} finally {
			isLoading = false;
		}
	}

	function filterAndSortDebts() {
		let filtered = [...debts]; // Create a copy to avoid mutating the original array

		// Apply search filter
		if (searchQuery) {
			filtered = filtered.filter(debt => 
				debt.contactName.toLowerCase().includes(searchQuery.toLowerCase()) ||
				debt.description?.toLowerCase().includes(searchQuery.toLowerCase()) ||
				debt.notes?.toLowerCase().includes(searchQuery.toLowerCase())
			);
		}

		// Apply status filter
		if (statusFilter !== 'all') {
			filtered = filtered.filter(debt => debt.status === statusFilter);
		}

		// Apply type filter
		if (typeFilter !== 'all') {
			filtered = filtered.filter(debt => debt.debt_type === typeFilter);
		}

		// Apply sorting
		filtered.sort((a, b) => {
			let aValue: any;
			let bValue: any;

			switch (sortBy) {
				case 'contactName':
					aValue = a.contactName.toLowerCase();
					bValue = b.contactName.toLowerCase();
					break;
				case 'totalAmount':
					aValue = parseFloat(a.total_amount);
					bValue = parseFloat(b.total_amount);
					break;
				case 'remainingBalance':
					aValue = a.remainingBalance;
					bValue = b.remainingBalance;
					break;
				case 'status':
					aValue = a.status;
					bValue = b.status;
					break;
				case 'dueDate':
					aValue = new Date(a.dueDate || '9999-12-31').getTime();
					bValue = new Date(b.dueDate || '9999-12-31').getTime();
					break;
				case 'debt_type':
					aValue = a.debt_type;
					bValue = b.debt_type;
					break;
				default:
					aValue = a[sortBy as keyof Debt];
					bValue = b[sortBy as keyof Debt];
			}

			if (sortOrder === 'asc') {
				return aValue < bValue ? -1 : aValue > bValue ? 1 : 0;
			} else {
				return aValue > bValue ? -1 : aValue < bValue ? 1 : 0;
			}
		});

		filteredDebts = filtered;
		totalPages = Math.ceil(filteredDebts.length / itemsPerPage);
		currentPage = Math.min(currentPage, totalPages || 1);
		
		console.log('Filtered debts:', filteredDebts, 'Total debts:', debts);
	}

	function handleSort(column: string) {
		if (sortBy === column) {
			sortOrder = sortOrder === 'asc' ? 'desc' : 'asc';
		} else {
			sortBy = column;
			sortOrder = 'asc';
		}
		filterAndSortDebts();
	}

	function formatCurrency(amount: number | string, currency: string = 'PHP'): string {
		const numAmount = typeof amount === 'string' ? parseFloat(amount) : amount;
		return new Intl.NumberFormat('en-PH', {
			style: 'currency',
			currency: currency
		}).format(numAmount);
	}

	function formatDate(dateString: string): string {
		if (!dateString) return 'N/A';
		const date = new Date(dateString);
		const now = new Date();
		const diffTime = date.getTime() - now.getTime();
		const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24));

		if (diffDays < 0) {
			return `Overdue by ${Math.abs(diffDays)} days`;
		} else if (diffDays === 0) {
			return 'Due today';
		} else if (diffDays <= 7) {
			return `Due in ${diffDays} days`;
		} else {
			return date.toLocaleDateString('en-US', {
				month: 'short',
				day: 'numeric',
				year: 'numeric'
			});
		}
	}

	function getStatusBadgeClass(status: string): string {
		switch (status) {
			case 'active':
				return 'bg-primary/10 text-primary';
			case 'settled':
				return 'bg-success/10 text-success';
			case 'archived':
				return 'bg-muted/50 text-muted-foreground';
			case 'overdue':
				return 'bg-destructive/10 text-destructive';
			default:
				return 'bg-muted/50 text-muted-foreground';
		}
	}

	function getTypeBadgeClass(type: string): string {
		return type === 'owed_to_me' 
			? 'bg-success/10 text-success' 
			: 'bg-destructive/10 text-destructive';
	}

	function getInstallmentText(plan: string): string {
		const plans = {
			'one_time': 'One-time',
			'weekly': 'Weekly',
			'biweekly': 'Bi-weekly',
			'monthly': 'Monthly',
			'quarterly': 'Quarterly',
			'yearly': 'Yearly'
		};
		return plans[plan as keyof typeof plans] || plan;
	}

	function viewDebt(debt: Debt) {
		selectedDebt = debt;
		showDetailsModal = true;
	}

	function editDebt(debt: Debt) {
		selectedDebt = debt;
		showEditModal = true;
	}

	function confirmDeleteDebt(debt: Debt) {
		debtToDelete = debt;
		showDeleteDialog = true;
	}

	async function deleteDebt() {
		if (debtToDelete) {
			try {
				await apiClient.deleteDebtList(debtToDelete.id);
				debts = debts.filter(d => d.id !== debtToDelete?.id);
				filterAndSortDebts();
				debtToDelete = null;
				showDeleteDialog = false;
			} catch (err) {
				error = err instanceof Error ? err.message : 'Failed to delete debt';
				console.error('Error deleting debt:', err);
			}
		}
	}

	async function markAsSettled(debt: Debt) {
		try {
			// Update the debt status to settled
			const updatedDebt = await apiClient.updateDebtList(debt.id, {
				// Note: The backend might not have a status field, so we'll handle this differently
				// For now, we'll just update the local state
			});
			
			const index = debts.findIndex(d => d.id === debt.id);
			if (index !== -1) {
				debts[index] = { ...debt, status: 'settled', remainingBalance: 0 };
				filterAndSortDebts();
			}
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to mark debt as settled';
			console.error('Error marking debt as settled:', err);
		}
	}

	function handleDebtUpdated(event: CustomEvent) {
		const updatedDebt = event.detail;
		const index = debts.findIndex(d => d.id === updatedDebt.id);
		if (index !== -1) {
			debts[index] = updatedDebt;
			filterAndSortDebts();
		}
		showEditModal = false;
	}

	// Initialize filteredDebts when debts change
	$: if (debts.length > 0) {
		filterAndSortDebts();
	} else {
		filteredDebts = [];
		totalPages = 1;
		currentPage = 1;
	}

	// React to filter changes
	$: if (searchQuery || statusFilter || typeFilter || sortBy || sortOrder) {
		filterAndSortDebts();
	}

	$: paginatedDebts = filteredDebts.slice(
		(currentPage - 1) * itemsPerPage,
		currentPage * itemsPerPage
	);
	
</script>

<div class="space-y-6">
	<!-- Loading State -->
	{#if isLoading}
		<div class="flex items-center justify-center py-12">
			<div class="flex items-center space-x-2">
				<div class="animate-spin rounded-full h-6 w-6 border-b-2 border-primary"></div>
				<span class="text-muted-foreground">Loading debts...</span>
			</div>
		</div>
	{/if}

	<!-- Error State -->
	{#if error && !isLoading}
		<div class="card p-4 bg-destructive/10 border border-destructive/20">
			<div class="flex items-center space-x-2">
				<svg class="w-5 h-5 text-destructive" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path>
				</svg>
				<span class="text-destructive font-medium">Error loading debts</span>
			</div>
			<p class="text-destructive/80 mt-1">{error}</p>
			<button on:click={loadDebts} class="btn-secondary mt-3">Try Again</button>
		</div>
	{/if}

	<!-- Main Content -->
	{#if !isLoading && !error}
		<!-- Header with Search and Filters -->
		<div class="flex flex-col lg:flex-row lg:items-center lg:justify-between gap-4">
			<div class="flex-1 max-w-md">
				<div class="relative">
					<svg class="absolute left-3 top-1/2 transform -translate-y-1/2 w-5 h-5 text-muted-foreground" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"></path>
					</svg>
					<input
						type="text"
						bind:value={searchQuery}
						placeholder="Search debts..."
						class="input pl-10"
					/>
				</div>
			</div>

			<div class="flex items-center space-x-4">
				<select bind:value={statusFilter} class="input">
					<option value="all">All Status</option>
					<option value="active">Active</option>
					<option value="settled">Settled</option>
					<option value="archived">Archived</option>
					<option value="overdue">Overdue</option>
				</select>

				<select bind:value={typeFilter} class="input">
					<option value="all">All Types</option>
					<option value="owed_to_me">Owed to Me</option>
					<option value="i_owe">I Owe</option>
				</select>

				<a href="/debts/new" class="btn-primary">
					<svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6"></path>
					</svg>
					Add Debt
				</a>
			</div>
		</div>

		<!-- Desktop Table -->
		<div class="hidden lg:block card overflow-hidden">
			<div class="overflow-x-auto">
				<table class="w-full">
					<thead class="bg-muted/50 border-b border-border">
						<tr>
							<th class="px-6 py-3 text-left text-xs font-medium text-muted-foreground uppercase tracking-wider cursor-pointer" on:click={() => handleSort('debt_type')}>
								Debt Type
								{#if sortBy === 'debt_type'}
									<span class="ml-1">{sortOrder === 'asc' ? '↑' : '↓'}</span>
								{/if}
							</th>
							<th class="px-6 py-3 text-left text-xs font-medium text-muted-foreground uppercase tracking-wider cursor-pointer" on:click={() => handleSort('contactName')}>
								Contact
								{#if sortBy === 'contactName'}
									<span class="ml-1">{sortOrder === 'asc' ? '↑' : '↓'}</span>
								{/if}
							</th>
							<th class="px-6 py-3 text-left text-xs font-medium text-muted-foreground uppercase tracking-wider cursor-pointer" on:click={() => handleSort('totalAmount')}>
								Total Amount
								{#if sortBy === 'totalAmount'}
									<span class="ml-1">{sortOrder === 'asc' ? '↑' : '↓'}</span>
								{/if}
							</th>
							<th class="px-6 py-3 text-left text-xs font-medium text-muted-foreground uppercase tracking-wider cursor-pointer" on:click={() => handleSort('remainingBalance')}>
								Remaining
								{#if sortBy === 'remainingBalance'}
									<span class="ml-1">{sortOrder === 'asc' ? '↑' : '↓'}</span>
								{/if}
							</th>
							<th class="px-6 py-3 text-left text-xs font-medium text-muted-foreground uppercase tracking-wider cursor-pointer" on:click={() => handleSort('status')}>
								Status
								{#if sortBy === 'status'}
									<span class="ml-1">{sortOrder === 'asc' ? '↑' : '↓'}</span>
								{/if}
							</th>
							<th class="px-6 py-3 text-left text-xs font-medium text-muted-foreground uppercase tracking-wider cursor-pointer" on:click={() => handleSort('dueDate')}>
								Due Date
								{#if sortBy === 'dueDate'}
									<span class="ml-1">{sortOrder === 'asc' ? '↑' : '↓'}</span>
								{/if}
							</th>
							<th class="px-6 py-3 text-left text-xs font-medium text-muted-foreground uppercase tracking-wider">
								Installment
							</th>
							<th class="px-6 py-3 text-left text-xs font-medium text-muted-foreground uppercase tracking-wider">
								Actions
							</th>
						</tr>
					</thead>
					<tbody class="bg-card divide-y divide-border">
						{#each paginatedDebts as debt (debt.id)}
							<tr class="hover:bg-muted/30 cursor-pointer transition-colors duration-200" on:click={() => viewDebt(debt)}>
								<td class="px-6 py-4 whitespace-nowrap">
									<span class="inline-flex px-2 py-1 text-xs font-medium rounded-full {getTypeBadgeClass(debt.debt_type)}">
										{debt.debt_type === 'owed_to_me' ? 'Owed to Me' : 'I Owe'}
									</span>
								</td>
								<td class="px-6 py-4 whitespace-nowrap">
									<div class="flex items-center">
										<div class="w-8 h-8 bg-primary rounded-full flex items-center justify-center mr-3">
											<span class="text-primary-foreground text-xs font-medium">
												{debt.contactName.split(' ').map(n => n[0]).join('')}
											</span>
										</div>
										<div>
											<div class="text-sm font-medium text-foreground">{debt.contactName}</div>
											{#if debt.description}
												<div class="text-sm text-muted-foreground truncate max-w-32">{debt.description}</div>
											{/if}
										</div>
									</div>
								</td>
								<td class="px-6 py-4 whitespace-nowrap text-sm font-medium text-foreground">
									{formatCurrency(debt.total_amount, debt.currency)}
								</td>
								<td class="px-6 py-4 whitespace-nowrap text-sm font-medium {debt.remainingBalance > 0 ? 'text-warning' : 'text-success'}">
									{formatCurrency(debt.remainingBalance, debt.currency)}
								</td>
								<td class="px-6 py-4 whitespace-nowrap">
									<span class="inline-flex px-2 py-1 text-xs font-medium rounded-full {getStatusBadgeClass(debt.status)}">
										{debt.status.charAt(0).toUpperCase() + debt.status.slice(1)}
									</span>
								</td>
								<td class="px-6 py-4 whitespace-nowrap text-sm text-foreground">
									{formatDate(debt.dueDate)}
								</td>
								<td class="px-6 py-4 whitespace-nowrap text-sm text-foreground">
									{getInstallmentText(debt.installment_plan)}
								</td>
								<td class="px-6 py-4 whitespace-nowrap text-sm font-medium">
									<div class="flex items-center space-x-2" on:click|stopPropagation>
										<button
											on:click={() => viewDebt(debt)}
											class="text-primary hover:text-primary/80 p-1"
											title="View Details"
										>
											<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
												<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"></path>
												<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"></path>
											</svg>
										</button>
										<button
											on:click={() => editDebt(debt)}
											class="text-secondary hover:text-secondary/80 p-1"
											title="Edit"
										>
											<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
												<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"></path>
											</svg>
										</button>
										{#if debt.status === 'active'}
											<button
												on:click={() => markAsSettled(debt)}
												class="text-success hover:text-success/80 p-1"
												title="Mark as Settled"
											>
												<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
													<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"></path>
												</svg>
											</button>
										{/if}
										<button
											on:click={() => confirmDeleteDebt(debt)}
											class="text-destructive hover:text-destructive/80 p-1"
											title="Delete"
										>
											<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
												<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"></path>
											</svg>
										</button>
									</div>
								</td>
							</tr>
						{/each}
					</tbody>
				</table>
			</div>
		</div>

		<!-- Mobile Card Layout -->
		<div class="lg:hidden space-y-4">
			{#each paginatedDebts as debt (debt.id)}
				<div class="card p-4" on:click={() => viewDebt(debt)}>
					<div class="flex items-center justify-between mb-3">
						<span class="inline-flex px-2 py-1 text-xs font-medium rounded-full {getTypeBadgeClass(debt.debt_type)}">
							{debt.debt_type === 'owed_to_me' ? 'Owed to Me' : 'I Owe'}
						</span>
						<span class="inline-flex px-2 py-1 text-xs font-medium rounded-full {getStatusBadgeClass(debt.status)}">
							{debt.status.charAt(0).toUpperCase() + debt.status.slice(1)}
						</span>
					</div>
					
					<div class="flex items-center mb-3">
						<div class="w-10 h-10 bg-primary rounded-full flex items-center justify-center mr-3">
							<span class="text-primary-foreground text-sm font-medium">
								{debt.contactName.split(' ').map(n => n[0]).join('')}
							</span>
						</div>
						<div class="flex-1">
							<div class="font-medium text-foreground">{debt.contactName}</div>
							{#if debt.description}
								<div class="text-sm text-muted-foreground">{debt.description}</div>
							{/if}
						</div>
					</div>

					<div class="grid grid-cols-2 gap-4 text-sm mb-3">
						<div>
							<span class="text-muted-foreground">Total:</span>
							<span class="font-medium ml-1">{formatCurrency(debt.total_amount, debt.currency)}</span>
						</div>
						<div>
							<span class="text-muted-foreground">Remaining:</span>
							<span class="font-medium ml-1 {debt.remainingBalance > 0 ? 'text-warning' : 'text-success'}">
								{formatCurrency(debt.remainingBalance, debt.currency)}
							</span>
						</div>
						<div>
							<span class="text-muted-foreground">Due:</span>
							<span class="ml-1">{formatDate(debt.dueDate)}</span>
						</div>
						<div>
							<span class="text-muted-foreground">Schedule:</span>
							<span class="ml-1">{getInstallmentText(debt.installment_plan)}</span>
						</div>
					</div>

					<div class="flex justify-end space-x-2" on:click|stopPropagation>
						<button on:click={() => viewDebt(debt)} class="btn-secondary text-xs px-3 py-1">View</button>
						<button on:click={() => editDebt(debt)} class="btn-secondary text-xs px-3 py-1">Edit</button>
						{#if debt.status === 'active'}
							<button on:click={() => markAsSettled(debt)} class="text-xs px-3 py-1 bg-success/10 text-success rounded-lg hover:bg-success/20">Settle</button>
						{/if}
					</div>
				</div>
			{/each}
		</div>

		<!-- Pagination -->
		{#if totalPages > 1}
			<div class="flex items-center justify-between">
				<div class="text-sm text-muted-foreground">
					Showing {(currentPage - 1) * itemsPerPage + 1} to {Math.min(currentPage * itemsPerPage, filteredDebts.length)} of {filteredDebts.length} debts
				</div>
				<div class="flex items-center space-x-2">
					<button
						on:click={() => currentPage = Math.max(1, currentPage - 1)}
						disabled={currentPage === 1}
						class="btn-secondary disabled:opacity-50 disabled:cursor-not-allowed"
					>
						Previous
					</button>
					<span class="text-sm text-muted-foreground">
						Page {currentPage} of {totalPages}
					</span>
					<button
						on:click={() => currentPage = Math.min(totalPages, currentPage + 1)}
						disabled={currentPage === totalPages}
						class="btn-secondary disabled:opacity-50 disabled:cursor-not-allowed"
					>
						Next
					</button>
				</div>
			</div>
		{/if}

		<!-- Empty State -->
		{#if filteredDebts.length === 0}
			<div class="text-center py-12">
				<svg class="mx-auto w-12 h-12 text-muted-foreground mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1"></path>
				</svg>
				<h3 class="text-lg font-medium text-foreground mb-2">No debts found</h3>
				<p class="text-muted-foreground mb-4">
					{searchQuery || statusFilter !== 'all' || typeFilter !== 'all' 
						? 'Try adjusting your filters or search query.'
						: 'Get started by adding your first debt.'}
				</p>
				<a href="/debts/new" class="btn-primary">Add First Debt</a>
			</div>
		{/if}
	{/if}
</div>

<!-- Modals -->
{#if showDetailsModal && selectedDebt}
	<DebtDetailsModal
		debt={selectedDebt}
		on:close={() => { showDetailsModal = false; selectedDebt = null; }}
		on:edit={() => selectedDebt && editDebt(selectedDebt)}
		on:delete={() => selectedDebt && confirmDeleteDebt(selectedDebt)}
	/>
{/if}

{#if showEditModal && selectedDebt}
	<EditDebtListModal
		debt={selectedDebt}
		on:close={() => { showEditModal = false; selectedDebt = null; }}
		on:debt-updated={handleDebtUpdated}
	/>
{/if}

{#if showDeleteDialog && debtToDelete}
	<DeleteDebtListModal
		debt={debtToDelete}
		on:confirm={deleteDebt}
		on:close={() => { showDeleteDialog = false; debtToDelete = null; }}
	/>
{/if}