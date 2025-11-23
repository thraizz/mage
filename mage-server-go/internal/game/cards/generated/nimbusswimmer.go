package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Nimbus Swimmer", NewNimbusSwimmer)
}

// NewNimbusSwimmer creates a Nimbus Swimmer
// {X}{G}{U} - CREATURE
// Flying
func NewNimbusSwimmer(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Nimbus Swimmer")
	card.ManaCost = "{X}{G}{U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"LEVIATHAN"}
	card.Power = "0"
	card.Toughness = "0"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	return card, nil
}
