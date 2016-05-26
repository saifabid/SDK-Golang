package recast

//An Entity represents a Recast entity and provides getter to access the entity fields
type Entity struct {
	data map[string]interface{}
	name string
}

// Field returns the value of the field gien by name, or nil if the field is not present
// The value returned is an interface{} and should be cast as follows:
// numbers: float64
// strings : string
// Refer to Recast.Ai manual for details about the entities
func (e *Entity) Field(name string) interface{} {
	for key, value := range e.data {
		if name == key {
			return value
		}
	}
	return nil
}

// Name returns the name of the entity
func (e *Entity) Name() string {
	return e.name
}

// Raw returns the raw value of the entity, as it was in the original text
func (e *Entity) Raw() string {
	return e.data["raw"].(string)
}
