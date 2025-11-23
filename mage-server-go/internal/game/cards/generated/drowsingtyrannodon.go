package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Drowsing Tyrannodon", NewDrowsingTyrannodon)
}

// NewDrowsingTyrannodon creates a Drowsing Tyrannodon
// {1}{G} - CREATURE
// Defender
func NewDrowsingTyrannodon(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Drowsing Tyrannodon")
	card.ManaCost = "{1}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"DINOSAUR"}
	card.Power = "3"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordDefender)
	card.AddAbility(ability0)
	return card, nil
}
