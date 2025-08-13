<script lang="ts">
	import { Button } from '$lib/components/ui/button/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '$lib/components/ui/card/index.js';
	import { authStore } from '$lib/stores/auth.svelte.js';
	import { Eye, EyeOff, Loader2 } from 'lucide-svelte';

	let firstName = $state('');
	let lastName = $state('');
	let email = $state('');
	let password = $state('');
	let confirmPassword = $state('');
	let showPassword = $state(false);
	let showConfirmPassword = $state(false);
	let error = $state('');

	async function handleSubmit(event: Event) {
		event.preventDefault();
		error = '';

		// Validation
		if (!firstName || !lastName || !email || !password || !confirmPassword) {
			error = 'Please fill in all fields';
			return;
		}

		if (password !== confirmPassword) {
			error = 'Passwords do not match';
			return;
		}

		if (password.length < 6) {
			error = 'Password must be at least 6 characters long';
			return;
		}

		const result = await authStore.register({ 
			firstName, 
			lastName, 
			email, 
			password 
		});
		
		if (!result.success) {
			error = result.error || 'Registration failed';
		}
	}

	function togglePasswordVisibility(field: 'password' | 'confirm') {
		if (field === 'password') {
			showPassword = !showPassword;
		} else {
			showConfirmPassword = !showConfirmPassword;
		}
	}
</script>

<Card class="w-full max-w-md">
	<CardHeader class="space-y-1">
		<CardTitle class="text-2xl font-bold">Create an account</CardTitle>
		<CardDescription>
			Enter your information to get started with Exit-Debt
		</CardDescription>
	</CardHeader>
	<CardContent>
		<form on:submit={handleSubmit} class="space-y-4">
			{#if error}
				<div class="rounded-md bg-destructive/15 p-3 text-sm text-destructive">
					{error}
				</div>
			{/if}

			<div class="grid grid-cols-2 gap-4">
				<div class="space-y-2">
					<Label for="firstName">First Name</Label>
					<Input
						id="firstName"
						type="text"
						placeholder="John"
						bind:value={firstName}
						disabled={authStore.isLoading}
						required
					/>
				</div>
				<div class="space-y-2">
					<Label for="lastName">Last Name</Label>
					<Input
						id="lastName"
						type="text"
						placeholder="Doe"
						bind:value={lastName}
						disabled={authStore.isLoading}
						required
					/>
				</div>
			</div>

			<div class="space-y-2">
				<Label for="email">Email</Label>
				<Input
					id="email"
					type="email"
					placeholder="john.doe@example.com"
					bind:value={email}
					disabled={authStore.isLoading}
					required
				/>
			</div>

			<div class="space-y-2">
				<Label for="password">Password</Label>
				<div class="relative">
					<Input
						id="password"
						type={showPassword ? 'text' : 'password'}
						placeholder="Create a password"
						bind:value={password}
						disabled={authStore.isLoading}
						required
					/>
					<Button
						type="button"
						variant="ghost"
						size="sm"
						class="absolute right-0 top-0 h-full px-3 py-2 hover:bg-transparent"
						on:click={() => togglePasswordVisibility('password')}
						disabled={authStore.isLoading}
					>
						{#if showPassword}
							<EyeOff class="h-4 w-4" />
							<span class="sr-only">Hide password</span>
						{:else}
							<Eye class="h-4 w-4" />
							<span class="sr-only">Show password</span>
						{/if}
					</Button>
				</div>
			</div>

			<div class="space-y-2">
				<Label for="confirmPassword">Confirm Password</Label>
				<div class="relative">
					<Input
						id="confirmPassword"
						type={showConfirmPassword ? 'text' : 'password'}
						placeholder="Confirm your password"
						bind:value={confirmPassword}
						disabled={authStore.isLoading}
						required
					/>
					<Button
						type="button"
						variant="ghost"
						size="sm"
						class="absolute right-0 top-0 h-full px-3 py-2 hover:bg-transparent"
						on:click={() => togglePasswordVisibility('confirm')}
						disabled={authStore.isLoading}
					>
						{#if showConfirmPassword}
							<EyeOff class="h-4 w-4" />
							<span class="sr-only">Hide password</span>
						{:else}
							<Eye class="h-4 w-4" />
							<span class="sr-only">Show password</span>
						{/if}
					</Button>
				</div>
			</div>

			<Button 
				type="submit" 
				class="w-full" 
				disabled={authStore.isLoading}
			>
				{#if authStore.isLoading}
					<Loader2 class="mr-2 h-4 w-4 animate-spin" />
				{/if}
				Create Account
			</Button>
		</form>

		<div class="mt-6 text-center text-sm text-muted-foreground">
			Already have an account?{' '}
			<a href="/auth/login" class="text-primary hover:underline">
				Sign in
			</a>
		</div>
	</CardContent>
</Card>
