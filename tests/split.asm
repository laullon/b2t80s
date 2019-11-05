;; This example shows how to make a vertical split/rupture.
;;
;; A vertical split/rupture is used to divide the display (in the vertical dimension)
;; into more than one block. Each block spans the entire width of the display,
;; can have a programmable height and it's own start address.
;; Therefore parts of the screen can be scrolled while other parts remain static.
;; The split is another method to make overscan, but the split must be refreshed
;; every frame for the effect to be maintained. The split can also require exact
;; timing.
;;
;; When making a split, be aware of the differences in the CRTC models used by Amstrad.
;; (HD6845S/UM6845 (type 0), UM6845R (type 1), MC6845 (type 2), CPC+ CRTC inside ASIC
;; (type 3), CPC CRTC inside Pre-ASIC in cost-down CPCs (type 4)).
;; You will certainly notice the difference between these CRTCs when you program splits.
;; 
;; The split is not a new effect, it has been used in some games (i.e. Mission Genocide, Octoplex
;; Prehistorik 2, Super Cauldron and Snowstrike, to name a few), and many demos.
;; This effect was made popular by demo groups such as Logon System.
;; 
;; This method is technical and requires good understanding of the operation
;; of the CRTC.
;;
;; Extensive comments have been included to explain the operation of the split
;; and the reasons for every CRTC register write.
;;
;; Abreviations used:
;;
;; VADJC = internal Vertical Adjust Count register of the CRTC
;; VC = internal Vertical Count register of the CRTC
;; RC = internal Raster Count register of the CRTC
;; HC = internal Horizontal Count register of the CRTC
;; HTOT = Horizontal total register of CRTC (register 0)
;; HDISP = Horizontal displayed register of CRTC (register 1)
;; VTOT = Vertical total register of CRTC (register 4)
;; VSYNCPOS = Vertical sync position register of CRTC (register 7)
;; VSYNC = The vertical sync signal. This can be monitored through PPI port B.
;; HSYNCWIDTH = Horizontal sync width (programmed by register 3 of the CRTC)
;; HSYNCPOS = Horizontal sync position register of CRTC (register 2)
;; HDISP = Horizontal displayed register of CRTC (register 1)
;; VADJ = Vertical adjust register of CRTC (register 5)
;; MR = Maximum raster register of CRTC (register 9)
;;
;; (c) Kevin Thacker, 2002
;;
;; This code has been released under the GNU Public License V2.

org &4000
nolist

;;------------------------------------------------
;; install a interrupt handler

di						;; disable interrupts
im 1					;; interrupt mode 1 (CPU will jump to &0038 when a interrupt occrs)
ld hl,&c9fb				;; C9 FB are the bytes for the Z80 opcodes EI:RET
ld (&0038),hl			;; setup interrupt handler

;; the interrupt handler we have setup is minimal:
;;
;; EI
;; RET
;; 
;; This will re-enable interrupts and then return to program control.
;; 
;; We know the CPU time taken for the interrupt handler because we have defined
;; the interrupt handler and the timing of the instructions is known.
;;

ei						;; re-enable interrupts
;;--------------------------------------------------------------------------


;;------------------------------------------------
;; define the horizontal and vertical sync widths
;; 
;; the vertical sync width is programmable on CRTC type 0, type 3 and type 4.
;; the vertical sync width is fixed on CRTC type 1 and type 2.
;;
;; the horizontal sync width is programmable on all CRTC types.
;;

ld bc,&bc03							;; select vertical and horizontal sync register of the CRTC
out (c),c
ld bc,&bd00+8						;; set vertical sync width = 16, horizontal sync width = 8
out (c),c

;;------------------------------------------------
;; define the horizontal sync position
;; 
;; With the vertical split, this is used to position the screen horizontally
;; within the display. Increasing this value will move the screen to the left
;; decreasing this value will move the screen to the right.
;; 
;; Note that on CRTC type 2, (HSYNCPOS+HSYNCWIDTH)<HTOT otherwise interrupts are not 
;; generated, and the split will not work!
;;
;; This is setup here once as the horizontal position will be the same
;; for all blocks.

ld bc,&bc02							;; select horizontal sync position register of the CRTC
out (c),c
ld bc,&bd00+48						;; set horizontal sync to 48
out (c),c

;;------------------------------------------------
;; define the horizontal displayed
;; 
;; The horizontal displayed defines the number of CRTC characters to display
;; on each CRTC scan line. (Each CRTC character is 2 bytes)
;; 
;; In this example this is setup here once as the horizontal displayed will 
;; be the same for all blocks.

ld bc,&bc01							;; select horizontal displayed register of the CRTC
out (c),c
ld bc,&bd00+48						;; set horizontal sync to 48
out (c),c

;;------------------------------------------------
;; Setup the vertical displayed so that it is larger than 
;; the vertical height of the tallest split block.
;;
;; In this example, this is setup here once as the horizontal displayed 
;; will be the same for all blocks.
;;
;; The vertical and horizontal displayed signals define when
;; the border is shown:
;;
;; - If VDISP<VTOT of any block then border colour will be shown 
;; between the block (where VDISP<VTOT) and the start of the next block.
;;
;; - If HDISP<HTOT of any block then border colour will be shown
;; on the sides of that block (where HDISP<HTOT).

ld bc,&bc06
out (c),c
ld bc,&bdff
out (c),c

;;------------------------------------------------
;; Setup the maximum raster so that it is the same for every split block.
;;
;; The following registers define the height of a split block:
;; - maximum raster register
;; - vertical total register
;; - vertical adjust register
;;
;; The height is computed as ((MR+1))*(VTOT+1))+VADJ  HTOT times.
;;
;; The time for a split block is defined as:
;; (((MR*HTOT)+1)*(VTOT+1))+(VADJ*HTOT)
;;
;; If HTOT is 64, then RC will increment once for every monitor scan-line.
;;
;; VERTICAL ADJUST:
;;
;; The vertical adjust, if defined to be a value other than 0, is activated
;; at the end of the split block it is defined for. The VADJC will increment
;; until it's value matches the programmed vertical adjust. Therefore the split block
;; is extended by (HTOT*VADJ) time.
;;
;; In this example where every HTOT is the duration of a monitor scan-line,
;; vertical adjust will count additional monitor scan-lines.
;;
;; MAXIMUM RASTER ADDRESS:
;; 
;; The maximum raster adjust defines the number of times RC must increment
;; for each CRTC character row.
;;
;; In this example where every HTOT is the duration of a monitor scan-line,
;; the maximum raster address defines the height of the CRTC character row in
;; monitor scan-lines.
 
ld bc,&bc05									;; select vertical adjust register of CRTC
out (c),c
ld bc,&bd00									;; 0 scan-lines of vertical adjust
out (c),c

ld bc,&bc09									;; select maximum raster register of CRTC
out (c),c
ld bc,&bd07									;; 8 scan-lines per CRTC character row
out (c),c

;;----------------------------------------------------
;; this main loop is executed every frame to maintain
;; the split. The timing required to maintain this split is not 
;; too important because the split is simple. If we are attempting
;; a more complex split, then the timing is often critical.


main_loop

;;----------------------------------------------------
;; wait for the start of the vsync
;;
;; the position of the vsync for the first time this code is executed,
;; and then subsequent executions will be different.
;;
;; The position of the first vsync is defined by the initial CRTC values
;; when this code is first entered.
;; 
;; The position of subsequent vsyncs is determined by our vertical split code.
;;
;; (Note: This code will also recognise when a VSYNC is already in progress.
;; It will only see the start of a VSYNC if the VSYNC was inactive before this
;; checking loop is started. 
;;
;; If the test is entered when the VSYNC has already started then our split
;; could be unstable as we may not always be synchronised to the CRTC at the same
;; point each frame. If the code requires the timing to be accurate, and we miss the 
;; start of the VSYNC, then the CRTC registers we program (to create the split) could
;; be setup at the wrong time and the split will break).

ld b,&f5				;; B = I/O address of PPI port B
vsync
in a,(c)				;; read PPI port B input
rra						;; transfer bit 0 into carry
jr nc,vsync				;; if carry=0 then vsync= 0 (inactive),
						;; if carry=1 then vsync=1 (active).

;;--------------------------------------------------
;; at this point we have seen the start of a vsync,
;; and are synchronised (but not exactly) to
;; the start of it.
;;
;;
;; (If we detected the start of the VSYNC, then we also know when the next interrupt
;; will occur.
;; 
;; The interrupt counter is updated every HSYNC.
;; The interrupt counter reset is synchronised to the start of the VSYNC.
;; A interrupt request is issued when the interrupt counter reaches 52.
;; 
;; The next interrupt could occur in two HSYNC times, assuming that
;; the previous interrupt was not serviced less than 32 lines ago.
;;
;; Otherwise the next interrupt will occur in 52+2 HSYNC times.
;;
;; A perfect split relies on predicting the position of the start of the VSYNC
;; and the position of the interrupts, as these are the signals we use to
;; synchronise with the display, and this means that we can setup the next split
;; block at the correct position).
;;
;; in our split example, the VSYNC is programmed to occur when 
;; VC = 0, RC = 0. (This trigger point is setup at the end of the main loop).
;; We know that when we have seen the start of the VSYNC the CRTC
;; will be processing VC = 0, RC=0, and HC will be somewhere on this line. 
;; (The value of HC will vary as the VSYNC testing check is accurate to about 8us).
;; 

;;-----------------------------------------------------------------
;; To stop the VSYNC from triggering again, until we want it
;; to, program the vsync position so that it is larger than 
;; the vertical height of the tallest split block.
;;
;; Now Vsync will not be triggered, because VC can never equal the 
;; vertical sync position we have programmed.
;;
;; The CRTC will continue to display the split blocks until we set
;; a valid VSYNC position which can be reached by VC. (e.g. VSYNCPOS<VTOT)

ld bc,&bc07				;; select vertical sync position register of the CRTC
out (c),c
ld bc,&bdff				;; set vertical sync position (255)
out (c),c

;;--------------------------------------------------------------
;; The height of each split block is defined by:
;;
;; - the vertical total register of the CRTC. The register must be reprogrammed for each
;; block.
;; - the maximum raster register of the CRTC.
;; - the vertical adjust register of the CRTC. 
;;
;;
;; In the case where HTOT is 64, the height of the split block is therefore
;; defined in complete monitor scan-lines.
;;
;; For compatibility reasons, the vertical total register should be written
;; when VC<new_VTOT, otherwise the CRTC may finish the current split block and
;; start a new one, and this may not be the effect you want.
;;
;; For a steady split, when HTOT is 64 the frame should use 312 scan-lines.
;; (64 microseconds * 312 crtc scan lines)=19968us.
;;
;; 1/50 second = 0.02 seconds
;; 1/1000000 second = 1 microsecond
;; 
;; 1 microsecond*19968 = 0.019968 seconds = approx 0.02 seconds
;;
;; If more scan-lines are used then the whole display may roll, or the display
;; will flicker.

;;---------------------------------------------------------------
;; set height of first split block. At this point VC<new_VTOT, so first split
;; block will become 5 char lines tall. The split block will
;; be 40 (5*8) scanlines.

ld bc,&bc04				;; select vertical total register of the CRTC
out (c),c
ld bc,&bd00+10-1		;; set height (height-1) of the first split block
out (c),c

;;----------------------------------------------------------------------------
;; HALT = wait for next interrupt to trigger, execute the interrupt handler
;; and continue with program execution. Assumes interrupts are ENABLED!
;;
;; NOTE: we use the interrupts for two purposes:
;; - to sync exactly with the CRTC to an exact point on the screen
;; (This position is determined by the HSYNCPOS, the CPU time required to acknowledge
;; the interrupt and the CPU time required to execute the interrupt handler)
;;
;; - to waste some time, so that we can wait until the appropiate position
;; to setup the split. When using the standard Amstrad CPC interrupts
;; the HALT provides a coarse positioning method.
;;
;; A HALT instruction will execute the equivalent of NOP instructions until
;; a interrupt request occurs. A NOP is one of the fastest instructions taking
;; 1us for each cycle. Therefore, when a interrupt request occurs, there will
;; always be the same time between interrupt request occuring and interrupt
;; request being acknowledge and serviced. This instruction works the best
;; when there is no interrupt request outstanding when the HALT is first executed.)

halt
;; At this point we are synchronised exactly to a predictable position on the screen.
;; We can now use software loops to delay to the exact position we want.

;;-------------------------------------------------------
;; In this example, the previous interrupt was processed greater than 32 lines
;; ago. 
;;
;; If IRQ_ACK is the time between the interrupt being triggered and then acknowledged
;; by the CPU, and IRQ_FUNC is the time for the interrupt handler to execute and
;; return control back to the program, then the interrupt will occur at time:
;;
;; HTOT+HSYNCPOS+IRQ_ACK+IRQ_FUNC
;; 
;; (This assumes that the HSYNCPOS and HTOT remain constant).
;; 
;; When VSYNCPOS=0 then the CRTC should be processing VC=0, RC=1 at the time the
;; interrupt is requested, and by the time we receive control the CRTC should
;; be processing VC=0, RC=2.

;; delay until VC=1.
;; Delay required = 8 scanlines = 8 * 64 us = 512 us
;; 
;; for DJNZ: if b-1==0, instruction will execute in 3us
;; if b-1<>0, instruction will execute in 4us.
;;

ld b,127						;; [2]
wait1 
djnz wait1						;; [3/4]

;; for this loop:
;; 
;; (14*4)+3+2 = 513 us

;;-----------------------------------------------------------------------------------
;; set the start address of the *next* split block
;;
;; - for CRTC type 0,2,3 and 4, this start address will take effect at the start
;; of the next split block.
;; - for CRTC type 1, it is possible to reprogram the start address of the current
;; split block if VC=0, otherwise this start address will take effect at the start
;; of the next split block.
;; 
;; For compatibility with other CRTC types, attempt to change the start position
;; when VC>0 and VC<(VTOT-1).
;;
;; NOTE: the start address of the first split block is defined at the end of the loop


ld hl,&1000							;; start address in CRTC form (&4000 - &7fff in RAM)
ld bc,&bc0c							;; start address high register of CRTC
out (c),c							;; select start address high register of CRTC
inc b								;; B = &BD
out (c),h							;; write data to start address low register

ld bc,&bc0d							;; select start address low register of CRTC
out (c),c							;; select start address low register of CRTC
inc b								;; B = &BD
out (c),l							;; write data to start address low register

;-----------------------------------------------------------------------------

halt
;; in this example, the CRTC has completed 2 + 52 monitor scanlines
;; since the start of the VSYNC.
;;
;; In this example HTOT and HSYNC remain constant, therefore the amount of time
;; that has passed is:
;;
;; (HTOT+HSYNCPOS)+(HTOT-HSYNCPOS)+(51*HTOT)+HSYNCPOS
;;
;; where (HTOT+HSYNCPOS) is the time to the first interrupt request,
;; (HTOT-HSYNCPOS) is the time from the first interrupt request to the first HTOT
;; after the interrupt request,
;; (51*HTOT)+HSYNCPOS is the time to the second interrupt request.
;;
;; approx 6.75 scanlines 

ld b,15
wait2 djnz wait2

ld bc,&bc04
out (c),c
ld bc,&bd00+24             
out (c),c

  ld bc, &7f5f
  out (c),c

;------------------------------------------------------------------------------
;blk3
halt

  ld bc, &7f56
  out (c),c

;------------------------------------------------------------------------------

halt
ld b,15
wait4 djnz wait4
ld bc,&bc0c
out (c),c
ld bc,&bd00
out (c),c
ld bc,&bc0d
out (c),c
ld bc,&bd00
out (c),c

halt

halt
ld b,15
wait6 djnz wait6

ld bc,&bc04
out (c),c
ld bc,&bd00+5 ;5+25+6=36 (nearly 39!)
out (c),c

ld bc,&bc0c
out (c),c
ld bc,&bd00+%00010000 ;top section of screen
out (c),c
ld bc,&bc0d
out (c),c
ld bc,&bd00
out (c),c


;;-------------------------------------------------------------
;; To maintain a steady split we must do the following: 
;; - ensure the register writes to the CRTC occur at the same position every frame
;; - force a VSYNC to be generated once per frame, once per 312 complete (64us) scan-lines
;;  (once every 19968us)
;;
;; If we want a steady split that will work on every CRTC type:
;; - ensure the register writes to the CRTC occur at a time that is compatible with
;; every CRTC

;; reprogram a new vertical sync position to force a VSYNC to be triggered.
;; 
;; We want our VSYNC to start on VC=0, RC=0, HC=0.
;;
;; - If VC<>0 then the VSYNC will be triggered at the start of the next split (VC=0, RC=0, HC=0).
;; - If VC==0, then on some CRTC's a VSYNC will be triggered immediatly (VC=0, RC!=0, HC=??), this
;; could cause a bad split if our split requires that the VSYNC must occur when VC=0, RC=0, HC=0)!
;;

ld bc,&bc07				;; vertical sync position register of the CRTC
out (c),c				;; select vertical sync position of the CRTC
ld bc,&bd00				;; vertical sync position = 0
out (c),c				;; set vertical sync position

;;----------------------------------------
;; continue to loop so that the split is maintained

jp main_loop

  end main_loop