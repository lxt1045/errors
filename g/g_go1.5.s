// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Assembly to get into package runtime without using exported symbols.
// See https://github.com/golang/go/blob/release-branch.go1.4/misc/cgo/test/backdoor/thunk.s


//go:build (386 || amd64 || amd64p32 || arm || arm64) && gc && go1.5


#include "textflag.h"

// func getg() *g
TEXT ·getg(SB),NOSPLIT,$0-8
#ifdef GOARCH_386
	MOVL (TLS), AX
	MOVL AX, ret+0(FP)
#endif
#ifdef GOARCH_amd64
	MOVQ (TLS), AX
	MOVQ AX, ret+0(FP)
#endif
#ifdef GOARCH_arm
	MOVW g, ret+0(FP)
#endif
#ifdef GOARCH_arm64
	MOVD g, ret+0(FP)
#endif
	RET
