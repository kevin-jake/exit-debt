<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import ContactSelector from './ContactSelector.svelte';
	import PaymentScheduleSection from './PaymentScheduleSection.svelte';
	import { apiClient, type Contact, type CreateDebtListRequest } from '../api';
	import { debtsStore } from '../stores/debts';

	// Form data
	let formData = {
		contact: null as Contact | null,
		debtType: 'i_owe' as 'i_owe' | 'owed_to_me',
		totalAmount: '',
		currency: 'PHP',
		description: '',
		paymentType: 'onetime' as 'onetime' | 'installment',
		dueDate: '',
		installmentPlan: 'monthly',
		numberOfPayments: 1,
		installmentCalculationMethod: 'number_of_payments' as 'number_of_payments' | 'due_date',
		finalDueDate: '',
		notes: ''
	};

	// Validation and state
	let errors: { [key: string]: string } = {};
	let isLoading = false;
	let contacts: any[] = [];

	// Character limits
	const DESCRIPTION_LIMIT = 500;
	const NOTES_LIMIT = 1000;

	onMount(() => {
		loadContacts();
	});

	async function loadContacts() {
		try {
			contacts = await apiClient.getContacts();
		} catch (error) {
			console.error('Error loading contacts:', error);
			errors.general = 'Failed to load contacts. Please try again.';
		}
	}

	function validateForm(): boolean {
		errors = {};

		// Contact validation
		if (!formData.contact) {
			errors.contact = 'Please select a contact';
		}

		// Amount validation
		const amount = parseFloat(formData.totalAmount);
		if (!formData.totalAmount || isNaN(amount) || amount <= 0) {
			errors.totalAmount = 'Please enter a valid amount greater than 0';
		}

		// Due date validation - only required for one-time payments or due_date method
		if (formData.paymentType === 'onetime' || formData.installmentCalculationMethod === 'due_date') {
			if (!formData.dueDate) {
				errors.dueDate = 'Please select a due date';
			} else {
				const selectedDate = new Date(formData.dueDate);
				const today = new Date();
				today.setHours(0, 0, 0, 0);
				
				if (selectedDate <= today) {
					errors.dueDate = 'Due date must be in the future';
				}
			}
		}

		// Installment validation
		if (formData.paymentType === 'installment') {
			// Note: Validation for numberOfPayments is now handled in PaymentScheduleSection
			// The component will ensure either numberOfPayments or dueDate is provided based on user selection
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
			// Prepare the debt data for the API
			const debtData: CreateDebtListRequest = {
				contact_id: formData.contact!.id,
				total_amount: formData.totalAmount.toString(),
				currency: formData.currency,
				debt_type: formData.debtType,
				installment_plan: formData.paymentType === 'installment' ? formData.installmentPlan : 'onetime',
				due_date: formData.dueDate ? new Date(formData.dueDate + 'T00:00:00').toISOString() : undefined,
				number_of_payments: formData.paymentType === 'installment' ? formData.numberOfPayments : undefined,
				description: formData.description || undefined,
				notes: formData.notes || undefined,
			};

			// Create the debt using the store
			const newDebt = await debtsStore.createDebt(debtData);

			console.log('Created debt:', newDebt);

			// Redirect to debts page
			goto('/debts');
			
		} catch (error) {
			console.error('Error creating debt:', error);
			errors.general = 'Failed to create debt. Please try again.';
		} finally {
			isLoading = false;
		}
	}

	function handleContactSelect(event: CustomEvent) {
		formData.contact = event.detail;
		if (errors.contact) delete errors.contact;
	}

	function handlePaymentScheduleChange(event: CustomEvent) {
		const { paymentType, dueDate, installmentPlan, numberOfPayments, installmentCalculationMethod, finalDueDate } = event.detail;
		formData.paymentType = paymentType;
		formData.dueDate = dueDate;
		formData.installmentPlan = installmentPlan;
		formData.numberOfPayments = numberOfPayments;
		formData.installmentCalculationMethod = installmentCalculationMethod || 'number_of_payments';
		formData.finalDueDate = finalDueDate || '';
		
		// Clear related errors
		if (errors.dueDate) delete errors.dueDate;
		if (errors.numberOfPayments) delete errors.numberOfPayments;
	}

	function formatCurrency(amount: string): string {
		const num = parseFloat(amount);
		if (isNaN(num)) return '';
		return new Intl.NumberFormat('en-PH', {
			style: 'currency',
			currency: formData.currency
		}).format(num);
	}

	$: displayAmount = formData.totalAmount ? formatCurrency(formData.totalAmount) : '';
</script>

<div class="max-w-4xl mx-auto">
	<form on:submit|preventDefault={handleSubmit} class="space-y-8">
		<!-- General Error -->
		{#if errors.general}
			<div class="bg-destructive/10 border border-destructive/20 text-destructive px-4 py-3 rounded-lg">
				{errors.general}
			</div>
		{/if}

		<!-- Contact Selection Section -->
		<div class="card p-6">
			<h2 class="text-xl font-semibold text-foreground mb-6">Contact Information</h2>
			<ContactSelector
				bind:selectedContact={formData.contact}
				{contacts}
				on:select={handleContactSelect}
			/>
			{#if errors.contact}
				<p class="mt-2 text-sm text-destructive">{errors.contact}</p>
			{/if}
		</div>

		<!-- Debt Information Section -->
		<div class="card p-6">
			<h2 class="text-xl font-semibold text-foreground mb-6">Debt Information</h2>
			
			<div class="space-y-6">
				<!-- Debt Type -->
				<div>
					<label for="debt-type" class="label">Debt Type *</label>
					<div class="grid grid-cols-1 md:grid-cols-2 gap-4">
						<label class="relative cursor-pointer">
							<input
								type="radio"
								bind:group={formData.debtType}
								value="i_owe"
								class="sr-only"
							/>
							<div class="border-2 rounded-lg p-4 transition-all duration-200 {formData.debtType === 'i_owe' ? 'border-destructive bg-destructive/10' : 'border-border hover:border-border/80'}">
								<div class="flex items-center space-x-3">
									<div class="w-10 h-10 bg-destructive rounded-full flex items-center justify-center">
										<svg class="w-5 h-5 text-destructive-foreground" fill="none" stroke="currentColor" viewBox="0 0 24 24">
											<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 17h8m0 0V9m0 8l-8-8-4 4-6-6"></path>
										</svg>
									</div>
									<div>
										<div class="font-medium text-foreground">I Owe</div>
										<div class="text-sm text-muted-foreground">Money you owe to someone</div>
									</div>
								</div>
							</div>
						</label>

						<label class="relative cursor-pointer">
							<input
								type="radio"
								bind:group={formData.debtType}
								value="owed_to_me"
								class="sr-only"
							/>
							<div class="border-2 rounded-lg p-4 transition-all duration-200 {formData.debtType === 'owed_to_me' ? 'border-success bg-success/10' : 'border-border hover:border-border/80'}">
								<div class="flex items-center space-x-3">
									<div class="w-10 h-10 bg-success rounded-full flex items-center justify-center">
										<svg class="w-5 h-5 text-success-foreground" fill="none" stroke="currentColor" viewBox="0 0 24 24">
											<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 7h8m0 0v8m0-8l-8 8-4-4-6 6"></path>
										</svg>
									</div>
									<div>
										<div class="font-medium text-foreground">Owed to Me</div>
										<div class="text-sm text-muted-foreground">Money someone owes you</div>
									</div>
								</div>
							</div>
						</label>
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
							class="input {errors.totalAmount ? 'border-destructive focus:border-destructive focus:ring-destructive' : ''}"
							placeholder="0.00"
							disabled={isLoading}
							required
						/>
						{#if errors.totalAmount}
							<p class="mt-1 text-sm text-destructive">{errors.totalAmount}</p>
						{/if}
						{#if displayAmount}
							<p class="mt-1 text-sm text-muted-foreground">Amount: {displayAmount}</p>
						{/if}
					</div>
					
					<div>
						<label for="currency" class="label">Currency</label>
						<select id="currency" bind:value={formData.currency} class="input" disabled={isLoading}>
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
						class="input resize-none {errors.description ? 'border-destructive focus:border-destructive focus:ring-destructive' : ''}"
						placeholder="What is this debt for? (e.g., Personal loan, Business expense, Equipment purchase...)"
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
		</div>

		<!-- Payment Schedule Section -->
		<div class="card p-6">
			<PaymentScheduleSection
				bind:paymentType={formData.paymentType}
				bind:dueDate={formData.dueDate}
				bind:installmentPlan={formData.installmentPlan}
				bind:numberOfPayments={formData.numberOfPayments}
				bind:installmentCalculationMethod={formData.installmentCalculationMethod}
				bind:finalDueDate={formData.finalDueDate}
				totalAmount={parseFloat(formData.totalAmount) || 0}
				on:change={handlePaymentScheduleChange}
			/>
			{#if errors.dueDate}
				<p class="mt-2 text-sm text-destructive">{errors.dueDate}</p>
			{/if}
			{#if errors.numberOfPayments}
				<p class="mt-2 text-sm text-destructive">{errors.numberOfPayments}</p>
			{/if}
		</div>

		<!-- Additional Information Section -->
		<div class="card p-6">
			<h2 class="text-xl font-semibold text-foreground mb-6">Additional Information</h2>
			
			<div>
				<label for="notes" class="label">Notes</label>
				<textarea
					id="notes"
					bind:value={formData.notes}
					rows="4"
					maxlength={NOTES_LIMIT}
					class="input resize-none {errors.notes ? 'border-destructive focus:border-destructive focus:ring-destructive' : ''}"
					placeholder="Any additional notes about this debt, payment terms, or special arrangements..."
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

		<!-- Form Actions -->
		<div class="flex items-center justify-between pt-6">
			<a href="/debts" class="btn-secondary">
				<svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 19l-7-7m0 0l7-7m-7 7h18"></path>
				</svg>
				Cancel
			</a>
			
			<button
				type="submit"
				class="btn-primary disabled:opacity-50 disabled:cursor-not-allowed"
				disabled={isLoading}
			>
				{#if isLoading}
					<svg class="animate-spin -ml-1 mr-3 h-5 w-5 text-primary-foreground" fill="none" viewBox="0 0 24 24">
						<circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
						<path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
					</svg>
					Creating Debt...
				{:else}
					<svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6"></path>
					</svg>
					Create Debt
				{/if}
			</button>
		</div>
	</form>
</div>
