/**
 * Toast notification types
 */

export type ToastType = 'success' | 'error' | 'info' | 'warning';

export interface Toast {
  id: string;
  type: ToastType;
  message: string;
  duration?: number; // Duration in milliseconds, undefined means no auto-dismiss
  dismissible?: boolean; // Whether the toast can be manually dismissed
}

export interface ToastOptions {
  type?: ToastType;
  duration?: number;
  dismissible?: boolean;
}
