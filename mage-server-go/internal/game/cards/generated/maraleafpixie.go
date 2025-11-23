package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Maraleaf Pixie", NewMaraleafPixie)
}

// NewMaraleafPixie creates a Maraleaf Pixie
// {G}{U} - CREATURE
// Flying
func NewMaraleafPixie(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Maraleaf Pixie")
	card.ManaCost = "{G}{U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"FAERIE"}
	card.Power = "2"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	ability1 := abilities.BuildSimpleManaAbility(card.ID, "G")
	card.AddAbility(ability1)
	ability2 := abilities.BuildSimpleManaAbility(card.ID, "U")
	card.AddAbility(ability2)
	return card, nil
}
