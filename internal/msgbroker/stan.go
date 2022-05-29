package msgbroker

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"time"

	"github.com/nats-io/stan.go"

	"github.com/yalagtyarzh/L0/internal/config"
	"github.com/yalagtyarzh/L0/internal/constants"
	"github.com/yalagtyarzh/L0/internal/models"
	"github.com/yalagtyarzh/L0/internal/repocache"
	"github.com/yalagtyarzh/L0/internal/repository"
	"github.com/yalagtyarzh/L0/internal/utils"
	"github.com/yalagtyarzh/L0/pkg/logging"
)

// STAN stands for posting and getting data from NATS streaming server and filling cache with it
type STAN struct {
	config *config.STAN
	cache  *repocache.Cache
	logger *logging.Logger
	db     repository.DatabaseRepo
}

// NewSTAN returns new stan object
func NewSTAN(cfg *config.STAN, c *repocache.Cache, l *logging.Logger, db repository.DatabaseRepo) *STAN {
	return &STAN{
		config: cfg,
		cache:  c,
		logger: l,
		db:     db,
	}
}

// SendMessages emulates getting data from message broker
func (s *STAN) SendMessages() {
	go func() {
		sc, err := stan.Connect(s.config.Cluster, s.config.Client)
		if err != nil {
			s.logger.Fatal(fmt.Sprintf("error to connect stan: %s", err))
		}
		defer sc.Close()

		sub, err := s.subscribe(&sc)
		if err != nil {
			s.logger.Fatal(fmt.Sprintf("error to subscribe: %s", err))
		}
		defer sub.Unsubscribe()

		delay := time.Duration(s.config.Delay) * time.Second

		checked := make([]string, 0)

		for {
			time.Sleep(delay)
			err = s.checkMsg(&checked, &sc)
			if err != nil {
				s.logger.Fatal(fmt.Sprintf("error to check msg: %s", err))
			}
		}
	}()
}

// subscribe subscribes to STAN channel and loads order data to cache and db
func (s *STAN) subscribe(sc *stan.Conn) (stan.Subscription, error) {
	sub, err := (*sc).Subscribe(
		s.config.Channel, func(m *stan.Msg) {
			var order models.Order
			json.Unmarshal(m.Data, &order)
			_, ok := s.cache.Load(order.OrderUID)
			if ok {
				return
			}

			s.cache.Store(order.OrderUID, order)
			err := s.db.InsertOrder(order)
			if err != nil {
				s.logger.Error(fmt.Sprintf("error to insert order: %s", err))
			}
		},
	)
	if err != nil {
		return sub, err
	}

	return sub, nil
}

// checkMsg
func (s *STAN) checkMsg(checked *[]string, sc *stan.Conn) error {
	files, err := ioutil.ReadDir(constants.Dir)
	if err != nil {
		return err
	}

	for _, file := range files {
		fileName := file.Name()
		if utils.IsInSlice(fileName, *checked) {
			continue
		}
		*checked = append(*checked, fileName)

		if !utils.CheckExt(fileName, constants.ValidExt) {
			s.logger.Error(fmt.Sprintf("Invalid file extension of file: %s", fileName))
			continue
		}

		f, err := os.Open(fmt.Sprintf("%s/%s", constants.Dir, fileName))
		if err != nil {
			s.logger.Error(fmt.Sprintf("Cannot open file: %s", err))
			continue
		}

		data, err := io.ReadAll(f)
		if err != nil {
			s.logger.Error(fmt.Sprintf("Couldn't read file: %s", err))
			continue
		}

		f.Close()

		var order models.Order
		err = json.Unmarshal(data, &order)
		if err != nil {
			s.logger.Error(fmt.Sprintf("Invalid json: %s", err))
			continue
		}

		_, ok := s.cache.Load(order.OrderUID)

		if ok {
			continue
		}

		jsonData, err := json.Marshal(order)
		if err != nil {
			s.logger.Error(fmt.Sprintf("Couldn't marshal: %s\n", err))
			continue
		}

		err = (*sc).Publish(s.config.Channel, jsonData)
		if err != nil {
			s.logger.Error(fmt.Sprintf("Error to send data: %s\n", err))
			continue
		}

		s.logger.Info("Message sent")
	}

	return nil
}
