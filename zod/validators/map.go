package validators

import "github.com/aymaneallaoui/zod-Go/zod"

type MapSchema struct {
	keySchema   zod.Schema
	valueSchema zod.Schema
	required    bool
}

func Map(keySchema, valueSchema zod.Schema) *MapSchema {
	return &MapSchema{keySchema: keySchema, valueSchema: valueSchema}
}

func (m *MapSchema) Required() *MapSchema {
	m.required = true
	return m
}

func (m *MapSchema) Validate(data interface{}) error {
	mapData, ok := data.(map[interface{}]interface{})
	if !ok {
		switch v := data.(type) {
		case map[string]interface{}:
			mapData = make(map[interface{}]interface{}, len(v))
			for k, val := range v {
				mapData[k] = val
			}
		default:
			return zod.NewValidationError("map", data, "invalid type, expected map")
		}
	}

	for key, value := range mapData {
		if err := m.keySchema.Validate(key); err != nil {
			return zod.NewValidationError("map key", key, err.Error())
		}
		if err := m.valueSchema.Validate(value); err != nil {
			return zod.NewValidationError("map value", value, err.Error())
		}
	}

	return nil
}
