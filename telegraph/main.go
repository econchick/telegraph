package main

import (
    // "flag"
    "fmt"
    "os"
    "time"
    "strconv"

    // "github.com/econchick/telegraph/commands/block"
    // "github.com/econchick/telegraph/commands/chain"
    // "github.com/econchick/telegraph/block"
)

type Block struct {
    PrevHash string
    payload Payload
    timestamp time.Time
}

type Payload = string

type Chain map[string]Block

func main() {
    // showCommand := flag.NewFlagSet("show", flag.ExitOnError)

    switch os.Args[1] {
    // case "block":
    //     switch os.Args[2] {
    //     case "show":
    //         commands.block.Show()
    //     case "create":
    //         commands.block.Create()
    //     default:
    //         fmt.Printf("%q is not valid command.\n", os.Args[1])
    //         os.Exit(2)
    //     }
    // case "chain":
    //     switch os.Args[2] {
    //     case "show":
    //         commands.chain.Show()
    //     case "create":
    //         commands.chain.Create()
    //     default:
    //         fmt.Printf("%q is not valid command.\n", os.Args[1])
    //         os.Exit(2)
    //     }
    case "demo":
        demo()
    default:
        fmt.Printf("%q is not valid command.\n", os.Args[1])
        os.Exit(2)
    }
}

func NewBlockChain() Chain {
    return Chain{}
}

func NewBlock(payload string) *Block {
    b := Block {
        PrevHash: "",
        payload: payload,
        timestamp: time.Now(),
    }
    return &b
}

func (c Chain) Show(start Block) {
    var b Block
    for b = start; b.PrevHash != ""; b = c[b.PrevHash] {
        b.Show()
    }
    b.Show()
}

func (b *Block) Show() {
    fmt.Printf("%#v\n", b)
}

func (b *Block) Hash() string {
    return fmt.Sprintf("%s:%s", b.payload, b.PrevHash)
}

func demo() {
    c := NewBlockChain()
    prevHash := ""
    var b *Block
    for i := 1; i <= 3; i++ {
        payload := strconv.Itoa(i)
        b = NewBlock(payload)
        b.PrevHash = prevHash
        curHash := b.Hash()
        c[curHash] = *b
        prevHash = curHash
    }
    c.Show(*b)
}
