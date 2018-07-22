package task

import "bitsync/services"

type PriceWatchTask struct {
}

// 启动价格监控
func (t *PriceWatchTask) Watch() {
	huobi := services.HuobiService{}
	go huobi.WatchHuobi()

	dragonex := services.DragonexService{}
	go dragonex.WatchDragonex()

	okex := services.OkexService{}
	go okex.WatchOkex()

	gate := services.GateService{}
	go gate.WatchGate()

	bithumb := services.BithumbService{}
	go bithumb.WatchBithumb()
}
