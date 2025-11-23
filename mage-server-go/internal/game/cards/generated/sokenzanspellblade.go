package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Sokenzan Spellblade", NewSokenzanSpellblade)
}

// NewSokenzanSpellblade creates a Sokenzan Spellblade
// {4}{R} - CREATURE
func NewSokenzanSpellblade(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Sokenzan Spellblade")
	card.ManaCost = "{4}{R}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"OGRE", "SAMURAI", "SHAMAN"}
	card.Power = "2"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewBoostEffect(CardsInControllerHandCount.ANY, StaticValue.get(0))).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
