<script lang="ts">
	import { createEventDispatcher } from 'svelte';

	export let paymentType: 'one_time' | 'installment' = 'one_time';
	export let dueDate: string = '';
	export let installmentPlan: string = 'monthly';
	export let numberOfPayments: number = 1;
	export let totalAmount: number = 0;

	const dispatch = createEventDispatcher();

	$: installmentAmount = totalAmount && numberOfPayments > 0 ? totalAmount / numberOfPayments : 0;

	$: {
		// Dispatch changes to parent
		dispatch('change', {
			paymentType,
			dueDate,
			installmentPlan,
			numberOfPayments
		});
	}

	function formatCurrency(amount: number): string {
		return new Intl.NumberFormat('en-PH', {
			style: 'currency',
			currency: 'PHP'
		}).format(amount);
	}

	function getTodayDate(): string {
		const today = new Date();
		today.setDate(today.getDate() + 1); // Tomorrow as minimum date
		return today.toISOString().split('T')[0];
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

	function calculateNextPaymentDates(): string[] {
		if (paymentType !== 'installment' || !dueDate || numberOfPayments <= 0) {
			return [];
		}

		const dates: string[] = [];
		const startDate = new Date(dueDate);
		
		for (let i = 0; i < Math.min(numberOfPayments, 5); i++) { // Show max 5 dates
			const paymentDate = new Date(startDate);
			
			switch (installmentPlan) {
				case 'weekly':
					paymentDate.setDate(startDate.getDate() + (i * 7));
					break;
				case 'biweekly':
					paymentDate.setDate(startDate.getDate() + (i * 14));
					break;
				case 'monthly':
					paymentDate.setMonth(startDate.getMonth() + i);
					break;
				case 'quarterly':
					paymentDate.setMonth(startDate.getMonth() + (i * 3));
					break;
				case 'yearly':
					paymentDate.setFullYear(startDate.getFullYear() + i);
					break;
			}
			
			dates.push(paymentDate.toLocaleDateString('en-US', {
				month: 'short',
				day: 'numeric',
				year: 'numeric'
			}));
		}
		
		return dates;
	}

	$: nextPaymentDates = calculateNextPaymentDates();
</script>

<div class="space-y-6">
	<div>
		<h3 class="text-lg font-medium text-foreground mb-4">Payment Schedule</h3>
	</div>

	<!-- Payment Type Selection -->
	<div>
		<label class="label">Payment Type *</label>
		<div class="space-y-3">
			<label class="flex items-center space-x-3 cursor-pointer">
				<input
					type="radio"
					bind:group={paymentType}
					value="one_time"
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
	{#if paymentType === 'one_time'}
		<div class="animate-fade-in">
			<label for="due-date" class="label">Due Date *</label>
			<input
				id="due-date"
				type="date"
				bind:value={dueDate}
				min={getTodayDate()}
				class="input"
				required
			/>
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

			<!-- Number of Payments -->
			<div>
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

			<!-- First Payment Due Date -->
			<div>
				<label for="first-due-date" class="label">First Payment Due Date *</label>
				<input
					id="first-due-date"
					type="date"
					bind:value={dueDate}
					min={getTodayDate()}
					class="input"
					required
				/>
				<p class="mt-1 text-sm text-muted-foreground">
					Date of the first installment payment
				</p>
			</div>

			<!-- Installment Amount Display -->
			{#if totalAmount > 0 && numberOfPayments > 0}
				<div class="bg-primary/10 border border-primary/20 rounded-lg p-4">
					<div class="flex items-center justify-between mb-2">
						<span class="text-sm font-medium text-primary">Installment Amount</span>
						<span class="text-lg font-semibold text-primary">
							{formatCurrency(installmentAmount)}
						</span>
					</div>
					<p class="text-sm text-primary/80">
						{formatCurrency(installmentAmount)} every {getInstallmentFrequencyText(installmentPlan)}
						for {numberOfPayments} payments
					</p>
				</div>
			{/if}

			<!-- Payment Schedule Preview -->
			{#if nextPaymentDates.length > 0}
				<div class="bg-muted/50 border border-border rounded-lg p-4">
					<h4 class="text-sm font-medium text-foreground mb-3">Payment Schedule Preview</h4>
					<div class="space-y-2">
						{#each nextPaymentDates as date, index}
							<div class="flex items-center justify-between text-sm">
								<span class="text-muted-foreground">Payment {index + 1}</span>
								<span class="font-medium text-foreground">{date}</span>
							</div>
						{/each}
						{#if numberOfPayments > 5}
							<div class="text-sm text-muted-foreground/60 text-center pt-2 border-t border-border">
								... and {numberOfPayments - 5} more payments
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
