package service

import (
	"encoding/hex"
	"regexp"

	"github.com/String-xyz/string-api/pkg/internal/common"
	"github.com/lmittmann/w3/module/debug"
)

func getAddressesFromTrace(trace debug.Trace) []string {

	// Map used for optimization
	var addresses map[string]int = make(map[string]int)

	// Slice used for return
	var returnAddresses []string

	// Addresses are a hexadecimal string of 40 characters left-padded by 24 zeros
	re := regexp.MustCompile(`000000000000000000000000[0-9a-fA-F]{40}`)
	for _, line := range trace.StructLogs {
		found := re.FindAllString(hex.EncodeToString(line.Memory), -1)
		for _, addr := range found {
			// Exclude strings with 36 or more prefix zeros, not sure how to do this with regex
			if addr[:36] == "000000000000000000000000000000000000" {
				continue
			}

			// discard padding, add prefix, and checksum
			addr = common.SanitizeChecksum("0x" + addr[24:])

			// check if addresses does not contain addr
			if _, ok := addresses[addr]; !ok {
				addresses[addr] = 1
				returnAddresses = append(returnAddresses, addr)
			}
		}
	}

	return returnAddresses
}
