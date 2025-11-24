package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Zombie Assassin", NewZombieAssassin)
}

// NewZombieAssassin creates a Zombie Assassin
// {4}{B} - CREATURE
func NewZombieAssassin(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Zombie Assassin")
	card.ManaCost = "{4}{B}"
	card.Types = []string{"CREATURE"}
	card.Power = "3"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddTapCost().
		AddEffect(abilities.NewDestroyEffect(true)).
		Build()
	card.AddAbility(ability0)
	return card, nil
}
