package tls

import (
	"context"
	"errors"
	"testing"
	"time"
)

// TestUQUICCanceledWaitingForData 验证 UQUICConn 在等待对端握手数据时关闭会结束握手并返回 close_notify
func TestUQUICCanceledWaitingForData(t *testing.T) {
	config := testConfigClient.Clone()
	config.MinVersion = VersionTLS13
	config.NextProtos = []string{"h3"}
	client := UQUICClient(&QUICConfig{TLSConfig: config}, HelloGolang)
	client.SetTransportParameters(nil)
	if err := client.Start(context.Background()); err != nil {
		t.Fatal(err)
	}
	for client.NextEvent().Kind != QUICNoEvent {
	}
	closed := make(chan error, 1)
	go func() { closed <- client.Close() }()
	select {
	case err := <-closed:
		if !errors.Is(err, alertCloseNotify) {
			t.Errorf("Close() = %v，应为 alertCloseNotify", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Close() 5 秒内没有返回")
	}
}
