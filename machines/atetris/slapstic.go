package atetris

const (
	DISABLED = iota
	ENABLED
	ALTERNATE1
	ALTERNATE2
	ALTERNATE3
	BITWISE1
	BITWISE2
	BITWISE3
	// ADDITIVE1
	// ADDITIVE2
	// ADDITIVE3
)

type mask_value struct {
	mask, value uint16
}

func MATCHES_MASK_VALUE(val uint16, maskval mask_value) bool {
	return val&maskval.mask == maskval.value
}

type slapstic struct {
	rom []byte

	state int

	// bankstart    uint16
	current_bank uint16
	alt_bank     uint16
	bit_bank     uint16
	bit_xor      uint16
	bank         [4]uint16

	alt1     mask_value
	alt2     mask_value
	alt3     mask_value
	alt4     mask_value
	altshift uint16

	bit1   mask_value
	bit2c0 mask_value
	bit2s0 mask_value
	bit2c1 mask_value
	bit2s1 mask_value
	bit3   mask_value

	// add1     mask_value
	// add2     mask_value
	// addplus1 mask_value
	// addplus2 mask_value
	// add3     mask_value
}

func newSlapstic(rom []uint8) *slapstic {
	return &slapstic{
		rom:   rom,
		state: DISABLED,

		// basic banking
		current_bank: 3,                                         // starting bank
		bank:         [4]uint16{0x0080, 0x0090, 0x00a0, 0x00b0}, // bank select values

		// alternate banking
		alt1:     mask_value{0x1fff, 0x1dfe}, // 1st mask/value in sequence  - 0xffff == UNKNOWN
		alt2:     mask_value{0x1fff, 0x1dff}, // 2nd mask/value in sequence
		alt3:     mask_value{0x1ffc, 0x1b5c}, // 3rd mask/value in sequence
		alt4:     mask_value{0x1fcf, 0x0080}, // 4th mask/value in sequence
		altshift: 0,                          // shift to get bank from 3rd

		// bitwise banking
		bit1:   mask_value{0x1ff0, 0x1540}, // 1st mask/value in sequence
		bit2c0: mask_value{0x1ff3, 0x1540}, // clear bit 0 value
		bit2s0: mask_value{0x1ff3, 0x1541}, //   set bit 0 value
		bit2c1: mask_value{0x1ff3, 0x1542}, // clear bit 1 value
		bit2s1: mask_value{0x1ff3, 0x1543}, //   set bit 1 value
		bit3:   mask_value{0x1ff8, 0x1550}, // final mask/value in sequence

		// // additive banking
		// NO_ADDITIVE
	}
}

func (slapstic *slapstic) ReadPort(addr uint16) (byte, bool) {
	pagedAddr := uint16(slapstic.current_bank&1) * 0x4000
	pagedAddr |= uint16(addr & 0x3fff)

	if (addr & 0x3fff) >= 0x2000 {
		slapstic.tweak(addr & 0x1fff)
	}
	return slapstic.rom[pagedAddr], false
}

func (slapstic *slapstic) WritePort(addr uint16, data byte) { panic(-1) }

func (slapstic *slapstic) tweak(offset uint16) {
	/* reset is universal */
	if offset == 0x0000 {
		slapstic.state = ENABLED
	} else { /* otherwise, use the state machine */
		switch slapstic.state {
		/* DISABLED state: everything is ignored except a reset */
		case DISABLED:

		/* ENABLED state: the chip has been activated and is ready for a bankswitch */
		case ENABLED:
			/* check for request to enter bitwise state */
			if MATCHES_MASK_VALUE(offset, slapstic.bit1) {
				slapstic.state = BITWISE1
				// } else if MATCHES_MASK_VALUE(offset, slapstic.add1) { /* check for request to enter additive state */
				// slapstic.state = ADDITIVE1
			} else if MATCHES_MASK_VALUE(offset, slapstic.alt1) { /* check for request to enter alternate state */
				slapstic.state = ALTERNATE1
			} else if MATCHES_MASK_VALUE(offset, slapstic.alt2) { /* special kludge for catching the second alternate address if he first one was missed (since it's usually an opcode fetch) */
				slapstic.state = ALTERNATE2
			} else if offset == slapstic.bank[0] { /* check for standard bankswitches */
				slapstic.state = DISABLED
				slapstic.current_bank = 0
			} else if offset == slapstic.bank[1] {
				slapstic.state = DISABLED
				slapstic.current_bank = 1
			} else if offset == slapstic.bank[2] {
				slapstic.state = DISABLED
				slapstic.current_bank = 2
			} else if offset == slapstic.bank[3] {
				slapstic.state = DISABLED
				slapstic.current_bank = 3
			}

		/* ALTERNATE1 state: look for alternate2 offset, or else fall back to ENABLED */
		case ALTERNATE1:
			if MATCHES_MASK_VALUE(offset, slapstic.alt2) {
				slapstic.state = ALTERNATE2
			} else {
				slapstic.state = ENABLED
			}

		/* ALTERNATE2 state: look for altbank offset, or else fall back to ENABLED */
		case ALTERNATE2:
			if MATCHES_MASK_VALUE(offset, slapstic.alt3) {
				slapstic.state = ALTERNATE3
				slapstic.alt_bank = (offset >> slapstic.altshift) & 3
			} else {
				slapstic.state = ENABLED
			}

		/* ALTERNATE3 state: wait for the final value to finish the transaction */
		case ALTERNATE3:
			if MATCHES_MASK_VALUE(offset, slapstic.alt4) {
				slapstic.state = DISABLED
				slapstic.current_bank = slapstic.alt_bank
			}

		/* BITWISE1 state: waiting for a bank to enter the BITWISE state */
		case BITWISE1:
			if offset == slapstic.bank[0] || offset == slapstic.bank[1] ||
				offset == slapstic.bank[2] || offset == slapstic.bank[3] {
				slapstic.state = BITWISE2
				slapstic.bit_bank = slapstic.current_bank
				slapstic.bit_xor = 0
			}

		/* BITWISE2 state: watch for twiddling and the escape mechanism */
		case BITWISE2:

			/* check for clear bit 0 case */
			if MATCHES_MASK_VALUE(offset^slapstic.bit_xor, slapstic.bit2c0) {
				slapstic.bit_bank &= ^uint16(1)
				slapstic.bit_xor ^= 3
			} else if MATCHES_MASK_VALUE(offset^slapstic.bit_xor, slapstic.bit2s0) { /* check for set bit 0 case */
				slapstic.bit_bank |= 1
				slapstic.bit_xor ^= 3
			} else if MATCHES_MASK_VALUE(offset^slapstic.bit_xor, slapstic.bit2c1) { /* check for clear bit 1 case */
				slapstic.bit_bank &= ^uint16(2)
				slapstic.bit_xor ^= 3
			} else if MATCHES_MASK_VALUE(offset^slapstic.bit_xor, slapstic.bit2s1) { /* check for set bit 1 case */
				slapstic.bit_bank |= 2
				slapstic.bit_xor ^= 3
			} else if MATCHES_MASK_VALUE(offset, slapstic.bit3) { /* check for escape case */
				slapstic.state = BITWISE3
			}

		/* BITWISE3 state: waiting for a bank to seal the deal */
		case BITWISE3:
			if offset == slapstic.bank[0] || offset == slapstic.bank[1] ||
				offset == slapstic.bank[2] || offset == slapstic.bank[3] {
				slapstic.state = DISABLED
				slapstic.current_bank = slapstic.bit_bank
			}

			/* ADDITIVE1 state: look for add2 offset, or else fall back to ENABLED */
			// case ADDITIVE1:
			// 	if MATCHES_MASK_VALUE(offset, slapstic.add2) {
			// 		slapstic.state = ADDITIVE2
			// 		add_bank = slapstic.current_bank
			// 	} else {
			// 		slapstic.state = ENABLED
			// 	}

			// /* ADDITIVE2 state: watch for twiddling and the escape mechanism */
			// case ADDITIVE2:
			// 	/* check for add 1 case -- can intermix */
			// 	if MATCHES_MASK_VALUE(offset, slapstic.addplus1) {
			// 		add_bank = (add_bank + 1) & 3
			// 	}

			// 	/* check for add 2 case -- can intermix */
			// 	if MATCHES_MASK_VALUE(offset, slapstic.addplus2) {
			// 		add_bank = (add_bank + 2) & 3
			// 	}

			// 	/* check for escape case -- can intermix with the above */
			// 	if MATCHES_MASK_VALUE(offset, slapstic.add3) {
			// 		slapstic.state = ADDITIVE3
			// 	}

			// /* ADDITIVE3 state: waiting for a bank to seal the deal */
			// case ADDITIVE3:
			// 	if offset == slapstic.bank[0] || offset == slapstic.bank[1] ||
			// 		offset == slapstic.bank[2] || offset == slapstic.bank[3] {
			// 		slapstic.state = DISABLED
			// 		slapstic.current_bank = add_bank
			// 	}

		}
	}
}
