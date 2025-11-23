package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Tyranid Harridan", NewTyranidHarridan)
}

// NewTyranidHarridan creates a Tyranid Harridan
// {4}{G}{U} - CREATURE
// Flying
func NewTyranidHarridan(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Tyranid Harridan")
	card.ManaCost = "{4}{G}{U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"TYRANID"}
	card.Power = "4"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	return card, nil
}
