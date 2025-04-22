package kite

import (
	"strings"

	"github.com/maxime-peim/gokite/pkg/bird/commands"
	"github.com/pkg/errors"
)

func (kite *BirdKite) Configure(file string, soft bool, timeout int) error {
	_, err := kite.bird.SendCommand(&commands.ConfigureCommand{
		File:    file,
		Soft:    soft,
		Timeout: timeout,
	})
	if err != nil {
		return errors.Wrap(err, "failed to send Configure command")
	}
	return nil
}

func (kite *BirdKite) ConfigureSoft(file string, timeout int) error {
	return kite.Configure(file, true, timeout)
}

func (kite *BirdKite) ConfigureHard(file string, timeout int) error {
	return kite.Configure(file, false, timeout)
}

func (kite *BirdKite) ConfigureCurrent(timeout int) error {
	return kite.Configure("", true, timeout)
}

func (kite *BirdKite) ConfigureConfirm() error {
	_, err := kite.bird.SendCommand(&commands.ConfigureConfirmCommand{})
	if err != nil {
		return errors.Wrap(err, "failed to send ConfigureConfirm command")
	}
	return nil
}

func (kite *BirdKite) ConfigureUndo() error {
	_, err := kite.bird.SendCommand(&commands.ConfigureUndoCommand{})
	if err != nil {
		return errors.Wrap(err, "failed to send ConfigureUndo command")
	}
	return nil
}

func (kite *BirdKite) ConfigureCheck(file string) error {
	reply, err := kite.bird.SendCommand(&commands.ConfigureCheckCommand{
		File: file,
	})
	if err != nil {
		return errors.Wrap(err, "failed to send ConfigureCheck command")
	} else if !strings.Contains(reply.String(), "OK") {
		return errors.Errorf("unexpected reply: %s", reply.String())
	}
	return nil
}

func (kite *BirdKite) ConfigureCheckCurrent() error {
	return kite.ConfigureCheck("")
}
