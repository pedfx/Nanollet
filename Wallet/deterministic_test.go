package Wallet

import (
	"bytes"
	"github.com/brokenbydefault/Nanollet/Util"
	"testing"
)

func TestReadSeedFY(t *testing.T) {
	fy, _ := ReadSeedFY("00000F15106351F426AC9B79321BF5421089835B2671F625C91272C5077EB2EA182D4E86CB")

	if fy.Type != 0 || fy.Version != 0 || fy.Thread != 16 || fy.Time != 15 || fy.Memory != 21 {
		t.Error("ReadWrong")
	}

}

func TestReadSeedFY2(t *testing.T) {
	p := SeedFY{}
	p.Version = 0
	p.Salt = []byte{0xF0, 0xF0, 0xF0, 0xF0, 0xF0, 0xF0, 0xF0}
	p.Memory = 22
	p.Thread = 24
	p.Time = 3

	if p.IsValid(V0, Nanollet) {
		t.Error("mininum salt light not verified")
	}
}

func TestReadSeedFY3(t *testing.T) {
	p := SeedFY{}
	p.Version = 0
	p.Salt = []byte{0xF0, 0xF0, 0xF0, 0xF0, 0xF0, 0xF0, 0xF0, 0xF0, 0xF0, 0xF0, 0xF0, 0xF0, 0xF0, 0xF0, 0xF0, 0xF0, 0xF0, 0xF0, 0xF0, 0xF0, 0xF0, 0xF0, 0xF0, 0xF0, 0xF0, 0xF0, 0xF0, 0xF0, 0xF0, 0xF0, 0xF0, 0xF0}
	p.Memory = 22
	p.Thread = 24
	p.Time = 3

	enc := p.Encode()

	readed, err := ReadSeedFY(enc)
	if readed.Type != p.Type || readed.Version != p.Version || readed.Thread != p.Thread || readed.Time != p.Time || readed.Memory != p.Memory {
		t.Error("ReadWrong")
	}

	if err != nil {
		t.Error(err)
	}
}

func TestNewSeedFY(t *testing.T) {

	var previous []byte
	for i := 0; i < 10; i++ {
		sf, err := NewSeedFY(V0, Nanollet)
		if err != nil {
			t.Error(t)
		}

		if bytes.Equal(previous, sf.Salt) {
			t.Error("Not random")
		}

		if !sf.IsValid(V0, Nanollet) {
			t.Error("Seed is not valid when generate from the default")
		}

	}

}

func TestSeed_CreateKeyPair(t *testing.T) {
	expectedsecretkey := map[uint32][]byte{
		1:       {0x26, 0x2B, 0x24, 0x50, 0x9D, 0x2C, 0x57, 0xD1, 0x18, 0x65, 0xB3, 0x4D, 0xDF, 0x21, 0x62, 0x63, 0x24, 0x5F, 0xEB, 0x79, 0x42, 0x61, 0x61, 0x9A, 0xFE, 0x81, 0x37, 0x70, 0x4E, 0xAB, 0xB4, 0xFE, 0x85, 0x3B, 0xF6, 0x25, 0xBF, 0x21, 0x6D, 0x4C, 0xA7, 0xA4, 0x89, 0x5C, 0xA6, 0x09, 0x59, 0x25, 0xF0, 0x15, 0xAB, 0x7B, 0x4E, 0x13, 0x07, 0x48, 0x4D, 0x50, 0x76, 0x6A, 0x58, 0xAA, 0x7D, 0x5B},
		1234567: {0x2F, 0xD2, 0x71, 0xC2, 0xA6, 0x75, 0xEA, 0x4E, 0x61, 0x8E, 0x47, 0x2B, 0x71, 0x86, 0xF1, 0xDA, 0x8E, 0x4D, 0xE2, 0xB8, 0xBF, 0xEC, 0x8A, 0xBA, 0xB4, 0xD0, 0xDB, 0x00, 0x3D, 0x67, 0x8A, 0xFC, 0x34, 0xCA, 0xA4, 0x29, 0x40, 0xFB, 0x62, 0xE9, 0x5E, 0xEB, 0x10, 0x9D, 0x3E, 0xBF, 0x89, 0x69, 0x56, 0x78, 0x65, 0x67, 0x1A, 0x76, 0xB8, 0xF8, 0x3D, 0x63, 0xE9, 0xAE, 0x98, 0xFE, 0xA7, 0x6D},
	}

	p := SeedFY{}
	p.Version = 0
	p.Salt = []byte{0xF0, 0xF0, 0xF0, 0xF0, 0xF0, 0xF0, 0xF0, 0xF0, 0xF0, 0xF0, 0xF0, 0xF0, 0xF0, 0xF0, 0xF0, 0xF0, 0xF0, 0xF0, 0xF0, 0xF0, 0xF0, 0xF0, 0xF0, 0xF0, 0xF0, 0xF0, 0xF0, 0xF0, 0xF0, 0xF0, 0xF0, 0xF0}
	p.Memory = 5
	p.Thread = 16
	p.Time = 10

	seed := p.RecoverSeed([]byte("myamazingpasswordthatisverygoodenough"), nil)

	for i, expected := range expectedsecretkey {
		_, sk, _ := seed.CreateKeyPair(Nano, i)

		if !bytes.Equal(sk, expected) {
			t.Errorf("recover seed failed, expecting %s and gives %s", Util.SecureHexEncode(expected), Util.SecureHexEncode(sk))
		}
	}
}

func TestRecoverCoinSeed(t *testing.T) {
	s, _ := Util.SecureHexDecode("15ac010515107c4ad61e4527cd8e43a9f6b4bd4fa6f1536361f3da686b0a88bc")
	seed := Seed(s)
	coin := Currency(0x00000001)

	if !bytes.Equal(Util.CreateKeyedXOFHash(32, seed, []byte{0x01, 0x00, 0x00, 0x00}), RecoverCoinSeed(seed, coin)) {
		t.Error("Coinseed is wrong")
	}

}
