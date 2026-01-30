/**
 * API Client for backend communication
 * Handles authentication tokens and standard API responses
 */

export interface ApiResponse<T = unknown> {
	success: boolean;
	data?: T;
	error?: {
		code: string;
		message: string;
		details?: Array<{ field: string; message: string }>;
	};
	meta?: {
		timestamp: string;
		requestId?: string;
	};
}

export interface User {
	id: string;
	email: string;
	name: string | null;
	role: 'user' | 'admin';
	isActive: boolean;
	lastLoginAt: string | null;
	createdAt: string;
	updatedAt: string;
}

export interface AuthTokens {
	accessToken: string;
	expiresIn: number;
}

export interface LoginCredentials {
	email: string;
	password: string;
}

export interface RegisterData {
	email: string;
	password: string;
	name?: string;
}

class ApiClient {
	private baseUrl: string;
	private accessToken: string | null = null;
	private isRefreshing: boolean = false;
	private refreshPromise: Promise<boolean> | null = null;

	constructor(baseUrl: string = '/api') {
		this.baseUrl = baseUrl;
	}

	setAccessToken(token: string | null) {
		this.accessToken = token;
	}

	getAccessToken(): string | null {
		return this.accessToken;
	}

	private async request<T>(
		endpoint: string,
		options: RequestInit = {},
		skipRefresh: boolean = false
	): Promise<ApiResponse<T>> {
		const url = `${this.baseUrl}${endpoint}`;

		const headers: HeadersInit = {
			'Content-Type': 'application/json',
			...options.headers
		};

		if (this.accessToken) {
			(headers as Record<string, string>)['Authorization'] = `Bearer ${this.accessToken}`;
		}

		try {
			const response = await fetch(url, {
				...options,
				headers,
				credentials: 'include'
			});

			const data: ApiResponse<T> = await response.json();

			if (response.status === 401 && endpoint !== '/auth/refresh' && !skipRefresh) {
				const refreshed = await this.refreshToken();
				if (refreshed) {
					return this.request<T>(endpoint, options, true);
				}
			}

			return data;
		} catch (error) {
			return {
				success: false,
				error: {
					code: 'NETWORK_ERROR',
					message: error instanceof Error ? error.message : 'Network error occurred'
				}
			};
		}
	}

	async get<T>(endpoint: string): Promise<ApiResponse<T>> {
		return this.request<T>(endpoint, { method: 'GET' });
	}

	async post<T>(endpoint: string, body?: unknown): Promise<ApiResponse<T>> {
		return this.request<T>(endpoint, {
			method: 'POST',
			body: body ? JSON.stringify(body) : undefined
		});
	}

	async put<T>(endpoint: string, body?: unknown): Promise<ApiResponse<T>> {
		return this.request<T>(endpoint, {
			method: 'PUT',
			body: body ? JSON.stringify(body) : undefined
		});
	}

	async delete<T>(endpoint: string): Promise<ApiResponse<T>> {
		return this.request<T>(endpoint, { method: 'DELETE' });
	}

	async login(credentials: LoginCredentials): Promise<ApiResponse<AuthTokens & { user: User }>> {
		const response = await this.post<AuthTokens & { user: User }>('/auth/login', credentials);
		if (response.success && response.data) {
			this.setAccessToken(response.data.accessToken);
		}
		return response;
	}

	async register(data: RegisterData): Promise<ApiResponse<AuthTokens & { user: User }>> {
		const response = await this.post<AuthTokens & { user: User }>('/auth/register', data);
		if (response.success && response.data) {
			this.setAccessToken(response.data.accessToken);
		}
		return response;
	}

	async logout(): Promise<ApiResponse<void>> {
		const response = await this.post<void>('/auth/logout');
		this.setAccessToken(null);
		return response;
	}

	async refreshToken(): Promise<boolean> {
		if (this.isRefreshing && this.refreshPromise) {
			return this.refreshPromise;
		}

		this.isRefreshing = true;
		this.refreshPromise = this.doRefresh();

		try {
			return await this.refreshPromise;
		} finally {
			this.isRefreshing = false;
			this.refreshPromise = null;
		}
	}

	private async doRefresh(): Promise<boolean> {
		const response = await this.post<AuthTokens>('/auth/refresh');
		if (response.success && response.data) {
			this.setAccessToken(response.data.accessToken);
			return true;
		}
		this.setAccessToken(null);
		return false;
	}

	async getMe(): Promise<ApiResponse<User>> {
		return this.get<User>('/auth/me');
	}
}

export const api = new ApiClient();
