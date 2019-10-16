package api

import (
	"bytes"
	"strconv"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestPrivateVxcCreateInput_toPayload(t *testing.T) {
	name := acctest.RandString(10)
	ref := acctest.RandString(10)
	rate := []uint64{100, 200, 300, 400, 500, 600, 700, 800, 900, 1000}[acctest.RandIntRange(0, 10)]
	rateString := strconv.FormatUint(rate, 10)
	vlanA := uint64(acctest.RandIntRange(1, 4094))
	vlanAString := strconv.FormatUint(vlanA, 10)
	vlanB := uint64(acctest.RandIntRange(1, 4094))
	vlanBString := strconv.FormatUint(vlanB, 10)
	uuidA := uuid.New().String()
	uuidB := uuid.New().String()
	testCases := []struct {
		i PrivateVxcCreateInput
		o []byte
	}{
		{
			PrivateVxcCreateInput{ // 0
				InvoiceReference: ref,
				Name:             name,
				ProductUidA:      uuidA,
				ProductUidB:      uuidB,
				RateLimit:        rate,
				VlanA:            vlanA,
				VlanB:            vlanB,
			},
			[]byte(`[{"productUid":"` + uuidA + `","associatedVxcs":[{"productName":"` + name + `","rateLimit":` + rateString + `,"costCentre":"` + ref + `","aEnd":{"vlan":` + vlanAString + `},"bEnd":{"productUid":"` + uuidB + `","vlan":` + vlanBString + `}}]}]`),
		},
		{
			PrivateVxcCreateInput{ // 1
				InvoiceReference: "",
				Name:             name,
				ProductUidA:      uuidA,
				ProductUidB:      uuidB,
				RateLimit:        rate,
				VlanA:            vlanA,
				VlanB:            vlanB,
			},
			[]byte(`[{"productUid":"` + uuidA + `","associatedVxcs":[{"productName":"` + name + `","rateLimit":` + rateString + `,"aEnd":{"vlan":` + vlanAString + `},"bEnd":{"productUid":"` + uuidB + `","vlan":` + vlanBString + `}}]}]`),
		},
		{
			PrivateVxcCreateInput{ // 2
				InvoiceReference: ref,
				Name:             "",
				ProductUidA:      uuidA,
				ProductUidB:      uuidB,
				RateLimit:        rate,
				VlanA:            vlanA,
				VlanB:            vlanB,
			},
			[]byte(`[{"productUid":"` + uuidA + `","associatedVxcs":[{"productName":"","rateLimit":` + rateString + `,"costCentre":"` + ref + `","aEnd":{"vlan":` + vlanAString + `},"bEnd":{"productUid":"` + uuidB + `","vlan":` + vlanBString + `}}]}]`),
		},
		{
			PrivateVxcCreateInput{ // 3
				InvoiceReference: ref,
				Name:             name,
				ProductUidA:      "",
				ProductUidB:      uuidB,
				RateLimit:        rate,
				VlanA:            vlanA,
				VlanB:            vlanB,
			},
			[]byte(`[{"productUid":"","associatedVxcs":[{"productName":"` + name + `","rateLimit":` + rateString + `,"costCentre":"` + ref + `","aEnd":{"vlan":` + vlanAString + `},"bEnd":{"productUid":"` + uuidB + `","vlan":` + vlanBString + `}}]}]`),
		},
		{
			PrivateVxcCreateInput{ // 4
				InvoiceReference: ref,
				Name:             name,
				ProductUidA:      uuidA,
				ProductUidB:      "",
				RateLimit:        rate,
				VlanA:            vlanA,
				VlanB:            vlanB,
			},
			[]byte(`[{"productUid":"` + uuidA + `","associatedVxcs":[{"productName":"` + name + `","rateLimit":` + rateString + `,"costCentre":"` + ref + `","aEnd":{"vlan":` + vlanAString + `},"bEnd":{"productUid":"","vlan":` + vlanBString + `}}]}]`),
		},
		{
			PrivateVxcCreateInput{ // 5
				InvoiceReference: ref,
				Name:             name,
				ProductUidA:      uuidA,
				ProductUidB:      uuidB,
				RateLimit:        0,
				VlanA:            vlanA,
				VlanB:            vlanB,
			},
			[]byte(`[{"productUid":"` + uuidA + `","associatedVxcs":[{"productName":"` + name + `","rateLimit":0,"costCentre":"` + ref + `","aEnd":{"vlan":` + vlanAString + `},"bEnd":{"productUid":"` + uuidB + `","vlan":` + vlanBString + `}}]}]`),
		},
		{
			PrivateVxcCreateInput{ // 6
				InvoiceReference: ref,
				Name:             name,
				ProductUidA:      uuidA,
				ProductUidB:      uuidB,
				RateLimit:        rate,
				VlanA:            0,
				VlanB:            vlanB,
			},
			[]byte(`[{"productUid":"` + uuidA + `","associatedVxcs":[{"productName":"` + name + `","rateLimit":` + rateString + `,"costCentre":"` + ref + `","bEnd":{"productUid":"` + uuidB + `","vlan":` + vlanBString + `}}]}]`),
		},
		{
			PrivateVxcCreateInput{ // 7
				InvoiceReference: ref,
				Name:             name,
				ProductUidA:      uuidA,
				ProductUidB:      uuidB,
				RateLimit:        rate,
				VlanA:            vlanA,
				VlanB:            0,
			},
			[]byte(`[{"productUid":"` + uuidA + `","associatedVxcs":[{"productName":"` + name + `","rateLimit":` + rateString + `,"costCentre":"` + ref + `","aEnd":{"vlan":` + vlanAString + `},"bEnd":{"productUid":"` + uuidB + `"}}]}]`),
		},
	}
	for i, tc := range testCases {
		p, err := tc.i.toPayload()
		if err != nil {
			t.Errorf("PrivateVxcCreateInput.toPayload (#%d): %w", i, err)
		}
		if !bytes.Equal(tc.o, p) {
			t.Errorf("PrivateVxcCreateInput.toPayload (#%d):\n\tgot      `%s`\n\texpected `%s`", i, p, tc.o)
		}
	}
}
