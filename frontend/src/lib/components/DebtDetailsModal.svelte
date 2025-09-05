<script lang="ts">
	import { createEventDispatcher, onMount } from 'svelte';
	import { apiClient, type Payment } from '../api';
	import ReceiptPhotoViewer from './ReceiptPhotoViewer.svelte';

	export let debt: any;

	const dispatch = createEventDispatcher();

	let payments: Payment[] = [];
	let showPaymentForm = false;
	let newPayment = {
		amount: '',
		payment_date: '',
		payment_method: 'cash',
		description: '',
		receiptPhoto: undefined as File | undefined
	};
	let isLoading = false;
	let error: string | null = null;

	// Receipt photo viewer state
	let showReceiptViewer = false;
	let selectedReceiptPhoto = '';

	console.log(debt);

	onMount(() => {
		loadPayments();
		// Prevent body scroll when modal is open
		document.body.style.overflow = 'hidden';
		return () => {
			document.body.style.overflow = 'auto';
		};
	});

	async function loadPayments() {
		if (!debt?.id) {
			console.error('No debt ID provided');
			error = 'No debt ID provided';
			payments = [];
			return;
		}

		isLoading = true;
		error = null;
		
		try {
			console.log('Loading payments for debt:', debt.id);
			payments = await apiClient.getPayments(debt.id);
			console.log('Loaded payments:', payments);
			
			// Filter out payments without valid IDs and ensure unique IDs
			payments = payments.filter(payment => payment.id && payment.id.trim() !== '');
			
			// Check for duplicate or undefined IDs
			const ids = payments.map(p => p.id);
			const uniqueIds = new Set(ids);
			if (ids.length !== uniqueIds.size) {
				console.warn('Duplicate payment IDs found:', ids);
				// Remove duplicates by keeping only the first occurrence
				const seen = new Set();
				payments = payments.filter(payment => {
					if (seen.has(payment.id)) {
						return false;
					}
					seen.add(payment.id);
					return true;
				});
			}
			if (ids.some(id => !id)) {
				console.warn('Some payments have undefined IDs:', payments.filter(p => !p.id));
			}
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to load payments';
			console.error('Error loading payments:', err);
			// Set empty array on error to prevent UI issues
			payments = [];
		} finally {
			isLoading = false;
		}
	}

	function formatCurrency(amount: number | string, currency: string = 'PHP'): string {
		const numAmount = typeof amount === 'string' ? parseFloat(amount) : amount;
		return new Intl.NumberFormat('en-PH', {
			style: 'currency',
			currency: currency
		}).format(numAmount);
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
			case 'verified':
			case 'completed':
				return 'bg-success/10 text-success';
			case 'pending':
				return 'bg-warning/10 text-warning';
			case 'overdue':
				return 'bg-destructive/10 text-destructive';
			case 'rejected':
				return 'bg-destructive/10 text-destructive';
			default:
				return 'bg-muted/50 text-muted-foreground';
		}
	}

	// File validation
	function validateFile(file: File): string | null {
		const maxSize = 5 * 1024 * 1024; // 5MB
		const allowedTypes = ['image/jpeg', 'image/jpg', 'image/png', 'image/gif', 'image/webp'];
		
		if (file.size > maxSize) return 'File size must be less than 5MB';
		if (!allowedTypes.includes(file.type)) return 'Only image files are allowed';
		return null;
	}

	function handleFileSelect(event: Event) {
		const target = event.target as HTMLInputElement;
		if (target.files && target.files[0]) {
			const file = target.files[0];
			const error = validateFile(file);
			
			if (error) {
				alert(error);
				return;
			}
			
			newPayment.receiptPhoto = file;
		}
	}

	function removeReceiptPhoto() {
		newPayment.receiptPhoto = undefined;
	}

	function viewReceiptPhoto(photoURL: string) {
		selectedReceiptPhoto = photoURL;
		showReceiptViewer = true;
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
		const totalAmount = parseFloat(debt.total_amount);
		if (totalAmount === 0) return 100;
		const paid = totalAmount - debt.remainingBalance;
		return Math.round((paid / totalAmount) * 100);
	}

	async function handleAddPayment() {
		if (parseFloat(newPayment.amount) > 0 && newPayment.payment_date) {
			try {
				console.log('Creating payment:', newPayment);
				// Create the payment via API
				const payment = await apiClient.createPayment(debt.id, {
					amount: newPayment.amount,
					payment_date: newPayment.payment_date,
					payment_method: newPayment.payment_method,
					description: newPayment.description
				});

				console.log('Payment created:', payment);

				// Upload receipt photo if provided
				if (newPayment.receiptPhoto) {
					console.log('Uploading receipt photo');
					await apiClient.uploadReceipt(payment.id, newPayment.receiptPhoto);
				}

				// Reload payments to get updated data
				await loadPayments();
				
				// Reset form
				newPayment = {
					amount: '',
					payment_date: '',
					payment_method: 'cash',
					description: '',
					receiptPhoto: undefined
				};
				showPaymentForm = false;
			} catch (err) {
				error = err instanceof Error ? err.message : 'Failed to add payment';
				console.error('Error adding payment:', err);
			}
		} else {
			error = 'Please provide both amount and payment date';
		}
	}

	async function verifyPayment(paymentId: string) {
		try {
			console.log('Verifying payment:', paymentId);
			await apiClient.verifyPayment(paymentId);
			console.log('Payment verified successfully');
			// Reload payments to get updated data
			await loadPayments();
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to verify payment';
			console.error('Error verifying payment:', err);
		}
	}

	function closeModal() {
		dispatch('close');
	}

	function editDebt() {
		dispatch('edit');
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
											{debt.debt_type === 'owed_to_me' ? 'Owes you money' : 'You owe money to'}
										</div>
									</div>
								</div>

								<div class="grid grid-cols-2 gap-4">
									<div>
										<label class="block text-sm font-medium text-muted-foreground mb-1">Debt Type</label>
										<span class="inline-flex px-3 py-1 text-sm font-medium rounded-full {debt.debt_type === 'owed_to_me' ? 'bg-success/10 text-success' : 'bg-destructive/10 text-destructive'}">
											{debt.debt_type === 'owed_to_me' ? 'Owed to Me' : 'I Owe'}
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
									<div class="text-sm text-foreground">{formatDate(debt.created_at)}</div>
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
										<div class="text-lg font-semibold text-foreground">{formatCurrency(debt.total_amount, debt.currency)}</div>
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
										<div class="text-sm text-foreground">{getInstallmentText(debt.installment_plan)}</div>
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
							{#if error}
								<div class="mb-4 p-3 bg-destructive/10 border border-destructive/20 rounded-lg">
									<p class="text-destructive text-sm">{error}</p>
								</div>
							{/if}
							<div class="grid grid-cols-1 md:grid-cols-2 gap-4">
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
									<label class="label">Payment Date</label>
									<input 
										type="date" 
										bind:value={newPayment.payment_date} 
										class="input"
									/>
								</div>
								<div>
									<label class="label">Payment Method</label>
									<select bind:value={newPayment.payment_method} class="input">
										<option value="cash">Cash</option>
										<option value="bank_transfer">Bank Transfer</option>
										<option value="check">Check</option>
										<option value="digital_wallet">Digital Wallet</option>
									</select>
								</div>
							</div>
							<div class="mt-4">
									<label class="label">Description (Optional)</label>
									<input 
										type="text" 
										bind:value={newPayment.description} 
										class="input"
										placeholder="Payment description"
									/>
							</div>

							<!-- Receipt Photo Upload -->
							<div class="mt-4">
								<label class="label">Receipt Photo (Optional)</label>
								<div class="border-2 border-dashed border-border rounded-lg p-4 text-center hover:border-primary/50 transition-colors">
									{#if newPayment.receiptPhoto}
										<!-- Photo preview -->
										<div class="relative inline-block">
											<img 
												src={URL.createObjectURL(newPayment.receiptPhoto)} 
												alt="Receipt preview" 
												class="w-24 h-24 object-cover rounded-lg"
											/>
											<button 
												type="button"
												on:click={removeReceiptPhoto}
												class="absolute -top-2 -right-2 bg-destructive text-destructive-foreground rounded-full w-6 h-6 flex items-center justify-center text-xs hover:bg-destructive/90"
											>
												×
											</button>
										</div>
									{:else}
										<!-- Upload prompt -->
										<div class="space-y-2">
											<input 
												type="file" 
												accept="image/*" 
												on:change={handleFileSelect}
												class="hidden" 
												id="receipt-upload"
											/>
											<label 
												for="receipt-upload" 
												class="cursor-pointer text-primary hover:text-primary/80"
											>
												Click to upload or drag and drop
											</label>
											<p class="text-xs text-muted-foreground">
												JPG, PNG, GIF, WebP up to 5MB
											</p>
										</div>
									{/if}
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
					{#if isLoading}
						<div class="card p-8 text-center">
							<div class="flex items-center justify-center space-x-2">
								<div class="animate-spin rounded-full h-6 w-6 border-b-2 border-primary"></div>
								<span class="text-muted-foreground">Loading payments...</span>
							</div>
						</div>
					{:else if error}
						<div class="card p-8 text-center">
							<div class="flex items-center justify-center space-x-2 mb-4">
								<svg class="w-8 h-8 text-destructive" fill="none" stroke="currentColor" viewBox="0 0 24 24">
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path>
								</svg>
								<span class="text-destructive font-medium">Failed to load payments</span>
							</div>
							<p class="text-muted-foreground mb-4">{error}</p>
							<button on:click={loadPayments} class="btn-secondary">
								<svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"></path>
								</svg>
								Try Again
							</button>
						</div>
					{:else if payments.length > 0}
						<div class="card overflow-hidden">
							<div class="overflow-x-auto">
								<table class="w-full">
									<thead class="bg-muted/50 border-b border-border">
										<tr>
											<th class="px-4 py-3 text-left text-xs font-medium text-muted-foreground uppercase tracking-wider">Date</th>
											<th class="px-4 py-3 text-left text-xs font-medium text-muted-foreground uppercase tracking-wider">Amount</th>
											<th class="px-4 py-3 text-left text-xs font-medium text-muted-foreground uppercase tracking-wider">Status</th>
											<th class="px-4 py-3 text-left text-xs font-medium text-muted-foreground uppercase tracking-wider">Receipt</th>
											{#if debt.debt_type === 'owed_to_me'}
												<th class="px-4 py-3 text-left text-xs font-medium text-muted-foreground uppercase tracking-wider">Actions</th>
											{/if}
										</tr>
									</thead>
									<tbody class="bg-card divide-y divide-border">
										{#each payments as payment, index (payment.id || `payment-${index}`)}
											<tr>
												<td class="px-4 py-3 text-sm text-foreground">{formatDate(payment.payment_date)}</td>
												<td class="px-4 py-3 text-sm font-medium text-foreground">{formatCurrency(payment.amount, payment.currency || debt.currency)}</td>
												<td class="px-4 py-3">
													<span class="inline-flex px-2 py-1 text-xs font-medium rounded-full {getPaymentStatusClass(payment.status)}">
														{payment.status.charAt(0).toUpperCase() + payment.status.slice(1)}
													</span>
												</td>
												<td class="px-4 py-3">
													{#if payment.receipt_photo_url}
														<button 
															on:click={() => viewReceiptPhoto(payment.receipt_photo_url!)}
															class="w-16 h-16 rounded-lg overflow-hidden border border-border hover:border-primary/50 transition-colors"
															title="View receipt"
														>
															<img 
																src={payment.receipt_photo_url} 
																alt="Receipt" 
																class="w-full h-full object-cover"
															/>
														</button>
													{:else}
														<div class="w-16 h-16 rounded-lg border border-border flex items-center justify-center text-muted-foreground">
															<svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
																<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" />
															</svg>
														</div>
													{/if}
												</td>
												{#if debt.debt_type === 'owed_to_me'}
													<td class="px-4 py-3">
														{#if payment.status !== 'verified' && payment.status !== 'completed'}
															<button 
																on:click={() => verifyPayment(payment.id)}
																class="inline-flex items-center px-3 py-1.5 text-xs font-medium rounded-md bg-blue-100 text-blue-700 hover:bg-blue-200 transition-colors"
																title={payment.status === 'rejected' ? 'Reverify payment' : 'Verify payment'}
															>
																<svg class="w-3 h-3 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
																	<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
																</svg>
																{payment.status === 'rejected' ? 'Reverify' : 'Verify'}
															</button>
														{:else}
															<span class="inline-flex items-center px-3 py-1.5 text-xs font-medium rounded-md bg-green-100 text-green-700">
																<svg class="w-3 h-3 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
																	<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
																</svg>
																Verified
															</span>
														{/if}
													</td>
												{/if}
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

<!-- Receipt Photo Viewer -->
<ReceiptPhotoViewer
	photoUrl={selectedReceiptPhoto}
	isOpen={showReceiptViewer}
	on:close={() => { showReceiptViewer = false; selectedReceiptPhoto = ''; }}
/>
