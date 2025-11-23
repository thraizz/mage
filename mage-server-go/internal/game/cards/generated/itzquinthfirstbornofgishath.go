package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Itzquinth Firstborn Of Gishath", NewItzquinthFirstbornOfGishath)
}

// NewItzquinthFirstbornOfGishath creates a Itzquinth Firstborn Of Gishath
// {R}{G} - CREATURE
// Haste
func NewItzquinthFirstbornOfGishath(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Itzquinth Firstborn Of Gishath")
	card.ManaCost = "{R}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"DINOSAUR"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "2"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordHaste)
	card.AddAbility(ability0)
	return card, nil
}
