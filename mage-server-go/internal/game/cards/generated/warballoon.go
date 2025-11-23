package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("War Balloon", NewWarBalloon)
}

// NewWarBalloon creates a War Balloon
// {2}{R} - ARTIFACT
// Flying
func NewWarBalloon(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "War Balloon")
	card.ManaCost = "{2}{R}"
	card.Types = []string{"ARTIFACT"}
	card.Subtypes = []string{"VEHICLE"}
	card.Power = "4"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	return card, nil
}
