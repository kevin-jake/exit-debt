export type DebtStatus = 'pending' | 'paid' | 'overdue' | 'cancelled';

export type DebtType = 'owed_by_me' | 'owed_to_me';

export interface Debt {
	id: string;
	title: string;
	description?: string;
	amount: number;
	currency: string;
	type: DebtType;
	status: DebtStatus;
	dueDate?: string;
	categoryId?: string;
	creditorId?: string;
	debtorId?: string;
	createdAt: string;
	updatedAt: string;
}

export interface DebtCategory {
	id: string;
	name: string;
	description?: string;
	color: string;
	userId: string;
	createdAt: string;
	updatedAt: string;
}

export interface DebtCreateRequest {
	title: string;
	description?: string;
	amount: number;
	currency: string;
	type: DebtType;
	dueDate?: string;
	categoryId?: string;
	creditorId?: string;
	debtorId?: string;
}

export interface DebtUpdateRequest extends Partial<DebtCreateRequest> {
	status?: DebtStatus;
}

export interface DebtFilters {
	status?: DebtStatus[];
	type?: DebtType[];
	categoryId?: string[];
	search?: string;
	sortBy?: 'amount' | 'dueDate' | 'createdAt' | 'title';
	sortOrder?: 'asc' | 'desc';
}
