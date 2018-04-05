// Auto-generated by avdl-compiler v1.3.22 (https://github.com/keybase/node-avdl-compiler)
//   Input file: avdl/keybase1/tlf_keys.avdl

package keybase1

import (
	"github.com/keybase/go-framed-msgpack-rpc/rpc"
	context "golang.org/x/net/context"
)

type TLFIdentifyBehavior int

const (
	TLFIdentifyBehavior_UNSET           TLFIdentifyBehavior = 0
	TLFIdentifyBehavior_CHAT_CLI        TLFIdentifyBehavior = 1
	TLFIdentifyBehavior_CHAT_GUI        TLFIdentifyBehavior = 2
	TLFIdentifyBehavior_CHAT_GUI_STRICT TLFIdentifyBehavior = 3
	TLFIdentifyBehavior_KBFS_REKEY      TLFIdentifyBehavior = 4
	TLFIdentifyBehavior_KBFS_QR         TLFIdentifyBehavior = 5
	TLFIdentifyBehavior_CHAT_SKIP       TLFIdentifyBehavior = 6
	TLFIdentifyBehavior_SALTPACK        TLFIdentifyBehavior = 7
	TLFIdentifyBehavior_CLI             TLFIdentifyBehavior = 8
	TLFIdentifyBehavior_GUI             TLFIdentifyBehavior = 9
	TLFIdentifyBehavior_DEFAULT_KBFS    TLFIdentifyBehavior = 10
	TLFIdentifyBehavior_PAGES           TLFIdentifyBehavior = 11
)

func (o TLFIdentifyBehavior) DeepCopy() TLFIdentifyBehavior { return o }

var TLFIdentifyBehaviorMap = map[string]TLFIdentifyBehavior{
	"UNSET":           0,
	"CHAT_CLI":        1,
	"CHAT_GUI":        2,
	"CHAT_GUI_STRICT": 3,
	"KBFS_REKEY":      4,
	"KBFS_QR":         5,
	"CHAT_SKIP":       6,
	"SALTPACK":        7,
	"CLI":             8,
	"GUI":             9,
	"DEFAULT_KBFS":    10,
	"PAGES":           11,
}

var TLFIdentifyBehaviorRevMap = map[TLFIdentifyBehavior]string{
	0:  "UNSET",
	1:  "CHAT_CLI",
	2:  "CHAT_GUI",
	3:  "CHAT_GUI_STRICT",
	4:  "KBFS_REKEY",
	5:  "KBFS_QR",
	6:  "CHAT_SKIP",
	7:  "SALTPACK",
	8:  "CLI",
	9:  "GUI",
	10: "DEFAULT_KBFS",
	11: "PAGES",
}

func (e TLFIdentifyBehavior) String() string {
	if v, ok := TLFIdentifyBehaviorRevMap[e]; ok {
		return v
	}
	return ""
}

type CanonicalTlfName string

func (o CanonicalTlfName) DeepCopy() CanonicalTlfName {
	return o
}

type CryptKey struct {
	KeyGeneration int     `codec:"KeyGeneration" json:"KeyGeneration"`
	Key           Bytes32 `codec:"Key" json:"Key"`
}

func (o CryptKey) DeepCopy() CryptKey {
	return CryptKey{
		KeyGeneration: o.KeyGeneration,
		Key:           o.Key.DeepCopy(),
	}
}

type TLFBreak struct {
	Breaks []TLFIdentifyFailure `codec:"breaks" json:"breaks"`
}

func (o TLFBreak) DeepCopy() TLFBreak {
	return TLFBreak{
		Breaks: (func(x []TLFIdentifyFailure) []TLFIdentifyFailure {
			if x == nil {
				return nil
			}
			var ret []TLFIdentifyFailure
			for _, v := range x {
				vCopy := v.DeepCopy()
				ret = append(ret, vCopy)
			}
			return ret
		})(o.Breaks),
	}
}

type TLFIdentifyFailure struct {
	User   User                 `codec:"user" json:"user"`
	Breaks *IdentifyTrackBreaks `codec:"breaks,omitempty" json:"breaks,omitempty"`
}

func (o TLFIdentifyFailure) DeepCopy() TLFIdentifyFailure {
	return TLFIdentifyFailure{
		User: o.User.DeepCopy(),
		Breaks: (func(x *IdentifyTrackBreaks) *IdentifyTrackBreaks {
			if x == nil {
				return nil
			}
			tmp := (*x).DeepCopy()
			return &tmp
		})(o.Breaks),
	}
}

type CanonicalTLFNameAndIDWithBreaks struct {
	TlfID         TLFID            `codec:"tlfID" json:"tlfID"`
	CanonicalName CanonicalTlfName `codec:"CanonicalName" json:"CanonicalName"`
	Breaks        TLFBreak         `codec:"breaks" json:"breaks"`
}

func (o CanonicalTLFNameAndIDWithBreaks) DeepCopy() CanonicalTLFNameAndIDWithBreaks {
	return CanonicalTLFNameAndIDWithBreaks{
		TlfID:         o.TlfID.DeepCopy(),
		CanonicalName: o.CanonicalName.DeepCopy(),
		Breaks:        o.Breaks.DeepCopy(),
	}
}

type GetTLFCryptKeysRes struct {
	NameIDBreaks CanonicalTLFNameAndIDWithBreaks `codec:"nameIDBreaks" json:"nameIDBreaks"`
	CryptKeys    []CryptKey                      `codec:"CryptKeys" json:"CryptKeys"`
}

func (o GetTLFCryptKeysRes) DeepCopy() GetTLFCryptKeysRes {
	return GetTLFCryptKeysRes{
		NameIDBreaks: o.NameIDBreaks.DeepCopy(),
		CryptKeys: (func(x []CryptKey) []CryptKey {
			if x == nil {
				return nil
			}
			var ret []CryptKey
			for _, v := range x {
				vCopy := v.DeepCopy()
				ret = append(ret, vCopy)
			}
			return ret
		})(o.CryptKeys),
	}
}

type TLFQuery struct {
	TlfName          string              `codec:"tlfName" json:"tlfName"`
	IdentifyBehavior TLFIdentifyBehavior `codec:"identifyBehavior" json:"identifyBehavior"`
}

func (o TLFQuery) DeepCopy() TLFQuery {
	return TLFQuery{
		TlfName:          o.TlfName,
		IdentifyBehavior: o.IdentifyBehavior.DeepCopy(),
	}
}

type GetTLFCryptKeysArg struct {
	Query TLFQuery `codec:"query" json:"query"`
}

type GetPublicCanonicalTLFNameAndIDArg struct {
	Query TLFQuery `codec:"query" json:"query"`
}

type TlfKeysInterface interface {
	// getTLFCryptKeys returns TLF crypt keys from all generations and the TLF ID.
	// TLF ID should not be cached or stored persistently.
	GetTLFCryptKeys(context.Context, TLFQuery) (GetTLFCryptKeysRes, error)
	// getPublicCanonicalTLFNameAndID return the canonical name and TLFID for tlfName.
	// TLF ID should not be cached or stored persistently.
	GetPublicCanonicalTLFNameAndID(context.Context, TLFQuery) (CanonicalTLFNameAndIDWithBreaks, error)
}

func TlfKeysProtocol(i TlfKeysInterface) rpc.Protocol {
	return rpc.Protocol{
		Name: "keybase.1.tlfKeys",
		Methods: map[string]rpc.ServeHandlerDescription{
			"getTLFCryptKeys": {
				MakeArg: func() interface{} {
					ret := make([]GetTLFCryptKeysArg, 1)
					return &ret
				},
				Handler: func(ctx context.Context, args interface{}) (ret interface{}, err error) {
					typedArgs, ok := args.(*[]GetTLFCryptKeysArg)
					if !ok {
						err = rpc.NewTypeError((*[]GetTLFCryptKeysArg)(nil), args)
						return
					}
					ret, err = i.GetTLFCryptKeys(ctx, (*typedArgs)[0].Query)
					return
				},
				MethodType: rpc.MethodCall,
			},
			"getPublicCanonicalTLFNameAndID": {
				MakeArg: func() interface{} {
					ret := make([]GetPublicCanonicalTLFNameAndIDArg, 1)
					return &ret
				},
				Handler: func(ctx context.Context, args interface{}) (ret interface{}, err error) {
					typedArgs, ok := args.(*[]GetPublicCanonicalTLFNameAndIDArg)
					if !ok {
						err = rpc.NewTypeError((*[]GetPublicCanonicalTLFNameAndIDArg)(nil), args)
						return
					}
					ret, err = i.GetPublicCanonicalTLFNameAndID(ctx, (*typedArgs)[0].Query)
					return
				},
				MethodType: rpc.MethodCall,
			},
		},
	}
}

type TlfKeysClient struct {
	Cli rpc.GenericClient
}

// getTLFCryptKeys returns TLF crypt keys from all generations and the TLF ID.
// TLF ID should not be cached or stored persistently.
func (c TlfKeysClient) GetTLFCryptKeys(ctx context.Context, query TLFQuery) (res GetTLFCryptKeysRes, err error) {
	__arg := GetTLFCryptKeysArg{Query: query}
	err = c.Cli.Call(ctx, "keybase.1.tlfKeys.getTLFCryptKeys", []interface{}{__arg}, &res)
	return
}

// getPublicCanonicalTLFNameAndID return the canonical name and TLFID for tlfName.
// TLF ID should not be cached or stored persistently.
func (c TlfKeysClient) GetPublicCanonicalTLFNameAndID(ctx context.Context, query TLFQuery) (res CanonicalTLFNameAndIDWithBreaks, err error) {
	__arg := GetPublicCanonicalTLFNameAndIDArg{Query: query}
	err = c.Cli.Call(ctx, "keybase.1.tlfKeys.getPublicCanonicalTLFNameAndID", []interface{}{__arg}, &res)
	return
}
