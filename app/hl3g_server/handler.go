package main

import (
	"github.com/rayjay214/parser/protocol/hl3g"
	"github.com/rayjay214/parser/server_base/hl3g_base"
	"github.com/rayjay214/parser/storage"
	log "github.com/sirupsen/logrus"
	"time"
)

func handleLK2(session *hl3g_base.Session, message *hl3g.Message) {
	entity := message.Body.(*hl3g.HL3G_LK2)
	log.Infof("%v:handle 01 %v, %v", session.ID(), message, entity)

	set, err := storage.SetStartTime(session.ID())
	if err == nil && set {
		storage.UpdateStartTime(session.ID())
	}

	info := map[string]interface{}{
		"state":     "4",
		"comm_time": time.Now(),
	}
	_ = storage.SetRunInfo(session.ID(), info)

	session.CommonReply(message.Header.Imei, message.Header.Proto)
}
