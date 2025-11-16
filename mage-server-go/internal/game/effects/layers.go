package effects

import (
	"strings"
	"sync"

	"github.com/google/uuid"
)

// Layer corresponds to the comprehensive rules layers for continuous effects.
type Layer int

const (
	LayerCopy Layer = 1 + iota
	LayerControl
	LayerText
	LayerType
	LayerColor
	LayerAbility
	LayerPowerToughness
)

var layerOrder = []Layer{
	LayerCopy,
	LayerControl,
	LayerText,
	LayerType,
	LayerColor,
	LayerAbility,
	LayerPowerToughness,
}

// Snapshot represents the mutable characteristics of a card when evaluating continuous effects.
type Snapshot struct {
	CardID        string
	ControllerID  string
	Types         []string
	BasePower     int
	BaseToughness int
	HasBasePower  bool
	HasBaseTough  bool
	Power         int
	Toughness     int
}

// NewSnapshot constructs a new snapshot for evaluation.
func NewSnapshot(cardID, controllerID string, types []string, basePower, baseToughness int, hasPower, hasToughness bool) *Snapshot {
	s := &Snapshot{
		CardID:        cardID,
		ControllerID:  controllerID,
		Types:         append([]string(nil), types...),
		BasePower:     basePower,
		BaseToughness: baseToughness,
		HasBasePower:  hasPower,
		HasBaseTough:  hasToughness,
	}
	s.Reset()
	return s
}

// Reset restores derived characteristics to their base values.
func (s *Snapshot) Reset() {
	if s.HasBasePower {
		s.Power = s.BasePower
	}
	if s.HasBaseTough {
		s.Toughness = s.BaseToughness
	}
}

// HasType returns true if the snapshot includes the provided type.
func (s *Snapshot) HasType(typeName string) bool {
	typeName = strings.ToLower(strings.TrimSpace(typeName))
	for _, t := range s.Types {
		if strings.ToLower(strings.TrimSpace(t)) == typeName {
			return true
		}
	}
	return false
}

// ContinuousEffect defines behaviour for modifying card characteristics.
type ContinuousEffect interface {
	ID() string
	Layer() Layer
	AppliesTo(*Snapshot) bool
	Apply(*Snapshot)
}

// LayerSystem manages registration and evaluation of continuous effects.
type LayerSystem struct {
	mu      sync.RWMutex
	effects map[Layer]map[string]ContinuousEffect
	index   map[string]Layer
}

// NewLayerSystem constructs an empty layer system.
func NewLayerSystem() *LayerSystem {
	return &LayerSystem{
		effects: make(map[Layer]map[string]ContinuousEffect),
		index:   make(map[string]Layer),
	}
}

// AddEffect registers a new continuous effect and returns its identifier.
func (ls *LayerSystem) AddEffect(effect ContinuousEffect) string {
	if effect == nil {
		return ""
	}
	ls.mu.Lock()
	defer ls.mu.Unlock()

	layer := effect.Layer()
	if layer == 0 {
		layer = LayerPowerToughness
	}

	id := effect.ID()
	if id == "" {
		id = uuid.NewString()
	}

	if _, ok := ls.effects[layer]; !ok {
		ls.effects[layer] = make(map[string]ContinuousEffect)
	}
	ls.effects[layer][id] = effect
	ls.index[id] = layer
	return id
}

// RemoveEffect removes a registered effect by ID.
func (ls *LayerSystem) RemoveEffect(id string) {
	if id == "" {
		return
	}
	ls.mu.Lock()
	defer ls.mu.Unlock()
	layer, ok := ls.index[id]
	if !ok {
		return
	}
	delete(ls.index, id)
	if layerMap, ok := ls.effects[layer]; ok {
		delete(layerMap, id)
		if len(layerMap) == 0 {
			delete(ls.effects, layer)
		}
	}
}

// Apply executes all relevant continuous effects across layers against the snapshot.
func (ls *LayerSystem) Apply(snapshot *Snapshot) {
	if snapshot == nil {
		return
	}
	ls.mu.RLock()
	defer ls.mu.RUnlock()

	snapshot.Reset()
	for _, layer := range layerOrder {
		layerEffects := ls.effects[layer]
		if len(layerEffects) == 0 {
			continue
		}
		for _, effect := range layerEffects {
			if effect.AppliesTo(snapshot) {
				effect.Apply(snapshot)
			}
		}
	}
}

// GetEffectsForCard returns all effects that apply to a specific card
func (ls *LayerSystem) GetEffectsForCard(cardID string) []ContinuousEffect {
	if ls == nil || cardID == "" {
		return nil
	}
	
	ls.mu.RLock()
	defer ls.mu.RUnlock()
	
	var result []ContinuousEffect
	snapshot := &Snapshot{CardID: cardID}
	
	for _, layerMap := range ls.effects {
		for _, effect := range layerMap {
			if effect.AppliesTo(snapshot) {
				result = append(result, effect)
			}
		}
	}
	
	return result
}

// HasEffectType checks if a card is affected by a specific effect type
// This is a helper for checking restrictions, requirements, etc.
func (ls *LayerSystem) HasEffectType(cardID string, checkFunc func(ContinuousEffect) bool) bool {
	if ls == nil || cardID == "" || checkFunc == nil {
		return false
	}
	
	ls.mu.RLock()
	defer ls.mu.RUnlock()
	
	snapshot := &Snapshot{CardID: cardID}
	
	for _, layerMap := range ls.effects {
		for _, effect := range layerMap {
			if effect.AppliesTo(snapshot) && checkFunc(effect) {
				return true
			}
		}
	}
	
	return false
}
