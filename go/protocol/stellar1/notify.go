// Auto-generated by avdl-compiler v1.3.24 (https://github.com/keybase/node-avdl-compiler)
//   Input file: avdl/stellar1/notify.avdl

package stellar1

import (
	"github.com/keybase/go-framed-msgpack-rpc/rpc"
	context "golang.org/x/net/context"
)

type PaymentNotificationArg struct {
	AccountID AccountID `codec:"accountID" json:"accountID"`
	PaymentID PaymentID `codec:"paymentID" json:"paymentID"`
}

type NotifyInterface interface {
	PaymentNotification(context.Context, PaymentNotificationArg) error
}

func NotifyProtocol(i NotifyInterface) rpc.Protocol {
	return rpc.Protocol{
		Name: "stellar.1.notify",
		Methods: map[string]rpc.ServeHandlerDescription{
			"paymentNotification": {
				MakeArg: func() interface{} {
					ret := make([]PaymentNotificationArg, 1)
					return &ret
				},
				Handler: func(ctx context.Context, args interface{}) (ret interface{}, err error) {
					typedArgs, ok := args.(*[]PaymentNotificationArg)
					if !ok {
						err = rpc.NewTypeError((*[]PaymentNotificationArg)(nil), args)
						return
					}
					err = i.PaymentNotification(ctx, (*typedArgs)[0])
					return
				},
				MethodType: rpc.MethodNotify,
			},
		},
	}
}

type NotifyClient struct {
	Cli rpc.GenericClient
}

func (c NotifyClient) PaymentNotification(ctx context.Context, __arg PaymentNotificationArg) (err error) {
	err = c.Cli.Notify(ctx, "stellar.1.notify.paymentNotification", []interface{}{__arg})
	return
}
