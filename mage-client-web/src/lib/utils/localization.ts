export function toPossessiveName(playerName: string): string {
  return playerName === 'You' ? 'Your' : `${playerName}'s`;
}
