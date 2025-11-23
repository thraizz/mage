package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Heart Of Kiran", NewHeartOfKiran)
}

// NewHeartOfKiran creates a Heart Of Kiran
// {2} - ARTIFACT
// Flying, Vigilance
func NewHeartOfKiran(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Heart Of Kiran")
	card.ManaCost = "{2}"
	card.Types = []string{"ARTIFACT"}
	card.Subtypes = []string{"VEHICLE"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "4"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	ability1 := abilities.NewKeywordAbility(card.ID, abilities.KeywordVigilance)
	card.AddAbility(ability1)
	return card, nil
}
