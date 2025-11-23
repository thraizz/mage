package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Boon Bringer Valkyrie", NewBoonBringerValkyrie)
}

// NewBoonBringerValkyrie creates a Boon Bringer Valkyrie
// {3}{W}{W} - CREATURE
// Flying, FirstStrike, Lifelink
func NewBoonBringerValkyrie(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Boon Bringer Valkyrie")
	card.ManaCost = "{3}{W}{W}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"ANGEL", "WARRIOR"}
	card.Power = "4"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	ability1 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFirstStrike)
	card.AddAbility(ability1)
	ability2 := abilities.NewKeywordAbility(card.ID, abilities.KeywordLifelink)
	card.AddAbility(ability2)
	return card, nil
}
