package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Scrapyard Steelbreaker", NewScrapyardSteelbreaker)
}

// NewScrapyardSteelbreaker creates a Scrapyard Steelbreaker
// {3}{R} - ARTIFACT CREATURE
func NewScrapyardSteelbreaker(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Scrapyard Steelbreaker")
	card.ManaCost = "{3}{R}"
	card.Types = []string{"ARTIFACT", "CREATURE"}
	card.Subtypes = []string{"HUMAN", "WARRIOR"}
	card.Power = "3"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddManaCost("{1}").
		AddEffect(abilities.NewBoostEffect(2, 1)).
		Build()
	card.AddAbility(ability0)
	return card, nil
}
