/**
 * Confirmation dialog store for promise-based confirmations
 */

import { writable } from 'svelte/store';

export interface ConfirmOptions {
	title?: string;
	message?: string;
	confirmText?: string;
	cancelText?: string;
	destructive?: boolean;
}

interface ConfirmState extends ConfirmOptions {
	open: boolean;
	resolve?: (value: boolean) => void;
}

const initialState: ConfirmState = {
	open: false,
	title: 'Confirm',
	message: 'Are you sure?',
	confirmText: 'Confirm',
	cancelText: 'Cancel',
	destructive: false
};

/**
 * Create the confirmation dialog store
 */
function createConfirmStore() {
	const { subscribe, set, update } = writable<ConfirmState>(initialState);

	return {
		subscribe,

		/**
		 * Show a confirmation dialog and return a promise
		 *
		 * @param options - Confirmation dialog options
		 * @returns Promise that resolves to true if confirmed, false if cancelled
		 */
		confirm(options: ConfirmOptions = {}): Promise<boolean> {
			return new Promise((resolve) => {
				set({
					open: true,
					title: options.title || 'Confirm',
					message: options.message || 'Are you sure?',
					confirmText: options.confirmText || 'Confirm',
					cancelText: options.cancelText || 'Cancel',
					destructive: options.destructive || false,
					resolve
				});
			});
		},

		/**
		 * Handle confirm action
		 */
		handleConfirm() {
			update((state) => {
				if (state.resolve) {
					state.resolve(true);
				}
				return { ...initialState };
			});
		},

		/**
		 * Handle cancel action
		 */
		handleCancel() {
			update((state) => {
				if (state.resolve) {
					state.resolve(false);
				}
				return { ...initialState };
			});
		}
	};
}

/**
 * Singleton confirm store instance
 */
export const confirm = createConfirmStore();
