package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Pramikon Sky Rampart", NewPramikonSkyRampart)
}

// NewPramikonSkyRampart creates a Pramikon Sky Rampart
// {U}{R}{W} - CREATURE
// Flying, Defender
func NewPramikonSkyRampart(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Pramikon Sky Rampart")
	card.ManaCost = "{U}{R}{W}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"WALL"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "1"
	card.Toughness = "5"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	ability1 := abilities.NewKeywordAbility(card.ID, abilities.KeywordDefender)
	card.AddAbility(ability1)
	return card, nil
}
