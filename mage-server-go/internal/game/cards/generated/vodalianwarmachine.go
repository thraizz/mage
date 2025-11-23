package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Vodalian War Machine", NewVodalianWarMachine)
}

// NewVodalianWarMachine creates a Vodalian War Machine
// {1}{U}{U} - CREATURE
// Defender
func NewVodalianWarMachine(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Vodalian War Machine")
	card.ManaCost = "{1}{U}{U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"WALL"}
	card.Power = "0"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordDefender)
	card.AddAbility(ability0)
	return card, nil
}
