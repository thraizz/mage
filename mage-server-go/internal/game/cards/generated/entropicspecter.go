package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Entropic Specter", NewEntropicSpecter)
}

// NewEntropicSpecter creates a Entropic Specter
// {3}{B}{B} - CREATURE
// Flying
func NewEntropicSpecter(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Entropic Specter")
	card.ManaCost = "{3}{B}{B}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"SPECTER", "SPIRIT"}
	card.Power = "0"
	card.Toughness = "0"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	// TODO: Implement spell ability with unmapped effects
	//   - ChooseOpponentEffect(Outcome.Detriment)
	// card.AddAbility(ability1)
	// TODO: Implement spell ability with unmapped effects
	//   - DiscardTargetEffect(1, false)
	// card.AddAbility(ability2)
	return card, nil
}
