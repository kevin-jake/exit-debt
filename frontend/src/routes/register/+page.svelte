<script lang="ts">
	import { goto } from '$app/navigation';

	let firstName = $state('');
	let lastName = $state('');
	let email = $state('');
	let phone = $state('');
	let password = $state('');
	let confirmPassword = $state('');
	let isLoading = $state(false);
	let error = $state('');

	async function handleRegister() {
		error = '';

		// Validation
		if (!firstName || !lastName || !email || !password) {
			error = 'Please fill in all required fields';
			return;
		}

		if (password.length < 6) {
			error = 'Password must be at least 6 characters long';
			return;
		}

		if (password !== confirmPassword) {
			error = 'Passwords do not match';
			return;
		}

		if (!/\S+@\S+\.\S+/.test(email)) {
			error = 'Please enter a valid email address';
			return;
		}

		isLoading = true;

		try {
			// TODO: Implement actual registration API call
			// Simulate API call
			await new Promise(resolve => setTimeout(resolve, 1500));
			
			// Mock registration success
			// TODO: Store JWT token
			localStorage.setItem('token', 'mock-jwt-token');
			goto('/');
		} catch (err) {
			error = 'An error occurred during registration. Please try again.';
		} finally {
			isLoading = false;
		}
	}

	function handleKeyPress(event: KeyboardEvent) {
		if (event.key === 'Enter') {
			handleRegister();
		}
	}
</script>

<svelte:head>
	<title>Register - DebtTracker</title>
</svelte:head>

<div class="min-h-screen flex items-center justify-center bg-gradient-to-br from-primary/10 to-muted px-4 py-12">
	<div class="max-w-md w-full space-y-8">
		<!-- Logo and Header -->
		<div class="text-center">
			<div class="mx-auto w-16 h-16 bg-primary rounded-xl flex items-center justify-center mb-6">
				<svg class="w-8 h-8 text-primary-foreground" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1"></path>
				</svg>
			</div>
			<h2 class="text-3xl font-bold text-foreground">Create account</h2>
			<p class="mt-2 text-sm text-muted-foreground">Join DebtTracker to manage your finances</p>
		</div>

		<!-- Registration Form -->
		<div class="card p-8">
			<form onsubmit={handleRegister} class="space-y-6">
				{#if error}
					<div class="bg-destructive/10 border border-destructive/20 text-destructive px-4 py-3 rounded-lg text-sm">
						{error}
					</div>
				{/if}

				<div class="grid grid-cols-2 gap-4">
					<div>
						<label for="firstName" class="label">First name</label>
						<input
							id="firstName"
							type="text"
							bind:value={firstName}
							onkeypress={handleKeyPress}
							class="input"
							placeholder="First name"
							required
						/>
					</div>
					<div>
						<label for="lastName" class="label">Last name</label>
						<input
							id="lastName"
							type="text"
							bind:value={lastName}
							onkeypress={handleKeyPress}
							class="input"
							placeholder="Last name"
							required
						/>
					</div>
				</div>

				<div>
					<label for="email" class="label">Email address</label>
					<input
						id="email"
						type="email"
						bind:value={email}
						onkeypress={handleKeyPress}
						class="input"
						placeholder="Enter your email"
						required
					/>
				</div>

				<div>
					<label for="phone" class="label">Phone number <span class="text-muted-foreground/60">(optional)</span></label>
					<input
						id="phone"
						type="tel"
						bind:value={phone}
						onkeypress={handleKeyPress}
						class="input"
						placeholder="Enter your phone number"
					/>
				</div>

				<div>
					<label for="password" class="label">Password</label>
					<input
						id="password"
						type="password"
						bind:value={password}
						onkeypress={handleKeyPress}
						class="input"
						placeholder="Create a password (min. 6 characters)"
						required
					/>
				</div>

				<div>
					<label for="confirmPassword" class="label">Confirm password</label>
					<input
						id="confirmPassword"
						type="password"
						bind:value={confirmPassword}
						onkeypress={handleKeyPress}
						class="input"
						placeholder="Confirm your password"
						required
					/>
				</div>

				<div class="flex items-center">
					<input
						id="terms"
						type="checkbox"
						class="h-4 w-4 text-primary focus:ring-primary border-input rounded"
						required
					/>
					<label for="terms" class="ml-2 block text-sm text-muted-foreground">
						I agree to the <a href="/terms" class="text-primary hover:text-primary/80">Terms of Service</a> and 
						<a href="/privacy" class="text-primary hover:text-primary/80">Privacy Policy</a>
					</label>
				</div>

				<button
					type="submit"
					disabled={isLoading}
					class="w-full btn-primary disabled:opacity-50 disabled:cursor-not-allowed"
				>
					{#if isLoading}
						<svg class="animate-spin -ml-1 mr-3 h-4 w-4 text-primary-foreground" fill="none" viewBox="0 0 24 24">
							<circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
							<path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
						</svg>
						Creating account...
					{:else}
						Create account
					{/if}
				</button>
			</form>
		</div>

		<!-- Login Link -->
		<div class="text-center">
			<p class="text-sm text-muted-foreground">
				Already have an account?
				<a href="/login" class="font-medium text-primary hover:text-primary/80">
					Sign in
				</a>
			</p>
		</div>
	</div>
</div>
