org &8000

start:
;;------------------------------------------------
;; install a interrupt handler

    di						;; disable interrupts
    im 1					;; interrupt mode 1 (CPU will jump to &0038 when a interrupt occrs)
    ld hl,&c9fb				;; C9 FB are the bytes for the Z80 opcodes EI:RET
    ld (&0038),hl			;; setup interrupt handler
    ei

    ld bc, &7f10
    out (c), c

loop:

    ld b,&f5
vsync
    in a,(c)
    rra
    jr nc,vsync

    halt

    ld bc, &7f10
    out (c), c
    ld bc, &7f56
    out (c),c

    halt

    ld bc, &7f10
    out (c), c
    ld bc, &7f5c
    out (c),c

    halt

    ld bc, &7f56
    out (c),c

    jp loop

    end start
