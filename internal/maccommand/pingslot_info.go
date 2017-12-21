package maccommand

import (
	"fmt"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/brocaar/loraserver/internal/common"
	"github.com/brocaar/loraserver/internal/storage"
	"github.com/brocaar/lorawan"
)

func handlePingSlotInfoReq(ds *storage.DeviceSession, block Block) error {
	if len(block.MACCommands) != 1 {
		return fmt.Errorf("exactly one mac-command expected, got: %d", len(block.MACCommands))
	}

	pl, ok := block.MACCommands[0].Payload.(*lorawan.PingSlotInfoReqPayload)
	if !ok {
		return fmt.Errorf("expected *lorawan.PingSlotInfoReqPayload, got: %T", block.MACCommands[0].Payload)
	}

	ds.PingSlotPeriodicity = int(pl.Periodicity)

	log.WithFields(log.Fields{
		"dev_eui":     ds.DevEUI,
		"periodicity": pl.Periodicity,
	}).Info("ping_slot_info_req request received")

	err := AddQueueItem(common.RedisPool, ds.DevEUI, Block{
		CID: lorawan.PingSlotInfoAns,
		MACCommands: []lorawan.MACCommand{
			{CID: lorawan.PingSlotInfoAns},
		},
	})
	if err != nil {
		return errors.Wrap(err, "add mac-command block to queue error")
	}

	return nil
}
