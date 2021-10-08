package errno

import "fmt"

var (
	UnknownError = fmt.Errorf("Unknow error")

	// 注册表驱动不可用
	RegisterDriverNotValid = fmt.Errorf("Register driver not valid")
)
