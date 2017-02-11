package model

import (
  "net/http"
  "encoding/json"
  "strconv"
)

type NameProvider interface {
  GetName(id int) string
  GetMove(id int) string
  GetNameID(name string) int
}

type RemoteNameProvider struct {
  nameURL string
  moveURL string
  nameToID map[int]string
  moveToID map[int]string
  idToName map[string]int
  idToMove map[string]int
}

func NewRemoteNameProvider(nameURL, moveURL string) (rnp RemoteNameProvider, err error) {
  rnp = RemoteNameProvider {
    nameURL: nameURL,
    moveURL: moveURL,
  }
  if err = rnp.LoadNames(); err != nil {
    return
  }
  if err = rnp.LoadMoves(); err != nil {
    return
  }
  return
}

func (r RemoteNameProvider) GetName(id int) string {
  return r.nameToID[id]
}

func (r RemoteNameProvider) GetMove(id int) string {
  return r.moveToID[id]
}

func (r RemoteNameProvider) GetNameID(name string) int {
  return r.idToName[name]
}


func (r *RemoteNameProvider) LoadNames() (err error) {
  r.nameToID, r.idToName, err = getAndUnmarshal(r.nameURL)
  return
}

func (r *RemoteNameProvider) LoadMoves() (err error) {
  r.moveToID, r.idToMove, err = getAndUnmarshal(r.moveURL)
  return
}

func getAndUnmarshal(url string) (idToValue map[int]string, valueToID map[string]int, err error) {
	// client := http.DefaultClient
	nameResp, err := http.Get(url)
	if err != nil {
		// handle error
	}
	defer nameResp.Body.Close()

	idToValue = make(map[int]string)
  valueToID = make(map[string]int)
	tmp := make(map[string]string)
	json.NewDecoder(nameResp.Body).Decode(&tmp)

	for k,v := range tmp {
		id, err := strconv.Atoi(k)
		if err != nil {
			break
		}
		idToValue[id] = v
    valueToID[v] = id
	}

	return
}
