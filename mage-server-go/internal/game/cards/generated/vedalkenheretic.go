package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Vedalken Heretic", NewVedalkenHeretic)
}

// NewVedalkenHeretic creates a Vedalken Heretic
// {G}{U} - CREATURE
func NewVedalkenHeretic(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Vedalken Heretic")
	card.ManaCost = "{G}{U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"VEDALKEN", "ROGUE"}
	card.Power = "1"
	card.Toughness = "1"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewDrawCardsEffect(1)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
