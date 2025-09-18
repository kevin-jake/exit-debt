<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import DatePicker from './DatePicker.svelte';

	export let paymentType: 'onetime' | 'installment' = 'onetime';
	export let dueDate: string = '';
	export let installmentPlan: string = 'monthly';
	export let numberOfPayments: number = 1;
	export let totalAmount: number = 0;
	export let installmentCalculationMethod: 'number_of_payments' | 'due_date' = 'number_of_payments';
	export let finalDueDate: string = '';

	const dispatch = createEventDispatcher();

	function calculateNumberOfPaymentsFromDueDate(): number {
		if (!finalDueDate) return 0;
		
		const today = new Date();
		const endDate = new Date(finalDueDate);
		
		// Ensure end date is in the future
		if (endDate <= today) return 0;
		
		const diffTime = endDate.getTime() - today.getTime();
		const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24));
		
		// Calculate number of payments based on frequency
		let payments: number;
		switch (installmentPlan) {
			case 'weekly':
				payments = Math.ceil(diffDays / 7);
				break;
			case 'biweekly':
				payments = Math.ceil(diffDays / 14);
				break;
			case 'monthly':
				// More accurate monthly calculation
				const monthsDiff = (endDate.getFullYear() - today.getFullYear()) * 12 + 
								  (endDate.getMonth() - today.getMonth());
				payments = Math.max(1, monthsDiff);
				break;
			case 'quarterly':
				// More accurate quarterly calculation (every 3 months)
				const quartersDiff = (endDate.getFullYear() - today.getFullYear()) * 4 + 
									Math.floor((endDate.getMonth() - today.getMonth()) / 3);
				payments = Math.max(1, quartersDiff);
				break;
			case 'yearly':
				// More accurate yearly calculation
				const yearsDiff = endDate.getFullYear() - today.getFullYear();
				payments = Math.max(1, yearsDiff);
				break;
			default:
				payments = 0;
		}
		
		// Ensure we have at least 1 payment
		return Math.max(1, payments);
	}

	// Calculate numberOfPayments when using due date method
	$: calculatedNumberOfPayments = (() => {
		if (installmentCalculationMethod === 'due_date' && finalDueDate) {
			return calculateNumberOfPaymentsFromDueDate();
		}
		return numberOfPayments;
	})();

	$: {
		// Dispatch changes to parent
		dispatch('change', {
			paymentType,
			dueDate: installmentCalculationMethod === 'number_of_payments' ? dueDate : finalDueDate,
			installmentPlan,
			numberOfPayments: calculatedNumberOfPayments,
			installmentCalculationMethod,
			finalDueDate
		});
	}

	// Reactive statement for installment amount - explicitly track all dependencies
	$: installmentAmount = (() => {
		if (paymentType !== 'installment' || totalAmount <= 0) {
			return 0;
		}

		return calculatedNumberOfPayments > 0 ? totalAmount / calculatedNumberOfPayments : 0;
	})();

	// Reactive statement for payment schedule preview - explicitly track all dependencies
	$: nextPaymentDates = (() => {
		if (paymentType !== 'installment' || calculatedNumberOfPayments <= 0) {
			return [];
		}

		const dates: string[] = [];
		
		// Calculate payment dates based on calculation method
		if (installmentCalculationMethod === 'number_of_payments') {
			// For number_of_payments method, show first few payments
			const firstPaymentDate = new Date(getTodayDate());
			
			for (let i = 0; i < Math.min(calculatedNumberOfPayments, 5); i++) {
				const paymentDate = new Date(firstPaymentDate);
				
				switch (installmentPlan) {
					case 'weekly':
						paymentDate.setDate(firstPaymentDate.getDate() + (i * 7));
						break;
					case 'biweekly':
						paymentDate.setDate(firstPaymentDate.getDate() + (i * 14));
						break;
					case 'monthly':
						paymentDate.setMonth(firstPaymentDate.getMonth() + i);
						break;
					case 'quarterly':
						paymentDate.setMonth(firstPaymentDate.getMonth() + (i * 3));
						break;
					case 'yearly':
						paymentDate.setFullYear(firstPaymentDate.getFullYear() + i);
						break;
				}
				
				dates.push(paymentDate.toLocaleDateString('en-US', {
					month: 'short',
					day: 'numeric',
					year: 'numeric'
				}));
			}
		} else {
			// For due_date method, work backwards from the final payment date
			if (!finalDueDate) return [];
			
			const finalPaymentDate = new Date(finalDueDate);
			
			// Calculate payment dates working backwards from the final date
			for (let i = 0; i < Math.min(calculatedNumberOfPayments, 5); i++) {
				const paymentDate = new Date(finalPaymentDate);
				
				switch (installmentPlan) {
					case 'weekly':
						// Go back by (calculatedNumberOfPayments - 1 - i) weeks
						paymentDate.setDate(finalPaymentDate.getDate() - ((calculatedNumberOfPayments - 1 - i) * 7));
						break;
					case 'biweekly':
						// Go back by (calculatedNumberOfPayments - 1 - i) * 2 weeks
						paymentDate.setDate(finalPaymentDate.getDate() - ((calculatedNumberOfPayments - 1 - i) * 14));
						break;
					case 'monthly':
						// Go back by (calculatedNumberOfPayments - 1 - i) months
						paymentDate.setMonth(finalPaymentDate.getMonth() - (calculatedNumberOfPayments - 1 - i));
						break;
					case 'quarterly':
						// Go back by (calculatedNumberOfPayments - 1 - i) * 3 months
						paymentDate.setMonth(finalPaymentDate.getMonth() - ((calculatedNumberOfPayments - 1 - i) * 3));
						break;
					case 'yearly':
						// Go back by (calculatedNumberOfPayments - 1 - i) years
						paymentDate.setFullYear(finalPaymentDate.getFullYear() - (calculatedNumberOfPayments - 1 - i));
						break;
				}
				
				dates.push(paymentDate.toLocaleDateString('en-US', {
					month: 'short',
					day: 'numeric',
					year: 'numeric'
				}));
			}
		}
		
		return dates;
	})();

	// Calculate the final payment date for display
	$: computedFinalDueDate = (() => {
		if (installmentCalculationMethod === 'number_of_payments' && calculatedNumberOfPayments > 0) {
			const firstPaymentDate = new Date(getTodayDate());
			const finalPaymentDate = new Date(firstPaymentDate);
			
			switch (installmentPlan) {
				case 'weekly':
					finalPaymentDate.setDate(firstPaymentDate.getDate() + ((calculatedNumberOfPayments - 1) * 7));
					break;
				case 'biweekly':
					finalPaymentDate.setDate(firstPaymentDate.getDate() + ((calculatedNumberOfPayments - 1) * 14));
					break;
				case 'monthly':
					finalPaymentDate.setMonth(firstPaymentDate.getMonth() + (calculatedNumberOfPayments - 1));
					break;
				case 'quarterly':
					finalPaymentDate.setMonth(firstPaymentDate.getMonth() + ((calculatedNumberOfPayments - 1) * 3));
					break;
				case 'yearly':
					finalPaymentDate.setFullYear(firstPaymentDate.getFullYear() + (calculatedNumberOfPayments - 1));
					break;
			}
			
			return finalPaymentDate.toLocaleDateString('en-US', {
				month: 'short',
				day: 'numeric',
				year: 'numeric'
			});
		}
		return '';
	})();

	function formatCurrency(amount: number): string {
		return new Intl.NumberFormat('en-PH', {
			style: 'currency',
			currency: 'PHP'
		}).format(amount);
	}

	function getTodayDate(): string {
		const today = new Date();
		today.setDate(today.getDate() + 1); // Tomorrow as minimum date
		return today.toDateString();
	}

	// Reactive minDate calculations based on payment frequency
	$: minDateForDueDate = (() => {
		const today = new Date();
		
		// Calculate minimum date based on payment frequency
		switch (installmentPlan) {
			case 'weekly':
				// Minimum is next week
				today.setDate(today.getDate() + 7);
				break;
			case 'biweekly':
				// Minimum is in 2 weeks
				today.setDate(today.getDate() + 14);
				break;
			case 'monthly':
				// Minimum is next month
				today.setMonth(today.getMonth() + 1);
				break;
			case 'quarterly':
				// Minimum is in 3 months
				today.setMonth(today.getMonth() + 3);
				break;
			case 'yearly':
				// Minimum is next year
				today.setFullYear(today.getFullYear() + 1);
				break;
			default:
				// Default to tomorrow
				today.setDate(today.getDate() + 1);
		}
		
		return today.toDateString();
	})();

	$: minDateForFinalDueDate = (() => {
		const today = new Date();
		
		// For final due date, we need to ensure it's at least one payment period away
		// and aligns with the payment frequency
		switch (installmentPlan) {
			case 'weekly':
				// Minimum is next week
				today.setDate(today.getDate() + 7);
				break;
			case 'biweekly':
				// Minimum is in 2 weeks
				today.setDate(today.getDate() + 14);
				break;
			case 'monthly':
				// Minimum is next month
				today.setMonth(today.getMonth() + 1);
				break;
			case 'quarterly':
				// Minimum is in 3 months
				today.setMonth(today.getMonth() + 3);
				break;
			case 'yearly':
				// Minimum is next year
				today.setFullYear(today.getFullYear() + 1);
				break;
			default:
				// Default to tomorrow
				today.setDate(today.getDate() + 1);
		}
		
		return today.toDateString();
	})();

	// Create a key that changes when minDate changes to force DatePicker re-render
	$: datePickerKey = `${minDateForDueDate}-${minDateForFinalDueDate}-${installmentPlan}`;

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

	$: isValidForInstallmentAmount = (() => {
		if (paymentType !== 'installment') return false;
		if (totalAmount <= 0) return false;
		if (calculatedNumberOfPayments <= 0) return false;
		
		// For number_of_payments method, only need numberOfPayments > 0
		// For due_date method, need finalDueDate to be set
		if (installmentCalculationMethod === 'number_of_payments') {
			return numberOfPayments > 0;
		} else if (installmentCalculationMethod === 'due_date') {
			return finalDueDate !== '' && finalDueDate.length > 0;
		}
		
		return false;
	})();

	$: isValidForPaymentSchedule = (() => {
		if (paymentType !== 'installment') return false;
		if (calculatedNumberOfPayments <= 0) return false;
		if (totalAmount <= 0) return false;
		
		return true;
	})();
</script>

<div class="space-y-6">
	<div>
		<h3 id="payment-schedule-heading" class="text-lg font-medium text-foreground mb-4">Payment Schedule</h3>
	</div>

	<!-- Payment Type Selection -->
	<div>
		<label for="payment-type" class="label">Payment Type *</label>
		<div class="space-y-3">
			<label class="flex items-center space-x-3 cursor-pointer">
				<input
					type="radio"
					bind:group={paymentType}
					value="onetime"
					class="w-4 h-4 text-primary focus:ring-primary border-border"
				/>
				<div class="flex-1">
					<div class="font-medium text-foreground">One-time Payment</div>
					<div class="text-sm text-muted-foreground">Single payment due on a specific date</div>
				</div>
			</label>
			
			<label class="flex items-center space-x-3 cursor-pointer">
				<input
					type="radio"
					bind:group={paymentType}
					value="installment"
					class="w-4 h-4 text-primary focus:ring-primary border-border"
				/>
				<div class="flex-1">
					<div class="font-medium text-foreground">Installment Plan</div>
					<div class="text-sm text-muted-foreground">Split into multiple payments over time</div>
				</div>
			</label>
		</div>
	</div>

	<!-- One-time Payment Fields -->
	{#if paymentType === 'onetime'}
		<div class="animate-fade-in">
			{#key datePickerKey}
				<DatePicker
					id="due-date"
					bind:value={dueDate}
					label="Due Date"
					placeholder="Select due date"
					required={true}
					minDate={getTodayDate()}
				/>
			{/key}
			<p class="mt-1 text-sm text-muted-foreground">
				The date when the full amount is due
			</p>
		</div>
	{/if}

	<!-- Installment Plan Fields -->
	{#if paymentType === 'installment'}
		<div class="animate-fade-in space-y-4">
			<!-- Installment Frequency -->
			<div>
				<label for="installment-plan" class="label">Payment Frequency *</label>
				<select id="installment-plan" bind:value={installmentPlan} class="input">
					<option value="weekly">Weekly</option>
					<option value="biweekly">Bi-weekly (Every 2 weeks)</option>
					<option value="monthly">Monthly</option>
					<option value="quarterly">Quarterly (Every 3 months)</option>
					<option value="yearly">Yearly</option>
				</select>
			</div>

			<!-- Calculation Method Selection -->
			<div>
				<label for="calculation-method" class="label">How would you like to set up your payment plan? *</label>
				<div class="space-y-3">
					<label class="flex items-center space-x-3 cursor-pointer">
						<input
							type="radio"
							bind:group={installmentCalculationMethod}
							value="number_of_payments"
							class="w-4 h-4 text-primary focus:ring-primary border-border"
						/>
						<div class="flex-1">
							<div class="font-medium text-foreground">Set Number of Payments</div>
							<div class="text-sm text-muted-foreground">Specify how many payments you want to make</div>
						</div>
					</label>
					
					<label class="flex items-center space-x-3 cursor-pointer">
						<input
							type="radio"
							bind:group={installmentCalculationMethod}
							value="due_date"
							class="w-4 h-4 text-primary focus:ring-primary border-border"
						/>
						<div class="flex-1">
							<div class="font-medium text-foreground">Set Final Due Date</div>
							<div class="text-sm text-muted-foreground">Specify when you want to finish paying</div>
						</div>
					</label>
				</div>
			</div>

			<!-- Number of Payments Field (when method is number_of_payments) -->
			{#if installmentCalculationMethod === 'number_of_payments'}
				<div class="animate-fade-in">
					<label for="number-of-payments" class="label">Number of Payments *</label>
					<input
						id="number-of-payments"
						type="number"
						bind:value={numberOfPayments}
						min="1"
						max="1000"
						class="input"
						placeholder="Enter number of payments"
						required
					/>
					<p class="mt-1 text-sm text-muted-foreground">
						Total number of installment payments
					</p>
				</div>
			{/if}

			<!-- Due Date Fields (when method is due_date) -->
			{#if installmentCalculationMethod === 'due_date'}
				<div class="animate-fade-in">
					{#key datePickerKey}
						<DatePicker
							id="final-due-date"
							bind:value={finalDueDate}
							label="Final Payment Date"
							placeholder="Select final payment date"
							required={true}
							minDate={minDateForFinalDueDate}
						/>
					{/key}
					<p class="mt-1 text-sm text-muted-foreground">
						Date when you want to finish paying
					</p>
				</div>
			{/if}

			<!-- Installment Amount Display -->
			{#if isValidForInstallmentAmount}
				<div class="bg-primary/10 border border-primary/20 rounded-lg p-4">
					<div class="flex items-center justify-between mb-2">
						<span class="text-sm font-medium text-primary">Installment Amount</span>
						<span class="text-lg font-semibold text-primary">
							{formatCurrency(installmentAmount)}
						</span>
					</div>
					<p class="text-sm text-primary/80">
						{formatCurrency(installmentAmount)} every {getInstallmentFrequencyText(installmentPlan)}
						for {calculatedNumberOfPayments} payments
					</p>
					{#if installmentCalculationMethod === 'due_date'}
						<p class="text-xs text-primary/60 mt-1">
							Calculated from your payment schedule
						</p>
					{/if}
				</div>
			{/if}

			<!-- Payment Schedule Preview -->
			{#if isValidForPaymentSchedule && nextPaymentDates.length > 0}
				<div class="bg-muted/50 border border-border rounded-lg p-4">
					<h4 class="text-sm font-medium text-foreground mb-3">Payment Schedule Preview</h4>
					<div class="space-y-2">
						{#each nextPaymentDates as date, index}
							<div class="flex items-center justify-between text-sm">
								<div class="flex flex-col">
									<span class="text-muted-foreground">Payment {index + 1}</span>
									<span class="font-medium text-foreground">{date}</span>
								</div>
								<span class="font-semibold text-primary">
									{formatCurrency(installmentAmount)}
								</span>
							</div>
						{/each}
						{#if calculatedNumberOfPayments > 5}
						<div class="text-sm text-muted-foreground/60 text-center pt-2 border-t border-border">
							... and {calculatedNumberOfPayments - 5} more payments of {formatCurrency(installmentAmount)} each
						</div>
						{/if}
						{#if computedFinalDueDate}
							<div class="text-sm text-muted-foreground/60 text-center pt-2 border-t border-border">
								Final payment date: {computedFinalDueDate}
							</div>
						{/if}
					</div>
				</div>
			{/if}
		</div>
	{/if}
</div>

<style>
	@keyframes fade-in {
		from {
			opacity: 0;
			transform: translateY(-10px);
		}
		to {
			opacity: 1;
			transform: translateY(0);
		}
	}

	.animate-fade-in {
		animation: fade-in 0.3s ease-out;
	}
</style>
