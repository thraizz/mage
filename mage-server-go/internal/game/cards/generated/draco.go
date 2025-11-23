package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Draco", NewDraco)
}

// NewDraco creates a Draco
// {16} - ARTIFACT CREATURE
// Flying
func NewDraco(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Draco")
	card.ManaCost = "{16}"
	card.Types = []string{"ARTIFACT", "CREATURE"}
	card.Subtypes = []string{"DRAGON"}
	card.Power = "9"
	card.Toughness = "9"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	return card, nil
}
