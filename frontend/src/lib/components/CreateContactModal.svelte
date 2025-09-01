<script lang="ts">
	import { createEventDispatcher, onMount } from 'svelte';

	const dispatch = createEventDispatcher();

	let formData = {
		name: '',
		email: '',
		phone: '',
		facebookId: '',
		notes: ''
	};

	let errors: { [key: string]: string } = {};
	let isLoading = false;

	onMount(() => {
		// Prevent body scroll when modal is open
		document.body.style.overflow = 'hidden';
		return () => {
			document.body.style.overflow = 'auto';
		};
	});

	function validateForm(): boolean {
		errors = {};

		// Required field validation
		if (!formData.name.trim()) {
			errors.name = 'Name is required';
		}

		// Email validation
		if (formData.email && !/\S+@\S+\.\S+/.test(formData.email)) {
			errors.email = 'Please enter a valid email address';
		}

		// Phone validation (basic)
		if (formData.phone && !/^[\d\s\-\+\(\)]+$/.test(formData.phone)) {
			errors.phone = 'Please enter a valid phone number';
		}

		return Object.keys(errors).length === 0;
	}

	async function handleSubmit() {
		if (!validateForm()) return;

		isLoading = true;

		try {
			// TODO: Replace with actual API call
			await new Promise(resolve => setTimeout(resolve, 1000));

			// Mock successful contact creation
			const newContact = {
				id: Date.now(),
				name: formData.name.trim(),
				email: formData.email.trim() || null,
				phone: formData.phone.trim() || null,
				facebookId: formData.facebookId.trim() || null,
				notes: formData.notes.trim() || null,
				createdAt: new Date().toISOString()
			};

			dispatch('contact-created', newContact);
			
		} catch (error) {
			console.error('Error creating contact:', error);
			errors.general = 'Failed to create contact. Please try again.';
		} finally {
			isLoading = false;
		}
	}

	function handleClose() {
		dispatch('close');
	}

	function handleKeydown(event: KeyboardEvent) {
		if (event.key === 'Escape') {
			handleClose();
		}
	}
</script>

<svelte:window on:keydown={handleKeydown} />

<!-- Modal Backdrop -->
<div class="fixed inset-0 bg-black/50 z-50 flex items-center justify-center p-4" on:click={handleClose}>
	<!-- Modal Content -->
	<div 
		class="bg-card rounded-xl shadow-medium max-w-md w-full max-h-[90vh] overflow-hidden"
		on:click|stopPropagation
	>
		<!-- Header -->
		<div class="px-6 py-4 border-b border-border flex items-center justify-between">
			<h2 class="text-xl font-semibold text-foreground">Create New Contact</h2>
			<button on:click={handleClose} class="text-muted-foreground hover:text-foreground">
				<svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
				</svg>
			</button>
		</div>

		<!-- Form -->
		<form on:submit|preventDefault={handleSubmit} class="p-6 space-y-4">
			<!-- General Error -->
			{#if errors.general}
				<div class="bg-destructive/10 border border-destructive/20 text-destructive px-4 py-3 rounded-lg text-sm">
					{errors.general}
				</div>
			{/if}

			<!-- Name Field -->
			<div>
				<label for="contact-name" class="label">Name *</label>
				<input
					id="contact-name"
					type="text"
					bind:value={formData.name}
					class="input {errors.name ? 'border-destructive focus:border-destructive focus:ring-destructive' : ''}"
					placeholder="Enter contact name"
					disabled={isLoading}
					required
				/>
				{#if errors.name}
					<p class="mt-1 text-sm text-destructive">{errors.name}</p>
				{/if}
			</div>

			<!-- Email Field -->
			<div>
				<label for="contact-email" class="label">Email</label>
				<input
					id="contact-email"
					type="email"
					bind:value={formData.email}
					class="input {errors.email ? 'border-destructive focus:border-destructive focus:ring-destructive' : ''}"
					placeholder="Enter email address"
					disabled={isLoading}
				/>
				{#if errors.email}
					<p class="mt-1 text-sm text-destructive">{errors.email}</p>
				{/if}
			</div>

			<!-- Phone Field -->
			<div>
				<label for="contact-phone" class="label">Phone</label>
				<input
					id="contact-phone"
					type="tel"
					bind:value={formData.phone}
					class="input {errors.phone ? 'border-destructive focus:border-destructive focus:ring-destructive' : ''}"
					placeholder="Enter phone number"
					disabled={isLoading}
				/>
				{#if errors.phone}
					<p class="mt-1 text-sm text-destructive">{errors.phone}</p>
				{/if}
			</div>

			<!-- Facebook ID Field -->
			<div>
				<label for="contact-facebook" class="label">Facebook ID</label>
				<input
					id="contact-facebook"
					type="text"
					bind:value={formData.facebookId}
					class="input"
					placeholder="Enter Facebook ID or username"
					disabled={isLoading}
				/>
			</div>

			<!-- Notes Field -->
			<div>
				<label for="contact-notes" class="label">Notes</label>
				<textarea
					id="contact-notes"
					bind:value={formData.notes}
					rows="3"
					class="input resize-none"
					placeholder="Additional notes about this contact..."
					disabled={isLoading}
				></textarea>
				<div class="mt-1 text-xs text-muted-foreground">
					{formData.notes.length}/500 characters
				</div>
			</div>

			<!-- Action Buttons -->
			<div class="flex items-center justify-end space-x-3 pt-4">
				<button
					type="button"
					on:click={handleClose}
					class="btn-secondary"
					disabled={isLoading}
				>
					Cancel
				</button>
				<button
					type="submit"
					class="btn-primary disabled:opacity-50 disabled:cursor-not-allowed"
					disabled={isLoading || !formData.name.trim()}
				>
					{#if isLoading}
						<svg class="animate-spin -ml-1 mr-3 h-4 w-4 text-primary-foreground" fill="none" viewBox="0 0 24 24">
							<circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
							<path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
						</svg>
						Creating...
					{:else}
						<svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6"></path>
						</svg>
						Create Contact
					{/if}
				</button>
			</div>
		</form>
	</div>
</div>
