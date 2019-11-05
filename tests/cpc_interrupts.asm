org &8000

start:

  di
  im 1

  ld a, &c3
  ld bc, int1
  ld (&38), a
  ld (&39), bc

  ld bc, &7f10
  out (c), c

  ei

  jp $

int1:
  push bc

  ld bc, &7f5c
  out (c),c

  ld bc, int2
  jp int_ret


int2:
  push bc

  ld bc, &7f55
  out (c),c

  ld bc, int3
  jp int_ret

int3:
  push bc
  ld bc, &7f4b
  out (c),c

  ld bc, int4
  jr int_ret

int4:
  push bc

  ld bc, &7f4d
  out (c),c

  ld bc, int5
  jr int_ret

int5:
  push bc


  ld bc, &7f56
  out (c),c

  ld bc, int6
  jr int_ret

int6:
  push bc

  ld bc, &7f5f
  out (c),c

  ld bc, int1
  jr int_ret

int_ret:
  ld (&39), bc

  pop bc
  ei
  ret

  end start