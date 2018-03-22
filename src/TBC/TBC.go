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

	//重置交易数据列表
	self.Datas = self.Datas[0:0:0]

	self.Chain = append(self.Chain, new(BlockChain))
	self.Chain[self.GetSize()-1] = block
}

//添加一个新交易信息
func (self *BlockChain) NewTransaction(sender, recipient, amount string) int {
	self.Datas = append(self.Datas, new(Transaction))
	data := &Transaction{sender, recipient, amount}
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

//是否有节点挖矿
func (self *BlockChain) CanMining() bool {
	return len(self.Datas) > 0
}

//打印所有区块信息
func (self *BlockChain) ShowAll() {
	fmt.Printf("---------------- ShowAll ----------------\n")
	for i := 0; i < len(self.Chain); i++ {
		fmt.Printf("Index:%v, Timestamp:%v, proof:%v, PrevHash:%X,Datasize:%v\n\t交易信息：",
			self.Chain[i].Index, self.Chain[i].Timestamp,
			self.Chain[i].Proof, self.Chain[i].PrevHash, len(self.Chain[i].Datas))

		for j := 0; j < len(self.Chain[i].Datas); j++ {
			fmt.Printf("[%v]Sender:%s, Recver:%s, Amount:%s ",
				j, self.Chain[i].Datas[j].Sender, self.Chain[i].Datas[j].Recver,
				self.Chain[i].Datas[j].Amount)
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
		if (string)(block.PrevHash) != (string)(self.GetLastBlock().Hash()) {
			return false
		}
		//校验工作量合法性
		if self.valid_proof(last_block.Proof, block.Proof, last_block.PrevHash) == false {
			return false
		}

		last_block = block
		current_index += 1
	}
	return true
}

//共识算法
func (self *BlockChain) ResolveConflicts() bool {
	//遍历所有我们的相邻节点，下载它们的链，并使用上述方法来验证它们。

	//:return: True if our chain was replaced, False if not
	/*
		neighbours = self.nodes
		var new_chain BlockChain

		//We're only looking for chains longer than ours
		max_length := self.GetSize()

		   //Grab and verify the chains from all the nodes in our network
		   for node in neighbours:
		       response = requests.get(f'http://{node}/chain')

		       if response.status_code == 200:
		           length = response.json()['length']
		           chain = response.json()['chain']

		           //Check if the length is longer and the chain is valid
		           if length > max_length and self.valid_chain(chain):
		               max_length = length
		               new_chain = chain

		   //Replace our chain if we discovered a new, valid chain longer than ours
		   if new_chain:
		       self.chain = new_chain
		       return true
	*/
	return false
}

//基本的工作量
func (self *BlockChain) Proof_of_work() int {
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
