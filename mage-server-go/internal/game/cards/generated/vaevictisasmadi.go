package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Vaevictis Asmadi", NewVaevictisAsmadi)
}

// NewVaevictisAsmadi creates a Vaevictis Asmadi
// {2}{B}{B}{R}{R}{G}{G} - CREATURE
// Flying
func NewVaevictisAsmadi(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Vaevictis Asmadi")
	card.ManaCost = "{2}{B}{B}{R}{R}{G}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"ELDER", "DRAGON"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "7"
	card.Toughness = "7"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	return card, nil
}
