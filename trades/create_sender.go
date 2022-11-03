package trades

import (
	"context"
	"fmt"
	"github.com/phoobynet/trade-ripper/configuration"
	qdb "github.com/questdb/go-questdb-client"
	"github.com/sirupsen/logrus"
)

func CreateSender(ctx context.Context, options configuration.Options) *qdb.LineSender {
	questDBAddress := fmt.Sprintf("%s:%d", options.DBHost, options.DBInfluxPort)
	logrus.Infof("Connecting to %s", questDBAddress)

	sender, err := qdb.NewLineSender(ctx, qdb.WithAddress(questDBAddress))

	if err != nil {
		logrus.Fatal("Error creating QuestDB line sender: ", err)
	}

	logrus.Infof("Connected to %s", questDBAddress)

	return sender
}
