/*
 * Copyright (C) HereweTech, Inc - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 */

/**
 * @file http.go
 * @package utils
 * @author Dr.NP <np@herewe.tech>
 * @since 02/25/2022
 */

package utils

import "time"

type Envelope struct {
	Code      int         `json:"code"`
	Status    int         `json:"status"`
	Timestamp time.Time   `json:"timestamp"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
}

func WrapResponse(data interface{}) *Envelope {
	e := &Envelope{
		Timestamp: time.Now(),
		Data:      data,
	}

	return e
}

func (e *Envelope) SetStatus(status int) *Envelope {
	e.Status = status

	return e
}

func (e *Envelope) SetMessage(msg string) *Envelope {
	e.Message = msg

	return e
}

/*
 * Local variables:
 * tab-width: 4
 * c-basic-offset: 4
 * End:
 * vim600: sw=4 ts=4 fdm=marker
 * vim<600: sw=4 ts=4
 */
