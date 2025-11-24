package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Godhead Of Awe", NewGodheadOfAwe)
}

// NewGodheadOfAwe creates a Godhead Of Awe
// {W/U}{W/U}{W/U}{W/U}{W/U} - CREATURE
// Flying
func NewGodheadOfAwe(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Godhead Of Awe")
	card.ManaCost = "{W/U}{W/U}{W/U}{W/U}{W/U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"SPIRIT", "AVATAR"}
	card.Power = "4"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	return card, nil
}