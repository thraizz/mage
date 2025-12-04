package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Campaign Of Vengeance", NewCampaignOfVengeance)
}

// NewCampaignOfVengeance creates a Campaign Of Vengeance
// {3}{W}{B} - ENCHANTMENT
func NewCampaignOfVengeance(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Campaign Of Vengeance")
	card.ManaCost = "{3}{W}{B}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
