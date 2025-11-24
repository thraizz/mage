package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Lurrus Of The Dream Den", NewLurrusOfTheDreamDen)
}

// NewLurrusOfTheDreamDen creates a Lurrus Of The Dream Den
// {1}{W/B}{W/B} - CREATURE
// Lifelink
func NewLurrusOfTheDreamDen(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Lurrus Of The Dream Den")
	card.ManaCost = "{1}{W/B}{W/B}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"CAT", "NIGHTMARE"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "3"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordLifelink)
	card.AddAbility(ability0)
	return card, nil
}