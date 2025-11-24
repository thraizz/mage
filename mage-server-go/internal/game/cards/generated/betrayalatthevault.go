package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Betrayal At The Vault", NewBetrayalAtTheVault)
}

// NewBetrayalAtTheVault creates a Betrayal At The Vault
// {4}{G}{G} - INSTANT
func NewBetrayalAtTheVault(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Betrayal At The Vault")
	card.ManaCost = "{4}{G}{G}"
	card.Types = []string{"INSTANT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}