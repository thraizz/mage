/**
 * Polling Utility
 * Provides smart periodic fetching with WebSocket awareness and visibility detection
 */

import { onMount, onDestroy } from 'svelte';
import { get } from 'svelte/store';
import { websocketStore } from '$lib/stores/websocket';

export interface PollingOptions {
  /**
   * Interval in milliseconds when tab is visible
   * @default 5000
   */
  interval?: number;

  /**
   * Interval in milliseconds when tab is hidden
   * @default 30000
   */
  intervalWhenHidden?: number;

  /**
   * Whether to poll when WebSocket is connected
   * @default false (only poll as fallback when disconnected)
   */
  pollWhenConnected?: boolean;

  /**
   * Whether to fetch immediately on mount
   * @default true
   */
  immediate?: boolean;

  /**
   * Whether polling is enabled
   * @default true
   */
  enabled?: boolean;
}

/**
 * Create a polling function that fetches data periodically with smart behavior
 *
 * Features:
 * - Only polls as fallback when WebSocket disconnected (by default)
 * - Reduces polling frequency when tab is hidden
 * - Automatically cleans up on component unmount
 * - Provides manual refresh capability
 *
 * @example
 * ```ts
 * const { refresh, isPolling } = usePolling(loadTables, {
 *   interval: 5000,
 *   intervalWhenHidden: 30000
 * });
 * ```
 */
export function usePolling(fetchFn: () => Promise<void>, options: PollingOptions = {}) {
  const {
    interval = 5000,
    intervalWhenHidden = 30000,
    pollWhenConnected = false,
    immediate = true,
    enabled = true
  } = options;

  let timerId: ReturnType<typeof setTimeout> | null = null;
  let isVisible = true;
  let isMounted = true;

  /**
   * Get current polling interval based on visibility
   */
  function getCurrentInterval(): number {
    return isVisible ? interval : intervalWhenHidden;
  }

  /**
   * Check if we should poll right now
   */
  function shouldPoll(): boolean {
    if (!enabled || !isMounted) return false;

    const wsState = get(websocketStore).state;
    const isWsConnected = wsState === 'connected';

    // If WebSocket is connected and we don't want to poll when connected, skip
    if (isWsConnected && !pollWhenConnected) {
      return false;
    }

    return true;
  }

  /**
   * Schedule next poll
   */
  function scheduleNext() {
    if (!shouldPoll()) return;

    if (timerId) {
      clearTimeout(timerId);
    }

    timerId = setTimeout(async () => {
      if (shouldPoll()) {
        try {
          await fetchFn();
        } catch (err) {
          console.error('[Polling] Fetch error:', err);
        }
      }
      scheduleNext();
    }, getCurrentInterval());
  }

  /**
   * Start polling
   */
  function start() {
    if (!shouldPoll()) return;
    scheduleNext();
  }

  /**
   * Stop polling
   */
  function stop() {
    if (timerId) {
      clearTimeout(timerId);
      timerId = null;
    }
  }

  /**
   * Manual refresh - fetch immediately and reset timer
   */
  async function refresh(): Promise<void> {
    stop();
    try {
      await fetchFn();
    } finally {
      start();
    }
  }

  /**
   * Handle visibility change
   */
  function handleVisibilityChange() {
    isVisible = !document.hidden;

    // If becoming visible and should poll, fetch immediately
    if (isVisible && shouldPoll()) {
      refresh();
    } else {
      // Reschedule with new interval
      stop();
      start();
    }
  }

  /**
   * Handle WebSocket state change
   */
  function handleWebSocketChange() {
    // When WebSocket state changes, re-evaluate if we should poll
    if (shouldPoll()) {
      start();
    } else {
      stop();
    }
  }

  // Setup on mount
  onMount(() => {
    // Listen for visibility changes
    document.addEventListener('visibilitychange', handleVisibilityChange);

    // Subscribe to WebSocket state changes
    const unsubscribeWs = websocketStore.subscribe(handleWebSocketChange);

    // Initial fetch if immediate
    if (immediate && shouldPoll()) {
      fetchFn().catch((err) => {
        console.error('[Polling] Initial fetch error:', err);
      });
    }

    // Start polling
    start();

    // Cleanup function
    return () => {
      isMounted = false;
      document.removeEventListener('visibilitychange', handleVisibilityChange);
      unsubscribeWs();
      stop();
    };
  });

  // Cleanup on destroy
  onDestroy(() => {
    isMounted = false;
    stop();
  });

  return {
    refresh,
    stop,
    start,
    get isPolling() {
      return timerId !== null;
    }
  };
}

/**
 * Simpler hook for components that just need periodic refresh
 * with sensible defaults
 *
 * @example
 * ```ts
 * usePeriodicRefresh(loadData);
 * ```
 */
export function usePeriodicRefresh(fetchFn: () => Promise<void>, intervalMs = 5000) {
  return usePolling(fetchFn, {
    interval: intervalMs,
    intervalWhenHidden: intervalMs * 6, // 6x slower when hidden
    immediate: true,
    enabled: true
  });
}
