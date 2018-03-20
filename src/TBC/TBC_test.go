//区块链相关
package TBC

import (
	"fmt"
	"testing"
	"time"
)

func TestInit(t *testing.T) {
	p := new(BlockChain)
	if p == nil {
		t.Error("Error new(BlockChain)")
	}

	p.Init()

	fmt.Println("BlockChain size:", p.GetSize())
	time.Sleep(time.Second * 1)
	p.NewBlock(1, p.GetLastBlock().Hash())
	fmt.Println("BlockChain size:", p.GetSize())

	p.ShowAll()
}
