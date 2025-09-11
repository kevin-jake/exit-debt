<script lang="ts">
	import { createEventDispatcher, onMount } from 'svelte';
	import DatePicker from './DatePicker.svelte';

	export let debt: any;

	const dispatch = createEventDispatcher();

	// Form data with current debt values
	let formData = {
		totalAmount: debt.total_amount?.toString() || debt.totalAmount?.toString() || '',
		currency: debt.currency || '',
		description: debt.description || '',
		dueDate: debt.due_date || debt.dueDate || '',
		installmentPlan: debt.installment_plan || debt.installmentPlan || 'one_time',
		numberOfPayments: debt.number_of_payments || debt.numberOfPayments || 1,
		notes: debt.notes || ''
	};

	// Store original data for comparison and reset
	let originalData = { ...formData };
	
	// State management
	let errors: { [key: string]: string } = {};
	let isLoading = false;
	let hasChanges = false;
	let showResetConfirmation = false;

	// Character limits
	const DESCRIPTION_LIMIT = 500;
	const NOTES_LIMIT = 1000;

	// Computed values
	$: installmentAmount = formData.totalAmount && formData.numberOfPayments > 0 
		? parseFloat(formData.totalAmount) / formData.numberOfPayments 
		: 0;

	$: isInstallmentPlan = (debt.installment_plan || debt.installmentPlan) !== 'one_time';

	onMount(() => {
		// Prevent body scroll when modal is open
		document.body.style.overflow = 'hidden';
		return () => {
			document.body.style.overflow = 'auto';
		};
	});

	function checkForChanges() {
		hasChanges = JSON.stringify(formData) !== JSON.stringify(originalData);
	}

	function validateForm(): boolean {
		errors = {};

		// Amount validation
		const amount = parseFloat(formData.totalAmount);
		if (!formData.totalAmount || isNaN(amount) || amount <= 0) {
			errors.totalAmount = 'Please enter a valid amount greater than 0';
		}

		// Due date validation
		if (!formData.dueDate) {
			errors.dueDate = 'Please select a due date';
		} else {
			const selectedDate = new Date(formData.dueDate);
			const today = new Date();
			today.setHours(0, 0, 0, 0);
			
			// Allow current date for existing debts
			if (selectedDate < today) {
				errors.dueDate = 'Due date cannot be in the past';
			}
		}

		// Installment validation
		if (isInstallmentPlan) {
			if (!formData.numberOfPayments || formData.numberOfPayments < 1) {
				errors.numberOfPayments = 'Number of payments must be at least 1';
			}
		}

		// Character limits
		if (formData.description.length > DESCRIPTION_LIMIT) {
			errors.description = `Description must be ${DESCRIPTION_LIMIT} characters or less`;
		}

		if (formData.notes.length > NOTES_LIMIT) {
			errors.notes = `Notes must be ${NOTES_LIMIT} characters or less`;
		}

		return Object.keys(errors).length === 0;
	}

	async function handleSubmit() {
		if (!validateForm()) {
			// Scroll to first error
			const firstErrorField = document.querySelector('.border-destructive');
			if (firstErrorField) {
				firstErrorField.scrollIntoView({ behavior: 'smooth', block: 'center' });
			}
			return;
		}

		isLoading = true;

		try {
			// TODO: Replace with actual API call
			await new Promise(resolve => setTimeout(resolve, 1500));

			// Mock successful debt update
			const updatedDebt = {
				...debt,
				total_amount: formData.totalAmount,
				currency: formData.currency,
				description: formData.description,
				due_date: formData.dueDate,
				installment_plan: formData.installmentPlan,
				number_of_payments: isInstallmentPlan ? formData.numberOfPayments : 1,
				notes: formData.notes,
				total_remaining_debt: calculateRemainingBalance().toString(),
				updated_at: new Date().toISOString()
			};

			dispatch('debt-updated', updatedDebt);
			
		} catch (error) {
			console.error('Error updating debt:', error);
			errors.general = 'Failed to update debt. Please try again.';
		} finally {
			isLoading = false;
		}
	}

	function calculateRemainingBalance(): number {
		// In a real app, this would calculate based on payments made
		// For now, we'll maintain the same ratio
		const totalAmount = debt.total_amount || debt.totalAmount;
		const remainingDebt = debt.total_remaining_debt || debt.remainingBalance;
		
		if (!totalAmount) return 0;
		
		const originalRatio = remainingDebt ? parseFloat(remainingDebt.toString()) / parseFloat(totalAmount.toString()) : 1;
		return parseFloat(formData.totalAmount) * originalRatio;
	}

	function handleReset() {
		if (hasChanges) {
			showResetConfirmation = true;
		}
	}

	function confirmReset() {
		formData = { ...originalData };
		errors = {};
		showResetConfirmation = false;
		checkForChanges();
	}

	function handleClose() {
		if (hasChanges) {
			const confirmed = confirm('You have unsaved changes. Are you sure you want to close?');
			if (!confirmed) return;
		}
		dispatch('close');
	}

	function formatCurrency(amount: number): string {
		return new Intl.NumberFormat('en-PH', {
			style: 'currency',
			currency: formData.currency
		}).format(amount);
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

	function getInstallmentFrequencyText(plan: string): string {
		const frequencies = {
			'weekly': 'week',
			'biweekly': '2 weeks',
			'monthly': 'month',
			'quarterly': '3 months',
			'yearly': 'year'
		};
		return frequencies[plan as keyof typeof frequencies] || plan;
	}

	function handleKeydown(event: KeyboardEvent) {
		if (event.key === 'Escape') {
			handleClose();
		}
	}

	// Watch for changes
	$: checkForChanges();

	// Recalculate installment amount when relevant fields change
	$: if (formData.totalAmount && formData.numberOfPayments) {
		installmentAmount = parseFloat(formData.totalAmount) / formData.numberOfPayments;
	}
</script>

<svelte:window on:keydown={handleKeydown} />

<!-- Modal Backdrop -->
<div class="fixed inset-0 bg-black/50 z-50 flex items-center justify-center p-4" on:click={handleClose} on:keydown={handleKeydown} role="dialog" aria-modal="true">
	<!-- Modal Content -->
	<div 
		class="bg-card rounded-xl shadow-medium max-w-2xl w-full max-h-[90vh] overflow-hidden flex flex-col"
		on:click|stopPropagation
		on:keydown|stopPropagation
		role="document"
	>
		<!-- Header -->
		<div class="px-6 py-4 border-b border-border">
			<div class="flex items-center justify-between">
				<div>
					<h2 class="text-xl font-semibold text-foreground">
						Edit Debt: {debt.description || `${(debt.debt_type === 'owed_to_me' || debt.type === 'owed_to_me') ? 'Money owed by' : 'Money owed to'} ${debt.contactName || debt.contact?.name || 'Unknown Contact'}`}
					</h2>
					<div class="flex items-center space-x-3 mt-1">
						<span class="text-sm text-muted-foreground">
							{(debt.debt_type === 'owed_to_me' || debt.type === 'owed_to_me') ? 'Owed to Me' : 'I Owe'}
						</span>
						<span class="inline-flex px-2 py-1 text-xs font-medium rounded-full {getStatusBadgeClass(debt.status || 'active')}">
							{(debt.status || 'active').charAt(0).toUpperCase() + (debt.status || 'active').slice(1)}
						</span>
												</div>
						</div>
						<button on:click={handleClose} class="text-muted-foreground hover:text-foreground" aria-label="Close modal">
					<svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
					</svg>
				</button>
			</div>
		</div>

		<!-- Form -->
		<div class="flex-1 overflow-y-auto">
			<form on:submit|preventDefault={handleSubmit} class="p-6 space-y-6">
				<!-- General Error -->
				{#if errors.general}
					<div class="bg-destructive/10 border border-destructive/20 text-destructive px-4 py-3 rounded-lg text-sm">
						{errors.general}
					</div>
				{/if}

				<!-- Change Indicator -->
				{#if hasChanges}
					<div class="bg-warning/10 border border-warning/20 text-warning px-4 py-3 rounded-lg text-sm flex items-center">
						<svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path>
						</svg>
						You have unsaved changes
					</div>
				{/if}

				<!-- Contact Information (Read-only) -->
				<div class="bg-muted/50 rounded-lg p-4">
					<h3 class="text-sm font-medium text-foreground mb-3">Contact Information</h3>
					<div class="flex items-center space-x-3">
						<div class="w-10 h-10 bg-primary rounded-full flex items-center justify-center">
							<span class="text-primary-foreground text-sm font-medium">
								{(debt.contactName || debt.contact?.name || 'Unknown').split(' ').map((n: string) => n[0]).join('')}
							</span>
						</div>
						<div>
							<div class="font-medium text-foreground">{debt.contactName || debt.contact?.name || 'Unknown Contact'}</div>
							<div class="text-sm text-muted-foreground">Contact cannot be changed after debt creation</div>
						</div>
					</div>
				</div>

				<!-- Debt Details -->
				<div class="space-y-4">
					<h3 class="text-sm font-medium text-foreground">Debt Details</h3>
					
					<!-- Debt Type (Read-only) -->
					<div class="bg-muted/50 rounded-lg p-4">
						<div class="block text-sm font-medium text-muted-foreground mb-2">Debt Type</div>
						<div class="flex items-center justify-between">
							<span class="inline-flex px-3 py-1 text-sm font-medium rounded-full {(debt.debt_type === 'owed_to_me' || debt.type === 'owed_to_me') ? 'bg-success/10 text-success' : 'bg-destructive/10 text-destructive'}">
								{(debt.debt_type === 'owed_to_me' || debt.type === 'owed_to_me') ? 'Owed to Me' : 'I Owe'}
							</span>
							<span class="text-sm text-muted-foreground">Debt type cannot be changed after creation</span>
						</div>
					</div>

					<!-- Amount and Currency -->
					<div class="grid grid-cols-1 md:grid-cols-3 gap-4">
						<div class="md:col-span-2">
							<label for="total-amount" class="label">Total Amount *</label>
							<input
								id="total-amount"
								type="number"
								bind:value={formData.totalAmount}
								step="0.01"
								min="0"
								class="input {errors.totalAmount ? 'border-destructive focus:border-destructive focus:ring-destructive' : ''} {formData.totalAmount !== originalData.totalAmount ? 'border-warning' : ''}"
								placeholder="0.00"
								disabled={isLoading}
								required
							/>
							{#if errors.totalAmount}
								<p class="mt-1 text-sm text-destructive">{errors.totalAmount}</p>
							{/if}
							{#if formData.totalAmount && !errors.totalAmount}
								<p class="mt-1 text-sm text-muted-foreground">Amount: {formatCurrency(parseFloat(formData.totalAmount))}</p>
							{/if}
						</div>
						
						<div>
							<label for="currency" class="label">Currency</label>
							<select 
								id="currency" 
								bind:value={formData.currency} 
								class="input {formData.currency !== originalData.currency ? 'border-warning' : ''}" 
								disabled={isLoading}
							>
								<option value="PHP">PHP (₱)</option>
								<option value="USD">USD ($)</option>
								<option value="EUR">EUR (€)</option>
								<option value="GBP">GBP (£)</option>
							</select>
						</div>
					</div>

					<!-- Description -->
					<div>
						<label for="description" class="label">Description</label>
						<textarea
							id="description"
							bind:value={formData.description}
							rows="3"
							maxlength={DESCRIPTION_LIMIT}
							class="input resize-none {errors.description ? 'border-destructive focus:border-destructive focus:ring-destructive' : ''} {formData.description !== originalData.description ? 'border-warning' : ''}"
							placeholder="What is this debt for?"
							disabled={isLoading}
						></textarea>
						<div class="mt-1 flex justify-between text-xs">
							<span class={errors.description ? 'text-destructive' : 'text-muted-foreground'}>
								{errors.description || 'Optional but recommended for tracking purposes'}
							</span>
							<span class="text-muted-foreground">
								{formData.description.length}/{DESCRIPTION_LIMIT}
							</span>
						</div>
					</div>
				</div>

				<!-- Payment Schedule -->
				<div class="space-y-4">
					<h3 class="text-sm font-medium text-foreground">Payment Schedule</h3>
					
					<!-- Payment Type (Read-only) -->
					<div class="bg-muted/50 rounded-lg p-4">
						<div class="block text-sm font-medium text-muted-foreground mb-2">Payment Type</div>
						<div class="flex items-center justify-between">
							<span class="font-medium text-foreground">
								{isInstallmentPlan ? 'Installment Plan' : 'One-time Payment'}
							</span>
							<span class="text-sm text-muted-foreground">Payment type cannot be changed after creation</span>
						</div>
					</div>

					{#if isInstallmentPlan}
						<!-- Installment Plan Fields -->
						<div class="grid grid-cols-1 md:grid-cols-2 gap-4">
							<div>
								<label for="installment-plan" class="label">Payment Frequency *</label>
								<select 
									id="installment-plan" 
									bind:value={formData.installmentPlan} 
									class="input {formData.installmentPlan !== originalData.installmentPlan ? 'border-warning' : ''}"
									disabled={isLoading}
								>
									<option value="weekly">Weekly</option>
									<option value="biweekly">Bi-weekly (Every 2 weeks)</option>
									<option value="monthly">Monthly</option>
									<option value="quarterly">Quarterly (Every 3 months)</option>
									<option value="yearly">Yearly</option>
								</select>
							</div>

							<div>
								<label for="number-of-payments" class="label">Number of Payments *</label>
								<input
									id="number-of-payments"
									type="number"
									bind:value={formData.numberOfPayments}
									min="1"
									max="1000"
									class="input {errors.numberOfPayments ? 'border-destructive focus:border-destructive focus:ring-destructive' : ''} {formData.numberOfPayments !== originalData.numberOfPayments ? 'border-warning' : ''}"
									placeholder="Enter number of payments"
									disabled={isLoading}
									required
								/>
								{#if errors.numberOfPayments}
									<p class="mt-1 text-sm text-destructive">{errors.numberOfPayments}</p>
								{/if}
							</div>
						</div>

						<!-- Installment Amount Display -->
						{#if formData.totalAmount && formData.numberOfPayments > 0}
							<div class="bg-primary/10 border border-primary/20 rounded-lg p-4">
								<div class="flex items-center justify-between mb-2">
									<span class="text-sm font-medium text-primary">Installment Amount</span>
									<span class="text-lg font-semibold text-primary">
										{formatCurrency(installmentAmount)}
									</span>
								</div>
								<p class="text-sm text-primary/80">
									{formatCurrency(installmentAmount)} every {getInstallmentFrequencyText(formData.installmentPlan)}
									for {formData.numberOfPayments} payments
								</p>
							</div>
						{/if}

						{#if formData.numberOfPayments !== originalData.numberOfPayments || formData.installmentPlan !== originalData.installmentPlan}
							<div class="bg-warning/10 border border-warning/20 text-warning px-4 py-3 rounded-lg text-sm">
								<div class="flex items-center">
									<svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
										<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.865-.833-2.635 0L4.178 16.5c-.77.833.192 2.5 1.732 2.5z"></path>
									</svg>
									Changing the payment schedule may affect existing payment records
								</div>
							</div>
						{/if}
					{/if}

					<!-- Due Date -->
					<div>
						<DatePicker
							id="due-date"
							bind:value={formData.dueDate}
							label={isInstallmentPlan ? 'Next Payment Due Date' : 'Due Date'}
							placeholder="Select due date"
							required={true}
							error={errors.dueDate}
							disabled={isLoading}
						/>
					</div>
				</div>

				<!-- Additional Information -->
				<div class="space-y-4">
					<h3 class="text-sm font-medium text-foreground">Additional Information</h3>
					
					<div>
						<label for="notes" class="label">Notes</label>
						<textarea
							id="notes"
							bind:value={formData.notes}
							rows="4"
							maxlength={NOTES_LIMIT}
							class="input resize-none {errors.notes ? 'border-destructive focus:border-destructive focus:ring-destructive' : ''} {formData.notes !== originalData.notes ? 'border-warning' : ''}"
							placeholder="Any additional notes about this debt..."
							disabled={isLoading}
						></textarea>
						<div class="mt-1 flex justify-between text-xs">
							<span class={errors.notes ? 'text-destructive' : 'text-muted-foreground'}>
								{errors.notes || 'Optional field for any additional context'}
							</span>
							<span class="text-muted-foreground">
								{formData.notes.length}/{NOTES_LIMIT}
							</span>
						</div>
					</div>
				</div>
			</form>
		</div>

		<!-- Footer Actions -->
		<div class="px-6 py-4 border-t border-border flex items-center justify-between">
			<button
				type="button"
				on:click={handleReset}
				class="btn-secondary"
				disabled={isLoading || !hasChanges}
			>
				<svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"></path>
				</svg>
				Reset
			</button>
			
			<div class="flex items-center space-x-3">
				<button
					type="button"
					on:click={handleClose}
					class="btn-secondary"
					disabled={isLoading}
				>
					Cancel
				</button>
				<button
					on:click={handleSubmit}
					class="btn-primary disabled:opacity-50 disabled:cursor-not-allowed"
					disabled={isLoading || !hasChanges}
				>
					{#if isLoading}
						<svg class="animate-spin -ml-1 mr-3 h-5 w-5 text-primary-foreground" fill="none" viewBox="0 0 24 24">
							<circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
							<path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
						</svg>
						Saving Changes...
					{:else}
						<svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"></path>
						</svg>
						Save Changes
					{/if}
				</button>
			</div>
		</div>
	</div>
</div>

<!-- Reset Confirmation Dialog -->
{#if showResetConfirmation}
	<div class="fixed inset-0 bg-black/50 z-[60] flex items-center justify-center p-4" on:click={() => showResetConfirmation = false} on:keydown={(e) => e.key === 'Escape' && (showResetConfirmation = false)} role="dialog" aria-modal="true">
		<div class="bg-card rounded-lg shadow-medium max-w-sm w-full p-6" on:click|stopPropagation on:keydown|stopPropagation role="document">
			<h3 class="text-lg font-semibold text-foreground mb-2">Reset Form?</h3>
			<p class="text-muted-foreground mb-4">This will discard all changes and restore the original values.</p>
			<div class="flex justify-end space-x-3">
				<button on:click={() => showResetConfirmation = false} class="btn-secondary">
					Keep Changes
				</button>
				<button on:click={confirmReset} class="btn-danger">
					Reset Form
				</button>
			</div>
		</div>
	</div>
{/if}
