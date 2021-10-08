package register

type Driver interface {
	Open(string) error
	Close()
	IsValid() bool

	Create(string, Record) error
	Update(string, Record) error
	Search(string) (map[string]Record, error)
	Delete(string) error
}
