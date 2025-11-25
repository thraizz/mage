package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Giants Ire", NewGiantsIre)
}

// NewGiantsIre creates a Giants Ire
// {3}{R} - KINDRED SORCERY
func NewGiantsIre(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Giants Ire")
	card.ManaCost = "{3}{R}"
	card.Types = []string{"KINDRED", "SORCERY"}
	card.Subtypes = []string{"GIANT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewDamageEffect(4)).
		AddEffect(abilities.NewDrawCardsEffect(1)).
		AddTarget(abilities.NewTargetRequirement(1, 1, abilities.NewAnyTargetFilter())).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
