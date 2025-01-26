package game

import (
	"log/slog"

	"github.com/Workiva/go-datastructures/queue"
	"github.com/spf13/viper"
)

var (
	iq = queue.New(viper.GetInt64("input_queue_capacity"))
)

func Enqueue(input string) {
	if err := iq.Put(input); err != nil {
		slog.Error("Failed to enqueue input",
			slog.String("input", input),
			slog.Any("error", err))
		return
	}
}

func ProcessQueue(fn func(string)) {
	go func() {
		for {
			item, err := iq.Get(1)
			if err != nil {
				slog.Error("Failed to dequeue input",
					slog.Any("error", err))
				continue
			}
			if len(item) > 0 {
				fn(item[0].(string))
			}
		}
	}()
}
