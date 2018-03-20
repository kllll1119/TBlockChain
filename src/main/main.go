package main

import (
	"TBC"
	"fmt"
)

var NodeIdentifier string = "矿工"

////////////////////////////////////////////////////////////////////////////
//新交易
//@app.route('/transactions/new', methods=['POST'])
func new_transaction(self *TBC.BlockChain) int {
	//values = request.get_json()

	//Check that the required fields are in the POST'ed data
	//required = ['sender', 'recipient', 'amount']
	//if not all(k in values for k in required):
	//    return 'Missing values', 400

	//Create a new Transaction
	//index = blockchain.new_transaction(values['sender'], values['recipient'], values['amount'])

	//response = {'message': f'Transaction will be added to Block {index}'}
	//return jsonify(response), 201
	return self.NewTransaction("aaa", "bbb", "10", false)
}

//挖矿函数
func mine(self *TBC.BlockChain) string {
	//1、计算工作量证明
	last_block := self.GetLastBlock()
	proof := self.Proof_of_work()

	//2、判断是否有矿可挖...
	//mining_index := self.CanMining()
	//if mining_index > 0 {
	//	fmt.Println("没有可挖矿...")
	//	return "没有可挖矿..."
	//}

	//3、通过一笔交易授予矿工（我们）代币，以作为奖励；
	self.NewTransaction("系统", NodeIdentifier, "1", true)

	//4、创造新区块，并将其添至区块链；
	previous_hash := last_block.Hash()
	self.NewBlock(proof, previous_hash)

	//5、移除挖走的记录
	//self.RemoveTransaction(mining_index)

	return "挖矿成功!"
}

////////////////////////////////////
//主函数
func main() {
	g_blockchain := new(TBC.BlockChain)
	g_blockchain.Init()

	index := new_transaction(g_blockchain) //第一笔交易
	fmt.Println("NewTransaction:", index)

	//index = new_transaction(g_blockchain) //第二笔交易
	//fmt.Println("NewTransaction:", index)

	g_blockchain.ShowAll() //打印所有节点

	fmt.Println(mine(g_blockchain)) //挖矿
	g_blockchain.ShowAll()

	//index = new_transaction(g_blockchain) //第二笔交易
	//fmt.Println("NewTransaction:", index)

	fmt.Println(mine(g_blockchain)) //挖矿
	g_blockchain.ShowAll()
	/*
		p := new(BlockChain)
		p.Init()

		fmt.Println("BlockChain size:", p.GetSize())
		time.Sleep(time.Second * 1)
		p.NewBlock([]byte(""), p.GetLastBlock().Hash())
		fmt.Println("BlockChain size:", p.GetSize())

		p.ShowAll()
	*/
}
