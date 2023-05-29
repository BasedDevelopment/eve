package auto

import (
	"net/http"
	"net/url"
)

func (a *Auto) WsReq(w http.ResponseWriter, r *http.Request, domid string) {
	wsUrl, err := url.Parse(a.Url)
	if err != nil {
		return
	}
	wsUrl.Path = "/libvirt/domains/" + domid + "/console"
	a.WSProxy(wsUrl, w, r)
}
