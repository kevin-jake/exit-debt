<script lang="ts">
	import { goto } from '$app/navigation';
	import { apiClient, tokenManager } from '$lib/api';
	import { authStore } from '$lib/stores/auth';

	let email = $state('');
	let password = $state('');
	let isLoading = $state(false);
	let error = $state('');

	async function handleLogin() {
		if (!email || !password) {
			error = 'Please fill in all fields';
			return;
		}

		isLoading = true;
		error = '';

		try {
			const response = await apiClient.login({ email, password });
			
			// Store the JWT token
			tokenManager.setToken(response.token);
			
			// Update auth store with user data
			authStore.setUser(response.user);
			
			// Redirect to dashboard
			goto('/');
		} catch (err) {
			if (err instanceof Error) {
				error = err.message;
			} else {
				error = 'An error occurred. Please try again.';
			}
		} finally {
			isLoading = false;
		}
	}

	function handleKeyPress(event: KeyboardEvent) {
		if (event.key === 'Enter') {
			handleLogin();
		}
	}
</script>

<svelte:head>
	<title>Login - DebtTracker</title>
</svelte:head>

<div class="min-h-screen flex items-center justify-center bg-gradient-to-br from-primary/10 to-muted px-4">
	<div class="max-w-md w-full space-y-8">
		<!-- Logo and Header -->
		<div class="text-center">
			<div class="mx-auto w-16 h-16 bg-primary rounded-xl flex items-center justify-center mb-6">
				<svg class="w-8 h-8 text-primary-foreground" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1"></path>
				</svg>
			</div>
			<h2 class="text-3xl font-bold text-foreground">Welcome back</h2>
			<p class="mt-2 text-sm text-muted-foreground">Sign in to your DebtTracker account</p>
		</div>

		<!-- Login Form -->
		<div class="card p-8">
			<form onsubmit={handleLogin} class="space-y-6">
				{#if error}
					<div class="bg-destructive/10 border border-destructive/20 text-destructive px-4 py-3 rounded-lg text-sm">
						{error}
					</div>
				{/if}

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
					<label for="password" class="label">Password</label>
					<input
						id="password"
						type="password"
						bind:value={password}
						onkeypress={handleKeyPress}
						class="input"
						placeholder="Enter your password"
						required
					/>
				</div>

				<div class="flex items-center justify-between">
					<div class="flex items-center">
						<input
							id="remember-me"
							type="checkbox"
							class="h-4 w-4 text-primary focus:ring-primary border-input rounded"
						/>
						<label for="remember-me" class="ml-2 block text-sm text-muted-foreground">
							Remember me
						</label>
					</div>

					<a href="/forgot-password" class="text-sm text-primary hover:text-primary/80">
						Forgot password?
					</a>
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
						Signing in...
					{:else}
						Sign in
					{/if}
				</button>
			</form>
		</div>

		<!-- Register Link -->
		<div class="text-center">
			<p class="text-sm text-muted-foreground">
				Don't have an account?
				<a href="/register" class="font-medium text-primary hover:text-primary/80">
					Sign up
				</a>
			</p>
		</div>
	</div>
</div>
