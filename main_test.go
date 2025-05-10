package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"testing"

	"golang.org/x/sync/errgroup"
)

func TestRun(t *testing.T) {
	l, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("failed to listen port: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		return run(ctx, l)
	})

	// テストリクエストを送信
	url := fmt.Sprintf("http://%s/%s", l.Addr().String(), "test")
	t.Logf("try request to %s", url)

	resp, err := http.Get(url)
	if err != nil {
		t.Fatalf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// レスポンスのステータスコードを確認
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	actual, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}

	// レスポンスの内容を確認
	expected := "Hello, test!"
	if string(actual) != expected {
		t.Errorf("expected response body %q, got %q", expected, actual)
	}

	// コンテキストをキャンセルしてサーバーを停止
	cancel()

	// run関数の終了を待機
	if err := eg.Wait(); err != nil {
		t.Errorf("run returned error: %v", err)
	}
}
