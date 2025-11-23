package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/effects"
)

func init() {
	cards.Register("Anikthea Hand Of Erebos", NewAniktheaHandOfErebos)
}

// NewAniktheaHandOfErebos creates a Anikthea Hand Of Erebos
// {2}{W}{B}{G} - ENCHANTMENT CREATURE
func NewAniktheaHandOfErebos(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Anikthea Hand Of Erebos")
	card.ManaCost = "{2}{W}{B}{G}"
	card.Types = []string{"ENCHANTMENT", "CREATURE"}
	card.Subtypes = []string{"DEMIGOD"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "4"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewGrantAbilityEffect(new MenaceAbility(false), filter, true)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}