package async

type (
	ResChan chan any
	ErrChan chan error

	Anyfunc func() (any, error)
)

func Await(rc ResChan, ec ErrChan) (any, error) {
	var (
		res any
		err error
	)
	select {
	case res = <-rc:
		return res, nil
	case err = <-ec:
		return nil, err
	}
}

func Async(f Anyfunc) (ResChan, ErrChan) {
	var (
		reschan ResChan = make(ResChan)
		errchan ErrChan = make(ErrChan)
	)
	go func() {
		defer close(reschan)
		defer close(errchan)

		res, err := f()
		if err != nil {
			errchan <- err
		}

		reschan <- res
	}()

	return reschan, errchan
}
