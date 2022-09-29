package cloud

type JianmuAuthorize struct {
	Authorization
}

func (auth *JianmuAuthorize) Authorize() ConnData {

	return ConnData{}
}
