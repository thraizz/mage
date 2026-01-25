export type LoadingSize = 'small' | 'medium' | 'large';

export interface LoadingSpinnerProps {
  /**
   * Size of the spinner
   * @default 'medium'
   */
  size?: LoadingSize;

  /**
   * Optional label text shown below spinner
   */
  label?: string;

  /**
   * Show as fullscreen overlay with backdrop
   * @default false
   */
  overlay?: boolean;

  /**
   * Custom color for the spinner (CSS color value)
   * @default '#667eea' (primary brand color)
   */
  color?: string;
}
