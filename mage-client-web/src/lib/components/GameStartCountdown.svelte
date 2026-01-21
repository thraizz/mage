<script lang="ts">
	import { onDestroy } from 'svelte';

	// Props
	let {
		show = $bindable(false),
		onComplete,
		onCancel
	}: {
		show?: boolean;
		onComplete: () => void;
		onCancel?: () => void;
	} = $props();

	// State
	let countdown = $state(5);
	let countdownInterval: ReturnType<typeof setInterval> | null = null;

	/**
	 * Start countdown
	 */
	function startCountdown(): void {
		countdown = 5;

		countdownInterval = setInterval(() => {
			countdown--;

			if (countdown <= 0) {
				stopCountdown();
				onComplete();
			}
		}, 1000);
	}

	/**
	 * Stop countdown
	 */
	function stopCountdown(): void {
		if (countdownInterval) {
			clearInterval(countdownInterval);
			countdownInterval = null;
		}
	}

	/**
	 * Handle cancel
	 */
	function handleCancel(): void {
		stopCountdown();
		show = false;
		if (onCancel) {
			onCancel();
		}
	}

	/**
	 * Start countdown when shown
	 */
	$effect(() => {
		if (show) {
			startCountdown();
		} else {
			stopCountdown();
		}
	});

	/**
	 * Cleanup on unmount
	 */
	onDestroy(() => {
		stopCountdown();
	});

	/**
	 * Get countdown message
	 */
	const countdownMessage = $derived.by(() => {
		if (countdown > 1) {
			return `Game starting in ${countdown}...`;
		} else if (countdown === 1) {
			return 'Game starting in 1...';
		} else {
			return 'Starting game!';
		}
	});
</script>

{#if show}
	<!-- Backdrop -->
	<!-- svelte-ignore a11y_click_events_have_key_events -->
	<!-- svelte-ignore a11y_no_static_element_interactions -->
	<div class="countdown-backdrop" onclick={handleCancel}>
		<!-- Modal -->
		<!-- svelte-ignore a11y_click_events_have_key_events -->
		<!-- svelte-ignore a11y_no_static_element_interactions -->
		<div class="countdown-modal" onclick={(e) => e.stopPropagation()}>
			<!-- Countdown Circle -->
			<div class="countdown-circle">
				<svg class="countdown-ring" viewBox="0 0 120 120">
					<circle
						class="countdown-ring-background"
						cx="60"
						cy="60"
						r="54"
						fill="none"
						stroke="#e5e7eb"
						stroke-width="8"
					/>
					<circle
						class="countdown-ring-progress"
						cx="60"
						cy="60"
						r="54"
						fill="none"
						stroke="#667eea"
						stroke-width="8"
						stroke-linecap="round"
						style="--progress: {(countdown / 5) * 100}%"
					/>
				</svg>
				<div class="countdown-number">{countdown > 0 ? countdown : '🎮'}</div>
			</div>

			<!-- Message -->
			<h2 class="countdown-title">{countdownMessage()}</h2>
			<p class="countdown-subtitle">Get ready!</p>

			<!-- Cancel Button (optional) -->
			{#if onCancel && countdown > 0}
				<button class="btn-cancel" onclick={handleCancel}> Cancel </button>
			{/if}
		</div>
	</div>
{/if}

<style>
	.countdown-backdrop {
		position: fixed;
		top: 0;
		left: 0;
		right: 0;
		bottom: 0;
		background: rgba(0, 0, 0, 0.75);
		display: flex;
		align-items: center;
		justify-content: center;
		z-index: 1000;
		animation: fadeIn 0.2s ease-out;
	}

	@keyframes fadeIn {
		from {
			opacity: 0;
		}
		to {
			opacity: 1;
		}
	}

	.countdown-modal {
		background: white;
		border-radius: 1.5rem;
		padding: 3rem 2rem;
		max-width: 400px;
		width: 90%;
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 1.5rem;
		box-shadow:
			0 20px 25px -5px rgba(0, 0, 0, 0.1),
			0 10px 10px -5px rgba(0, 0, 0, 0.04);
		animation: scaleIn 0.3s ease-out;
	}

	@keyframes scaleIn {
		from {
			transform: scale(0.9);
			opacity: 0;
		}
		to {
			transform: scale(1);
			opacity: 1;
		}
	}

	.countdown-circle {
		position: relative;
		width: 160px;
		height: 160px;
		display: flex;
		align-items: center;
		justify-content: center;
	}

	.countdown-ring {
		position: absolute;
		width: 100%;
		height: 100%;
		transform: rotate(-90deg);
	}

	.countdown-ring-progress {
		stroke-dasharray: 339.292; /* 2 * π * 54 */
		stroke-dashoffset: calc(339.292 * (1 - var(--progress) / 100));
		transition:
			stroke-dashoffset 1s linear,
			stroke 0.3s;
	}

	.countdown-number {
		font-size: 4rem;
		font-weight: 800;
		color: #667eea;
		animation: pulse 1s ease-in-out infinite;
		user-select: none;
	}

	@keyframes pulse {
		0%,
		100% {
			transform: scale(1);
		}
		50% {
			transform: scale(1.1);
		}
	}

	.countdown-title {
		margin: 0;
		font-size: 1.5rem;
		font-weight: 700;
		color: #111827;
		text-align: center;
	}

	.countdown-subtitle {
		margin: 0;
		font-size: 1.125rem;
		color: #6b7280;
		text-align: center;
	}

	.btn-cancel {
		margin-top: 1rem;
		padding: 0.75rem 1.5rem;
		background: white;
		color: #6b7280;
		border: 1px solid #d1d5db;
		border-radius: 0.5rem;
		font-size: 0.9375rem;
		font-weight: 600;
		cursor: pointer;
		transition: all 0.2s;
	}

	.btn-cancel:hover {
		background: #f3f4f6;
		border-color: #9ca3af;
		color: #374151;
	}

	/* Responsive */
	@media (max-width: 640px) {
		.countdown-modal {
			padding: 2rem 1.5rem;
		}

		.countdown-circle {
			width: 140px;
			height: 140px;
		}

		.countdown-number {
			font-size: 3.5rem;
		}

		.countdown-title {
			font-size: 1.25rem;
		}

		.countdown-subtitle {
			font-size: 1rem;
		}
	}
</style>
