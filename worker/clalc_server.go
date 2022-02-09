package worker

type Handle struct{}

func (h *Handle) Name() string {
	return "calc"
}

func (h *Handle) Calc(req ClacReq, resp *ClacResp) error {
	//fmt.Println("Calc deal")
	resp.Result = req.A<<2 + req.B<<1
	return nil
}
