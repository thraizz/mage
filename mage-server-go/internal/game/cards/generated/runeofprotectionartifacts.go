package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Rune Of Protection Artifacts", NewRuneOfProtectionArtifacts)
}

// NewRuneOfProtectionArtifacts creates a Rune Of Protection Artifacts
// {1}{W} - ENCHANTMENT
func NewRuneOfProtectionArtifacts(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Rune Of Protection Artifacts")
	card.ManaCost = "{1}{W}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
