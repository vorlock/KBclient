// Copyright 2016 Keybase Inc. All rights reserved.
// Use of this source code is governed by a BSD
// license that can be found in the LICENSE file.

package libkbfs

import (
	"testing"

	"github.com/keybase/client/go/libkb"
	keybase1 "github.com/keybase/client/go/protocol"
	"golang.org/x/net/context"
)

// Test that Put/Get works for TLF crypt key server halves.
func TestKeyServerLocalTLFCryptKeyServerHalves(t *testing.T) {
	// simulate two users
	var userName1, userName2 libkb.NormalizedUsername = "u1", "u2"
	config1, uid1, ctx := kbfsOpsConcurInit(t, userName1, userName2)
	defer CheckConfigAndShutdown(t, config1)

	config2 := ConfigAsUser(config1.(*ConfigLocal), userName2)
	defer CheckConfigAndShutdown(t, config2)
	_, uid2, err := config2.KBPKI().GetCurrentUserInfo(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	publicKey1, err := config1.KBPKI().GetCurrentCryptPublicKey(ctx)
	if err != nil {
		t.Fatal(err)
	}
	publicKey2, err := config2.KBPKI().GetCurrentCryptPublicKey(ctx)
	if err != nil {
		t.Fatal(err)
	}

	serverHalf1 := MakeTLFCryptKeyServerHalf([32]byte{1})
	serverHalf2 := MakeTLFCryptKeyServerHalf([32]byte{2})
	serverHalf3 := MakeTLFCryptKeyServerHalf([32]byte{3})
	serverHalf4 := MakeTLFCryptKeyServerHalf([32]byte{4})

	// write 1
	keyHalves := make(map[keybase1.UID]map[keybase1.KID]TLFCryptKeyServerHalf)
	deviceHalves := make(map[keybase1.KID]TLFCryptKeyServerHalf)
	deviceHalves[publicKey1.kid] = serverHalf1
	keyHalves[uid1] = deviceHalves

	err = config1.KeyOps().PutTLFCryptKeyServerHalves(ctx, keyHalves)
	if err != nil {
		t.Fatal(err)
	}

	// write 2
	keyHalves = make(map[keybase1.UID]map[keybase1.KID]TLFCryptKeyServerHalf)
	deviceHalves = make(map[keybase1.KID]TLFCryptKeyServerHalf)
	deviceHalves[publicKey1.kid] = serverHalf2
	keyHalves[uid1] = deviceHalves

	err = config1.KeyOps().PutTLFCryptKeyServerHalves(ctx, keyHalves)
	if err != nil {
		t.Fatal(err)
	}

	// write 3 and 4 together
	keyHalves = make(map[keybase1.UID]map[keybase1.KID]TLFCryptKeyServerHalf)
	deviceHalves1 := make(map[keybase1.KID]TLFCryptKeyServerHalf)
	deviceHalves2 := make(map[keybase1.KID]TLFCryptKeyServerHalf)
	deviceHalves1[publicKey1.kid] = serverHalf3
	keyHalves[uid1] = deviceHalves1
	deviceHalves2[publicKey2.kid] = serverHalf4
	keyHalves[uid2] = deviceHalves2

	err = config1.KeyOps().PutTLFCryptKeyServerHalves(ctx, keyHalves)
	if err != nil {
		t.Fatal(err)
	}

	serverHalfID1, err :=
		config1.Crypto().GetTLFCryptKeyServerHalfID(uid1, publicKey1.kid, serverHalf1)
	if err != nil {
		t.Fatal(err)
	}

	serverHalfID2, err :=
		config1.Crypto().GetTLFCryptKeyServerHalfID(uid1, publicKey1.kid, serverHalf2)
	if err != nil {
		t.Fatal(err)
	}

	serverHalfID3, err :=
		config1.Crypto().GetTLFCryptKeyServerHalfID(uid1, publicKey1.kid, serverHalf3)
	if err != nil {
		t.Fatal(err)
	}

	serverHalfID4, err :=
		config1.Crypto().GetTLFCryptKeyServerHalfID(uid2, publicKey2.kid, serverHalf4)
	if err != nil {
		t.Fatal(err)
	}

	half1, err := config1.KeyOps().GetTLFCryptKeyServerHalf(ctx, serverHalfID1, publicKey1)
	if err != nil {
		t.Error(err)
	}

	if half1 != serverHalf1 {
		t.Errorf("Expected %v, got %v", serverHalf1, half1)
	}

	half2, err := config1.KeyOps().GetTLFCryptKeyServerHalf(ctx, serverHalfID2, publicKey1)
	if err != nil {
		t.Error(err)
	}

	if half2 != serverHalf2 {
		t.Errorf("Expected %v, got %v", serverHalf2, half2)
	}

	half3, err := config1.KeyOps().GetTLFCryptKeyServerHalf(ctx, serverHalfID3, publicKey1)
	if err != nil {
		t.Error(err)
	}

	if half3 != serverHalf3 {
		t.Errorf("Expected %v, got %v", serverHalf3, half3)
	}

	half4, err := config1.KeyOps().GetTLFCryptKeyServerHalf(ctx, serverHalfID4, publicKey1)
	_, unauthorized := err.(MDServerErrorUnauthorized)
	if !unauthorized {
		t.Errorf("Expected unauthorized, got %v", err)
	}

	// try to get uid2's key now as uid2
	half4, err = config2.KeyOps().GetTLFCryptKeyServerHalf(ctx, serverHalfID4, publicKey2)
	if err != nil {
		t.Error(err)
	}

	if half4 != serverHalf4 {
		t.Errorf("Expected %v, got %v", serverHalf4, half4)
	}

	serverHalfIDNope, err := config1.Crypto().GetTLFCryptKeyServerHalfID(uid1, publicKey1.kid, serverHalf4)
	if err != nil {
		t.Error(err)
	}

	_, err = config1.KeyOps().GetTLFCryptKeyServerHalf(ctx, serverHalfIDNope, publicKey1)
	if err == nil {
		t.Error("GetTLFCryptKeyServerHalf(id2, keyGen2, publicKey2) unexpectedly succeeded")
	}
}
