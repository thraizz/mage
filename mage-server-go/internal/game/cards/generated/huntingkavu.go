package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Hunting Kavu", NewHuntingKavu)
}

// NewHuntingKavu creates a Hunting Kavu
// {1}{R}{G} - CREATURE
func NewHuntingKavu(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Hunting Kavu")
	card.ManaCost = "{1}{R}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"KAVU"}
	card.Power = "2"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddTapCost().
		AddEffect(abilities.NewExileSourceEffect()).
		AddEffect(abilities.NewExileTargetEffect()).
		Build()
	card.AddAbility(ability0)
	return card, nil
}
