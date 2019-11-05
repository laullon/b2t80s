org &4000
nolist

scr_set_mode equ &bc0e
txt_output equ &bb5a
km_read_char equ &bb09
start:
call set_crtc


ld a,2
call scr_set_mode

ld bc,24*80
ld d,' '
l2:
inc d
ld a,d
cp &7f
jr nz,no_char_reset
ld d,' '
no_char_reset:
ld a,d
call txt_output
dec bc
ld a,b
or c
jr nz,l2

;;---------------------------------------------

loop1:
ld b,&f5
l1:
in a,(c)
rra
jr nc,l1
call check_keys

ld bc,&bc0c
out (c),c
ld hl,(scrl2)
inc b
out (c),h
ld bc,&bc0d
out (c),c
inc b
out (c),l

halt
halt
halt
jp loop1


set_crtc:
ld bc,&bc00
set_crtc_vals:
out (c),c
inc b
ld a,(hl)
out (c),a
dec b
inc hl
inc c
ld a,c
cp 14
jr nz,set_crtc_vals
ret

crtc_vals:
defb &3f
defb 48
defb 49
defb &89
defb 38
defb 0
defb 35
defb 35
defb 0
defb 7
defb 0
defb 0
defb &0c
defb 208


;; check if a key has been pressed and perform action if it has
check_keys:
call km_read_char
ret nc
;; A = ascii char of key pressed
;; we check for both upper and lower case chars
cp '8'
jp z,scroll_up
cp '2'
jp z,scroll_down
cp '4'
jp z,scroll_left
cp '6'
jp z,scroll_right
cp '7'
jp z,scroll_up_left
cp '9'
jp z,scroll_up_right
cp '1'
jp z,scroll_down_left
cp '3'
jp z,scroll_down_right

ret


scroll_down_right:
ld c,1
call scroll_down
call scroll_right
ret


scroll_down_left:
ld c,1
call scroll_down
call scroll_left
ret


scroll_up_left:
ld c,1
call scroll_up
call scroll_left
ret

scroll_up_right:
ld c,1
call scroll_up
call scroll_right
ret

scroll_right:
ld hl,(scrl2)
inc hl
ld a,h
and &3
or &c
ld h,a
ld (scrl2),hl
ret


scroll_left:
ld hl,(scrl2)
dec hl
ld a,h
and &3
or &c
ld h,a
ld (scrl2),hl
ret
;;---------------------------------------------

;; update these seperate of display
scroll_up:
ld hl,(scrl2)
ld bc,48
add hl,bc
ld a,h
and &3
or &c
ld h,a
ld (scrl2),hl
ret

;; update these seperate of display
scroll_down:
ld hl,(scrl2)
or a
ld bc,48
sbc hl,bc
ld a,h
and &3
or &c
ld h,a
ld (scrl2),hl
ret


scrl2:
defw 0

end start