package app

import "github.com/rei0721/go-scaffold/pkg/dbtx"

func (app *App) initDBTx() error {
	tx, err := dbtx.NewManager(app.DB.DB(), app.Logger)
	if err != nil {
		return err
	}
	app.DBTx = tx
	return nil
}
