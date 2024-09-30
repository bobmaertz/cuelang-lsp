package serve

import (
	"bufio"
	"log"
	"os"

	"github.com/bobmaertz/cuelang-lsp/pkg/lsp/handler"
	"github.com/bobmaertz/cuelang-lsp/pkg/lsp/rpc"
)

func Serve(l *log.Logger) {
	// state := analysis.NewState()

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(rpc.SplitFunc)

	for scanner.Scan() {
		method, contents, err := rpc.DecodeMessage(scanner.Bytes())
		if err != nil {
			l.Printf("error decoding message: %v\n", err)
			continue
		}
		handler.HandleMessage(l, nil, method, contents)
	}
}
