package counters

// Counter represents a counter on a permanent or player.
// Mirrors Java Counter class.
type Counter struct {
	Name  string
	Count int
}

// NewCounter creates a new counter with the given name and count.
func NewCounter(name string, count int) *Counter {
	if count <= 0 {
		count = 1
	}
	return &Counter{
		Name:  name,
		Count: count,
	}
}

// Add adds the specified amount to the counter.
func (c *Counter) Add(amount int) {
	if amount > 0 {
		c.Count += amount
	}
}

// Remove removes the specified amount from the counter.
// Will not allow count to go below 0.
func (c *Counter) Remove(amount int) {
	if amount > 0 {
		if c.Count >= amount {
			c.Count -= amount
		} else {
			c.Count = 0
		}
	}
}

// Copy creates a deep copy of the counter.
func (c *Counter) Copy() *Counter {
	return &Counter{
		Name:  c.Name,
		Count: c.Count,
	}
}

// BoostCounter represents a power/toughness boost counter (e.g., +1/+1, -1/-1).
type BoostCounter struct {
	*Counter
	Power     int
	Toughness int
}

// NewBoostCounter creates a new boost counter.
func NewBoostCounter(power, toughness, count int) *BoostCounter {
	name := getBoostCounterName(power, toughness)
	return &BoostCounter{
		Counter:   NewCounter(name, count),
		Power:     power,
		Toughness: toughness,
	}
}

// Copy creates a deep copy of the boost counter.
func (bc *BoostCounter) Copy() *BoostCounter {
	return &BoostCounter{
		Counter:   bc.Counter.Copy(),
		Power:     bc.Power,
		Toughness: bc.Toughness,
	}
}

// getBoostCounterName generates a name for a boost counter (e.g., "+1/+1", "-1/-1").
func getBoostCounterName(power, toughness int) string {
	powerStr := formatBoost(power)
	toughnessStr := formatBoost(toughness)
	return powerStr + "/" + toughnessStr
}

func formatBoost(value int) string {
	if value > 0 {
		return "+" + formatInt(value)
	} else if value < 0 {
		return formatInt(value)
	}
	return "±0"
}

func formatInt(value int) string {
	if value < 0 {
		return formatInt(-value)
	}
	if value < 10 {
		return string(rune('0' + value))
	}
	// For values >= 10, convert to string
	result := ""
	for value > 0 {
		result = string(rune('0'+(value%10))) + result
		value /= 10
	}
	return result
}

// Counters manages a collection of counters.
// Mirrors Java Counters class.
type Counters struct {
	Counters map[string]*Counter
}

// NewCounters creates a new Counters collection.
func NewCounters() *Counters {
	return &Counters{
		Counters: make(map[string]*Counter),
	}
}

// AddCounter adds a counter to the collection.
// If a counter with the same name already exists, adds to its count.
func (cs *Counters) AddCounter(counter *Counter) {
	if counter == nil {
		return
	}
	if existing, ok := cs.Counters[counter.Name]; ok {
		existing.Add(counter.Count)
	} else {
		cs.Counters[counter.Name] = counter.Copy()
	}
}

// RemoveCounter removes the specified amount of counters of the given name.
// Returns true if any counters were removed.
func (cs *Counters) RemoveCounter(name string, amount int) bool {
	if amount <= 0 {
		return false
	}
	if counter, ok := cs.Counters[name]; ok {
		counter.Remove(amount)
		if counter.Count == 0 {
			delete(cs.Counters, name)
		}
		return true
	}
	return false
}

// GetCount returns the count of counters with the given name.
func (cs *Counters) GetCount(name string) int {
	if counter, ok := cs.Counters[name]; ok {
		return counter.Count
	}
	return 0
}

// HasCounter returns true if there are any counters with the given name.
func (cs *Counters) HasCounter(name string) bool {
	return cs.GetCount(name) > 0
}

// GetTotalCount returns the total number of all counters.
func (cs *Counters) GetTotalCount() int {
	total := 0
	for _, counter := range cs.Counters {
		total += counter.Count
	}
	return total
}

// GetAll returns all counters as a map.
func (cs *Counters) GetAll() map[string]*Counter {
	result := make(map[string]*Counter)
	for name, counter := range cs.Counters {
		result[name] = counter.Copy()
	}
	return result
}

// GetBoostCounters returns all boost counters (power/toughness modifying counters).
// Checks counter names for boost counter patterns (e.g., "+1/+1", "-1/-1").
func (cs *Counters) GetBoostCounters() []*BoostCounter {
	var boostCounters []*BoostCounter
	for _, counter := range cs.Counters {
		// Check if counter name matches boost counter pattern (e.g., "+1/+1", "-1/-1")
		if power, toughness, ok := parseBoostCounterName(counter.Name); ok {
			boostCounters = append(boostCounters, NewBoostCounter(power, toughness, counter.Count))
		}
	}
	return boostCounters
}

// parseBoostCounterName parses a boost counter name (e.g., "+1/+1") into power/toughness deltas.
// Returns power, toughness, and true if parsing succeeded.
func parseBoostCounterName(name string) (int, int, bool) {
	// Simple pattern matching for boost counters like "+1/+1", "-1/-1", "+2/+2", etc.
	// This is a simplified parser - full implementation would handle all edge cases
	if len(name) < 3 {
		return 0, 0, false
	}
	// Look for pattern: [+-]?[0-9]+/[+-]?[0-9]+
	// For now, handle common cases
	power, toughness := 0, 0
	parts := splitBoostName(name)
	if len(parts) != 2 {
		return 0, 0, false
	}
	var ok bool
	power, ok = parseBoostValue(parts[0])
	if !ok {
		return 0, 0, false
	}
	toughness, ok = parseBoostValue(parts[1])
	if !ok {
		return 0, 0, false
	}
	return power, toughness, true
}

func splitBoostName(name string) []string {
	// Split on "/" to get power and toughness parts
	for i := 0; i < len(name); i++ {
		if name[i] == '/' {
			return []string{name[:i], name[i+1:]}
		}
	}
	return nil
}

func parseBoostValue(s string) (int, bool) {
	if len(s) == 0 {
		return 0, false
	}
	negative := false
	start := 0
	if s[0] == '+' {
		start = 1
	} else if s[0] == '-' {
		negative = true
		start = 1
	}
	if start >= len(s) {
		return 0, false
	}
	value := 0
	for i := start; i < len(s); i++ {
		if s[i] >= '0' && s[i] <= '9' {
			value = value*10 + int(s[i]-'0')
		} else {
			return 0, false
		}
	}
	if negative {
		value = -value
	}
	return value, true
}

// Copy creates a deep copy of the Counters collection.
func (cs *Counters) Copy() *Counters {
	copy := NewCounters()
	for name, counter := range cs.Counters {
		copy.Counters[name] = counter.Copy()
	}
	return copy
}

// ToView converts counters to the view format.
func (cs *Counters) ToView() []CounterView {
	var views []CounterView
	for name, counter := range cs.Counters {
		views = append(views, CounterView{
			Name:  name,
			Count: counter.Count,
		})
	}
	return views
}

// CounterView represents a counter in the view format.
type CounterView struct {
	Name  string
	Count int
}
