package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"strconv"
	"time"
	//"unsafe"
)

//交易数据
type Transaction struct {
	Sender string //发送者
	Recver string //接收者
	Amount string //数量
}

//区块链结构体
type BlockChain struct {
	Index     int64 //序列号
	Timestamp int64 //时间戳

	proof    []byte //工作量证明
	PrevHash []byte //前一个hash

	Datas []*Transaction //交易数据，可为多个

	Chain []*BlockChain //区块链列表
}

//初始化创世区块
func (self *BlockChain) Init() {
	self.NewBlock([]byte(""), []byte(""))
}

//添加区块方法
func (self *BlockChain) NewBlock(proof []byte, prevHash []byte) {
	self.Index++

	//fmt.Println(self.Index, time.Now().Unix(), proof, prevHash)

	block := &BlockChain{self.Index, time.Now().Unix(), proof, prevHash, self.Datas, self.Chain}

	self.Chain = append(self.Chain, new(BlockChain))
	self.Chain[len(self.Chain)-1] = block
}

//获取整个区块链大小
func (self *BlockChain) GetSize() int64 {
	return int64(len(self.Chain))
}

//计算当前节点hash值
func (self *BlockChain) Hash() []byte {
	timestamp := []byte(strconv.FormatInt(self.Timestamp, 10))
	headers := bytes.Join([][]byte{[]byte(self.PrevHash), timestamp, self.PrevHash}, []byte{})
	hash := sha256.Sum256(headers)

	//slice := hash[:]
	slice := hash[:4] //简化取4位
	return slice
}

//打印所有区块信息
func (self *BlockChain) ShowAll() {
	for i := 0; i < len(self.Chain); i++ {
		fmt.Printf("Index:%v, Timestamp:%v, proof:%X, PrevHash:%X, Datas:%v\n",
			self.Chain[i].Index, self.Chain[i].Timestamp,
			self.Chain[i].proof, self.Chain[i].PrevHash, self.Chain[i].Datas)
	}
}

func main() {
	p := new(BlockChain)
	p.Init()

	fmt.Println("BlockChain size:", p.GetSize())
	time.Sleep(time.Second * 1)
	p.NewBlock([]byte(""), p.Hash())
	fmt.Println("BlockChain size:", p.GetSize())

	p.ShowAll()

}
