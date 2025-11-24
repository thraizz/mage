package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Mnemonic Deluge", NewMnemonicDeluge)
}

// NewMnemonicDeluge creates a Mnemonic Deluge
// {6}{U}{U}{U} - SORCERY
func NewMnemonicDeluge(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Mnemonic Deluge")
	card.ManaCost = "{6}{U}{U}{U}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
