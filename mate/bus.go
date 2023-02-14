package mate

func NewBus() *Bus {
	return &Bus{
		data: make(map[string]interface{}),
	}
}

type Bus struct {
	data map[string]interface{}
}

func (b *Bus) Set(key string, value interface{}) {
	b.data[key] = value
}

func (b *Bus) Get(key string) interface{} {
	return b.data[key]
}

func (b *Bus) Clone() map[string]interface{} {
	res := make(map[string]interface{})
	for k, v := range b.data {
		res[k] = v
	}
	return res
}
