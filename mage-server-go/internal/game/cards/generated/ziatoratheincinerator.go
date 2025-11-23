package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Ziatora The Incinerator", NewZiatoraTheIncinerator)
}

// NewZiatoraTheIncinerator creates a Ziatora The Incinerator
// {3}{B}{R}{G} - CREATURE
// Flying
func NewZiatoraTheIncinerator(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Ziatora The Incinerator")
	card.ManaCost = "{3}{B}{R}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"DEMON", "DRAGON"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "6"
	card.Toughness = "6"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	return card, nil
}
