package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Gluntch The Bestower", NewGluntchTheBestower)
}

// NewGluntchTheBestower creates a Gluntch The Bestower
// {1}{G}{W} - CREATURE
// Flying
func NewGluntchTheBestower(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Gluntch The Bestower")
	card.ManaCost = "{1}{G}{W}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"JELLYFISH"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "0"
	card.Toughness = "5"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	return card, nil
}
