package projectIBC

import (
	"bytes"
	sha256 "crypto/sha256"
	"fmt"
	"reflect"
)

type Transaction struct {
	To     string
	From   string
	Bcoins float64
}

func (a Transaction) print() {
	fmt.Printf(a.To + " received ")
	fmt.Print(a.Bcoins)
	fmt.Println(" Bcoins from  " + a.From)

}

type Block struct {
	Hash          []byte
	Transactions  []Transaction
	PrevPointer   *Block
	PrevBlockHash []byte
}

func (a *Block) DeriveHash() {
	/*buf := &bytes.Buffer{}
	err := binary.Write(buf, binary.BigEndian, &a)
	if err != nil {
		panic(err)
	}
	bytess := buf.Bytes()
	hash := sha256.Sum256(bytess)*/
	h := sha256.New()
	s := fmt.Sprintf("%v", &a.Transactions)
	sum := h.Sum([]byte(s))
	//fmt.Printf("%s hashes to %x", s, sum)
	a.Hash = sum[:]
}

func (a Block) AddTransaction(trans Transaction) {
	a.Transactions = append(a.Transactions, trans)
}

func createBlock(data []Transaction, prevBlock *Block, prevhash []byte) *Block {
	block := &Block{[]byte{}, data, prevBlock, prevhash}
	block.DeriveHash()
	return block
}

func InsertBlock(transaction []Transaction, chainHead *Block) *Block {
	var newHead *Block
	if chainHead == nil {
		block := createBlock(transaction, nil, nil)
		newHead = block
	} else {
		block := createBlock(transaction, chainHead, chainHead.Hash)
		newHead = block
	}

	return newHead
}

func ListBlocks(chainHead *Block) {
	fmt.Println("Current Block Chain is")
	var temp *Block
	temp = chainHead
	for chainHead.PrevPointer != nil {
		i := 0
		for i < len(chainHead.Transactions) {
			chainHead.Transactions[i].print()
			i = i + 1
		}
		fmt.Println("Hash  ", chainHead.Hash)
		fmt.Println("previous Hash  ", chainHead.PrevBlockHash)
		chainHead = chainHead.PrevPointer
		fmt.Println("   ")
	}
	i := 0
	for i < len(chainHead.Transactions) {
		chainHead.Transactions[i].print()
		i = i + 1
	}
	fmt.Println("Hash   ", chainHead.Hash)
	fmt.Println("previous Hash  ", chainHead.PrevBlockHash)
	fmt.Println("   ")
	chainHead = temp
}

func ChangeBlock(oldTrans []Transaction, newTrans []Transaction, chainHead *Block) {
	var temp *Block
	temp = chainHead
	for chainHead != nil {
		if reflect.DeepEqual(chainHead.Transactions, oldTrans) {
			chainHead.Transactions = newTrans
			break
		}
		chainHead = chainHead.PrevPointer
	}

	chainHead = temp
}

func VerifyChain(chainHead *Block) {
	var modified bool
	modified = false
	var temp *Block
	temp = chainHead
	for chainHead.PrevPointer != nil {
		chainHead.PrevPointer.DeriveHash()
		newhash := chainHead.PrevPointer.Hash
		oldhash := chainHead.PrevBlockHash
		result := bytes.Compare(newhash, oldhash)
		if result != 0 {
			modified = true
			fmt.Println("Block " + chainHead.PrevPointer.Transactions[0].To + "  has been changed")
		}
		chainHead = chainHead.PrevPointer
	}

	if modified == false {
		fmt.Println("Block chain is not modified")
	} else {
		fmt.Println("block chain has been modified")
	}
	chainHead = temp

}
