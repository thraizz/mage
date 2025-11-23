package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/effects"
)

func init() {
	cards.Register("Torrent Of Souls", NewTorrentOfSouls)
}

// NewTorrentOfSouls creates a Torrent Of Souls
// {4}{B/R} - SORCERY
func NewTorrentOfSouls(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Torrent Of Souls")
	card.ManaCost = "{4}{B/R}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - ReturnFromGraveyardToBattlefieldTargetEffect()
	//
	// Targets:
	//   - abilities.NewPlayerTargetFilter()
	// card.AddAbility(ability0)
	return card, nil
}
