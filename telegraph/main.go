package main

import (
    // "flag"
    "crypto/md5"
    "encoding/hex"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "os"
    "strconv"
    "strings"
    "time"
)

func main() {
    // showCommand := flag.NewFlagSet("show", flag.ExitOnError)

    what := os.Args[1]
    // command := os.Args[2]

    switch what {
    // case "block":
    //     switch command {
    //     case "show":
    //         block.Show()
    //     case "create":
    //         block.Create()
    //     default:
    //         fmt.Printf("%q is not valid command.\n", os.Args[1])
    //         os.Exit(2)
    //     }
    // case "chain":
    //     switch command {
    //     case "show":
    //         chain.Show()
    //     case "create":
    //         chain.Create()
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

// Block ...?
type Block struct {
    PrevHash  string  `json:"prevHash"`
    Payload   Payload `json:"payload"`
    Timestamp string  `json:"timestamp"`
    next      *Block  `json:"next"`
}

// Payload ...?
type Payload = string

// Chain ...?
type Chain map[string]Block

// LinkedChain ...?
type LinkedChain struct {
    head *Block
    tail *Block
}

// NewBlockChain ...?
func NewBlockChain() Chain {
    return Chain{}
}

// NewLinkedBlockChain ...?
func NewLinkedBlockChain() LinkedChain {
    return LinkedChain{}
}

// NewBlock ...?
func NewBlock(payload string) *Block {
    t := time.Now()
    b := Block{
        PrevHash:  "",
        Payload:   payload,
        Timestamp: t.Format(time.RFC3339),
    }
    return &b
}

// Show ...?
func (c Chain) Show(start Block) {
    var b Block
    for b = start; b.PrevHash != ""; b = c[b.PrevHash] {
        b.Show()
    }
    b.Show()
}

// Show ...?
func (b Block) Show() {
    fmt.Printf("%#v\n", b)
}

// Hash ...?
func (b *Block) Hash() string {
    bytes := tojson(b)
    h := md5.New()
    h.Write([]byte(bytes))
    return hex.EncodeToString(h.Sum(nil))
}

// AddBlock ...?
func (c *LinkedChain) AddBlock(b Block) {
    block := &Block{
        next:      c.head,
        Payload:   b.Payload,
        PrevHash:  b.PrevHash,
        Timestamp: b.Timestamp,
    }
    if c.head != nil {
        c.head.PrevHash = block.PrevHash
    }
    c.head = block

    head := c.head
    for head.next != nil {
        head = head.next
    }
    c.tail = head
}

// Reverse ...?
func (c *LinkedChain) Reverse() {
    curr := c.head
    var prev *Block
    c.tail = c.head
    for curr != nil {
        next := curr.next
        curr.next = prev
        prev = curr
        curr = next
    }
    c.head = prev
}

func tojson(b interface{}) []byte {
    bytes, err := json.MarshalIndent(b, "", "    ")
    if err != nil {
        fmt.Println(err.Error())
        os.Exit(1)
    }
    return bytes
}

// Dump ...?
func (b Block) Dump(path string) {
    jsonData := tojson(b)

    file, err := os.Create(path)
    if err != nil {
        fmt.Println("Cannot create file", err)
        os.Exit(1)
    }
    defer file.Close()

    file.Write(jsonData)
    file.Close()
    fmt.Println("JSON data written to ", file.Name())
}

// Dumps ...?
func (b Block) Dumps() string {
    j := tojson(b)
    return string(j)
}

// Load ...?
func Load(path string) Block {
    fmt.Println("Loading JSON data from: ", path)
    raw, err := ioutil.ReadFile(path)

    if err != nil {
        fmt.Println(err.Error())
    }
    var b Block
    json.Unmarshal(raw, &b)
    return b
}

// PrettyPrint ...?
func (b Block) PrettyPrint() {
    j := string(tojson(b))
    fmt.Println(j)
}

// PrettyPrint ...?
func (c Chain) PrettyPrint(start Block) {
    var b Block
    for b = start; b.PrevHash != ""; b = c[b.PrevHash] {
        b.PrettyPrint()
    }
    b.PrettyPrint()
}

// PrettyPrint ...?
func (c LinkedChain) PrettyPrint() {
    block := c.head
    var nodeSlice []string
    for block != nil {
        nodeSlice = append(nodeSlice, block.Payload)
        block = block.next
    }
    output := strings.Join(nodeSlice, " -> ")
    fmt.Println(output)
}

func demo() {
    c := NewLinkedBlockChain()
    prevHash := ""
    var b *Block
    for i := 1; i <= 5; i++ {
        payload := strconv.Itoa(i)
        b = NewBlock(payload)
        b.PrevHash = prevHash
        curHash := b.Hash()
        c.AddBlock(*b)
        prevHash = curHash
        out := fmt.Sprintf("%d.json", i)
        b.Dump(out)
    }
    fmt.Println("Pretty printing the chain", c)
    c.PrettyPrint()

    fmt.Println("Flip that and reverse it")
    c.Reverse()
    c.PrettyPrint()
}
