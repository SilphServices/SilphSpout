package model

import (
  //"io"
  "time"
  "fmt"
  //"crypto/md5"
  "github.com/SilphServices/SilphSpout/ttlcache"
)

type DedupeFilter struct {
  cache *ttlcache.Cache
}

func NewDedupeFilter() (DedupeFilter) {
  return DedupeFilter {
    cache: ttlcache.NewCache(time.Hour),
  }
}

func (f DedupeFilter) Filter (spawn Spawn) bool {
  hash := hashSpawn(spawn)
  return f.cache.GetAndSet(hash)
}

func hashSpawn(spawn Spawn) string {
  /*
  digest := md5.New()
  idString := strconv.Itoa(spawn.NameID)
  latString := strconv.FormatFloat(spawn.Latitude, 'f', -1, 64)
	lngString := strconv.FormatFloat(spawn.Longitude, 'f', -1, 64)

  io.WriteString(digest, idString)
  io.WriteString(digest, latString)
  io.WriteString(digest, lngString)
  hash := digest.Sum(nil)
  return string(hash[:])
  */
  return fmt.Sprintf("%d", spawn.EncounterID)
}
