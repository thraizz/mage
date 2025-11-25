<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { auth } from '$lib/stores/auth';
	import { toast } from '$lib/stores/toast';
	import {
		fetchUserProfile,
		fetchUserStats,
		fetchMatchHistory,
		changePassword
	} from '$lib/api/profile';
	import type { UserProfile, UserStats, MatchHistory } from '$lib/types/profile';

	// State
	let profile = $state<UserProfile | null>(null);
	let stats = $state<UserStats | null>(null);
	let matches = $state<MatchHistory[]>([]);
	let loading = $state(true);
	let error = $state<string | null>(null);
	let showChangePassword = $state(false);

	// Change password form
	let currentPassword = $state('');
	let newPassword = $state('');
	let confirmPassword = $state('');
	let changingPassword = $state(false);
	let passwordError = $state<string | null>(null);

	/**
	 * Load profile data
	 */
	async function loadProfile(): Promise<void> {
		loading = true;
		error = null;

		try {
			// Load all profile data in parallel
			const [profileData, statsData, matchesData] = await Promise.all([
				fetchUserProfile(),
				fetchUserStats(),
				fetchMatchHistory(10)
			]);

			profile = profileData;
			stats = statsData;
			matches = matchesData;
		} catch (err) {
			console.error('Failed to load profile:', err);
			error = err instanceof Error ? err.message : 'Failed to load profile';

			// If session expired, redirect to login
			if (error.toLowerCase().includes('session')) {
				toast.error('Session expired - please login again');
				auth.logout();
				goto('/login');
			}
		} finally {
			loading = false;
		}
	}

	/**
	 * Handle password change
	 */
	async function handleChangePassword(): Promise<void> {
		// Validate inputs
		if (!currentPassword || !newPassword || !confirmPassword) {
			passwordError = 'All fields are required';
			return;
		}

		if (newPassword !== confirmPassword) {
			passwordError = 'New passwords do not match';
			return;
		}

		if (newPassword.length < 8) {
			passwordError = 'New password must be at least 8 characters';
			return;
		}

		if (newPassword === currentPassword) {
			passwordError = 'New password must be different from current password';
			return;
		}

		changingPassword = true;
		passwordError = null;

		try {
			await changePassword({ currentPassword, newPassword });
			toast.success('Password changed successfully');

			// Reset form
			currentPassword = '';
			newPassword = '';
			confirmPassword = '';
			showChangePassword = false;
		} catch (err) {
			console.error('Failed to change password:', err);
			passwordError = err instanceof Error ? err.message : 'Failed to change password';
		} finally {
			changingPassword = false;
		}
	}

	/**
	 * Cancel password change
	 */
	function cancelChangePassword(): void {
		showChangePassword = false;
		currentPassword = '';
		newPassword = '';
		confirmPassword = '';
		passwordError = null;
	}

	/**
	 * Format timestamp to readable date
	 */
	function formatDate(timestamp: number): string {
		return new Date(timestamp).toLocaleDateString('en-US', {
			year: 'numeric',
			month: 'long',
			day: 'numeric'
		});
	}

	/**
	 * Format duration in seconds to readable string
	 */
	function formatDuration(seconds: number): string {
		const minutes = Math.floor(seconds / 60);
		if (minutes < 60) {
			return `${minutes}m`;
		}
		const hours = Math.floor(minutes / 60);
		const remainingMinutes = minutes % 60;
		return `${hours}h ${remainingMinutes}m`;
	}

	/**
	 * Format time ago
	 */
	function formatTimeAgo(timestamp: number): string {
		const now = Date.now();
		const diff = now - timestamp;
		const seconds = Math.floor(diff / 1000);
		const minutes = Math.floor(seconds / 60);
		const hours = Math.floor(minutes / 60);
		const days = Math.floor(hours / 24);

		if (days > 0) return `${days}d ago`;
		if (hours > 0) return `${hours}h ago`;
		if (minutes > 0) return `${minutes}m ago`;
		return 'Just now';
	}

	// Load profile on mount
	onMount(() => {
		loadProfile();
	});

	// Derived username for display
	const username = $derived($auth.user?.username || 'Unknown');
</script>

<svelte:head>
	<title>Profile - MAGE</title>
</svelte:head>

<div class="container">
	<header>
		<h1>Player Profile</h1>
	</header>

	{#if loading}
		<div class="loading-state">
			<div class="spinner"></div>
			<p>Loading profile...</p>
		</div>
	{:else if error}
		<div class="error-state">
			<p class="error-message">{error}</p>
			<button class="btn-primary" onclick={loadProfile}>Retry</button>
		</div>
	{:else if profile && stats}
		<div class="profile-content">
			<!-- Profile Card (T036) -->
			<div class="profile-card">
				<div class="avatar">
					<div class="avatar-placeholder">{username.charAt(0).toUpperCase()}</div>
				</div>
				<h2>{profile.username}</h2>
				<p class="email">{profile.email}</p>
				<p class="joined">Joined: {formatDate(profile.createdAt)}</p>
				{#if profile.lastLogin}
					<p class="last-login">Last login: {formatTimeAgo(profile.lastLogin)}</p>
				{/if}
			</div>

			<!-- Statistics (T036) -->
			<div class="stats-grid">
				<div class="stat-card">
					<h3>Games Played</h3>
					<p class="stat-value">{stats.gamesPlayed}</p>
				</div>

				<div class="stat-card">
					<h3>Wins</h3>
					<p class="stat-value win">{stats.wins}</p>
				</div>

				<div class="stat-card">
					<h3>Losses</h3>
					<p class="stat-value loss">{stats.losses}</p>
				</div>

				<div class="stat-card">
					<h3>Win Rate</h3>
					<p class="stat-value">{stats.winRate.toFixed(1)}%</p>
				</div>

				{#if stats.draws > 0}
					<div class="stat-card">
						<h3>Draws</h3>
						<p class="stat-value">{stats.draws}</p>
					</div>
				{/if}

				{#if stats.totalPlayTime}
					<div class="stat-card">
						<h3>Play Time</h3>
						<p class="stat-value">{formatDuration(stats.totalPlayTime)}</p>
					</div>
				{/if}
			</div>

			<!-- Recent Matches (T038) -->
			<div class="section">
				<h3>Recent Matches</h3>
				{#if matches.length === 0}
					<div class="empty-state">
						<p>No matches played yet</p>
						<button class="btn-primary" onclick={() => goto('/lobby')}> Find a Game </button>
					</div>
				{:else}
					<div class="matches-list">
						{#each matches as match (match.id)}
							<div class="match-item {match.result}">
								<div class="match-info">
									<strong>vs {match.opponent}</strong>
									<span class="match-meta">
										{match.format}
										{#if match.duration}
											• {formatDuration(match.duration)}
										{/if}
										• {formatTimeAgo(match.timestamp)}
									</span>
								</div>
								<span class="result">{match.result.toUpperCase()}</span>
							</div>
						{/each}
					</div>
				{/if}
			</div>

			<!-- Settings (T037) -->
			<div class="section">
				<h3>Account Settings</h3>

				{#if !showChangePassword}
					<button class="btn-secondary" onclick={() => (showChangePassword = true)}>
						Change Password
					</button>
				{:else}
					<div class="password-form">
						<h4>Change Password</h4>

						{#if passwordError}
							<div class="error-banner">{passwordError}</div>
						{/if}

						<div class="form-group">
							<label for="current-password">Current Password</label>
							<input
								id="current-password"
								type="password"
								bind:value={currentPassword}
								disabled={changingPassword}
								placeholder="Enter current password"
							/>
						</div>

						<div class="form-group">
							<label for="new-password">New Password</label>
							<input
								id="new-password"
								type="password"
								bind:value={newPassword}
								disabled={changingPassword}
								placeholder="Enter new password (min 8 characters)"
							/>
						</div>

						<div class="form-group">
							<label for="confirm-password">Confirm New Password</label>
							<input
								id="confirm-password"
								type="password"
								bind:value={confirmPassword}
								disabled={changingPassword}
								placeholder="Confirm new password"
								onkeydown={(e) => e.key === 'Enter' && handleChangePassword()}
							/>
						</div>

						<div class="form-actions">
							<button
								class="btn-primary"
								onclick={handleChangePassword}
								disabled={changingPassword}
							>
								{changingPassword ? 'Changing...' : 'Change Password'}
							</button>
							<button
								class="btn-secondary"
								onclick={cancelChangePassword}
								disabled={changingPassword}
							>
								Cancel
							</button>
						</div>
					</div>
				{/if}
			</div>
		</div>
	{/if}
</div>

<style>
	.container {
		max-width: 1200px;
		margin: 0 auto;
		padding: 2rem;
	}

	header {
		margin-bottom: 2rem;
	}

	h1 {
		margin: 0;
		font-size: 2.5rem;
		color: #333;
	}

	/* Loading State */
	.loading-state {
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		padding: 4rem 2rem;
		gap: 1rem;
	}

	.spinner {
		width: 48px;
		height: 48px;
		border: 4px solid #e5e7eb;
		border-top-color: #667eea;
		border-radius: 50%;
		animation: spin 0.8s linear infinite;
	}

	@keyframes spin {
		to {
			transform: rotate(360deg);
		}
	}

	.loading-state p {
		color: #6b7280;
		font-size: 1rem;
	}

	/* Error State */
	.error-state {
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 1rem;
		padding: 2rem;
		background: #fef2f2;
		border: 1px solid #fecaca;
		border-radius: 8px;
	}

	.error-message {
		color: #dc2626;
		margin: 0;
	}

	/* Profile Content */
	.profile-content {
		display: flex;
		flex-direction: column;
		gap: 2rem;
	}

	.profile-card {
		background: white;
		border: 1px solid #ddd;
		border-radius: 8px;
		padding: 2rem;
		text-align: center;
	}

	.avatar {
		margin: 0 auto 1rem;
		width: 100px;
		height: 100px;
	}

	.avatar-placeholder {
		width: 100%;
		height: 100%;
		background: #667eea;
		border-radius: 50%;
		display: flex;
		align-items: center;
		justify-content: center;
		font-size: 3rem;
		font-weight: 700;
		color: white;
	}

	.profile-card h2 {
		margin: 0 0 0.5rem 0;
		color: #333;
		font-size: 1.75rem;
	}

	.profile-card p {
		margin: 0.25rem 0;
		color: #666;
	}

	.email {
		font-size: 1rem;
		color: #667eea;
	}

	.joined,
	.last-login {
		font-size: 0.875rem;
		color: #999;
	}

	/* Statistics Grid */
	.stats-grid {
		display: grid;
		grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
		gap: 1rem;
	}

	.stat-card {
		background: white;
		border: 1px solid #ddd;
		border-radius: 8px;
		padding: 1.5rem;
		text-align: center;
	}

	.stat-card h3 {
		margin: 0 0 0.5rem 0;
		font-size: 1rem;
		color: #666;
		font-weight: 500;
	}

	.stat-value {
		margin: 0;
		font-size: 2rem;
		font-weight: 700;
		color: #667eea;
	}

	.stat-value.win {
		color: #10b981;
	}

	.stat-value.loss {
		color: #ef4444;
	}

	/* Section */
	.section {
		background: white;
		border: 1px solid #ddd;
		border-radius: 8px;
		padding: 1.5rem;
	}

	.section h3 {
		margin: 0 0 1rem 0;
		color: #333;
		font-size: 1.5rem;
	}

	/* Empty State */
	.empty-state {
		text-align: center;
		padding: 2rem;
		color: #6b7280;
	}

	.empty-state p {
		margin: 0 0 1rem 0;
	}

	/* Matches List */
	.matches-list {
		display: flex;
		flex-direction: column;
		gap: 0.75rem;
	}

	.match-item {
		display: flex;
		justify-content: space-between;
		align-items: center;
		padding: 1rem;
		border-radius: 4px;
		background: #f9fafb;
	}

	.match-item.win {
		border-left: 4px solid #10b981;
	}

	.match-item.loss {
		border-left: 4px solid #ef4444;
	}

	.match-item.draw {
		border-left: 4px solid #6b7280;
	}

	.match-info {
		display: flex;
		flex-direction: column;
		gap: 0.25rem;
	}

	.match-info strong {
		color: #333;
		font-size: 1rem;
	}

	.match-meta {
		font-size: 0.875rem;
		color: #666;
	}

	.result {
		font-weight: 700;
		font-size: 0.875rem;
	}

	.match-item.win .result {
		color: #10b981;
	}

	.match-item.loss .result {
		color: #ef4444;
	}

	.match-item.draw .result {
		color: #6b7280;
	}

	/* Password Form */
	.password-form {
		padding: 1.5rem;
		background: #f9fafb;
		border: 1px solid #e5e7eb;
		border-radius: 8px;
	}

	.password-form h4 {
		margin: 0 0 1rem 0;
		color: #333;
		font-size: 1.125rem;
	}

	.form-group {
		margin-bottom: 1rem;
	}

	.form-group label {
		display: block;
		margin-bottom: 0.5rem;
		color: #374151;
		font-weight: 500;
		font-size: 0.875rem;
	}

	.form-group input {
		width: 100%;
		padding: 0.75rem;
		border: 1px solid #d1d5db;
		border-radius: 4px;
		font-size: 1rem;
		transition: border-color 0.2s;
	}

	.form-group input:focus {
		outline: none;
		border-color: #667eea;
	}

	.form-group input:disabled {
		background: #f3f4f6;
		cursor: not-allowed;
	}

	.form-actions {
		display: flex;
		gap: 1rem;
		margin-top: 1.5rem;
	}

	.error-banner {
		padding: 0.75rem;
		background: #fef2f2;
		border: 1px solid #fecaca;
		border-radius: 4px;
		color: #dc2626;
		font-size: 0.875rem;
		margin-bottom: 1rem;
	}

	/* Buttons */
	.btn-primary {
		padding: 0.75rem 1.5rem;
		background: #667eea;
		color: white;
		border: none;
		border-radius: 4px;
		font-size: 1rem;
		font-weight: 500;
		cursor: pointer;
		transition: background 0.2s;
	}

	.btn-primary:hover:not(:disabled) {
		background: #5568d3;
	}

	.btn-primary:disabled {
		opacity: 0.6;
		cursor: not-allowed;
	}

	.btn-secondary {
		padding: 0.75rem 1.5rem;
		background: white;
		color: #667eea;
		border: 1px solid #667eea;
		border-radius: 4px;
		font-size: 1rem;
		font-weight: 500;
		cursor: pointer;
		transition: all 0.2s;
	}

	.btn-secondary:hover:not(:disabled) {
		background: #667eea;
		color: white;
	}

	.btn-secondary:disabled {
		opacity: 0.6;
		cursor: not-allowed;
	}

	/* Responsive */
	@media (max-width: 768px) {
		.container {
			padding: 1rem;
		}

		h1 {
			font-size: 2rem;
		}

		.stats-grid {
			grid-template-columns: repeat(2, 1fr);
		}

		.form-actions {
			flex-direction: column;
		}

		.form-actions button {
			width: 100%;
		}
	}
</style>
