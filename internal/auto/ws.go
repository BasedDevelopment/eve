package auto

import (
	"net/http"

	eUtil "github.com/BasedDevelopment/eve/pkg/util"
	"github.com/BasedDevelopment/eve/pkg/wsproxy"
)

func (a *Auto) WsReq(w http.ResponseWriter, r *http.Request, domid string) http.Handler {
	wsurl := a.Url + "/libvirt/domains/" + domid + "/console"
	conn, err := a.getWSConn(wsurl)
	if err != nil {
		eUtil.WriteError(w, r, err, http.StatusInternalServerError, "Internal Server Error")
	}
	return wsproxy.WsProxy(w, r, conn)
}
