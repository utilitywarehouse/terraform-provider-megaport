package api

import (
	"bytes"
	"strconv"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
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
	emptyString := ""
	emptyUint := uint64(0)
	testCases := []struct {
		i PrivateVxcCreateInput
		o []byte
	}{
		{ // 0
			PrivateVxcCreateInput{
				InvoiceReference: &ref,
				Name:             &name,
				ProductUidA:      &uuidA,
				ProductUidB:      &uuidB,
				RateLimit:        &rate,
				VlanA:            &vlanA,
				VlanB:            &vlanB,
			},
			[]byte(`[{"productUid":"` + uuidA + `","associatedVxcs":[{"productName":"` + name + `","rateLimit":` + rateString + `,"costCentre":"` + ref + `","aEnd":{"vlan":` + vlanAString + `},"bEnd":{"productUid":"` + uuidB + `","vlan":` + vlanBString + `}}]}]`),
		},
		{ // 1
			PrivateVxcCreateInput{
				InvoiceReference: &emptyString,
				Name:             &emptyString,
				ProductUidA:      &emptyString,
				ProductUidB:      &emptyString,
				RateLimit:        &emptyUint,
				VlanA:            &emptyUint,
				VlanB:            &emptyUint,
			},
			[]byte(`[{"productUid":"","associatedVxcs":[{"productName":"","rateLimit":0,"costCentre":"","aEnd":{"vlan":0},"bEnd":{"productUid":"","vlan":0}}]}]`),
		},
		{ // 2
			PrivateVxcCreateInput{},
			[]byte(`[{}]`),
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

func TestPrivateVxcUpdateInput_toPayload(t *testing.T) {
	name := acctest.RandString(10)
	ref := acctest.RandString(10)
	rate := []uint64{100, 200, 300, 400, 500, 600, 700, 800, 900, 1000}[acctest.RandIntRange(0, 10)]
	rateString := strconv.FormatUint(rate, 10)
	vlanA := uint64(acctest.RandIntRange(1, 4094))
	vlanAString := strconv.FormatUint(vlanA, 10)
	vlanB := uint64(acctest.RandIntRange(1, 4094))
	vlanBString := strconv.FormatUint(vlanB, 10)
	uuidB := uuid.New().String()
	emptyString := ""
	emptyUint := uint64(0)
	testCases := []struct {
		i PrivateVxcUpdateInput
		o []byte
	}{
		{ // 0
			PrivateVxcUpdateInput{
				InvoiceReference: &ref,
				Name:             &name,
				ProductUid:       &uuidB,
				RateLimit:        &rate,
				VlanA:            &vlanA,
				VlanB:            &vlanB,
			},
			[]byte(`{"aEndVlan":` + vlanAString + `,"bEndVlan":` + vlanBString + `,"costCentre":"` + ref + `","name":"` + name + `","rateLimit":` + rateString + `}`),
		},
		{ // 1
			PrivateVxcUpdateInput{
				InvoiceReference: &emptyString,
				Name:             &emptyString,
				ProductUid:       &emptyString,
				RateLimit:        &emptyUint,
				VlanA:            &emptyUint,
				VlanB:            &emptyUint,
			},
			[]byte(`{"aEndVlan":0,"bEndVlan":0,"costCentre":"","name":"","rateLimit":0}`),
		},
		{ // 2
			PrivateVxcUpdateInput{},
			[]byte(`{}`),
		},
	}
	for i, tc := range testCases {
		p, err := tc.i.toPayload()
		if err != nil {
			t.Errorf("PrivateVxcUpdateInput.toPayload (#%d): %w", i, err)
		}
		if !bytes.Equal(tc.o, p) {
			t.Errorf("PrivateVxcUpdateInput.toPayload (#%d):\n\tgot      `%s`\n\texpected `%s`", i, p, tc.o)
		}
	}
}
