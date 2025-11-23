package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Spoils Of The Vault", NewSpoilsOfTheVault)
}

// NewSpoilsOfTheVault creates a Spoils Of The Vault
// {B} - INSTANT
func NewSpoilsOfTheVault(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Spoils Of The Vault")
	card.ManaCost = "{B}"
	card.Types = []string{"INSTANT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
