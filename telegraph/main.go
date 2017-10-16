package main

import (
    // "flag"
    "fmt"
    "os"

    // "github.com/econchick/telegraph/commands/block"
    // "github.com/econchick/telegraph/commands/chain"
    "github.com/econchick/telegraph/block"
)

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

func demo() {
    c := block.NewChain()
    prevHash := ""
    var b block.Block
    for i := 1; i <= 3; i++ {
        b = block.NewBlock(fmt.Sprintf("%d", i))
        b.PrevHash = prevHash
        curHash := b.Hash()
        c[curHash] = b
        prevHash = curHash
    }
    c.Show(b)
}
