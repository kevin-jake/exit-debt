<script lang="ts">
	import { createEventDispatcher, onMount } from 'svelte';

	export let photoUrl: string;
	export let isOpen: boolean = false;

	const dispatch = createEventDispatcher<{ close: void }>();

	let isLoading = true;
	let hasError = false;

	onMount(() => {
		// Prevent body scroll when modal is open
		if (isOpen) {
			document.body.style.overflow = 'hidden';
		}
		return () => {
			document.body.style.overflow = 'auto';
		};
	});

	// Alternative approach: use a more reliable image loading method
	function loadImage() {
		if (!photoUrl) return;
		
		const img = new Image();
		img.onload = () => {
			isLoading = false;
			hasError = false;
		};
		img.onerror = () => {
			isLoading = false;
			hasError = true;
		};
		img.src = photoUrl;
	}

	function handleImageLoad() {
		isLoading = false;
		hasError = false;
	}

	function handleImageError() {
		isLoading = false;
		hasError = true;
	}

	function closeViewer() {
		dispatch('close');
	}

	function handleKeydown(event: KeyboardEvent) {
		if (event.key === 'Escape') {
			closeViewer();
		}
	}

	function handleBackdropClick(event: MouseEvent) {
		if (event.target === event.currentTarget) {
			closeViewer();
		}
	}

	// Watch for changes in isOpen to manage body scroll
	$: if (isOpen) {
		document.body.style.overflow = 'hidden';
		// Reset loading state when modal opens
		isLoading = true;
		hasError = false;
		// Use the new image loading method
		loadImage();
	} else {
		document.body.style.overflow = 'auto';
		// Reset state when modal closes
		isLoading = false;
		hasError = false;
	}

	// Watch for photoUrl changes to reload image
	$: if (photoUrl && isOpen) {
		loadImage();
	}
</script>

<svelte:window on:keydown={handleKeydown} />

{#if isOpen}
	<!-- Modal Backdrop -->
	<div 
		class="fixed inset-0 bg-black/90 z-[70] flex items-center justify-center p-4" 
		role="dialog" 
		aria-modal="true"
		on:click={handleBackdropClick}
	>
		<!-- Modal Content -->
		<div 
			class="relative max-w-[95vw] max-h-[95vh] flex items-center justify-center"
			on:click|stopPropagation
		>
			<!-- Close Button -->
			<button
				on:click={closeViewer}
				class="absolute top-4 right-4 z-10 w-10 h-10 bg-black/50 hover:bg-black/70 text-white rounded-full flex items-center justify-center transition-colors duration-200"
				aria-label="Close photo viewer"
			>
				<svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
				</svg>
			</button>

			<!-- Loading State -->
			{#if isLoading}
				<div class="flex flex-col items-center space-y-4">
					<svg class="animate-spin h-12 w-12 text-white" fill="none" viewBox="0 0 24 24">
						<circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
						<path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
					</svg>
					<p class="text-white text-lg">Loading receipt photo...</p>
				</div>
			{/if}

			<!-- Error State -->
			{#if hasError}
				<div class="flex flex-col items-center space-y-4 text-center">
					<svg class="w-16 h-16 text-red-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.865-.833-2.635 0L4.178 16.5c-.77.833.192 2.5 1.732 2.5z"></path>
					</svg>
					<div class="space-y-2">
						<h3 class="text-white text-xl font-semibold">Failed to load photo</h3>
						<p class="text-white/80">The receipt photo could not be loaded. Please try again.</p>
					</div>
					<button
						on:click={() => { isLoading = true; hasError = false; }}
						class="px-4 py-2 bg-white/20 hover:bg-white/30 text-white rounded-lg transition-colors duration-200"
					>
						Retry
					</button>
				</div>
			{/if}

			<!-- Image Display -->
			{#if !isLoading && !hasError}
				<img
					src={photoUrl}
					alt="Receipt photo"
					class="max-w-full max-h-full object-contain rounded-lg shadow-2xl"
					on:load={handleImageLoad}
					on:error={handleImageError}
				/>
			{/if}
		</div>
	</div>
{/if}
