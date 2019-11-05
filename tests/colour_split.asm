;; firmware functions we use in this example
kl_new_fast_ticker equ &bce0
;;mc_set_mode        equ &bd1c
mc_set_inks        equ &bd25
mc_wait_flyback    equ &bd19
kl_new_frame_fly   equ &bcd7
kl_init_event equ &bcef
km_wait_char equ &bb06

;;txt_set_pen equ &bb90
txt_set_cursor equ &bb75
txt_output equ &bb5a
scr_set_mode equ &bc0e

;; this code must be in the range &4000-&bfff
;; to work correctly.
;;
;; assemble then from BASIC:
;;
;; call &8000
org &8000

start:
;; set mode
ld a,1
call scr_set_mode

;; draw text in each of the pens
ld h,8
ld l,1
ld b,4
tl1:
push bc
push hl
call txt_set_cursor

ld hl,text
call print_msg
pop hl
ld a,l
add a,8
ld l,a
pop bc
djnz tl1

call km_wait_char

;; wait for a screen refresh
;; we do this to synchronise our effect
call mc_wait_flyback

ld a,6
ld (ticker_counter),a
ld hl,colours
ld (current_colour_pointer),hl

;; install interrupt
ld hl,ticker_event_block
ld b,%10000010
ld c,&80
ld de,ticker_function
call kl_new_fast_ticker

;; return to BASIC
loop1:
jp loop1

text:
defb 15,0,"PEN 0",13,10
defb 15,1,"PEN 1",13,10
defb 15,2,"PEN 2",13,10
defb 15,3,"PEN 3","$"

print_msg:
ld a,(hl)
inc hl
cp "$"
ret z
call txt_output
jr print_msg



;; this is initialised by
;; the firmware; holds runtime state of ticker interrupt
ticker_event_block:
defs 10

;; this is the function called each 1/300th of a second
ticker_function:
push af
push hl

;; The 1/300th of a second interrupt effectively splits
;; the screen into 6 sections of equal height. Each section
;; spans the entire width of the screen.
;;
;; We want to ensure that the effect is stationary so we reset
;; every 6 calls of this function.
ld a,(ticker_counter)
dec a
ld (ticker_counter),a
or a
jr nz,ticker_function2
ld a,6
ld (ticker_counter),a
ld hl,colours
ld (current_colour_pointer),hl

ticker_function2:

;; setting the colours will occur immeditately.

;; get pointer to current colours
ld de,(current_colour_pointer)
call mc_set_inks

;; update colours pointer
ld hl,(current_colour_pointer)
ld bc,17
add hl,bc
ld (current_colour_pointer),hl


pop af
pop hl
ret

ticker_counter: defb 0

current_colour_pointer: defw colours

;; The 1/300th of a second interrupt effectively splits
;; the screen into 6 sections of equal height. Each section
;; spans the entire width of the screen.
;;
colours:

;; colours for 1st section
;; border colour, followed by pen 0, pen 1 up to pen 15
;; NOTE: These are hardware colour numbers
defb &04,&14,&19,&13,&0c,&0b,&14,&15,&0d,&06,&1e,&1f,&07,&12,&19,&04,&17
;; colours for 2nd section
defb &03,&01,&1f,&1e,&0c,&0b,&14,&15,&0d,&06,&1e,&1f,&07,&12,&19,&04,&17
;; colours for 3rd section
defb &11,&05,&02,&0b,&0c,&0b,&14,&15,&0d,&06,&1e,&1f,&07,&12,&19,&04,&17
;; colours for 4th section
defb &0c,&08,&03,&1b,&0c,&0b,&14,&15,&0d,&06,&1e,&1f,&07,&12,&19,&04,&17
;; colours for 5th section
defb &12,&09,&0c,&06,&0c,&0b,&14,&15,&0d,&06,&1e,&1f,&07,&12,&19,&04,&17
;; colours for 6th section
defb &1f,&0a,&06,&16,&0c,&0b,&14,&15,&0d,&06,&1e,&1f,&07,&12,&19,&04,&17

end start