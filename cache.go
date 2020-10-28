package deep_cache

type (
	ItemKey      interface{} // Ключ элемента (item_id)
	Key          interface{} // Значение категории (category_id, flow_id и т.д.)
	Value        interface{}
	CategoryKey  interface{} // Ключ категории ("category_id", "flow_id" и т.д.)
	CategoryData struct {
		key   CategoryKey
		value Key
	}
	node struct {
		parent    *node
		key       Key
		value     Value
		isDeleted bool
		child     map[Key]*node
	}
	Cache struct {
		tree       *node
		categories []CategoryKey
		items      map[ItemKey]*node
	}
)

func NewCache(categories []CategoryKey) *Cache {
	root := &node{
		child: make(map[Key]*node),
	}

	return &Cache{
		tree:       root,
		categories: categories,
		items:      make(map[ItemKey]*node),
	}
}

func (c *Cache) Add(itemKey ItemKey, value Value, categoriesInfo []CategoryData) {
	// Add item to the tree
	prev := c.tree
	for _, cat := range categoriesInfo {
		v, ok := prev.child[cat.value]
		if !ok {
			v = &node{
				parent: prev,
				key:    cat.key,
				value:  cat.value,
				child:  make(map[Key]*node),
			}
			prev.child[cat.value] = v
		}
		prev = v
	}

	n := &node{
		parent: prev,
		key:    itemKey,
		value:  value,
		child:  nil,
	}

	prev.child[itemKey] = n

	// Add item to the items
	c.items[itemKey] = n
}

func (c *Cache) GetByItemKey(key ItemKey) (Value, bool) {
	v, ok := c.items[key]
	if !ok || v.isDeleted {
		return nil, false
	}
	return v.value, true
}

func (c *Cache) GetByCategories(categoriesInfo []CategoryData) []Value {
	prev := c.tree
	for _, cat := range categoriesInfo {
		v, ok := prev.child[cat.value]
		if !ok {
			return []Value{}
		}

		prev = v
	}

	return getItemsFromNode(prev)
}

func (c *Cache) Remove(key ItemKey) {
	v, ok := c.items[key]
	if !ok {
		// TODO: Should return error?
		return
	}
	v.isDeleted = true
}

func getItemsFromNode(n *node) []Value {
	if n.value == nil {
		if n.isDeleted {
			return []Value{}
		}
		return []Value{n.value}
	}

	values := make([]Value, 0, len(n.child))
	for _, v := range n.child {
		values = append(values, getItemsFromNode(v)...)
	}
	return values
}
