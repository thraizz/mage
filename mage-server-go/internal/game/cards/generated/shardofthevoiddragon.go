package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Shard Of The Void Dragon", NewShardOfTheVoidDragon)
}

// NewShardOfTheVoidDragon creates a Shard Of The Void Dragon
// {4}{B}{B}{B} - CREATURE
// Flying
func NewShardOfTheVoidDragon(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Shard Of The Void Dragon")
	card.ManaCost = "{4}{B}{B}{B}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"CTAN"}
	card.Power = "7"
	card.Toughness = "7"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	return card, nil
}
