import { env } from '$env/dynamic/public';

export const API_CONFIG = {
	BASE_URL: env.PUBLIC_API_BASE_URL || 'http://localhost:8080/api/v1',
	ENDPOINTS: {
		// Auth endpoints
		AUTH: {
			LOGIN: '/auth/login',
			REGISTER: '/auth/register',
			LOGOUT: '/auth/logout',
			REFRESH: '/auth/refresh',
			PROFILE: '/auth/profile'
		},
		// Debt endpoints
		DEBTS: {
			LIST: '/debts',
			CREATE: '/debts',
			UPDATE: (id: string) => `/debts/${id}`,
			DELETE: (id: string) => `/debts/${id}`,
			GET: (id: string) => `/debts/${id}`
		},
		// Category endpoints
		CATEGORIES: {
			LIST: '/categories',
			CREATE: '/categories',
			UPDATE: (id: string) => `/categories/${id}`,
			DELETE: (id: string) => `/categories/${id}`
		},
		// User endpoints
		USERS: {
			PROFILE: '/users/profile',
			SETTINGS: '/users/settings',
			UPDATE: '/users/profile'
		}
	}
} as const;

export const APP_CONFIG = {
	NAME: env.PUBLIC_APP_NAME || 'Exit-Debt',
	VERSION: env.PUBLIC_APP_VERSION || '1.0.0'
} as const;
