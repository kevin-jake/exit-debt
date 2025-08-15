<script lang="ts">
	import { createEventDispatcher, onMount } from 'svelte';

	export let debt: any;

	const dispatch = createEventDispatcher();

	type Payment = {
		id: number;
		date: string;
		amount: number;
		method: 'cash' | 'bank_transfer' | 'check' | 'digital_wallet' | 'other';
		status: 'completed' | 'pending' | 'failed' | 'refunded';
		description?: string;
	};

	let payments: Payment[] = [];
	let showPaymentForm = false;
	let newPayment = {
		amount: 0,
		method: 'cash' as const,
		description: ''
	};

	onMount(() => {
		loadPayments();
		// Prevent body scroll when modal is open
		document.body.style.overflow = 'hidden';
		return () => {
			document.body.style.overflow = 'auto';
		};
	});

	function loadPayments() {
		// Mock payment data - replace with actual API call
		payments = [
			{
				id: 1,
				date: '2024-01-01',
				amount: 500.00,
				method: 'bank_transfer',
				status: 'completed',
				description: 'First installment payment'
			},
			{
				id: 2,
				date: '2024-01-15',
				amount: 200.00,
				method: 'cash',
				status: 'completed',
				description: 'Partial payment'
			}
		];
	}

	function formatCurrency(amount: number, currency: string = 'PHP'): string {
		return new Intl.NumberFormat('en-PH', {
			style: 'currency',
			currency: currency
		}).format(amount);
	}

	function formatDate(dateString: string): string {
		return new Date(dateString).toLocaleDateString('en-US', {
			month: 'short',
			day: 'numeric',
			year: 'numeric'
		});
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

	function getPaymentStatusClass(status: string): string {
		switch (status) {
			case 'completed':
				return 'bg-success/10 text-success';
			case 'pending':
				return 'bg-warning/10 text-warning';
			case 'failed':
				return 'bg-destructive/10 text-destructive';
			case 'refunded':
				return 'bg-muted/50 text-muted-foreground';
			default:
				return 'bg-muted/50 text-muted-foreground';
		}
	}

	function getMethodText(method: string): string {
		const methods = {
			'cash': 'Cash',
			'bank_transfer': 'Bank Transfer',
			'check': 'Check',
			'digital_wallet': 'Digital Wallet',
			'other': 'Other'
		};
		return methods[method as keyof typeof methods] || method;
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

	function calculateProgress(): number {
		if (debt.totalAmount === 0) return 100;
		const paid = debt.totalAmount - debt.remainingBalance;
		return Math.round((paid / debt.totalAmount) * 100);
	}

	function handleAddPayment() {
		if (newPayment.amount > 0) {
			const payment: Payment = {
				id: payments.length + 1,
				date: new Date().toISOString().split('T')[0],
				amount: newPayment.amount,
				method: newPayment.method,
				status: 'completed',
				description: newPayment.description
			};
			payments = [payment, ...payments];
			
			// Reset form
			newPayment = {
				amount: 0,
				method: 'cash',
				description: ''
			};
			showPaymentForm = false;
		}
	}

	function closeModal() {
		dispatch('close');
	}

	function editDebt() {
		dispatch('edit');
		closeModal();
	}

	function deleteDebt() {
		dispatch('delete');
		closeModal();
	}

	function handleKeydown(event: KeyboardEvent) {
		if (event.key === 'Escape') {
			closeModal();
		}
	}
</script>

<svelte:window on:keydown={handleKeydown} />

<!-- Modal Backdrop -->
<div class="fixed inset-0 bg-black/50 z-50 flex items-center justify-center p-4" on:click={closeModal}>
	<!-- Modal Content -->
	<div 
		class="bg-card rounded-xl shadow-medium max-w-4xl w-full max-h-[90vh] overflow-hidden flex flex-col"
		on:click|stopPropagation
	>
		<!-- Header -->
		<div class="px-6 py-4 border-b border-border flex items-center justify-between">
			<div class="flex items-center space-x-4">
				<h2 class="text-xl font-semibold text-foreground">Debt Details</h2>
				<span class="inline-flex px-3 py-1 text-sm font-medium rounded-full {getStatusBadgeClass(debt.status)}">
					{debt.status.charAt(0).toUpperCase() + debt.status.slice(1)}
				</span>
			</div>
			<button on:click={closeModal} class="text-muted-foreground hover:text-foreground">
				<svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
				</svg>
			</button>
		</div>

		<!-- Content -->
		<div class="flex-1 overflow-y-auto">
			<div class="p-6 space-y-8">
				<!-- Debt Information -->
				<div class="grid grid-cols-1 lg:grid-cols-2 gap-8">
					<!-- Basic Details -->
					<div class="space-y-6">
						<div>
							<h3 class="text-lg font-medium text-foreground mb-4">Basic Information</h3>
							<div class="space-y-4">
								<div class="flex items-center space-x-4">
									<div class="w-12 h-12 bg-primary rounded-full flex items-center justify-center">
										<span class="text-primary-foreground font-medium">
											{debt.contactName.split(' ').map((n: string) => n[0]).join('')}
										</span>
									</div>
									<div>
										<div class="font-medium text-foreground">{debt.contactName}</div>
										<div class="text-sm text-muted-foreground">
											{debt.type === 'owed_to_me' ? 'Owes you money' : 'You owe money to'}
										</div>
									</div>
								</div>

								<div class="grid grid-cols-2 gap-4">
									<div>
										<label class="block text-sm font-medium text-muted-foreground mb-1">Debt Type</label>
										<span class="inline-flex px-3 py-1 text-sm font-medium rounded-full {debt.type === 'owed_to_me' ? 'bg-success/10 text-success' : 'bg-destructive/10 text-destructive'}">
											{debt.type === 'owed_to_me' ? 'Owed to Me' : 'I Owe'}
										</span>
									</div>
									<div>
										<label class="block text-sm font-medium text-muted-foreground mb-1">Currency</label>
										<div class="text-sm text-foreground">{debt.currency}</div>
									</div>
								</div>

								{#if debt.description}
									<div>
										<label class="block text-sm font-medium text-muted-foreground mb-1">Description</label>
										<div class="text-sm text-foreground">{debt.description}</div>
									</div>
								{/if}

								<div>
									<label class="block text-sm font-medium text-muted-foreground mb-1">Created</label>
									<div class="text-sm text-foreground">{formatDate(debt.createdAt)}</div>
								</div>
							</div>
						</div>
					</div>

					<!-- Financial Details -->
					<div class="space-y-6">
						<div>
							<h3 class="text-lg font-medium text-foreground mb-4">Financial Details</h3>
							<div class="space-y-4">
								<div class="grid grid-cols-2 gap-4">
									<div>
										<label class="block text-sm font-medium text-muted-foreground mb-1">Total Amount</label>
										<div class="text-lg font-semibold text-foreground">{formatCurrency(debt.totalAmount, debt.currency)}</div>
									</div>
									<div>
										<label class="block text-sm font-medium text-muted-foreground mb-1">Remaining Balance</label>
										<div class="text-lg font-semibold {debt.remainingBalance > 0 ? 'text-warning' : 'text-success'}">
											{formatCurrency(debt.remainingBalance, debt.currency)}
										</div>
									</div>
								</div>

								<div>
									<label class="block text-sm font-medium text-muted-foreground mb-2">Payment Progress</label>
									<div class="w-full bg-muted rounded-full h-2">
										<div 
											class="bg-primary h-2 rounded-full transition-all duration-300" 
											style="width: {calculateProgress()}%"
										></div>
									</div>
									<div class="text-sm text-muted-foreground mt-1">{calculateProgress()}% paid</div>
								</div>

								<div class="grid grid-cols-2 gap-4">
									<div>
										<label class="block text-sm font-medium text-muted-foreground mb-1">Due Date</label>
										<div class="text-sm text-foreground">{formatDate(debt.dueDate)}</div>
									</div>
									<div>
										<label class="block text-sm font-medium text-muted-foreground mb-1">Installment Plan</label>
										<div class="text-sm text-foreground">{getInstallmentText(debt.installmentPlan)}</div>
									</div>
								</div>

								{#if debt.nextPayment}
									<div>
										<label class="block text-sm font-medium text-muted-foreground mb-1">Next Payment</label>
										<div class="text-sm text-foreground">{formatDate(debt.nextPayment)}</div>
									</div>
								{/if}
							</div>
						</div>
					</div>
				</div>

				<!-- Payment History -->
				<div>
					<div class="flex items-center justify-between mb-4">
						<h3 class="text-lg font-medium text-foreground">Payment History</h3>
						<button 
							on:click={() => showPaymentForm = !showPaymentForm}
							class="btn-primary text-sm"
						>
							<svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6"></path>
							</svg>
							Add Payment
						</button>
					</div>

					<!-- Add Payment Form -->
					{#if showPaymentForm}
						<div class="card p-4 mb-4">
							<h4 class="font-medium text-foreground mb-4">Record New Payment</h4>
							<div class="grid grid-cols-1 md:grid-cols-3 gap-4">
								<div>
									<label class="label">Amount</label>
									<input 
										type="number" 
										bind:value={newPayment.amount} 
										step="0.01" 
										min="0"
										class="input"
										placeholder="0.00"
									/>
								</div>
								<div>
									<label class="label">Payment Method</label>
									<select bind:value={newPayment.method} class="input">
										<option value="cash">Cash</option>
										<option value="bank_transfer">Bank Transfer</option>
										<option value="check">Check</option>
										<option value="digital_wallet">Digital Wallet</option>
										<option value="other">Other</option>
									</select>
								</div>
								<div>
									<label class="label">Description (Optional)</label>
									<input 
										type="text" 
										bind:value={newPayment.description} 
										class="input"
										placeholder="Payment description"
									/>
								</div>
							</div>
							<div class="flex justify-end space-x-3 mt-4">
								<button on:click={() => showPaymentForm = false} class="btn-secondary">
									Cancel
								</button>
								<button on:click={handleAddPayment} class="btn-primary">
									Add Payment
								</button>
							</div>
						</div>
					{/if}

					<!-- Payments Table -->
					{#if payments.length > 0}
						<div class="card overflow-hidden">
							<div class="overflow-x-auto">
								<table class="w-full">
									<thead class="bg-muted/50 border-b border-border">
										<tr>
											<th class="px-4 py-3 text-left text-xs font-medium text-muted-foreground uppercase tracking-wider">Date</th>
											<th class="px-4 py-3 text-left text-xs font-medium text-muted-foreground uppercase tracking-wider">Amount</th>
											<th class="px-4 py-3 text-left text-xs font-medium text-muted-foreground uppercase tracking-wider">Method</th>
											<th class="px-4 py-3 text-left text-xs font-medium text-muted-foreground uppercase tracking-wider">Status</th>
											<th class="px-4 py-3 text-left text-xs font-medium text-muted-foreground uppercase tracking-wider">Description</th>
										</tr>
									</thead>
									<tbody class="bg-card divide-y divide-border">
										{#each payments as payment (payment.id)}
											<tr>
												<td class="px-4 py-3 text-sm text-foreground">{formatDate(payment.date)}</td>
												<td class="px-4 py-3 text-sm font-medium text-foreground">{formatCurrency(payment.amount, debt.currency)}</td>
												<td class="px-4 py-3 text-sm text-foreground">{getMethodText(payment.method)}</td>
												<td class="px-4 py-3">
													<span class="inline-flex px-2 py-1 text-xs font-medium rounded-full {getPaymentStatusClass(payment.status)}">
														{payment.status.charAt(0).toUpperCase() + payment.status.slice(1)}
													</span>
												</td>
												<td class="px-4 py-3 text-sm text-muted-foreground">{payment.description || 'N/A'}</td>
											</tr>
										{/each}
									</tbody>
								</table>
							</div>
						</div>
					{:else}
						<div class="card p-8 text-center">
							<svg class="mx-auto w-12 h-12 text-muted-foreground mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"></path>
							</svg>
							<h4 class="text-lg font-medium text-foreground mb-2">No payments recorded</h4>
							<p class="text-muted-foreground mb-4">Start tracking payments by adding the first payment record.</p>
							<button on:click={() => showPaymentForm = true} class="btn-primary">
								Add First Payment
							</button>
						</div>
					{/if}
				</div>
			</div>
		</div>

		<!-- Footer Actions -->
		<div class="px-6 py-4 border-t border-border flex items-center justify-between">
			<div class="flex items-center space-x-3">
				<button on:click={editDebt} class="btn-secondary">
					<svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"></path>
					</svg>
					Edit Debt
				</button>
				<button on:click={deleteDebt} class="btn-danger">
					<svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"></path>
					</svg>
					Delete Debt
				</button>
			</div>
			<button on:click={closeModal} class="btn-primary">
				Close
			</button>
		</div>
	</div>
</div>
