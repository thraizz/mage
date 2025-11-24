package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Biomantic Mastery", NewBiomanticMastery)
}

// NewBiomanticMastery creates a Biomantic Mastery
// {4}{G/U}{G/U}{G/U} - SORCERY
func NewBiomanticMastery(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Biomantic Mastery")
	card.ManaCost = "{4}{G/U}{G/U}{G/U}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}