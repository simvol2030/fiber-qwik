/**
 * Admin Store - Plain TypeScript module
 * Module-level singleton state for admin panel
 */

let _sidebarCollapsed = false;
let _sidebarMobileOpen = false;
let _theme: 'light' | 'dark' = 'light';

interface Toast {
  id: string;
  type: 'success' | 'error' | 'info' | 'warning';
  message: string;
  duration?: number;
}

let _toasts: Toast[] = [];
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

export function initTheme() {
  if (typeof window !== 'undefined') {
    const saved = localStorage.getItem('admin-theme') as
      | 'light'
      | 'dark'
      | null;
    if (saved) {
      _theme = saved;
      document.documentElement.setAttribute('data-theme', saved);
    }
  }
}

export function toggleTheme() {
  _theme = _theme === 'light' ? 'dark' : 'light';
  if (typeof window !== 'undefined') {
    localStorage.setItem('admin-theme', _theme);
    document.documentElement.setAttribute('data-theme', _theme);
  }
  notify();
}

export function toggleSidebar() {
  _sidebarCollapsed = !_sidebarCollapsed;
  if (typeof window !== 'undefined') {
    localStorage.setItem(
      'admin-sidebar-collapsed',
      String(_sidebarCollapsed)
    );
  }
  notify();
}

export function toggleMobileSidebar() {
  _sidebarMobileOpen = !_sidebarMobileOpen;
  notify();
}

export function closeMobileSidebar() {
  _sidebarMobileOpen = false;
  notify();
}

export function initSidebar() {
  if (typeof window !== 'undefined') {
    const saved = localStorage.getItem('admin-sidebar-collapsed');
    if (saved === 'true') {
      _sidebarCollapsed = true;
    }
  }
}

function generateId(): string {
  return Math.random().toString(36).substring(2, 9);
}

function addToast(
  type: Toast['type'],
  message: string,
  duration = 5000
): string {
  const id = generateId();
  _toasts = [..._toasts, { id, type, message, duration }];
  notify();
  if (duration > 0) {
    setTimeout(() => removeToast(id), duration);
  }
  return id;
}

function removeToast(id: string) {
  _toasts = _toasts.filter((t) => t.id !== id);
  notify();
}

export const toast = {
  success: (message: string, duration?: number) =>
    addToast('success', message, duration),
  error: (message: string, duration?: number) =>
    addToast('error', message, duration),
  info: (message: string, duration?: number) =>
    addToast('info', message, duration),
  warning: (message: string, duration?: number) =>
    addToast('warning', message, duration),
  remove: removeToast,
};

export function getSidebarCollapsed() {
  return _sidebarCollapsed;
}
export function getSidebarMobileOpen() {
  return _sidebarMobileOpen;
}
export function getTheme() {
  return _theme;
}
export function getToasts(): Toast[] {
  return _toasts;
}
