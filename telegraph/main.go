// Copyright (c) 2017 Nick Platt, Lynn Root
package main

import (
    "crypto/md5"
    "encoding/hex"
    "encoding/json"
    "flag"
    "fmt"
    "io/ioutil"
    "os"
    "strconv"
    "strings"
    "time"
)

// Create a new type for a list of Strings (block files)
type stringList []string

// flag.Value interface for createCommand.Var
func (s *stringList) String() string {
    return fmt.Sprintf("%v", *s)
}

func (s *stringList) Set(value string) error {
    // *s = strings.Split(value, ",")
    *s = append(*s, value)
    return nil
}

func main() {
    // subcommands
    showCommand := flag.NewFlagSet("show", flag.ExitOnError)
    createCommand := flag.NewFlagSet("create", flag.ExitOnError)

    // showCommand flags
    jsonFile := showCommand.String("name", "", "JSON filename to show. (Required)")
    showWhat := showCommand.String("what", "", "choices {block|chain} (Required)")

    // createCommand flags
    createWhat := createCommand.String("what", "", "choices {block|chain} (Required)")
    payload := createCommand.String("payload", "", "payload text (Required for --what block)")
    prevBlock := createCommand.String("prev", "", "path to previous block, if any")

    var blockList stringList
    createCommand.Var(&blockList, "blocks", "A comma seperated list of paths to JSON blocks for --what chain.")

    command := os.Args[1]

    switch command {
    case "show":
        showCommand.Parse(os.Args[2:])
    case "create":
        createCommand.Parse(os.Args[2:])
    case "demo":
        demo()
    default:
        fmt.Printf("%q is not valid command.\n", os.Args[1])
        os.Exit(1)
    }

    // Handle commands
    if showCommand.Parsed() {
        // Required Flags
        if *jsonFile == "" {
            showCommand.PrintDefaults()
            os.Exit(1)
        }
        if *showWhat == "" {
            showCommand.PrintDefaults()
            os.Exit(1)
        }
        switch *showWhat {
        case "block":
            ShowBlock(*jsonFile)
        // TODO: case chain
        case "default":
            fmt.Printf("%q is not valid subcommand.\n", *showWhat)
            os.Exit(1)
        }
    }
    if createCommand.Parsed() {
        if *createWhat == "" {
            createCommand.PrintDefaults()
            os.Exit(1)
        }
        switch *createWhat {
        case "block":
            if *payload == "" {
                createCommand.PrintDefaults()
                os.Exit(1)
            }
            CreateBlock(*payload, *prevBlock)
        case "chain":
            if (&blockList).String() == "[]" {
                createCommand.PrintDefaults()
                os.Exit(1)
            }
            CreateChain(&blockList)
        case "default":
            fmt.Printf("%q is not valid subcommand.\n", *showWhat)
            os.Exit(1)
        }
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

// ShowBlock ...?
func ShowBlock(path string) {
    b := Load(path)
    b.PrettyPrint()
}

// CreateBlock ...?
func CreateBlock(payload string, prevBlock string) {
    var b *Block
    b = NewBlock(payload)

    if prevBlock != "" {
        p := Load(prevBlock)
        b.PrevHash = p.Hash()
    }
    hash := b.Hash()
    b.Dump(hash)
}

// CreateChain ...?
func CreateChain(blocks *stringList) {
    // TODO: this is b0rked - CLI order should not matter
    c := NewLinkedBlockChain()

    prevHash := ""
    for _, blockFile := range *blocks {
        b := Load(blockFile)
        b.PrevHash = prevHash
        curHash := b.Hash()
        c.AddBlock(b)
        prevHash = curHash
    }
    c.PrettyPrint()
}
