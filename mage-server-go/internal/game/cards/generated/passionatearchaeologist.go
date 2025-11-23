package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/effects"
)

func init() {
	cards.Register("Passionate Archaeologist", NewPassionateArchaeologist)
}

// NewPassionateArchaeologist creates a Passionate Archaeologist
// {1}{R} - ENCHANTMENT
func NewPassionateArchaeologist(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Passionate Archaeologist")
	card.ManaCost = "{1}{R}"
	card.Types = []string{"ENCHANTMENT"}
	card.Subtypes = []string{"BACKGROUND"}
	card.Supertypes = []string{"LEGENDARY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewGrantAbilityEffect()).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
