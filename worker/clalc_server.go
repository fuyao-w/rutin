package worker

import (
	"fmt"
	"time"
)

type Handle struct{}

func (h *Handle) Name() string {
	return "calc"
}

func (h *Handle) Calc(req ClacReq, resp *ClacResp) error {
	fmt.Println("Calc deal")
	time.Sleep(500 * time.Millisecond)
	resp.Result = req.A<<2 + req.B<<1
	return nil
}
