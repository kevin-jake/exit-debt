// API client for the DebtTracker backend
const API_BASE_URL =
  import.meta.env.VITE_API_URL || "http://localhost:8080/api/v1";

// Types for API requests and responses
export interface LoginRequest {
  email: string;
  password: string;
}

export interface LoginResponse {
  token: string;
  user: {
    id: string;
    email: string;
    first_name: string;
    last_name: string;
    phone?: string;
    created_at: string;
    updated_at: string;
  };
}

export interface RegisterRequest {
  email: string;
  password: string;
  first_name: string;
  last_name: string;
  phone?: string;
}

export interface RegisterResponse {
  user: {
    id: string;
    email: string;
    first_name: string;
    last_name: string;
    phone?: string;
    created_at: string;
    updated_at: string;
  };
}

export interface ApiResponse<T> {
  message: string;
  data?: T;
  request_id: string;
  timestamp: string;
}

export interface ApiError {
  error: string;
  details?: string;
  request_id: string;
  timestamp: string;
}

// Contact Management Types
export interface Contact {
  id: string;
  name: string;
  email?: string;
  phone?: string;
  notes?: string;
  created_at: string;
  updated_at: string;
}

export interface CreateContactRequest {
  name: string;
  email?: string;
  phone?: string;
  notes?: string;
}

export interface UpdateContactRequest {
  name?: string;
  email?: string;
  phone?: string;
  notes?: string;
}

// Debt Management Types
export interface DebtList {
  id: string;
  contact_id: string;
  total_amount: string;
  currency: string;
  debt_type: string;
  installment_plan: string;
  description?: string;
  notes?: string;
  created_at: string;
  updated_at: string;
  // Additional fields from backend response
  installment_amount?: string;
  total_payments_made?: string;
  total_remaining_debt?: string;
  status?: string;
  due_date?: string;
  next_payment_date?: string;
  number_of_payments?: number;
  contact?: Contact;
}

export interface CreateDebtListRequest {
  contact_id: string;
  total_amount: string;
  currency: string;
  debt_type: string;
  installment_plan: string;
  description?: string;
  notes?: string;
}

export interface UpdateDebtListRequest {
  contact_id?: string;
  total_amount?: string;
  currency?: string;
  debt_type?: string;
  installment_plan?: string;
  description?: string;
  notes?: string;
}

// Payment Management Types
export interface Payment {
  id: string;
  debt_list_id: string;
  amount: string;
  currency: string;
  payment_date: string;
  payment_method: string;
  description?: string;
  status:
    | "pending"
    | "paid"
    | "overdue"
    | "verified"
    | "rejected"
    | "completed";
  receipt_photo_url?: string;
  verified_by?: string;
  verified_at?: string;
  verification_notes?: string;
  created_at: string;
  updated_at: string;
}

export interface CreatePaymentRequest {
  amount: string;
  payment_date: string;
  payment_method?: string;
  description?: string;
}

export interface PaymentSchedule {
  debt_list_id: string;
  payments: Payment[];
  total_amount: string;
  remaining_amount: string;
  next_payment_date?: string;
}

// API client class
class ApiClient {
  private baseUrl: string;

  constructor(baseUrl: string) {
    this.baseUrl = baseUrl;
  }

  private async request<T>(
    endpoint: string,
    options: RequestInit = {}
  ): Promise<T> {
    const url = `${this.baseUrl}${endpoint}`;

    const config: RequestInit = {
      headers: {
        "Content-Type": "application/json",
        ...options.headers,
      },
      ...options,
    };

    // Add authorization header if token exists
    const token = localStorage.getItem("token");
    if (token) {
      config.headers = {
        ...config.headers,
        Authorization: `Bearer ${token}`,
      };
    }

    try {
      const response = await fetch(url, config);

      if (!response.ok) {
        const errorData: ApiError = await response.json();
        throw new Error(errorData.error || `HTTP ${response.status}`);
      }

      const data: ApiResponse<T> = await response.json();
      return data.data as T;
    } catch (error) {
      if (error instanceof Error) {
        throw error;
      }
      throw new Error("An unexpected error occurred");
    }
  }

  // Authentication methods
  async login(credentials: LoginRequest): Promise<LoginResponse> {
    const response = await this.request<{
      token: string;
      user: {
        ID: string;
        Email: string;
        FirstName: string;
        LastName: string;
        Phone?: string;
        CreatedAt: string;
        UpdatedAt: string;
      };
    }>("/auth/login", {
      method: "POST",
      body: JSON.stringify(credentials),
    });

    // Transform the response to match our frontend interface
    return {
      token: response.token,
      user: {
        id: response.user.ID,
        email: response.user.Email,
        first_name: response.user.FirstName,
        last_name: response.user.LastName,
        phone: response.user.Phone,
        created_at: response.user.CreatedAt,
        updated_at: response.user.UpdatedAt,
      },
    };
  }

  async register(userData: RegisterRequest): Promise<RegisterResponse> {
    const response = await this.request<{
      user: {
        ID: string;
        Email: string;
        FirstName: string;
        LastName: string;
        Phone?: string;
        CreatedAt: string;
        UpdatedAt: string;
      };
    }>("/auth/register", {
      method: "POST",
      body: JSON.stringify(userData),
    });

    // Transform the response to match our frontend interface
    return {
      user: {
        id: response.user.ID,
        email: response.user.Email,
        first_name: response.user.FirstName,
        last_name: response.user.LastName,
        phone: response.user.Phone,
        created_at: response.user.CreatedAt,
        updated_at: response.user.UpdatedAt,
      },
    };
  }

  // Health check
  async healthCheck(): Promise<{
    status: string;
    service: string;
    version: string;
    timestamp: string;
  }> {
    return this.request("/health");
  }

  // Contact Management methods
  async createContact(contactData: CreateContactRequest): Promise<Contact> {
    const response = await this.request<{
      ID: string;
      Name: string;
      Email?: string;
      Phone?: string;
      Notes?: string;
      IsUser: boolean;
      UserIDRef?: string;
      CreatedAt: string;
      UpdatedAt: string;
    }>("/contacts", {
      method: "POST",
      body: JSON.stringify(contactData),
    });

    // Transform the response to match our frontend interface
    return {
      id: response.ID,
      name: response.Name,
      email: response.Email,
      phone: response.Phone,
      notes: response.Notes,
      created_at: response.CreatedAt,
      updated_at: response.UpdatedAt,
    };
  }

  async getContacts(): Promise<Contact[]> {
    const response = await this.request<
      {
        ID: string;
        Name: string;
        Email?: string;
        Phone?: string;
        Notes?: string;
        IsUser: boolean;
        UserIDRef?: string;
        CreatedAt: string;
        UpdatedAt: string;
      }[]
    >("/contacts");

    // Transform the response to match our frontend interface
    return response.map((contact) => ({
      id: contact.ID,
      name: contact.Name,
      email: contact.Email,
      phone: contact.Phone,
      notes: contact.Notes,
      created_at: contact.CreatedAt,
      updated_at: contact.UpdatedAt,
    }));
  }

  async getContact(id: string): Promise<Contact> {
    const response = await this.request<{
      ID: string;
      Name: string;
      Email?: string;
      Phone?: string;
      Notes?: string;
      IsUser: boolean;
      UserIDRef?: string;
      CreatedAt: string;
      UpdatedAt: string;
    }>(`/contacts/${id}`);

    // Transform the response to match our frontend interface
    return {
      id: response.ID,
      name: response.Name,
      email: response.Email,
      phone: response.Phone,
      notes: response.Notes,
      created_at: response.CreatedAt,
      updated_at: response.UpdatedAt,
    };
  }

  async updateContact(
    id: string,
    contactData: UpdateContactRequest
  ): Promise<Contact> {
    const response = await this.request<{
      ID: string;
      Name: string;
      Email?: string;
      Phone?: string;
      Notes?: string;
      IsUser: boolean;
      UserIDRef?: string;
      CreatedAt: string;
      UpdatedAt: string;
    }>(`/contacts/${id}`, {
      method: "PUT",
      body: JSON.stringify(contactData),
    });

    // Transform the response to match our frontend interface
    return {
      id: response.ID,
      name: response.Name,
      email: response.Email,
      phone: response.Phone,
      notes: response.Notes,
      created_at: response.CreatedAt,
      updated_at: response.UpdatedAt,
    };
  }

  async deleteContact(id: string): Promise<void> {
    return this.request<void>(`/contacts/${id}`, {
      method: "DELETE",
    });
  }

  async getDebtLists(): Promise<DebtList[]> {
    const response = await this.request<
      {
        id: string;
        user_id: string;
        contact_id: string;
        debt_type: string;
        total_amount: string;
        installment_amount: string;
        total_payments_made: string;
        total_remaining_debt: string;
        currency: string;
        status: string;
        due_date: string;
        next_payment_date: string;
        installment_plan: string;
        number_of_payments: number;
        description?: string;
        notes?: string;
        created_at: string;
        updated_at: string;
        contact: {
          ID: string;
          Name: string;
          Email: string;
          Phone: string;
          Notes: string;
          IsUser: boolean;
          UserIDRef: string | null;
          CreatedAt: string;
          UpdatedAt: string;
        };
      }[]
    >("/debts");

    return response.map((debt) => ({
      id: debt.id,
      contact_id: debt.contact_id,
      total_amount: debt.total_amount,
      currency: debt.currency,
      debt_type: debt.debt_type,
      installment_plan: debt.installment_plan,
      description: debt.description,
      notes: debt.notes,
      created_at: debt.created_at,
      updated_at: debt.updated_at,
      // Additional fields from backend response
      installment_amount: debt.installment_amount,
      total_payments_made: debt.total_payments_made,
      total_remaining_debt: debt.total_remaining_debt,
      status: debt.status,
      due_date: debt.due_date,
      next_payment_date: debt.next_payment_date,
      number_of_payments: debt.number_of_payments,
      contact: {
        id: debt.contact.ID,
        name: debt.contact.Name,
        email: debt.contact.Email,
        phone: debt.contact.Phone,
        notes: debt.contact.Notes,
        created_at: debt.contact.CreatedAt,
        updated_at: debt.contact.UpdatedAt,
      },
    }));
  }

  async getDebtList(id: string): Promise<DebtList> {
    const response = await this.request<{
      id: string;
      contact_id: string;
      total_amount: string;
      currency: string;
      debt_type: string;
      installment_plan: string;
      description?: string;
      notes?: string;
      created_at: string;
      updated_at: string;
    }>(`/debts/${id}`);

    return {
      id: response.id,
      contact_id: response.contact_id,
      total_amount: response.total_amount,
      currency: response.currency,
      debt_type: response.debt_type,
      installment_plan: response.installment_plan,
      description: response.description,
      notes: response.notes,
      created_at: response.created_at,
      updated_at: response.updated_at,
    };
  }

  async createDebtList(debtData: CreateDebtListRequest): Promise<DebtList> {
    const response = await this.request<{
      id: string;
      contact_id: string;
      total_amount: string;
      currency: string;
      debt_type: string;
      installment_plan: string;
      description?: string;
      notes?: string;
      created_at: string;
      updated_at: string;
    }>("/debts", {
      method: "POST",
      body: JSON.stringify(debtData),
    });

    return {
      id: response.id,
      contact_id: response.contact_id,
      total_amount: response.total_amount,
      currency: response.currency,
      debt_type: response.debt_type,
      installment_plan: response.installment_plan,
      description: response.description,
      notes: response.notes,
      created_at: response.created_at,
      updated_at: response.updated_at,
    };
  }

  async updateDebtList(
    id: string,
    debtData: UpdateDebtListRequest
  ): Promise<DebtList> {
    const response = await this.request<{
      id: string;
      contact_id: string;
      total_amount: string;
      currency: string;
      debt_type: string;
      installment_plan: string;
      description?: string;
      notes?: string;
      created_at: string;
      updated_at: string;
    }>(`/debts/${id}`, {
      method: "PUT",
      body: JSON.stringify(debtData),
    });

    return {
      id: response.id,
      contact_id: response.contact_id,
      total_amount: response.total_amount,
      currency: response.currency,
      debt_type: response.debt_type,
      installment_plan: response.installment_plan,
      description: response.description,
      notes: response.notes,
      created_at: response.created_at,
      updated_at: response.updated_at,
    };
  }

  async deleteDebtList(id: string): Promise<void> {
    return this.request<void>(`/debts/${id}`, {
      method: "DELETE",
    });
  }

  // Payment Management methods
  async createPayment(
    debtId: string,
    paymentData: CreatePaymentRequest
  ): Promise<Payment> {
    const response = await this.request<{
      ID: string;
      DebtListID: string;
      Amount: string;
      Currency: string;
      PaymentDate: string;
      PaymentMethod: string;
      Description?: string;
      Status:
        | "pending"
        | "paid"
        | "overdue"
        | "verified"
        | "rejected"
        | "completed";
      ReceiptPhotoURL?: string;
      VerifiedBy?: string;
      VerifiedAt?: string;
      VerificationNotes?: string;
      CreatedAt: string;
      UpdatedAt: string;
    }>("/debts/payments", {
      method: "POST",
      body: JSON.stringify({
        debt_list_id: debtId,
        amount: paymentData.amount,
        payment_date: paymentData.payment_date,
        payment_method: paymentData.payment_method || "cash",
        description: paymentData.description,
      }),
    });

    return {
      id: response.ID,
      debt_list_id: response.DebtListID,
      amount: response.Amount,
      currency: response.Currency,
      payment_date: response.PaymentDate,
      payment_method: response.PaymentMethod,
      description: response.Description,
      status: response.Status,
      receipt_photo_url: response.ReceiptPhotoURL,
      verified_by: response.VerifiedBy,
      verified_at: response.VerifiedAt,
      verification_notes: response.VerificationNotes,
      created_at: response.CreatedAt,
      updated_at: response.UpdatedAt,
    };
  }

  async getPayments(debtId: string): Promise<Payment[]> {
    const response = await this.request<
      {
        ID: string;
        DebtListID: string;
        Amount: string;
        Currency: string;
        PaymentDate: string;
        PaymentMethod: string;
        Description?: string;
        Status:
          | "pending"
          | "paid"
          | "overdue"
          | "verified"
          | "rejected"
          | "completed";
        ReceiptPhotoURL?: string;
        VerifiedBy?: string;
        VerifiedAt?: string;
        VerificationNotes?: string;
        CreatedAt: string;
        UpdatedAt: string;
      }[]
    >(`/debts/${debtId}/payments`);

    return response.map((payment) => ({
      id: payment.ID,
      debt_list_id: payment.DebtListID,
      amount: payment.Amount,
      currency: payment.Currency,
      payment_date: payment.PaymentDate,
      payment_method: payment.PaymentMethod,
      description: payment.Description,
      status: payment.Status,
      receipt_photo_url: payment.ReceiptPhotoURL,
      verified_by: payment.VerifiedBy,
      verified_at: payment.VerifiedAt,
      verification_notes: payment.VerificationNotes,
      created_at: payment.CreatedAt,
      updated_at: payment.UpdatedAt,
    }));
  }

  async verifyPayment(paymentId: string): Promise<void> {
    return this.request<void>(`/debts/payments/${paymentId}/verify`, {
      method: "POST",
    });
  }

  async rejectPayment(paymentId: string): Promise<void> {
    return this.request<void>(`/debts/payments/${paymentId}/reject`, {
      method: "POST",
    });
  }

  async uploadReceipt(paymentId: string, file: File): Promise<void> {
    const formData = new FormData();
    formData.append("receipt", file);

    const url = `${this.baseUrl}/debts/payments/${paymentId}/receipt`;
    const token = localStorage.getItem("token");

    const config: RequestInit = {
      method: "POST",
      headers: {
        ...(token ? { Authorization: `Bearer ${token}` } : {}),
      },
      body: formData,
    };

    const response = await fetch(url, config);

    if (!response.ok) {
      const errorData: ApiError = await response.json();
      throw new Error(errorData.error || `HTTP ${response.status}`);
    }
  }

  // Analytics methods
  async getOverdueItems(): Promise<Payment[]> {
    const response = await this.request<
      {
        ID: string;
        DebtListID: string;
        Amount: string;
        Currency: string;
        PaymentDate: string;
        PaymentMethod: string;
        Description?: string;
        Status:
          | "pending"
          | "paid"
          | "overdue"
          | "verified"
          | "rejected"
          | "completed";
        ReceiptPhotoURL?: string;
        VerifiedBy?: string;
        VerifiedAt?: string;
        VerificationNotes?: string;
        CreatedAt: string;
        UpdatedAt: string;
      }[]
    >("/debts/overdue");

    return response.map((payment) => ({
      id: payment.ID,
      debt_list_id: payment.DebtListID,
      amount: payment.Amount,
      currency: payment.Currency,
      payment_date: payment.PaymentDate,
      payment_method: payment.PaymentMethod,
      description: payment.Description,
      status: payment.Status,
      receipt_photo_url: payment.ReceiptPhotoURL,
      verified_by: payment.VerifiedBy,
      verified_at: payment.VerifiedAt,
      verification_notes: payment.VerificationNotes,
      created_at: payment.CreatedAt,
      updated_at: payment.UpdatedAt,
    }));
  }

  async getDueSoonItems(): Promise<Payment[]> {
    const response = await this.request<
      {
        ID: string;
        DebtListID: string;
        Amount: string;
        Currency: string;
        PaymentDate: string;
        PaymentMethod: string;
        Description?: string;
        Status:
          | "pending"
          | "paid"
          | "overdue"
          | "verified"
          | "rejected"
          | "completed";
        ReceiptPhotoURL?: string;
        VerifiedBy?: string;
        VerifiedAt?: string;
        VerificationNotes?: string;
        CreatedAt: string;
        UpdatedAt: string;
      }[]
    >("/debts/due-soon");

    return response.map((payment) => ({
      id: payment.ID,
      debt_list_id: payment.DebtListID,
      amount: payment.Amount,
      currency: payment.Currency,
      payment_date: payment.PaymentDate,
      payment_method: payment.PaymentMethod,
      description: payment.Description,
      status: payment.Status,
      receipt_photo_url: payment.ReceiptPhotoURL,
      verified_by: payment.VerifiedBy,
      verified_at: payment.VerifiedAt,
      verification_notes: payment.VerificationNotes,
      created_at: payment.CreatedAt,
      updated_at: payment.UpdatedAt,
    }));
  }

  async getPaymentSchedule(debtId: string): Promise<PaymentSchedule> {
    const response = await this.request<{
      debt_list_id: string;
      payments: {
        ID: string;
        DebtListID: string;
        Amount: string;
        Currency: string;
        PaymentDate: string;
        PaymentMethod: string;
        Description?: string;
        Status:
          | "pending"
          | "paid"
          | "overdue"
          | "verified"
          | "rejected"
          | "completed";
        ReceiptPhotoURL?: string;
        VerifiedBy?: string;
        VerifiedAt?: string;
        VerificationNotes?: string;
        CreatedAt: string;
        UpdatedAt: string;
      }[];
      total_amount: string;
      remaining_amount: string;
      next_payment_date?: string;
    }>(`/debts/${debtId}/schedule`);

    return {
      debt_list_id: response.debt_list_id,
      payments: response.payments.map((payment) => ({
        id: payment.ID,
        debt_list_id: payment.DebtListID,
        amount: payment.Amount,
        currency: payment.Currency,
        payment_date: payment.PaymentDate,
        payment_method: payment.PaymentMethod,
        description: payment.Description,
        status: payment.Status,
        receipt_photo_url: payment.ReceiptPhotoURL,
        verified_by: payment.VerifiedBy,
        verified_at: payment.VerifiedAt,
        verification_notes: payment.VerificationNotes,
        created_at: payment.CreatedAt,
        updated_at: payment.UpdatedAt,
      })),
      total_amount: response.total_amount,
      remaining_amount: response.remaining_amount,
      next_payment_date: response.next_payment_date,
    };
  }

  async getUpcomingPayments(): Promise<Payment[]> {
    const response = await this.request<
      {
        ID: string;
        DebtListID: string;
        Amount: string;
        Currency: string;
        PaymentDate: string;
        PaymentMethod: string;
        Description?: string;
        Status:
          | "pending"
          | "paid"
          | "overdue"
          | "verified"
          | "rejected"
          | "completed";
        ReceiptPhotoURL?: string;
        VerifiedBy?: string;
        VerifiedAt?: string;
        VerificationNotes?: string;
        CreatedAt: string;
        UpdatedAt: string;
      }[]
    >("/upcoming-payments");

    return response.map((payment) => ({
      id: payment.ID,
      debt_list_id: payment.DebtListID,
      amount: payment.Amount,
      currency: payment.Currency,
      payment_date: payment.PaymentDate,
      payment_method: payment.PaymentMethod,
      description: payment.Description,
      status: payment.Status,
      receipt_photo_url: payment.ReceiptPhotoURL,
      verified_by: payment.VerifiedBy,
      verified_at: payment.VerifiedAt,
      verification_notes: payment.VerificationNotes,
      created_at: payment.CreatedAt,
      updated_at: payment.UpdatedAt,
    }));
  }
}

// Create and export the API client instance
export const apiClient = new ApiClient(API_BASE_URL);

// Utility functions for token management
export const tokenManager = {
  getToken(): string | null {
    return localStorage.getItem("token");
  },

  setToken(token: string): void {
    localStorage.setItem("token", token);
  },

  removeToken(): void {
    localStorage.removeItem("token");
  },

  isAuthenticated(): boolean {
    return !!this.getToken();
  },
};
