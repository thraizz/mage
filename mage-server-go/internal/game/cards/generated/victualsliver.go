package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/effects"
)

func init() {
	cards.Register("Victual Sliver", NewVictualSliver)
}

// NewVictualSliver creates a Victual Sliver
// {G}{W} - CREATURE
func NewVictualSliver(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Victual Sliver")
	card.ManaCost = "{G}{W}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"SLIVER"}
	card.Power = "2"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddSacrificeSourceCost().
		AddManaCost("{2}").
		AddEffect(abilities.NewGainLifeEffect(4)).
		Build()
	card.AddAbility(ability0)
	ability1, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewGrantAbilityEffect({2}, Sacrifice permanent: You gain 4 life.\"")).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability1)
	return card, nil
}