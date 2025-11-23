package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Faridehs Fireball", NewFaridehsFireball)
}

// NewFaridehsFireball creates a Faridehs Fireball
// {3}{R}{R} - INSTANT
func NewFaridehsFireball(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Faridehs Fireball")
	card.ManaCost = "{3}{R}{R}"
	card.Types = []string{"INSTANT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewDamageEffect(5)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
