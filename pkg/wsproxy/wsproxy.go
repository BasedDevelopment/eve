/*
Modified from https://github.com/gobwas/ws-examples/blob/master/src/proxy/proxy.go
MIT License
*/
package wsproxy

import (
	"io"
	"net"
	"net/http"

	eUtil "github.com/BasedDevelopment/eve/pkg/util"
	"github.com/rs/zerolog/log"
)

func WsProxy(w http.ResponseWriter, r *http.Request, peer net.Conn) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.Write(peer); err != nil {
			eUtil.WriteError(w, r, err, http.StatusBadGateway, "Failed to write to peer")
			return
		}
		hj, ok := w.(http.Hijacker)
		if !ok {
			eUtil.WriteError(w, r, nil, http.StatusInternalServerError, "Failed to hijack connection")
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

	})
}
