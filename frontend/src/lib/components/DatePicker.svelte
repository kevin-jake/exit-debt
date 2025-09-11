<script lang="ts">
	import { createEventDispatcher, onMount, onDestroy, afterUpdate } from 'svelte';
	import AirDatepicker from 'air-datepicker';
	import 'air-datepicker/air-datepicker.css';
	import localeEn from 'air-datepicker/locale/en';

	export let value: string = '';
	export let placeholder: string = 'Select date';
	export let label: string = '';
	export let required: boolean = false;
	export let disabled: boolean = false;
	export let minDate: string = '';
	export let maxDate: string = '';
	export let error: string = '';
	export let id: string = '';

	const dispatch = createEventDispatcher();

	let inputElement: HTMLInputElement;
	let containerElement: HTMLDivElement | undefined;
	let datepicker: any = null;

	onMount(() => {
		// Wait for next tick to ensure DOM is ready
		setTimeout(() => {
			if (inputElement && containerElement) {
				console.log('Initializing Air Datepicker on:', inputElement, 'container:', containerElement);

				try {
					// Initialize Air Datepicker with container
					datepicker = new AirDatepicker(inputElement, {
						container: containerElement,
						locale: localeEn,
						dateFormat: 'yyyy-MM-dd',
						autoClose: true,
						minDate: minDate ? minDate : undefined,
						maxDate: maxDate ? maxDate : undefined,
						onRenderCell: ({ date, cellType }) => {
							// Disable dates before minDate
							if (cellType === 'day' && minDate) {
								const minDateObj = new Date(minDate);
								if (date < minDateObj) {
									return {
										disabled: true
									};
								}
							}
							// Disable dates after maxDate
							if (cellType === 'day' && maxDate) {
								const maxDateObj = new Date(maxDate);
								if (date > maxDateObj) {
									return {
										disabled: true
									};
								}
							}
							return {};
						},
						onSelect: ({ date }: { date: any }) => {
							console.log('Date selected:', date);
							if (date) {
								// Handle both single date and array of dates
								const selectedDate = Array.isArray(date) ? date[0] : date;
								// Convert to YYYY-MM-DD format
								const year = selectedDate.getFullYear();
								const month = String(selectedDate.getMonth() + 1).padStart(2, '0');
								const day = String(selectedDate.getDate()).padStart(2, '0');
								const dateString = `${year}-${month}-${day}`;
								value = dateString;
								dispatch('change', dateString);
							}
						},
						onShow: () => {
							console.log('Datepicker shown');
						},
						onHide: () => {
							console.log('Datepicker hidden');
						}
					});

					console.log('Datepicker instance created:', datepicker);

					// Set initial value if provided
					if (value) {
						datepicker.selectDate(new Date(value));
					}
				} catch (error) {
					console.error('Error initializing Air Datepicker:', error);
				}
			} else {
				console.error('Input element or container not found');
			}
		}, 0);
	});

	onDestroy(() => {
		if (datepicker) {
			datepicker.destroy();
		}
	});


	// Update datepicker when value changes externally
	$: if (datepicker && value) {
		const currentValue = inputElement?.value;
		if (currentValue !== value) {
			datepicker.selectDate(new Date(value));
		}
	}

	// Update min/max dates when props change
	$: if (datepicker) {
		if (minDate) {
			datepicker.update('minDate', new Date(minDate));
		}
		if (maxDate) {
			datepicker.update('maxDate', new Date(maxDate));
		}
	}

	function handleInputClick() {
		console.log('Input clicked, datepicker:', datepicker, 'disabled:', disabled);
		if (datepicker && !disabled) {
			console.log('Calling datepicker.show()');
			try {
				datepicker.show();
			} catch (error) {
				console.error('Error showing datepicker:', error);
				// Try alternative method
				try {
					datepicker.toggle();
				} catch (toggleError) {
					console.error('Error toggling datepicker:', toggleError);
				}
			}
		}
	}

	function handleInputFocus() {
		console.log('Input focused, datepicker:', datepicker, 'disabled:', disabled);
		if (datepicker && !disabled) {
			console.log('Calling datepicker.show() on focus');
			try {
				datepicker.show();
			} catch (error) {
				console.error('Error showing datepicker on focus:', error);
			}
		}
	}
</script>

<div class="w-full">
	{#if label}
		<label for={id} class="label">
			{label}
			{#if required}
				<span class="text-destructive">*</span>
			{/if}
		</label>
	{/if}

	<div class="relative" bind:this={containerElement}>
		<input
			bind:this={inputElement}
			{id}
			type="text"
			{value}
			{placeholder}
			{required}
			{disabled}
			on:focus={handleInputFocus}
			on:click={handleInputClick}
			class="input {error ? 'border-destructive focus:border-destructive focus:ring-destructive' : ''} {disabled ? 'opacity-50 cursor-not-allowed' : 'cursor-pointer'}"
			readonly
			style="cursor: pointer;"
		/>
		<div class="absolute right-3 top-1/2 transform -translate-y-1/2 pointer-events-none">
			<svg class="w-5 h-5 text-muted-foreground" fill="none" stroke="currentColor" viewBox="0 0 24 24">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"></path>
			</svg>
		</div>
	</div>

	{#if error}
		<p class="mt-1 text-sm text-destructive">{error}</p>
	{/if}
</div>

<style>
	/* Basic Air Datepicker styling to match your theme */
	:global(.air-datepicker) {
		background: hsl(var(--card));
		border: 1px solid hsl(var(--border));
		border-radius: 0.5rem;
		box-shadow: 0 10px 15px -3px rgba(0, 0, 0, 0.1), 0 4px 6px -2px rgba(0, 0, 0, 0.05);
		z-index: 1000;
	}

	:global(.air-datepicker-nav) {
		background: hsl(var(--muted) / 0.5);
		border-bottom: 1px solid hsl(var(--border));
	}

	:global(.air-datepicker-nav--title) {
		color: hsl(var(--foreground));
		font-weight: 500;
	}
	:global(.air-datepicker-nav--title:hover) {
		color: hsl(var(--foreground));
		background: hsl(var(--muted) / 0.5);
		font-weight: 500;
	}

	:global(.air-datepicker-nav--action) {
		color: hsl(var(--muted-foreground));
	}

	:global(.air-datepicker-nav--action:hover) {
		color: hsl(var(--foreground));
		background: hsl(var(--muted) / 0.5);
		border-radius: 0.25rem;
	}

	:global(.air-datepicker-body--day-name) {
		color: hsl(var(--muted-foreground));
		font-weight: 500;
	}

	:global(.air-datepicker-cell) {
		color: hsl(var(--foreground));
	}

	:global(.air-datepicker-cell:hover) {
		background: hsl(var(--muted) / 0.5);
		border-radius: 0.25rem;
	}

	:global(.air-datepicker-cell.-day-.-other-month-) {
		color: hsl(var(--muted-foreground) / 0.5);
	}

	:global(.air-datepicker-cell.-day-.-current-) {
		color: hsl(var(--primary));
		font-weight: 500;
		border-radius: 0.25rem;
	}

	:global(.air-datepicker-cell.-day-.-selected-) {
		background: hsl(var(--primary));
		color: hsl(var(--primary-foreground));
		font-weight: 500;
		border-radius: 0.25rem;
	}

	:global(.air-datepicker-cell.-day-.-disabled-) {
		color: hsl(var(--muted-foreground) / 0.3);
		cursor: not-allowed;
		background: transparent;
	}

	:global(.air-datepicker-cell.-day-.-disabled-:hover) {
		background: transparent;
		color: hsl(var(--muted-foreground) / 0.3);
	}

	:global(.air-datepicker-buttons) {
		border-top: 1px solid hsl(var(--border));
		background: hsl(var(--muted) / 0.3);
	}

	:global(.air-datepicker-button) {
		color: hsl(var(--primary));
		font-weight: 500;
		padding: 0.5rem 0.75rem;
		border-radius: 0.25rem;
		transition: all 0.2s;
	}

	:global(.air-datepicker-button:hover) {
		color: hsl(var(--primary) / 0.8);
		background: hsl(var(--muted) / 0.5);
	}
</style>