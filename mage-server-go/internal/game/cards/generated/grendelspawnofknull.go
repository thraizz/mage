package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Grendel Spawn Of Knull", NewGrendelSpawnOfKnull)
}

// NewGrendelSpawnOfKnull creates a Grendel Spawn Of Knull
// {3}{B} - CREATURE
// Flying, Deathtouch
func NewGrendelSpawnOfKnull(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Grendel Spawn Of Knull")
	card.ManaCost = "{3}{B}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"SYMBIOTE", "DRAGON", "VILLAIN"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "3"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	ability1 := abilities.NewKeywordAbility(card.ID, abilities.KeywordDeathtouch)
	card.AddAbility(ability1)
	return card, nil
}
