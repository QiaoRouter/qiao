package ripng

import "errors"

var (
	ErrLength             = errors.New("RipngErrorCode::ERR_LENGTH")
	ErrPortNotRipng       = errors.New("RipngErrorCode::ERR_UDP_PORT_NOT_RIPNG")
	ErrVersion            = errors.New("RipngErrorCode::ERR_Version")
	ErrBadZero            = errors.New("RipngErrorCode::ERR_RIPNG_BAD_ZERO")
	ErrBadPrefixLen       = errors.New("RipngErrorCode::ERR_RIPNG_BAD_PREFIX_LEN")
	ErrInconsistentPrefix = errors.New("RipngErrorCode::ERR_RIPNG_INCONSISTENT_PREFIX_LENGTH")
)
