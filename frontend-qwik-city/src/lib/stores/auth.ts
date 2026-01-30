/**
 * Auth Store - Plain TypeScript module
 * Module-level singleton state for authentication
 */

import { api, type User } from '../api/client';

let _user: User | null = null;
let _isLoading = true;
let _isInitialized = false;
let _initPromise: Promise<void> | null = null;
let _listeners: Array<() => void> = [];

function notify() {
  _listeners.forEach((fn) => fn());
}

export function subscribe(fn: () => void): () => void {
  _listeners.push(fn);
  return () => {
    _listeners = _listeners.filter((l) => l !== fn);
  };
}

export async function initAuth(): Promise<void> {
  if (_isInitialized) return;
  if (_initPromise) return _initPromise;
  _initPromise = doInitAuth();
  await _initPromise;
}

async function doInitAuth(): Promise<void> {
  _isLoading = true;
  notify();
  try {
    const response = await api.getMe();
    if (response.success && response.data) {
      _user = response.data;
    } else {
      _user = null;
    }
  } catch {
    _user = null;
  } finally {
    _isLoading = false;
    _isInitialized = true;
    notify();
  }
}

export async function login(
  email: string,
  password: string
): Promise<{ success: boolean; error?: string }> {
  _isLoading = true;
  notify();
  try {
    const response = await api.login({ email, password });
    if (response.success && response.data) {
      _user = response.data.user;
      notify();
      return { success: true };
    }
    return {
      success: false,
      error: response.error?.message || 'Login failed',
    };
  } catch (error) {
    return {
      success: false,
      error: error instanceof Error ? error.message : 'Login failed',
    };
  } finally {
    _isLoading = false;
    notify();
  }
}

export async function register(
  email: string,
  password: string,
  name?: string
): Promise<{ success: boolean; error?: string }> {
  _isLoading = true;
  notify();
  try {
    const response = await api.register({ email, password, name });
    if (response.success && response.data) {
      _user = response.data.user;
      notify();
      return { success: true };
    }
    return {
      success: false,
      error: response.error?.message || 'Registration failed',
    };
  } catch (error) {
    return {
      success: false,
      error: error instanceof Error ? error.message : 'Registration failed',
    };
  } finally {
    _isLoading = false;
    notify();
  }
}

export async function logout(): Promise<void> {
  await api.logout();
  _user = null;
  notify();
}

export function getUser(): User | null {
  return _user;
}
export function getIsLoading(): boolean {
  return _isLoading;
}
export function getIsAuthenticated(): boolean {
  return _user !== null;
}
