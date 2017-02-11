package model

import (
	"encoding/json"
	"os"
	"errors"
)

type IVFilter struct {
  defaultIV int
  idToIV map[int]int
}

func (f IVFilter) Filter (spawn Spawn) bool {
  minIV, ok := f.idToIV[spawn.NameID]
  if !ok {
    minIV = f.defaultIV
  }
  return spawn.IVPercent() < minIV
}

func LoadFilter(pathToFilterJSON string, nameProvider NameProvider) (filter IVFilter, err error) {
  var rawIVFilter struct {
		DefaultMinIV int
		MinIV map[string]int
	}

  fd, err := os.Open(pathToFilterJSON)
	decoder := json.NewDecoder(fd)
	err = decoder.Decode(&rawIVFilter)

	filter.defaultIV = rawIVFilter.DefaultMinIV
	filter.idToIV = make(map[int]int)
	for name, iv := range rawIVFilter.MinIV {
		id := nameProvider.GetNameID(name)
		if id == 0 {
			err = errors.New("Unrecognized Pokemon Name " + name)
			return
		}
		filter.idToIV[id] = iv
	}

  return
}
