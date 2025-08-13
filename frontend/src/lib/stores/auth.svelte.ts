import { browser } from '$app/environment';
import { goto } from '$app/navigation';
import type {
	User,
	AuthState,
	LoginRequest,
	RegisterRequest,
	AuthResponse
} from '$lib/types/index.js';
import { apiClient } from '$lib/utils/api-client.js';
import { API_CONFIG } from '$lib/config/api.js';

class AuthStore {
	private state = $state<AuthState>({
		user: null,
		token: null,
		isAuthenticated: false,
		isLoading: false
	});

	constructor() {
		if (browser) {
			this.initializeAuth();
		}
	}

	get user(): User | null {
		return this.state.user;
	}

	get token(): string | null {
		return this.state.token;
	}

	get isAuthenticated(): boolean {
		return this.state.isAuthenticated;
	}

	get isLoading(): boolean {
		return this.state.isLoading;
	}

	async login(credentials: LoginRequest): Promise<{ success: boolean; error?: string }> {
		this.state.isLoading = true;

		try {
			const response = await apiClient.post<AuthResponse>(
				API_CONFIG.ENDPOINTS.AUTH.LOGIN,
				credentials
			);

			if (response.success && response.data) {
				this.setAuthData(response.data);
				goto('/dashboard');
				return { success: true };
			}

			return { success: false, error: response.message || 'Login failed' };
		} catch (error: any) {
			return {
				success: false,
				error: error.message || 'An error occurred during login'
			};
		} finally {
			this.state.isLoading = false;
		}
	}

	async register(userData: RegisterRequest): Promise<{ success: boolean; error?: string }> {
		this.state.isLoading = true;

		try {
			const response = await apiClient.post<AuthResponse>(
				API_CONFIG.ENDPOINTS.AUTH.REGISTER,
				userData
			);

			if (response.success && response.data) {
				this.setAuthData(response.data);
				goto('/dashboard');
				return { success: true };
			}

			return { success: false, error: response.message || 'Registration failed' };
		} catch (error: any) {
			return {
				success: false,
				error: error.message || 'An error occurred during registration'
			};
		} finally {
			this.state.isLoading = false;
		}
	}

	async logout(): Promise<void> {
		try {
			if (this.state.token) {
				await apiClient.post(API_CONFIG.ENDPOINTS.AUTH.LOGOUT, {}, this.state.token);
			}
		} catch (error) {
			console.error('Error during logout:', error);
		} finally {
			this.clearAuthData();
			goto('/auth/login');
		}
	}

	async refreshProfile(): Promise<void> {
		if (!this.state.token) return;

		try {
			const response = await apiClient.get<User>(
				API_CONFIG.ENDPOINTS.AUTH.PROFILE,
				this.state.token
			);

			if (response.success && response.data) {
				this.state.user = response.data;
			}
		} catch (error) {
			console.error('Error refreshing profile:', error);
			// Token might be invalid, logout user
			this.logout();
		}
	}

	private setAuthData(authData: AuthResponse): void {
		this.state.user = authData.user;
		this.state.token = authData.token;
		this.state.isAuthenticated = true;

		// Store in localStorage
		if (browser) {
			localStorage.setItem('exit-debt-token', authData.token);
			localStorage.setItem('exit-debt-user', JSON.stringify(authData.user));
		}
	}

	private clearAuthData(): void {
		this.state.user = null;
		this.state.token = null;
		this.state.isAuthenticated = false;

		// Clear from localStorage
		if (browser) {
			localStorage.removeItem('exit-debt-token');
			localStorage.removeItem('exit-debt-user');
		}
	}

	private initializeAuth(): void {
		const token = localStorage.getItem('exit-debt-token');
		const userStr = localStorage.getItem('exit-debt-user');

		if (token && userStr) {
			try {
				const user = JSON.parse(userStr);
				this.state.user = user;
				this.state.token = token;
				this.state.isAuthenticated = true;

				// Optionally refresh user profile
				this.refreshProfile();
			} catch (error) {
				console.error('Error parsing stored user data:', error);
				this.clearAuthData();
			}
		}
	}
}

export const authStore = new AuthStore();
