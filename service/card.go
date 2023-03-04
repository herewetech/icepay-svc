/*
 * Copyright (C) HereweTech, Inc - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 */

/**
 * @file card.go
 * @package service
 * @author Dr.NP <np@herewe.tech>
 * @since 02/26/2023
 */

package service

type Card struct{}

func NewCard() *Card {
	s := new(Card)

	return s
}

/* {{{ [Methods] */
func (s *Card) Create(owner_id, owner_type, number string) error {
	return nil
}

/* }}} */

/*
 * Local variables:
 * tab-width: 4
 * c-basic-offset: 4
 * End:
 * vim600: sw=4 ts=4 fdm=marker
 * vim<600: sw=4 ts=4
 */
