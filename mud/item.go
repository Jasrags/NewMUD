package mud

import "sync"

type Item struct {
	ID          string `yaml:"id"`
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
}

type ItemManager struct {
	mu    sync.RWMutex
	items map[string]*Item // Keyed by item ID
}

// NewItemManager creates and initializes an ItemManager
func NewItemManager() *ItemManager {
	return &ItemManager{
		items: make(map[string]*Item),
	}
}

// AddItem adds an item to the manager
func (im *ItemManager) AddItem(item *Item) {
	im.mu.Lock()
	defer im.mu.Unlock()
	im.items[item.ID] = item
}

// GetItem retrieves an item by its ID
func (im *ItemManager) GetItem(id string) *Item {
	im.mu.RLock()
	defer im.mu.RUnlock()
	return im.items[id]
}

// LoadItemsFromYAML loads items from a YAML file
func (im *ItemManager) LoadItemsFromYAML(filePath string) error {
	var items []Item
	if err := LoadYAML(filePath, &items); err != nil {
		return err
	}

	im.mu.Lock()
	defer im.mu.Unlock()
	for _, item := range items {
		im.items[item.ID] = &item
	}
	return nil
}
