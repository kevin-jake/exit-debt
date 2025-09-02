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
