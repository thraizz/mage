package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Mossbridge Troll", NewMossbridgeTroll)
}

// NewMossbridgeTroll creates a Mossbridge Troll
// {5}{G}{G} - CREATURE
func NewMossbridgeTroll(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Mossbridge Troll")
	card.ManaCost = "{5}{G}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"TROLL"}
	card.Power = "5"
	card.Toughness = "5"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddEffect(abilities.NewBoostEffect(20, 20)).
		Build()
	card.AddAbility(ability0)
	return card, nil
}
