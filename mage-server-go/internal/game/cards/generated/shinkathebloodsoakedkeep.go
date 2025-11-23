package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/effects"
)

func init() {
	cards.Register("Shinka The Bloodsoaked Keep", NewShinkaTheBloodsoakedKeep)
}

// NewShinkaTheBloodsoakedKeep creates a Shinka The Bloodsoaked Keep
//   - LAND
func NewShinkaTheBloodsoakedKeep(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Shinka The Bloodsoaked Keep")
	card.ManaCost = ""
	card.Types = []string{"LAND"}
	card.Supertypes = []string{"LEGENDARY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.BuildSimpleManaAbility(card.ID, "R")
	card.AddAbility(ability0)
	ability1 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddTapCost().
		AddEffect(abilities.NewGrantAbilityEffect("FirstStrikeAbility", effects.DurationEndOfTurn)).
		Build()
	card.AddAbility(ability1)
	return card, nil
}
