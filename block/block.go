package block

import (
    "fmt"
    "time"
)

func NewBlock(payload string) Block {
    return Block {
        PrevHash: "",
        payload: payload,
        timestamp: time.Now(),
    }
}

func (b Block) Show() {
    fmt.Printf("%#v\n", b)
}

func (b Block) Hash() string {
    return fmt.Sprintf("%s:%s", b.payload, b.PrevHash)
}
