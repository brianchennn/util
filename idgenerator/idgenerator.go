// idgenerator is used for generating ID from minValue to maxValue.
// It will allocate IDs in range [minValue, maxValue]
// It is thread-safe when allocating IDs
package idgenerator

import (
	"errors"
	"sync"
)

type IDGenerator struct {
	Lock       sync.Mutex
	MinValue   int64
	MaxValue   int64
	ValueRange int64
	Offset     int64
	UsedMap    map[int64]bool
}

// Initialize an IDGenerator with minValue and maxValue.
func NewGenerator(minValue, maxValue int64) *IDGenerator {
	idGenerator := &IDGenerator{}
	idGenerator.init(minValue, maxValue)
	return idGenerator
}

func (idGenerator *IDGenerator) init(minValue, maxValue int64) {
	idGenerator.MinValue = minValue
	idGenerator.MaxValue = maxValue
	idGenerator.ValueRange = maxValue - minValue + 1
	idGenerator.Offset = 0
	idGenerator.UsedMap = make(map[int64]bool)
}

// Allocate and return an id in range [minValue, maxValue]
func (idGenerator *IDGenerator) Allocate() (int64, error) {
	idGenerator.Lock.Lock()
	defer idGenerator.Lock.Unlock()

	offsetBegin := idGenerator.Offset
	for {
		if _, ok := idGenerator.UsedMap[idGenerator.Offset]; ok {
			idGenerator.updateOffset()

			if idGenerator.Offset == offsetBegin {
				return 0, errors.New("No available value range to allocate id")
			}
		} else {
			break
		}
	}
	idGenerator.UsedMap[idGenerator.Offset] = true
	id := idGenerator.Offset + idGenerator.MinValue
	idGenerator.updateOffset()
	return id, nil
}

// param:
//   - id: id to free
func (idGenerator *IDGenerator) FreeID(id int64) {
	if id < idGenerator.MinValue || id > idGenerator.MaxValue {
		return
	}
	idGenerator.Lock.Lock()
	defer idGenerator.Lock.Unlock()
	delete(idGenerator.UsedMap, id-idGenerator.MinValue)
}

func (idGenerator *IDGenerator) updateOffset() {
	idGenerator.Offset++
	idGenerator.Offset = idGenerator.Offset % idGenerator.ValueRange
}
