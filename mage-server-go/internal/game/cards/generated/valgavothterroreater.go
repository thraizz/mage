package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Valgavoth Terror Eater", NewValgavothTerrorEater)
}

// NewValgavothTerrorEater creates a Valgavoth Terror Eater
// {6}{B}{B}{B} - CREATURE
// Flying, Lifelink
func NewValgavothTerrorEater(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Valgavoth Terror Eater")
	card.ManaCost = "{6}{B}{B}{B}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"ELDER", "DEMON"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "9"
	card.Toughness = "9"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	ability1 := abilities.NewKeywordAbility(card.ID, abilities.KeywordLifelink)
	card.AddAbility(ability1)
	return card, nil
}
