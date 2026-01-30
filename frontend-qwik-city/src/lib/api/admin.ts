/**
 * Admin API Client
 * Handles admin panel API requests
 */

import { api, type ApiResponse } from './client';

export interface AdminUser {
	id: string;
	email: string;
	name: string | null;
	role: 'user' | 'admin';
	isActive: boolean;
	lastLoginAt: string | null;
	createdAt: string;
	updatedAt: string;
}

export interface DashboardStats {
	totalUsers: number;
	activeUsers: number;
	adminUsers: number;
	newUsersToday: number;
	newUsersThisWeek: number;
	newUsersThisMonth: number;
	recentUsers: Array<{
		id: string;
		email: string;
		name: string | null;
		createdAt: string;
	}>;
	recentActivity: Array<{
		type: string;
		message: string;
		timestamp: string;
	}>;
}

export interface ListParams {
	page?: number;
	pageSize?: number;
	search?: string;
	sortBy?: string;
	sortDir?: 'asc' | 'desc';
	role?: string;
	isActive?: boolean;
}

export interface ListResult<T> {
	items: T[];
	total: number;
	page: number;
	pageSize: number;
	totalPages: number;
}

export interface CreateUserInput {
	email: string;
	password: string;
	name?: string;
	role?: 'user' | 'admin';
	isActive?: boolean;
}

export interface UpdateUserInput {
	email?: string;
	password?: string;
	name?: string;
	role?: 'user' | 'admin';
	isActive?: boolean;
}

export interface FileInfo {
	name: string;
	path: string;
	size: number;
	isDir: boolean;
	modTime: string;
	extension: string;
	mimeType: string;
}

export interface FilesResult {
	files: FileInfo[];
	total: number;
	totalSize: number;
	currentDir: string;
}

export interface AppSetting {
	id: string;
	key: string;
	value: string;
	type: 'string' | 'number' | 'boolean' | 'json';
	label: string;
	group: string;
	updatedAt: string;
}

export const adminApi = {
	getDashboard: (): Promise<ApiResponse<DashboardStats>> => {
		return api.get<DashboardStats>('/admin/dashboard');
	},

	getUsers: (params: ListParams = {}): Promise<ApiResponse<ListResult<AdminUser>>> => {
		const searchParams = new URLSearchParams();
		if (params.page) searchParams.set('page', String(params.page));
		if (params.pageSize) searchParams.set('pageSize', String(params.pageSize));
		if (params.search) searchParams.set('search', params.search);
		if (params.sortBy) searchParams.set('sortBy', params.sortBy);
		if (params.sortDir) searchParams.set('sortDir', params.sortDir);
		if (params.role) searchParams.set('role', params.role);
		if (params.isActive !== undefined) searchParams.set('isActive', String(params.isActive));

		const query = searchParams.toString();
		return api.get<ListResult<AdminUser>>(`/admin/users${query ? `?${query}` : ''}`);
	},

	getUser: (id: string): Promise<ApiResponse<AdminUser>> => {
		return api.get<AdminUser>(`/admin/users/${id}`);
	},

	createUser: (data: CreateUserInput): Promise<ApiResponse<AdminUser>> => {
		return api.post<AdminUser>('/admin/users', data);
	},

	updateUser: (id: string, data: UpdateUserInput): Promise<ApiResponse<AdminUser>> => {
		return api.put<AdminUser>(`/admin/users/${id}`, data);
	},

	deleteUser: (id: string): Promise<ApiResponse<{ message: string }>> => {
		return api.delete<{ message: string }>(`/admin/users/${id}`);
	},

	getFiles: (dir?: string): Promise<ApiResponse<FilesResult>> => {
		const query = dir ? `?dir=${encodeURIComponent(dir)}` : '';
		return api.get<FilesResult>(`/admin/files${query}`);
	},

	deleteFile: (path: string): Promise<ApiResponse<{ message: string }>> => {
		return api.delete<{ message: string }>(`/admin/files/${encodeURIComponent(path)}`);
	},

	getSettings: (group?: string): Promise<ApiResponse<AppSetting[]>> => {
		const query = group ? `?group=${encodeURIComponent(group)}` : '';
		return api.get<AppSetting[]>(`/admin/settings${query}`);
	},

	getSetting: (key: string): Promise<ApiResponse<AppSetting>> => {
		return api.get<AppSetting>(`/admin/settings/${key}`);
	},

	updateSetting: (key: string, value: string): Promise<ApiResponse<AppSetting>> => {
		return api.put<AppSetting>(`/admin/settings/${key}`, { value });
	},

	updateSettings: (
		settings: Array<{ key: string; value: string }>
	): Promise<ApiResponse<AppSetting[]>> => {
		return api.put<AppSetting[]>('/admin/settings', { settings });
	}
};
