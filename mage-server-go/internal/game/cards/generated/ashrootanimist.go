package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Ashroot Animist", NewAshrootAnimist)
}

// NewAshrootAnimist creates a Ashroot Animist
// {2}{R}{G} - CREATURE
// Trample
func NewAshrootAnimist(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Ashroot Animist")
	card.ManaCost = "{2}{R}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"LIZARD", "DRUID"}
	card.Power = "4"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordTrample)
	card.AddAbility(ability0)
	return card, nil
}
