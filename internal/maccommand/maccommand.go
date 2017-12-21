package maccommand

import (
	"fmt"

	"github.com/brocaar/loraserver/internal/models"
	"github.com/brocaar/loraserver/internal/storage"
	"github.com/brocaar/lorawan"
)

// Handle handles a MACCommand sent by a node.
func Handle(ds *storage.DeviceSession, block Block, pending *Block, rxInfoSet models.RXInfoSet) error {
	var err error
	switch block.CID {
	case lorawan.LinkADRAns:
		err = handleLinkADRAns(ds, block, pending)
	case lorawan.LinkCheckReq:
		err = handleLinkCheckReq(ds, rxInfoSet)
	case lorawan.DevStatusAns:
		err = handleDevStatusAns(ds, block)
	case lorawan.PingSlotInfoReq:
		err = handlePingSlotInfoReq(ds, block)
	default:
		err = fmt.Errorf("undefined CID %d", block.CID)

	}
	return err
}
