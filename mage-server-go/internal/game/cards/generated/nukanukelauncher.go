package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Nuka Nuke Launcher", NewNukaNukeLauncher)
}

// NewNukaNukeLauncher creates a Nuka Nuke Launcher
// {2} - ARTIFACT
func NewNukaNukeLauncher(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Nuka Nuke Launcher")
	card.ManaCost = "{2}"
	card.Types = []string{"ARTIFACT"}
	card.Subtypes = []string{"EQUIPMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewEquipAbility(card.ID, "{3}", true)
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	ability1, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewBoostEquippedEffect(3, 0)).
		AddEffect(abilities.NewGainAbilityAttachedEffect(abilities.NewKeywordAbility(card.ID, abilities.KeywordIntimidate), abilities.AttachmentTypeEquipment, abilities.DurationWhileOnBattlefield, "")).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability1)
	return card, nil
}
