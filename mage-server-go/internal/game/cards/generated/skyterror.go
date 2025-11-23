package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Sky Terror", NewSkyTerror)
}

// NewSkyTerror creates a Sky Terror
// {R}{W} - CREATURE
// Flying
func NewSkyTerror(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Sky Terror")
	card.ManaCost = "{R}{W}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"DINOSAUR"}
	card.Power = "2"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	return card, nil
}
