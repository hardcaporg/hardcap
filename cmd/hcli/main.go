package main

import (
    "github.com/hardcaporg/hardcap/internal/model"
    "github.com/hardcaporg/hardcap/internal/rpc/client"
    "github.com/hardcaporg/hardcap/internal/rpc/common"
    "log"
    "net/rpc"
    "os"

    "github.com/k0kubun/pp/v3"
    "github.com/urfave/cli/v2"
)

var rpcClient *rpc.Client

func conn(ctx *cli.Context) *rpc.Client {
	if rpcClient != nil {
		return rpcClient
	}

	var err error
	host := ctx.String("host")
	rpcClient, err = rpc.DialHTTP("tcp", host)
	if err != nil {
		log.Fatal("cannot connect: ", err)
	}

	return rpcClient
}

func main() {
	var host string

	app := &cli.App{
		EnableBashCompletion: true,

		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "host",
				Value:       ":8007",
				Usage:       "RPC host",
				Destination: &host,
			},
		},

		Commands: []*cli.Command{
			{
				Name:    "appliances",
				Aliases: []string{"app"},
				Usage:   "appliances",
				Subcommands: []*cli.Command{
					{
						Name:  "register",
						Usage: "register hardware or virtual appliance",
						Flags: []cli.Flag{
                            &cli.StringFlag{Name: "name", Required: true},
                            &cli.StringFlag{Name: "url", Required: true},
                            },
						Action: func(ctx *cli.Context) error {
							reg := client.ApplianceClient(conn(ctx))

							reply := &common.ApplianceRegisterReply{}
							app := &model.Appliance{
                                Name: ctx.String("name"),
                                URL: ctx.String("url"),
                            }
							err := reg.Register(common.ApplianceRegisterArgs{Appliance: app}, reply)
							if err != nil {
								return err
							}

							pp.Println(reply.Appliance)
							return nil
						},
					},
					{
						Name:  "list",
						Usage: "list appliances",
						Action: func(ctx *cli.Context) error {
							reg := client.ApplianceClient(conn(ctx))

							reply := &common.ApplianceListReply{}
							err := reg.List(common.ApplianceListArgs{}, reply)
							if err != nil {
								return err
							}

							pp.Println(reply.List)
							return nil
						},
					},
				},
			},
			{
				Name:    "registrations",
				Aliases: []string{"reg"},
				Usage:   "host registrations",
				Subcommands: []*cli.Command{
					{
						Name:  "list",
						Usage: "list registrations",
						Action: func(ctx *cli.Context) error {
							reg := client.RegistrationClient(conn(ctx))

							reply := &common.RegistrationListReply{}
							err := reg.List(common.RegistrationListArgs{}, reply)
							if err != nil {
								return err
							}

							pp.Println(reply.List)
							return nil
						},
					},
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
