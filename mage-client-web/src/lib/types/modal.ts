/**
 * Modal dialog types
 */

export type ModalSize = 'small' | 'medium' | 'large';

export interface ModalProps {
	/**
	 * Whether the modal is currently open
	 */
	open: boolean;

	/**
	 * Modal title (optional)
	 */
	title?: string;

	/**
	 * Modal size
	 * @default 'medium'
	 */
	size?: ModalSize;

	/**
	 * Whether clicking the backdrop closes the modal
	 * @default true
	 */
	closeOnBackdrop?: boolean;

	/**
	 * Whether pressing ESC closes the modal
	 * @default true
	 */
	closeOnEsc?: boolean;

	/**
	 * Whether to show the close button (X icon)
	 * @default true
	 */
	showCloseButton?: boolean;

	/**
	 * Callback when the modal is closed
	 */
	onClose?: () => void;
}
