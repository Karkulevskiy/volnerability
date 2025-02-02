package db


type Storage struct {
	// TODO выбрать окончательно какая будет бд
}

func New(storagePath string) (*Storage, error) {
	//TODO добавить логику инициализации БД

	return  &Storage{}, nil
}

func (s *Storage) Login() error {
	// TODO implement me
	return nil
}

func (s *Storage) Register() error {
	// TODO implement me
	return nil
}