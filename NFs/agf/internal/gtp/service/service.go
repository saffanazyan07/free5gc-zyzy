package service

import (
	"context"
	"errors"
	"net"

	gtp "github.com/wmnsk/go-gtp/gtpv1"
	gtpMsg "github.com/wmnsk/go-gtp/gtpv1/message"

	"github.com/free5gc/agf/internal/gtp/handler"
	"github.com/free5gc/agf/internal/logger"
	agfContext "github.com/free5gc/agf/pkg/context"
)

var gtpContext context.Context = context.TODO()

// SetupGTPTunnelWithUPF set up GTP connection with UPF
// return *gtp.UPlaneConn, net.Addr and error
func SetupGTPTunnelWithUPF(upfIPAddr string) (*gtp.UPlaneConn, net.Addr, error) {
	agfSelf := agfContext.W-AGFSelf()

	// Set up GTP connection
	upfUDPAddr := upfIPAddr + gtp.GTPUPort

	remoteUDPAddr, err := net.ResolveUDPAddr("udp", upfUDPAddr)
	if err != nil {
		logger.GTPLog.Errorf("Resolve UDP address %s failed: %+v", upfUDPAddr, err)
		return nil, nil, errors.New("Resolve Address Failed")
	}

	agfUDPAddr := agfSelf.GTPBindAddress + gtp.GTPUPort

	localUDPAddr, err := net.ResolveUDPAddr("udp", agfUDPAddr)
	if err != nil {
		logger.GTPLog.Errorf("Resolve UDP address %s failed: %+v", agfUDPAddr, err)
		return nil, nil, errors.New("Resolve Address Failed")
	}

	// Dial to UPF
	userPlaneConnection, err := gtp.DialUPlane(gtpContext, localUDPAddr, remoteUDPAddr)
	if err != nil {
		logger.GTPLog.Errorf("Dial to UPF failed: %+v", err)
		return nil, nil, errors.New("Dial failed")
	}

	// Overwrite T-PDU handler for supporting extension header containing QoS parameters
	userPlaneConnection.AddHandler(gtpMsg.MsgTypeTPDU, handler.HandleQoSTPDU)

	return userPlaneConnection, remoteUDPAddr, nil
}
