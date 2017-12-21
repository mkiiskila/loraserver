package maccommand

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/brocaar/loraserver/internal/common"
	"github.com/brocaar/loraserver/internal/models"
	"github.com/brocaar/loraserver/internal/storage"
	"github.com/brocaar/lorawan"
)

func handleLinkCheckReq(ds *storage.DeviceSession, rxInfoSet models.RXInfoSet) error {
	if len(rxInfoSet) == 0 {
		return errors.New("rx info-set contains zero items")
	}

	requiredSNR, ok := common.SpreadFactorToRequiredSNRTable[rxInfoSet[0].DataRate.SpreadFactor]
	if !ok {
		return fmt.Errorf("sf %d not in sf to required snr table", rxInfoSet[0].DataRate.SpreadFactor)
	}

	margin := rxInfoSet[0].LoRaSNR - requiredSNR
	if margin < 0 {
		margin = 0
	}

	block := Block{
		CID: lorawan.LinkCheckAns,
		MACCommands: MACCommands{
			{
				CID: lorawan.LinkCheckAns,
				Payload: &lorawan.LinkCheckAnsPayload{
					Margin: uint8(margin),
					GwCnt:  uint8(len(rxInfoSet)),
				},
			},
		},
	}

	if err := AddQueueItem(common.RedisPool, ds.DevEUI, block); err != nil {
		return errors.Wrap(err, "add mac-command block to queue error")
	}
	return nil
}
