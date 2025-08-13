export interface User {
	id: string;
	email: string;
	firstName: string;
	lastName: string;
	createdAt: string;
	updatedAt: string;
}

export interface LoginRequest {
	email: string;
	password: string;
}

export interface RegisterRequest {
	email: string;
	password: string;
	firstName: string;
	lastName: string;
}

export interface AuthResponse {
	token: string;
	user: User;
}

export interface AuthState {
	user: User | null;
	token: string | null;
	isAuthenticated: boolean;
	isLoading: boolean;
}
