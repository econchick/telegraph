package block

import "time"

type Block struct {
    PrevHash string
    payload Payload
    timestamp time.Time
}

type Chain map[string]Block

type Payload = string

