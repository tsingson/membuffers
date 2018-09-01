package e2e

import (
	"testing"
	"reflect"
	"github.com/orbs-network/membuffers/go/e2e/types"
)

var transactionBuf = []byte{
	// Transaction
		// TransactionData
		0x58,0x00,0x00,0x00, // size
		0x33,0x00,0x00,0x00,
		0x08,0x07,0x06,0x05, 0x04,0x03,0x02,0x01,
		0x40,0x00,0x00,0x00,
			// TransactionSender
			0x1c,0x00,0x00,0x00, // size
			0x03,0x00,0x00,0x00, '1','2','3',0x00,
			0x10,0x00,0x00,0x00, 0x03,0x00,0x00,0x00, '4','5','6',0x00, 0x03,0x00,0x00,0x00, '7','8','9',0x00,
			// TransactionSender
			0x1c,0x00,0x00,0x00, // size
			0x03,0x00,0x00,0x00, 'a','b','c',0x00,
			0x10,0x00,0x00,0x00, 0x03,0x00,0x00,0x00, 'd','e','f',0x00, 0x03,0x00,0x00,0x00, 'g','h','i',0x00,
		0x18,0x17,0x16,0x15, 0x14,0x13,0x12,0x11,
	0x06,0x00,0x00,0x00, 0x11,0x22,0x33,0x44, 0x55,0x66,
	0x01,0x00,
}

func TestReadPrimitives(t *testing.T) {
	transaction := types.TransactionReader(transactionBuf)
	sig := transaction.Signature()
	if !reflect.DeepEqual(sig, []byte{0x11,0x22,0x33,0x44,0x55,0x66}) {
		t.Fatalf("Signature: instead of expected got %v", sig)
	}
	pv := transaction.Data().ProtocolVersion()
	if pv != 0x33 {
		t.Fatalf("ProtocolVersion: instead of expected got %v", pv)
	}
	vc := transaction.Data().VirtualChain()
	if vc != 0x0102030405060708 {
		t.Fatalf("VirtualChain: instead of expected got %v", vc)
	}
	ts := transaction.Data().TimeStamp()
	if ts != 0x1112131415161718 {
		t.Fatalf("TimeStamp: instead of expected got %v", ts)
	}
	tt := transaction.Type()
	if tt != types.NETWORK_TYPE_TEST_NET {
		t.Fatalf("Type: instead of expected got %v", tt)
	}
}

func TestReadArrays(t *testing.T) {
	transaction := types.TransactionReader(transactionBuf)
	names := []string{}
	friends := []string{}
	for i:= transaction.Data().SenderIterator(); i.HasNext(); {
		sender := i.NextSender()
		names = append(names, sender.Name())
		for j := sender.FriendIterator(); j.HasNext(); {
			friends = append(friends, j.NextFriend())
		}
	}
	if !reflect.DeepEqual(names, []string{"123","abc"}) {
		t.Fatalf("Names: instead of expected got %v", names)
	}
	if !reflect.DeepEqual(friends, []string{"456","789","def","ghi"}) {
		t.Fatalf("Friends: instead of expected got %v", friends)
	}
}

func TestReadRawBuffers(t *testing.T) {
	transaction := types.TransactionReader(transactionBuf)
	transaction_rb := transaction.Raw()
	if !reflect.DeepEqual(transaction_rb, transactionBuf) {
		t.Fatalf("transaction_rb: instead of expected got %v", transaction_rb)
	}
	data_rb := transaction.RawData()
	if !reflect.DeepEqual(data_rb, []byte{0x33,0x00,0x00,0x00,0x08,0x07,0x06,0x05,0x04,0x03,0x02,0x01,0x40,0x00,0x00,0x00,0x1c,0x00,0x00,0x00,0x03,0x00,0x00,0x00, '1','2','3',0x00,0x10,0x00,0x00,0x00,0x03,0x00,0x00,0x00,'4','5','6',0x00, 0x03,0x00,0x00,0x00,'7','8','9',0x00,0x1c,0x00,0x00,0x00,0x03,0x00,0x00,0x00,'a','b','c',0x00, 0x10,0x00,0x00,0x00,0x03,0x00,0x00,0x00, 'd','e','f',0x00, 0x03,0x00,0x00,0x00, 'g','h','i',0x00,0x18,0x17,0x16,0x15, 0x14,0x13,0x12,0x11,}) {
		t.Fatalf("data_rb: instead of expected got %v", data_rb)
	}
	senders_rb := transaction.Data().RawSenderArray()
	if !reflect.DeepEqual(senders_rb, []byte{0x1c,0x00,0x00,0x00,0x03,0x00,0x00,0x00,'1','2','3',0x00,0x10,0x00,0x00,0x00,0x03,0x00,0x00,0x00,'4','5','6',0x00,0x03,0x00,0x00,0x00,'7','8','9',0x00,0x1c,0x00,0x00,0x00,0x03,0x00,0x00,0x00,'a','b','c',0x00,0x10,0x00,0x00,0x00,0x03,0x00,0x00,0x00,'d','e','f',0x00,0x03,0x00,0x00,0x00,'g','h','i',0x00,}) {
		t.Fatalf("senders_rb: instead of expected got %v", senders_rb)
	}
	vc_rb := transaction.Data().RawVirtualChain()
	if !reflect.DeepEqual(vc_rb, []byte{0x08,0x07,0x06,0x05,0x04,0x03,0x02,0x01,}) {
		t.Fatalf("vc_rb: instead of expected got %v", vc_rb)
	}
}

var methodBuf  = []byte{
	// Method
	0x03,0x00,0x00,0x00, 'm','e','t',0x00,
	0x1c,0x00,0x00,0x00,
		// MethodCallArgument
		0x08,0x00,0x00,0x00, // size
		0x00,0x00,0x00,0x00, 0x04,0x03,0x02,0x01,
		// MethodCallArgument
		0x0c,0x00,0x00,0x00, // size
		0x01,0x00,0x00,0x00, 0x03,0x00,0x00,0x00, 'a','r','g', 0x00,
}

func TestReadUnion(t *testing.T) {
	method := types.MethodReader(methodBuf)
	iteration := 0
	for i := method.ArgIterator(); i.HasNext(); {
		arg := i.NextArg()
		if iteration == 0 {
			if !arg.IsTypeNum() {
				t.Fatalf("first union is not a num although it is")
			}
			if arg.IsTypeData() {
				t.Fatalf("first union is a data although it is not")
			}
			if arg.Type() != types.METHOD_CALL_ARGUMENT_TYPE_NUM {
				t.Fatalf("first union is incorrect enum value")
			}
			num := arg.Num()
			if num != 0x01020304 {
				t.Fatalf("first union instead of expected got %v", num)
			}
		}
		if iteration == 1 {
			if !arg.IsTypeStr() {
				t.Fatalf("second union is not a str although it is")
			}
			if arg.IsTypeData() {
				t.Fatalf("second union is a data although it is not")
			}
			if arg.Type() != types.METHOD_CALL_ARGUMENT_TYPE_STR {
				t.Fatalf("second union is incorrect enum value")
			}
			str := arg.Str()
			if str != "arg" {
				t.Fatalf("second union instead of expected got %v", str)
			}
		}
		iteration++
	}
}

func TestReadRawBuffersWithHeader(t *testing.T) {
	transaction := types.TransactionReader(transactionBuf)
	data_rb := transaction.RawDataWithHeader()
	if !reflect.DeepEqual(data_rb, []byte{0x58,0x00,0x00,0x00, 0x33,0x00,0x00,0x00,0x08,0x07,0x06,0x05,0x04,0x03,0x02,0x01,0x40,0x00,0x00,0x00,0x1c,0x00,0x00,0x00,0x03,0x00,0x00,0x00, '1','2','3',0x00,0x10,0x00,0x00,0x00,0x03,0x00,0x00,0x00,'4','5','6',0x00, 0x03,0x00,0x00,0x00,'7','8','9',0x00,0x1c,0x00,0x00,0x00,0x03,0x00,0x00,0x00,'a','b','c',0x00, 0x10,0x00,0x00,0x00,0x03,0x00,0x00,0x00, 'd','e','f',0x00, 0x03,0x00,0x00,0x00, 'g','h','i',0x00,0x18,0x17,0x16,0x15, 0x14,0x13,0x12,0x11,}) {
		t.Fatalf("data_rb: instead of expected got %v", data_rb)
	}
	senders_rb := transaction.Data().RawSenderArrayWithHeader()
	if !reflect.DeepEqual(senders_rb, []byte{0x40,0x00,0x00,0x00, 0x1c,0x00,0x00,0x00,0x03,0x00,0x00,0x00,'1','2','3',0x00,0x10,0x00,0x00,0x00,0x03,0x00,0x00,0x00,'4','5','6',0x00,0x03,0x00,0x00,0x00,'7','8','9',0x00,0x1c,0x00,0x00,0x00,0x03,0x00,0x00,0x00,'a','b','c',0x00,0x10,0x00,0x00,0x00,0x03,0x00,0x00,0x00,'d','e','f',0x00,0x03,0x00,0x00,0x00,'g','h','i',0x00,}) {
		t.Fatalf("senders_rb: instead of expected got %v", senders_rb)
	}
	method := types.MethodReader(methodBuf)
	args_rb := method.RawArgArrayWithHeader()
	if !reflect.DeepEqual(args_rb, []byte{0x1c,0x00,0x00,0x00, 0x08,0x00,0x00,0x00, 0x00,0x00,0x00,0x00, 0x04,0x03,0x02,0x01, 0x0c,0x00,0x00,0x00, 0x01,0x00,0x00,0x00, 0x03,0x00,0x00,0x00, 'a','r','g', 0x00,}) {
		t.Fatalf("args_rb: instead of expected got %v", args_rb)
	}
}