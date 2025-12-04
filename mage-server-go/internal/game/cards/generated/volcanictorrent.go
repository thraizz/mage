package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Volcanic Torrent", NewVolcanicTorrent)
}

// NewVolcanicTorrent creates a Volcanic Torrent
// {4}{R} - SORCERY
func NewVolcanicTorrent(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Volcanic Torrent")
	card.ManaCost = "{4}{R}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - DamageAllEffect(VolcanicTorrentValue.instance, filter)
	// card.AddAbility(ability0)
	return card, nil
}
