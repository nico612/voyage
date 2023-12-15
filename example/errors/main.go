package main

import (
	"fmt"
	"github.com/nico612/voyage/internal/pkg/code"
	"github.com/nico612/voyage/pkg/errors"
	"github.com/nico612/voyage/pkg/log"
)

func main() {

	if err := getUser(); err != nil {
		fmt.Printf("%+v\n", err)
	}
}

func getUser() error {
	if err := queryDatabase(); err != nil {
		//return errors.Wrap(err, "get sysuser failed.")
		return errors.WrapC(err, code.ErrEncodingFailed, "编码错误")
	}
	return nil
}

func queryDatabase() error {
	opts := &log.Options{
		OutputPaths:      []string{"test.log", "stdout"},
		ErrorOutputPaths: []string{},
		Level:            "info",
		Format:           "json",
		EnableColor:      true,
		Development:      true,
	}

	log.Init(opts)
	defer log.Flush()

	err := errors.WithCode(code.ErrDatabase, "sysuser xxx not found.")
	if err != nil {
		log.Errorf("%+v", err)
	}

	return err

}
