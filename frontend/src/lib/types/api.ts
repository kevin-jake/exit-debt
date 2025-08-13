export interface ApiResponse<T = any> {
	data: T;
	message?: string;
	success: boolean;
}

export interface ApiError {
	message: string;
	status: number;
	details?: any;
}

export interface PaginatedResponse<T> {
	data: T[];
	pagination: {
		page: number;
		limit: number;
		total: number;
		totalPages: number;
	};
}

export interface PaginationParams {
	page?: number;
	limit?: number;
}
