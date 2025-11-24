package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/effects"
)

func init() {
	cards.Register("Dagobah Maw Slug", NewDagobahMawSlug)
}

// NewDagobahMawSlug creates a Dagobah Maw Slug
// {R}{G}{W} - CREATURE
func NewDagobahMawSlug(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Dagobah Maw Slug")
	card.ManaCost = "{R}{G}{W}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"SLUG"}
	card.Power = "4"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewGrantAbilityEffect("SwampwalkAbility")).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
