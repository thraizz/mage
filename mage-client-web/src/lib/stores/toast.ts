/**
 * Toast notification store for displaying user feedback messages
 */

import { writable } from 'svelte/store';
import type { Toast, ToastOptions } from '$lib/types/toast';

// Default toast duration (3 seconds)
const DEFAULT_DURATION = 3000;

/**
 * Generate a unique ID for a toast
 */
function generateId(): string {
  return `toast-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`;
}

/**
 * Create the toast notification store
 */
function createToastStore() {
  const { subscribe, update } = writable<Toast[]>([]);

  return {
    subscribe,

    /**
     * Add a toast notification
     *
     * @param message - The message to display
     * @param options - Optional configuration (type, duration, dismissible)
     * @returns The toast ID
     */
    add(message: string, options: ToastOptions = {}): string {
      const id = generateId();
      const toast: Toast = {
        id,
        type: options.type || 'info',
        message,
        duration: options.duration !== undefined ? options.duration : DEFAULT_DURATION,
        dismissible: options.dismissible !== undefined ? options.dismissible : true
      };

      update((toasts) => [...toasts, toast]);

      // Auto-dismiss if duration is set
      if (toast.duration && toast.duration > 0) {
        setTimeout(() => {
          this.dismiss(id);
        }, toast.duration);
      }

      return id;
    },

    /**
     * Show a success toast
     *
     * @param message - Success message
     * @param duration - Optional duration override
     */
    success(message: string, duration?: number): string {
      return this.add(message, { type: 'success', duration });
    },

    /**
     * Show an error toast
     *
     * @param message - Error message
     * @param duration - Optional duration override
     */
    error(message: string, duration?: number): string {
      return this.add(message, { type: 'error', duration });
    },

    /**
     * Show an info toast
     *
     * @param message - Info message
     * @param duration - Optional duration override
     */
    info(message: string, duration?: number): string {
      return this.add(message, { type: 'info', duration });
    },

    /**
     * Show a warning toast
     *
     * @param message - Warning message
     * @param duration - Optional duration override
     */
    warning(message: string, duration?: number): string {
      return this.add(message, { type: 'warning', duration });
    },

    /**
     * Dismiss a specific toast
     *
     * @param id - The toast ID to dismiss
     */
    dismiss(id: string): void {
      update((toasts) => toasts.filter((t) => t.id !== id));
    },

    /**
     * Dismiss all toasts
     */
    dismissAll(): void {
      update(() => []);
    }
  };
}

/**
 * Singleton toast store instance
 */
export const toast = createToastStore();
