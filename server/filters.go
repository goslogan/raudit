package server

type AuthStatus uint

const (
	AUTHENTICATION_FAILED AuthStatus = iota
	AUTHENTICATION_FAILED_TOO_LONG
	AUTHENTICATION_NOT_REQUIRED
	AUTHENTICATION_DIRECTORY_PENDING
	AUTHENTICATION_DIRECTORY_ERROR
	AUTHENTICATION_SYNCER_IN_PROGRESS
	AUTHENTICATION_SYNCER_FAILED
	AUTHENTICATION_SYNCER_OK
	AUTHENTICATION_OK
	AUTH_STATUS_MIN = AUTHENTICATION_FAILED
	AUTH_STATUS_MAX = AUTHENTICATION_OK
)

func ExcludeNewConnection() Filter {
	return func(s *Server, m map[string]any) bool {
		_, ok := m["new_conn"]
		return !ok
	}
}

func ExcludeInternalConnection() Filter {
	return func(s *Server, m map[string]any) bool {
		_, ok := m["new_int_conn"]
		return !ok
	}
}

func ExcludeCloseConnection() Filter {
	return func(s *Server, m map[string]any) bool {
		_, ok := m["close_conn"]
		return !ok
	}
}

func ExcludeAuthConnection() Filter {
	return func(s *Server, m map[string]any) bool {
		val, ok := m["action"]
		return !ok || val != "auth"
	}
}

func ExcludeAuthConnectionStatus(exclude ...AuthStatus) Filter {
	return func(s *Server, m map[string]any) bool {

		i, ok := m["status"]
		if !ok {
			return true
		}

		if s, ok := i.(float64); !ok {
			return true
		} else {
			status := AuthStatus(s)
			for _, n := range exclude {
				if status == n {
					return false
				}
			}
			return true
		}

	}
}
