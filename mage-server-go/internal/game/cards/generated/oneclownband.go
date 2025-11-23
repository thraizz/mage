package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("One Clown Band", NewOneClownBand)
}

// NewOneClownBand creates a One Clown Band
// {1}{R} - ARTIFACT CREATURE
func NewOneClownBand(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "One Clown Band")
	card.ManaCost = "{1}{R}"
	card.Types = []string{"ARTIFACT", "CREATURE"}
	card.Subtypes = []string{"CLOWN", "ROBOT", "BARD"}
	card.Power = "1"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddEffect(abilities.NewBoostEffect(2, 0)).
		Build()
	card.AddAbility(ability0)
	return card, nil
}
