package e2e

import (
	"testing"
	"bytes"
	"github.com/orbs-network/membuffers/go/membufc/e2e/protos"
	"crypto/md5"
	"github.com/orbs-network/membuffers/go/membufc/e2e/protos/dep2"
	"github.com/orbs-network/membuffers/go/membufc/e2e/protos/dep1/dep11"
)

func hashBytes(buffer []byte) []byte {
	hash := md5.Sum(buffer)
	return hash[:]
}

func TestReadWriteTransaction(t *testing.T) {
	// write
	builder := &types.TransactionBuilder{
		Data: &types.TransactionDataBuilder{
			ProtocolVersion: 0x01,
			VirtualChain: 0x11223344,
			Sender: []*types.TransactionSenderBuilder{
				&types.TransactionSenderBuilder{
					Name: "johnny",
					Friend: []string{"billy","jeff","alex"},
				},
				&types.TransactionSenderBuilder{
					Name: "rachel",
					Friend: []string{"jessica","sara"},
				},
			},
			TimeStamp: 0x445566778899,
		},
		Signature: []byte{0x22,0x22,0x22,0x22,0x22,0x22,0x22,0x22,0x22,0x22,0x22,0x22,0x22,0x22,0x22,0x22,},
		Type: types.NetworkType_RESERVED,
	}
	buf := make([]byte, builder.CalcRequiredSize())
	err := builder.Write(buf)
	if err != nil {
		t.Fatalf("error while writing in builder")
	}

	// read
	transaction := types.TransactionReader(buf)
	if transaction.Data().ProtocolVersion() != 0x01 {
		t.Fatalf("ProtocolVersion: instead of expected got %v", transaction.Data().ProtocolVersion())
	}
	if transaction.Data().VirtualChain() != 0x11223344 {
		t.Fatalf("VirtualChain: instead of expected got %v", transaction.Data().VirtualChain())
	}
	if transaction.Data().TimeStamp() != 0x445566778899 {
		t.Fatalf("TimeStamp: instead of expected got %v", transaction.Data().TimeStamp())
	}
	if !bytes.Equal(transaction.Signature(), []byte{0x22,0x22,0x22,0x22,0x22,0x22,0x22,0x22,0x22,0x22,0x22,0x22,0x22,0x22,0x22,0x22,}) {
		t.Fatalf("Signature: instead of expected got %v", transaction.Signature())
	}
	if transaction.Type() != types.NetworkType_RESERVED {
		t.Fatalf("Type: instead of expected got %v", transaction.Type())
	}
	senderCount := 0
	for i := transaction.Data().SenderIterator(); i.HasNext(); {
		sender := i.NextSender()
		if senderCount == 0 {
			if sender.Name() != "johnny" {
				t.Fatalf("Sender0.Name: instead of expected got %v", sender.Name())
			}
			j := sender.FriendIterator()
			if j.NextFriend() != "billy" {
				t.Fatalf("Sender0.Friend0: did not get as expected")
			}
			if j.NextFriend() != "jeff" {
				t.Fatalf("Sender0.Friend1: did not get as expected")
			}
			if j.NextFriend() != "alex" {
				t.Fatalf("Sender0.Friend2: did not get as expected")
			}
		}
		if senderCount == 1 {
			if sender.Name() != "rachel" {
				t.Fatalf("Sender1.Name: instead of expected got %v", sender.Name())
			}
			j := sender.FriendIterator()
			if j.NextFriend() != "jessica" {
				t.Fatalf("Sender1.Friend0: did not get as expected")
			}
			if j.NextFriend() != "sara" {
				t.Fatalf("Sender1.Friend1: did not get as expected")
			}
		}
		senderCount++
	}

	// mutate
	hash := hashBytes(transaction.RawData())
	err = transaction.MutateSignature(hash)
	if err != nil {
		t.Fatalf("Signature: mutate to hash value failed: %v", err.Error())
	}
	if bytes.Equal(transaction.Signature(), []byte{0x22,0x22,0x22,0x22,0x22,0x22,0x22,0x22,0x22,0x22,0x22,0x22,0x22,0x22,0x22,0x22,}) {
		t.Fatalf("Signature: value did not change after mutation")
	}
}

func TestReadWriteMethod(t *testing.T) {
	// write
	builder := &types.MethodBuilder{
		Name: "MyMethod",
		Arg: []*types.MethodCallArgumentBuilder{
			&types.MethodCallArgumentBuilder{
				Type: types.MethodCallArgumentTypeNum,
				Num:  0x17,
			},
			&types.MethodCallArgumentBuilder{
				Type: types.MethodCallArgumentTypeStr,
				Str:  "flower",
			},
			&types.MethodCallArgumentBuilder{
				Type: types.MethodCallArgumentTypeData,
				Data: []byte{0x01,0x02,0x03},
			},
		},
	}
	buf := make([]byte, builder.CalcRequiredSize())
	err := builder.Write(buf)
	if err != nil {
		t.Fatalf("error while writing in builder")
	}

	// read
	method := types.MethodReader(buf)
	if method.Name() != "MyMethod" {
		t.Fatalf("Name: instead of expected got %v", method.Name())
	}
	i := method.ArgIterator()
	arg0 := i.NextArg()
	if !arg0.IsTypeNum() {
		t.Fatalf("Arg0: is type is not num")
	}
	if arg0.IsTypeStr() {
		t.Fatalf("Arg0: is type is str")
	}
	if arg0.Type() != types.MethodCallArgumentTypeNum {
		t.Fatalf("Arg0: type is not num")
	}
	if arg0.TypeNum() != 0x17 {
		t.Fatalf("Arg0.Num: instead of expected got %v", arg0.TypeNum())
	}
	arg1 := i.NextArg()
	if !arg1.IsTypeStr() {
		t.Fatalf("Arg1: is type is not str")
	}
	if arg1.IsTypeNum() {
		t.Fatalf("Arg1: is type is num")
	}
	if arg1.Type() != types.MethodCallArgumentTypeStr {
		t.Fatalf("Arg1: type is not str")
	}
	if arg1.TypeStr() != "flower" {
		t.Fatalf("Arg1.Str: instead of expected got %v", arg1.TypeStr())
	}
}

func TestReadWriteWithImports(t *testing.T) {
	// write
	builder := &dep2.DependentBuilder{}
	buf := make([]byte, builder.CalcRequiredSize())
	err := builder.Write(buf)
	if err != nil {
		t.Fatalf("error while writing in builder")
	}

	// read
	dependent := dep2.DependentReader(buf)
	if dependent.A().Field() != 0 {
		t.Fatalf("A.Field: is not empty")
	}
	if dependent.B() != dep11.DependencyEnum_OPTION_A {
		t.Fatalf("B: is not DependencyEnum_OPTION_A")
	}
}

func TestQuickBuildTransaction(t *testing.T) {
	transaction := (&types.TransactionBuilder{
		Data: &types.TransactionDataBuilder{
			ProtocolVersion: 0x01,
			VirtualChain: 0x11223344,
			Sender: []*types.TransactionSenderBuilder{
				&types.TransactionSenderBuilder{
					Name: "johnny",
					Friend: []string{"billy","jeff","alex"},
				},
				&types.TransactionSenderBuilder{
					Name: "rachel",
					Friend: []string{"jessica","sara"},
				},
			},
			TimeStamp: 0x445566778899,
		},
		Signature: []byte{0x22,0x22,0x22,0x22,0x22,0x22,0x22,0x22,0x22,0x22,0x22,0x22,0x22,0x22,0x22,0x22,},
		Type: types.NetworkType_RESERVED,
	}).Build()

	if transaction.Data().ProtocolVersion() != 0x01 {
		t.Fatalf("ProtocolVersion: instead of expected got %v", transaction.Data().ProtocolVersion())
	}
	if transaction.Data().VirtualChain() != 0x11223344 {
		t.Fatalf("VirtualChain: instead of expected got %v", transaction.Data().VirtualChain())
	}
	if transaction.Data().TimeStamp() != 0x445566778899 {
		t.Fatalf("TimeStamp: instead of expected got %v", transaction.Data().TimeStamp())
	}
	if !bytes.Equal(transaction.Signature(), []byte{0x22,0x22,0x22,0x22,0x22,0x22,0x22,0x22,0x22,0x22,0x22,0x22,0x22,0x22,0x22,0x22,}) {
		t.Fatalf("Signature: instead of expected got %v", transaction.Signature())
	}
	if transaction.Type() != types.NetworkType_RESERVED {
		t.Fatalf("Type: instead of expected got %v", transaction.Type())
	}
}
