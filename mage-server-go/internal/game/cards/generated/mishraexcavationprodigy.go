package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Mishra Excavation Prodigy", NewMishraExcavationProdigy)
}

// NewMishraExcavationProdigy creates a Mishra Excavation Prodigy
// {2}{R} - CREATURE
// Haste
func NewMishraExcavationProdigy(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Mishra Excavation Prodigy")
	card.ManaCost = "{2}{R}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"HUMAN", "ARTIFICER"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "2"
	card.Toughness = "1"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordHaste)
	card.AddAbility(ability0)
	ability1 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddManaCost("{1}").
		AddTapCost().
		AddEffect(abilities.NewDrawCardsEffect(1)).
		Build()
	card.AddAbility(ability1)
	return card, nil
}
