package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Swans Of Bryn Argoll", NewSwansOfBrynArgoll)
}

// NewSwansOfBrynArgoll creates a Swans Of Bryn Argoll
// {2}{W/U}{W/U} - CREATURE
// Flying
func NewSwansOfBrynArgoll(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Swans Of Bryn Argoll")
	card.ManaCost = "{2}{W/U}{W/U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"BIRD", "SPIRIT"}
	card.Power = "4"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	return card, nil
}