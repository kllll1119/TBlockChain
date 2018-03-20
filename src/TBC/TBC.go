//区块链相关
package TBC

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"strconv"
	"time"
)

//交易数据
type Transaction struct {
	Sender string //发送者
	Recver string //接收者
	Amount string //数量
	Flag   bool   //挖取标识
}

//区块链结构体
type BlockChain struct {
	Index     int            //序列号
	Timestamp int64          //时间戳
	Proof     int            //工作量证明
	PrevHash  []byte         //前一个hash
	Datas     []*Transaction //交易数据列表
	Chain     []*BlockChain  //区块链列表
}

//初始化创世区块
func (self *BlockChain) Init() {
	self.NewBlock(0, []byte(""))
}

//获取最后一个区块
func (self *BlockChain) GetLastBlock() *BlockChain {
	length := len(self.Chain)
	if length <= 0 {
		return nil
	}
	return self.Chain[length-1]
}

//添加一个新区块
func (self *BlockChain) NewBlock(proof int, prevHash []byte) {
	self.Index++
	block := &BlockChain{self.Index, time.Now().Unix(), proof, prevHash, self.Datas, self.Chain}

	self.Chain = append(self.Chain, new(BlockChain))
	self.Chain[self.GetSize()-1] = block
}

//添加一个新交易信息
func (self *BlockChain) NewTransaction(sender, recipient, amount string, flag bool) int {
	self.Datas = append(self.Datas, new(Transaction))
	data := &Transaction{sender, recipient, amount, flag}
	self.Datas[len(self.Datas)-1] = data
	return self.Index + 1
}

//获取整个区块链大小
func (self *BlockChain) GetSize() int {
	return len(self.Chain)
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

/*
//是否有节点挖矿
func (self *BlockChain) CanMining() int {
	for i := 0; i < len(self.Datas); i++ {
		if self.Datas[i].Flag == false {
			return i
		}
	}
	return -1
}


//移除挖走的交易数据
func (self *BlockChain) RemoveTransaction(index int) {
	fmt.Println("RemoveTransaction size:", len(self.Datas))

	var Datas []*Transaction
	for i := range self.Datas {
		if i != index {
			Datas = append(Datas, new(Transaction))
			Datas[len(Datas)-1] = self.Datas[i]
		}
	}
	self.Datas = Datas
	fmt.Println("RemoveTransaction size:", len(self.Datas))
	return
}
*/

//打印所有区块信息
func (self *BlockChain) ShowAll() {
	fmt.Printf("---------------- ShowAll ----------------\n")
	for i := 0; i < len(self.Chain); i++ {
		fmt.Printf("Index:%v, Timestamp:%v, proof:%v, PrevHash:%X,Datasize:%v\n\t交易信息：",
			self.Chain[i].Index, self.Chain[i].Timestamp,
			self.Chain[i].Proof, self.Chain[i].PrevHash, len(self.Chain[i].Datas))

		for j := 0; j < len(self.Chain[i].Datas); j++ {
			fmt.Printf("[%v]Sender:%s, Recver:%s, Amount:%s ,Flag:%v",
				j, self.Chain[i].Datas[j].Sender, self.Chain[i].Datas[j].Recver,
				self.Chain[i].Datas[j].Amount, self.Chain[i].Datas[j].Flag)
		}

		fmt.Printf("\n")
	}
}

//检验新区块合法性
func (self *BlockChain) valid_chain(chain BlockChain) bool {

	last_block := chain.Chain[0]
	current_index := 1

	for current_index < chain.GetSize() {
		block := chain.Chain[current_index]
		//print('{last_block}')
		//print('{block}')
		//print("\n-----------\n")
		//# Check that the hash of the block is correct
		if (string)(block.PrevHash) != (string)(self.GetLastBlock().Hash()) {
			return false
		}
		//# Check that the Proof of Work is correct
		if self.valid_proof(last_block.Proof, block.Proof, last_block.PrevHash) == false {
			return false
		}

		last_block = block
		current_index += 1
	}
	return true
}

//基本的工作量证明
func (self *BlockChain) Proof_of_work() int {

	// Find a number p' such that hash(pp') contains leading 4 zeroes
	//Where p is the previous proof, and p' is the new proof

	last_proof := self.GetLastBlock().Proof
	last_hash := self.GetLastBlock().Hash()

	proof := 0
	for self.valid_proof(last_proof, proof, last_hash) == false {
		proof += 1
	}
	return proof
}

//验证证明
func (self *BlockChain) valid_proof(last_proof int, proof int, last_hash []byte) bool {
	headers := bytes.Join([][]byte{IntToBytes(last_proof), IntToBytes(proof), last_hash}, []byte{})
	hash := sha256.Sum256(headers)

	//fmt.Println(hash)
	return (string)(hash[:1]) == "0"
}

//////////////// tool //////////////
//整形转换成字节
func IntToBytes(n int) []byte {
	tmp := int32(n)
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, tmp)
	return bytesBuffer.Bytes()
}

//字节转换成整形
func BytesToInt(b []byte) int {
	bytesBuffer := bytes.NewBuffer(b)
	var tmp int32
	binary.Read(bytesBuffer, binary.BigEndian, &tmp)
	return int(tmp)
}
