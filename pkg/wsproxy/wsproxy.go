/*
Modified from https://github.com/gobwas/ws-examples/blob/master/src/proxy/proxy.go
MIT License
*/
package wsproxy

import (
	"io"
	"net"
	"net/http"

	"github.com/rs/zerolog/log"
)

func WsProxy(w http.ResponseWriter, r *http.Request, wsUrl string) http.Handler {
	peer, err := net.Dial("tcp", wsUrl)
	if err != nil {
		eUtil.WriteError(w, r, err, http.StatusInternalServerError, "Failed to dial TCP connection to ws")
		return
	}

	if err := r.Write(peer); err != nil {
		eUtil.WriteError(w, r, err, http.StatusBadGateway, "Failed to write to peer")
		return
	}
	hj, ok := w.(http.Hijacker)
	if !ok {
		eUtil.WriteError(w, r, err, http.StatusInternalServerError, "Failed to hijack connection")
		return
	}
	conn, _, err := hj.Hijack()
	if err != nil {
		eUtil.WriteError(w, r, err, http.StatusInternalServerError, "Failed to hijack connection")
		return
	}

	log.Info().
		Str("remote", peer.RemoteAddr().String()).
		Str("local", peer.LocalAddr().String()).
		Str("remote", conn.RemoteAddr().String()).
		Str("local", conn.LocalAddr().String()).
		Msg("proxying websocket")

	go func() {
		defer peer.Close()
		defer conn.Close()
		io.Copy(peer, conn)
	}()
	go func() {
		defer peer.Close()
		defer conn.Close()
		io.Copy(conn, peer)
	}()

}
