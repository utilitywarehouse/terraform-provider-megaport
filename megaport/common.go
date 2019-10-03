package megaport

import (
	"log"
	"time"

	"github.com/google/uuid"
)

var (
	uuidSpace uuid.UUID
)

func init() {
	var err error
	uuidSpace, err = uuid.NewRandom()
	if err != nil {
		log.Fatal(err)
	}
}

func newUUID(data string) string {
	t, _ := time.Now().MarshalBinary() // Should be safe to ignore the error here
	return uuid.NewSHA1(uuidSpace, append(t, data...)).String()
}
