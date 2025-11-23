package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Dracosaur Auxiliary", NewDracosaurAuxiliary)
}

// NewDracosaurAuxiliary creates a Dracosaur Auxiliary
// {4}{R}{R} - CREATURE
// Flying, Haste
func NewDracosaurAuxiliary(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Dracosaur Auxiliary")
	card.ManaCost = "{4}{R}{R}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"DINOSAUR", "DRAGON", "MOUNT"}
	card.Power = "4"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	ability1 := abilities.NewKeywordAbility(card.ID, abilities.KeywordHaste)
	card.AddAbility(ability1)
	ability2, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewDamageEffect(2)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability2)
	return card, nil
}
