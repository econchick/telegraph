package block

// add/append to chain

func NewChain() Chain {
    return make(map[string]Block)
}

func (c Chain) Show(start Block) {
    var b Block
    for b = start; b.PrevHash != ""; b = c[b.PrevHash] {
        b.Show()
    }
    b.Show()
}
