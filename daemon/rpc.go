package main

import (
	"../core"
	"../events"
	"context"
	"fmt"
	"net"
	"net/rpc"
	"time"
)

type RPCService struct {
	daemon *Daemon
}

func RunRPCService(ctx context.Context, d *Daemon, netType string, addr string) error {
	l, err := net.Listen(netType, addr)
	if err != nil {
		return err
	}
	s := RPCService{d}

	rc := rpc.NewServer()
	err = rc.RegisterName(core.RPCName, s)
	if err != nil {
		return err
	}
	go func() {
		for {
			select {
			case <-ctx.Done():
				l.Close()
				return
			case <-time.After(time.Second):
			}
		}
	}()
	rc.Accept(l)
	return nil
}

func (s RPCService) ListConfig(_ bool, out *[]string) error {
	var ret []string
	ret = append(ret, "ID\tHitCount")
	for _, cfg := range s.daemon.cfgs {
		v := fmt.Sprintf("%s\t%d", cfg.Id, s.daemon.history.Usage(cfg.Id))
		ret = append(ret, v)
	}
	*out = ret
	return nil
}

func (s RPCService) SwitchUserSession(env map[string]string, out *bool) error {
	return s.daemon.SwitchUserSession(env)
}

func (s RPCService) Schedule(_ string, out *bool) error {
	err := s.daemon.Schedule(s.daemon.ctx)
	if err != nil {
		return err
	}
	return nil
}
func (s RPCService) SchedulePendings(_ string, out *map[string][]string) error {
	var ret = make(map[string][]string)
	for _, s := range events.Scopes() {
		p := events.Pendings(s)
		if len(p) != 0 {
			ret[s] = p
		}
	}
	*out = ret
	return nil
}
