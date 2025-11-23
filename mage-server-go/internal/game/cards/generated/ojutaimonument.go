package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Ojutai Monument", NewOjutaiMonument)
}

// NewOjutaiMonument creates a Ojutai Monument
// {3} - ARTIFACT
// Flying
func NewOjutaiMonument(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Ojutai Monument")
	card.ManaCost = "{3}"
	card.Types = []string{"ARTIFACT"}
	card.Subtypes = []string{"DRAGON"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	ability1 := abilities.BuildSimpleManaAbility(card.ID, "W")
	card.AddAbility(ability1)
	ability2 := abilities.BuildSimpleManaAbility(card.ID, "U")
	card.AddAbility(ability2)
	return card, nil
}
