<script lang="ts">
  import PlayerSeat from './PlayerSeat.svelte';

  interface Player {
    id: string;
    username: string;
    isHost: boolean;
    isReady?: boolean;
  }

  interface Props {
    players: Player[];
    maxPlayers: number;
    currentUsername: string;
    isHost: boolean;
    onkick?: (playerId: string) => void;
  }

  let { players, maxPlayers, currentUsername, isHost, onkick }: Props = $props();

  const columns = $derived(maxPlayers > 4 ? 3 : 2);
</script>

<div class="seat-layout" style="--columns: {columns}">
  {#each Array(maxPlayers) as _, i}
    {@const player = players[i]}
    {@const isCurrentUser = player?.username === currentUsername}
    <PlayerSeat
      {player}
      seatNumber={i + 1}
      {isCurrentUser}
      canKick={isHost && !player?.isHost}
      {onkick}
    />
  {/each}
</div>

<style>
  .seat-layout {
    display: grid;
    grid-template-columns: repeat(var(--columns), 1fr);
    gap: var(--space-4);
  }

  @media (max-width: 768px) {
    .seat-layout {
      grid-template-columns: 1fr;
    }
  }
</style>
