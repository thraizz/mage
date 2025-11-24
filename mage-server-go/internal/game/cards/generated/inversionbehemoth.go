package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Inversion Behemoth", NewInversionBehemoth)
}

// NewInversionBehemoth creates a Inversion Behemoth
// {2}{C}{C} - CREATURE
func NewInversionBehemoth(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Inversion Behemoth")
	card.ManaCost = "{2}{C}{C}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"ELDRAZI"}
	card.Power = "2"
	card.Toughness = "9"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}