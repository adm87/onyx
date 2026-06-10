package assert

func Fatal(err error) {
	if err != nil {
		panic(err)
	}
}

func True(condition bool, message string) {
	if !condition {
		panic(message)
	}
}

func False(condition bool, message string) {
	if condition {
		panic(message)
	}
}

func NotNil(value any, message string) {
	if value == nil {
		panic(message)
	}
}

func Nil(value any, message string) {
	if value != nil {
		panic(message)
	}
}
