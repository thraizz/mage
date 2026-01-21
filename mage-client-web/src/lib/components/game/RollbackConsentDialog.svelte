<script lang="ts">
	import type { PendingRollbackRequest } from '$lib/types/game';

	export let request: PendingRollbackRequest;
	export let onApprove: () => void;
	export let onDeny: () => void;
</script>

<div class="rollback-overlay">
	<div class="rollback-dialog">
		<div class="dialog-header">
			<span class="icon">⏪</span>
			<h3>Rollback Request</h3>
		</div>

		<div class="dialog-content">
			<p class="requester">
				<strong>{request.requestingPlayerName}</strong> wants to rollback the game
			</p>
			<div class="target-message">
				<span class="label">Rollback to:</span>
				<span class="message">"{request.targetMessageText}"</span>
			</div>
			<p class="note">This will undo all game actions since that point for all players.</p>
		</div>

		<div class="dialog-actions">
			<button class="deny-btn" onclick={onDeny}> Deny </button>
			<button class="approve-btn" onclick={onApprove}> Approve </button>
		</div>
	</div>
</div>

<style>
	.rollback-overlay {
		position: fixed;
		top: 0;
		left: 0;
		right: 0;
		bottom: 0;
		background: rgba(0, 0, 0, 0.7);
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

	.rollback-dialog {
		background: #1a1a2e;
		border: 1px solid #3d3d5c;
		border-radius: 12px;
		padding: 0;
		max-width: 400px;
		width: 90%;
		box-shadow: 0 8px 32px rgba(0, 0, 0, 0.5);
		animation: slideUp 0.2s ease-out;
	}

	@keyframes slideUp {
		from {
			opacity: 0;
			transform: translateY(20px);
		}
		to {
			opacity: 1;
			transform: translateY(0);
		}
	}

	.dialog-header {
		display: flex;
		align-items: center;
		gap: 10px;
		padding: 16px 20px;
		border-bottom: 1px solid #3d3d5c;
		background: #252545;
		border-radius: 12px 12px 0 0;
	}

	.dialog-header .icon {
		font-size: 24px;
	}

	.dialog-header h3 {
		margin: 0;
		font-size: 18px;
		color: #e0e0e0;
		font-weight: 600;
	}

	.dialog-content {
		padding: 20px;
	}

	.requester {
		margin: 0 0 16px 0;
		font-size: 15px;
		color: #c0c0c0;
	}

	.requester strong {
		color: #64b5f6;
	}

	.target-message {
		background: #0d0d1a;
		border-radius: 8px;
		padding: 12px;
		margin-bottom: 16px;
	}

	.target-message .label {
		display: block;
		font-size: 12px;
		color: #808080;
		margin-bottom: 4px;
		text-transform: uppercase;
	}

	.target-message .message {
		display: block;
		font-size: 14px;
		color: #ffcc80;
		font-style: italic;
		word-break: break-word;
	}

	.note {
		margin: 0;
		font-size: 13px;
		color: #909090;
		line-height: 1.5;
	}

	.dialog-actions {
		display: flex;
		gap: 12px;
		padding: 16px 20px;
		border-top: 1px solid #3d3d5c;
		background: #151528;
		border-radius: 0 0 12px 12px;
	}

	.dialog-actions button {
		flex: 1;
		padding: 12px 16px;
		border-radius: 8px;
		font-size: 14px;
		font-weight: 600;
		cursor: pointer;
		transition: all 0.15s ease;
	}

	.deny-btn {
		background: #2a2a3e;
		border: 1px solid #4a4a6e;
		color: #c0c0c0;
	}

	.deny-btn:hover {
		background: #3a3a4e;
		border-color: #ef5350;
		color: #ef5350;
	}

	.approve-btn {
		background: #2e7d32;
		border: 1px solid #4caf50;
		color: white;
	}

	.approve-btn:hover {
		background: #388e3c;
	}
</style>
