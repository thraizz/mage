package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Flame Wreathed Phoenix", NewFlameWreathedPhoenix)
}

// NewFlameWreathedPhoenix creates a Flame Wreathed Phoenix
// {2}{R}{R} - CREATURE
// Flying
func NewFlameWreathedPhoenix(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Flame Wreathed Phoenix")
	card.ManaCost = "{2}{R}{R}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"PHOENIX"}
	card.Power = "3"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	return card, nil
}
