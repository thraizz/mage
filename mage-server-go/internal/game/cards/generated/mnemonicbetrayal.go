package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Mnemonic Betrayal", NewMnemonicBetrayal)
}

// NewMnemonicBetrayal creates a Mnemonic Betrayal
// {1}{U}{B} - SORCERY
func NewMnemonicBetrayal(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Mnemonic Betrayal")
	card.ManaCost = "{1}{U}{B}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}