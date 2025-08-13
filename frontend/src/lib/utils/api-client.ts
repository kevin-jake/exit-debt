import { API_CONFIG } from '$lib/config/api.js';
import type { ApiResponse, ApiError } from '$lib/types/index.js';

interface RequestOptions extends RequestInit {
	token?: string;
}

class ApiClient {
	private baseURL: string;

	constructor(baseURL: string = API_CONFIG.BASE_URL) {
		this.baseURL = baseURL;
	}

	private async request<T>(
		endpoint: string,
		options: RequestOptions = {}
	): Promise<ApiResponse<T>> {
		const { token, ...requestOptions } = options;

		const config: RequestInit = {
			headers: {
				'Content-Type': 'application/json',
				...(token && { Authorization: `Bearer ${token}` }),
				...requestOptions.headers
			},
			...requestOptions
		};

		try {
			const response = await fetch(`${this.baseURL}${endpoint}`, config);
			const data = await response.json();

			if (!response.ok) {
				const error: ApiError = {
					message: data.message || 'An error occurred',
					status: response.status,
					details: data
				};
				throw error;
			}

			return data;
		} catch (error) {
			if (error instanceof Error && 'status' in error) {
				throw error;
			}
			throw {
				message: 'Network error occurred',
				status: 0,
				details: error
			} as ApiError;
		}
	}

	async get<T>(endpoint: string, token?: string): Promise<ApiResponse<T>> {
		return this.request<T>(endpoint, { method: 'GET', token });
	}

	async post<T>(endpoint: string, data?: any, token?: string): Promise<ApiResponse<T>> {
		return this.request<T>(endpoint, {
			method: 'POST',
			body: data ? JSON.stringify(data) : undefined,
			token
		});
	}

	async put<T>(endpoint: string, data?: any, token?: string): Promise<ApiResponse<T>> {
		return this.request<T>(endpoint, {
			method: 'PUT',
			body: data ? JSON.stringify(data) : undefined,
			token
		});
	}

	async delete<T>(endpoint: string, token?: string): Promise<ApiResponse<T>> {
		return this.request<T>(endpoint, { method: 'DELETE', token });
	}
}

export const apiClient = new ApiClient();
