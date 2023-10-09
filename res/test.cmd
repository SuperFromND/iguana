; This is an Iguana test CMD file.
; The purpose of this file is to be used for testing Iguana's basic functionality.
; This file will NOT function as an actual CMD file within MUGEN, IKEMEN Go, or any other engine that supports this format.

; Button Remaps (not yet supported by Iguana)
[Remap]
x = x
y = y
z = z
a = a
b = b
c = c
s = s

; Command time default
[Defaults]
command.time = 15

; Single buttons
[Command]
name = "a"
command = a

[Command]
name = "b"
command = b

[Command]
name = "c"
command = c

[Command]
name = "x"
command = x

[Command]
name = "y"
command = y

[Command]
name = "z"
command = z

; Double-tapping
[Command]
name = "RunForward"
command = F, F

[Command]
name = "RunBack"
command = B, B

; Common motion inputs
[Command]
name = "QuarterCircleForward"
command = D, DF, F, a

[Command]
name = "QuarterCircleBack"
command = D, DB, B, a

[Command]
name = "HalfCircleForward"
command = B, DB, D, DF, F, a

[Command]
name = "HalfCircleBack"
command = F, DF, D, DB, B, a

[Command]
name = "DragonPunchForward"
command = F, D, DF, a

[Command]
name = "DragonPunchBack"
command = B, D, DB, a

[Command]
name = "360Forward"
command = F, D, B, U, a

[Command]
name = "360Back"
command = F, U, B, D, a

; Defines the start of move definitions
[Statedef -1]

[State -1, A Button]
type = ChangeState
triggerall = command = "a"

[State -1, B Button]
type = ChangeState
triggerall = command = "b"

[State -1, C Button]
type = ChangeState
triggerall = command = "c"

[State -1, X Button]
type = ChangeState
triggerall = command = "x"

[State -1, Y Button]
type = ChangeState
triggerall = command = "y"

[State -1, Z Button]
type = ChangeState
triggerall = command = "z"

[State -1, Dash Forward]
type = ChangeState
triggerall = command = "RunForward"

[State -1, Dash Back]
type = ChangeState
triggerall = command = "RunBack"

[State -1, Multiple Buttons]
type = ChangeState
triggerall = command = "a"
triggerall = command = "b"
triggerall = command = "c"

[State -1, Dragon Punch]
type = ChangeState
triggerall = command = "DragonPunchForward"
triggerall = command = "DragonPunchBack"

[State -1, Move with Power]
type = ChangeState
triggerall = command = "HalfCircleForward"
triggerall = command = "HalfCircleBack"
triggerall = power > 1000

[State -1, 360 Motions]
type = ChangeState
triggerall = command = "360Forward"
triggerall = command = "360Back"

[State -1, Multiple Moves via &&]
type = ChangeState
triggerall = command = "x" && command = "y"