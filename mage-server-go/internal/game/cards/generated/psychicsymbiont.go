package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Psychic Symbiont", NewPsychicSymbiont)
}

// NewPsychicSymbiont creates a Psychic Symbiont
// {4}{U}{B} - CREATURE
// Flying
func NewPsychicSymbiont(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Psychic Symbiont")
	card.ManaCost = "{4}{U}{B}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"NIGHTMARE", "HORROR"}
	card.Power = "3"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	// TODO: Implement spell ability with unmapped effects
	//   - DiscardTargetEffect(1)
	// card.AddAbility(ability1)
	return card, nil
}
